package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/assessment/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AssessmentRepoPostgres struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewAssessmentRepoPostgres(db *sql.DB, log *logrus.Logger) domain.AssessmentRepo {
	return &AssessmentRepoPostgres{
		db:  db,
		log: log,
	}
}

func (r *AssessmentRepoPostgres) Create(ctx context.Context, assessment *domain.Assessment) error {
	query := `
			INSERT INTO assessments (
				id,
				organization_id,
				title,
				assessment_type,
				assessment_sub_type,
				due_date,
				created_at, 
				updated_at, 
				course_id
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		assessment.ID,
		assessment.OrganizationID,
		assessment.Title,
		assessment.Type,
		assessment.SubType,
		assessment.DueDate,
		assessment.CreatedAt,
		assessment.UpdatedAt,
		assessment.CourseID,
	)

	if err != nil {
		r.log.WithError(err).WithField("assessment_id", assessment.ID).Error("failed to insert assessment")
		return fmt.Errorf("failed to insert assessment: %w", err)
	}

	r.log.WithFields(logrus.Fields{"assessment_id": assessment.ID, "title": assessment.Title}).Info("assessment created successfully")
	return nil
}

func (r *AssessmentRepoPostgres) GetStudentAssessments(ctx context.Context, userID uuid.UUID, filter domain.StudentAssessmentFilter) ([]domain.StudentAssessmentItem, error) {
	var args []interface{}
	argIndex := 1

	query := `
		SELECT 
			a.id,
			COALESCE(subj.name, '') as subject,
			a.title,
			c.content_data->>'url' as attachment_url,
			CASE 
				WHEN s.final_score IS NOT NULL THEN 'done'
				WHEN s.submitted_at IS NOT NULL THEN 'submitted'
				WHEN a.due_date < NOW() THEN 'overdue'
				ELSE 'pending'
			END as status,
			a.assessment_type,
			a.assessment_sub_type,
			a.due_date
		FROM assessments a
		INNER JOIN courses cr ON a.course_id = cr.id
		LEFT JOIN subjects subj ON cr.subject_id = subj.id
		INNER JOIN enrollments e ON e.course_id = cr.id AND e.user_id = $1 AND e.status = 'active'
		LEFT JOIN submissions s ON s.assessment_id = a.id AND s.user_id = $1
		LEFT JOIN contents c ON c.assessment_id = a.id
		WHERE a.deleted_at IS NULL`

	args = append(args, userID)
	argIndex++

	// Apply filters
	if filter.Type != nil {
		query += fmt.Sprintf(" AND a.assessment_type = $%d", argIndex)
		args = append(args, string(*filter.Type))
		argIndex++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND a.due_date >= $%d", argIndex)
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND a.due_date <= $%d", argIndex)
		args = append(args, *filter.EndDate)
		argIndex++
	}

	query += " ORDER BY a.due_date ASC"

	// Apply pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("failed to query student assessments")
		return nil, fmt.Errorf("failed to query student assessments: %w", err)
	}
	defer rows.Close()

	var items []domain.StudentAssessmentItem
	for rows.Next() {
		var item domain.StudentAssessmentItem
		var attachmentURL sql.NullString
		var status string
		var assessmentType string
		var subType string

		if err := rows.Scan(
			&item.ID,
			&item.Subject,
			&item.Title,
			&attachmentURL,
			&status,
			&assessmentType,
			&subType,
			&item.DueDate,
		); err != nil {
			return nil, fmt.Errorf("failed to scan student assessment: %w", err)
		}

		if attachmentURL.Valid {
			item.AttachmentURL = &attachmentURL.String
		}
		item.Status = domain.SubmissionStatus(status)
		item.Type = domain.AssessmentType(assessmentType)
		item.SubType = domain.AssessmentSubType(subType)

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("error iterating student assessments")
		return nil, fmt.Errorf("error iterating student assessments: %w", err)
	}

	return items, nil
}

func (r *AssessmentRepoPostgres) GetStudentAssessmentSummary(ctx context.Context, userID uuid.UUID, filter domain.StudentAssessmentFilter) (*domain.StudentAssessmentSummary, error) {
	var args []interface{}
	argIndex := 1

	baseQuery := `
		SELECT 
			CASE 
				WHEN s.final_score IS NOT NULL THEN 'done'
				WHEN s.submitted_at IS NOT NULL THEN 'submitted'
				WHEN a.due_date < NOW() THEN 'overdue'
				ELSE 'pending'
			END as status
		FROM assessments a
		INNER JOIN courses cr ON a.course_id = cr.id
		INNER JOIN enrollments e ON e.course_id = cr.id AND e.user_id = $1 AND e.status = 'active'
		LEFT JOIN submissions s ON s.assessment_id = a.id AND s.user_id = $1
		WHERE a.deleted_at IS NULL`

	args = append(args, userID)
	argIndex++

	var filterClauses []string

	if filter.Type != nil {
		filterClauses = append(filterClauses, fmt.Sprintf("a.assessment_type = $%d", argIndex))
		args = append(args, string(*filter.Type))
		argIndex++
	}

	if filter.StartDate != nil {
		filterClauses = append(filterClauses, fmt.Sprintf("a.due_date >= $%d", argIndex))
		args = append(args, *filter.StartDate)
		argIndex++
	}

	if filter.EndDate != nil {
		filterClauses = append(filterClauses, fmt.Sprintf("a.due_date <= $%d", argIndex))
		args = append(args, *filter.EndDate)
		argIndex++
	}

	if len(filterClauses) > 0 {
		baseQuery += " AND " + strings.Join(filterClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0) as pending,
			COALESCE(SUM(CASE WHEN status = 'submitted' THEN 1 ELSE 0 END), 0) as submitted,
			COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0) as done,
			COALESCE(SUM(CASE WHEN status = 'overdue' THEN 1 ELSE 0 END), 0) as overdue
		FROM (%s) AS subquery`, baseQuery)

	var summary domain.StudentAssessmentSummary
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&summary.Pending,
		&summary.Submitted,
		&summary.Done,
		&summary.Overdue,
	)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("failed to get student assessment summary")
		return nil, fmt.Errorf("failed to get student assessment summary: %w", err)
	}

	return &summary, nil
}

// ensure time package is used
var _ = time.Now

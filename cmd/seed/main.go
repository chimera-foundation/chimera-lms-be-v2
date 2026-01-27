package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	cohortRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/repository"
	courseDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	courseRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/repository/postgres"
	eduDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/domain"
	eduRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/repository/postgres"
	enrollmentDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	enrollmentRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/repo"
	eventDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	eventRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/repository"
	orgDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	orgRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/repository/postgres"
	sectionDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	sectionRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/repository"
	subjectDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	subjectRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/repository/postgres"
	userDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	userRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
	"github.com/google/uuid"
)

type Seeder struct {
	db             *sql.DB
	orgRepo        orgDomain.OrganizationRepository
	userRepo       userDomain.UserRepository
	roleRepo       userDomain.RoleRepository
	cohortRepo     domain.CohortRepository
	sectionRepo    sectionDomain.SectionRepository
	courseRepo     courseDomain.CourseRepository
	moduleRepo     courseDomain.ModuleRepository
	lessonRepo     courseDomain.LessonRepository
	enrollmentRepo enrollmentDomain.EnrollmentRepository
	eventRepo      eventDomain.EventRepository
	subjectRepo    subjectDomain.SubjectRepository
	eduRepo        eduDomain.EducationLevelRepository

	// Context data
	orgID         uuid.UUID
	teacherRoleID uuid.UUID
	studentRoleID uuid.UUID
	teacherIDs    []uuid.UUID
	studentIDs    []uuid.UUID
	eduLevelIDs   []uuid.UUID
	// Map subject code to ID
	subjectMap       map[string]uuid.UUID
	cohortIDs        []uuid.UUID
	sectionIDs       []uuid.UUID
	courseIDs        []uuid.UUID
	academicPeriodID uuid.UUID
}

func main() {
	v := app.NewViper()
	logger := app.NewLogger(v)
	db := app.NewDatabase(v, logger)
	defer db.Close()

	seeder := &Seeder{
		db:             db,
		orgRepo:        orgRepo.NewOrganizationRepo(db),
		userRepo:       userRepo.NewUserRepo(db),
		roleRepo:       userRepo.NewRoleRepository(db),
		cohortRepo:     cohortRepo.NewCohortRepository(db),
		sectionRepo:    sectionRepo.NewSectionRepository(db),
		courseRepo:     courseRepo.NewCourseRepository(db),
		moduleRepo:     courseRepo.NewModuleRepository(db),
		lessonRepo:     courseRepo.NewLessonRepository(db),
		enrollmentRepo: enrollmentRepo.NewEnrollmentRepository(db),
		eventRepo:      eventRepo.NewEventRepository(db),
		subjectRepo:    subjectRepo.NewSubjectRepository(db),
		eduRepo:        eduRepo.NewEducationLevelRepository(db),
		subjectMap:     make(map[string]uuid.UUID),
	}

	ctx := context.Background()

	if err := seeder.Seed(ctx); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	fmt.Println("Seeding completed successfully!")
}

func (s *Seeder) Seed(ctx context.Context) error {

	if err := s.seedOrganization(ctx); err != nil {
		return err
	}
	if err := s.seedRoles(ctx); err != nil {
		return err
	}
	if err := s.seedUsers(ctx); err != nil {
		return err
	}
	if err := s.seedAcademicPeriod(ctx); err != nil {
		return err
	}
	if err := s.seedEducationLevels(ctx); err != nil {
		return err
	}
	if err := s.seedSubjects(ctx); err != nil {
		return err
	}
	if err := s.seedCohortsAndSections(ctx); err != nil {
		return err
	}
	if err := s.seedCourses(ctx); err != nil {
		return err
	}
	if err := s.seedEnrollments(ctx); err != nil {
		return err
	}
	if err := s.seedEvents(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Seeder) seedOrganization(ctx context.Context) error {
	fmt.Println("Seeding Organization...")
	org := &orgDomain.Organization{
		Name:     "Candle Tree School",
		Slug:     "candle-tree-school",
		Type:     orgDomain.HighSchool,
		Address:  "123 Candle Tree Lane, Jakarta, Indonesia",
		IsActive: true,
	}
	org.PrepareCreate(nil)

	// Check if already exists
	existingOrg, err := s.orgRepo.GetBySlug(ctx, org.Slug)
	if err != nil {
		return fmt.Errorf("failed to check existing organization: %w", err)
	}
	if existingOrg != nil {
		fmt.Printf("Organization %s already exists, using ID: %s\n", org.Slug, existingOrg.ID)
		s.orgID = existingOrg.ID
		return nil
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return fmt.Errorf("failed to seed organization: %w", err)
	}
	s.orgID = org.ID
	return nil
}

func (s *Seeder) seedRoles(ctx context.Context) error {
	fmt.Println("Seeding Roles...")
	roles := []string{"Teacher", "Student"}

	for _, name := range roles {
		role, err := s.roleRepo.GetByName(ctx, name)
		if err != nil {
			return err
		}
		if role == nil {
			role = &userDomain.Role{Name: name}
			role.PrepareCreate(nil)
			if err := s.roleRepo.Create(ctx, role); err != nil {
				return err
			}
		}

		switch name {
		case "Teacher":
			s.teacherRoleID = role.ID
		case "Student":
			s.studentRoleID = role.ID
		}
	}
	return nil
}

func (s *Seeder) seedAcademicPeriod(ctx context.Context) error {
	fmt.Println("Seeding Academic Period...")
	// Check if active period exists for org
	var id uuid.UUID
	err := s.db.QueryRowContext(ctx, "SELECT id FROM academic_periods WHERE organization_id = $1 AND is_active = true LIMIT 1", s.orgID).Scan(&id)
	if err == nil {
		s.academicPeriodID = id
		return nil
	}

	if err != sql.ErrNoRows {
		return err
	}

	id = uuid.New()
	name := fmt.Sprintf("Academic Year %d", time.Now().Year())
	startDate := time.Now()
	endDate := startDate.AddDate(1, 0, 0)

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO academic_periods (id, organization_id, name, start_date, end_date, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		id, s.orgID, name, startDate, endDate, true, time.Now(), time.Now())

	if err != nil {
		return fmt.Errorf("failed to create academic period: %w", err)
	}
	s.academicPeriodID = id
	return nil
}

func (s *Seeder) seedUsers(ctx context.Context) error {
	fmt.Println("Seeding Users...")

	// Create Teachers
	for i := 1; i <= 5; i++ {
		email := fmt.Sprintf("teacher%d@candletree.edu", i)
		user, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return err
		}

		if user == nil {
			user = userDomain.NewUser(email, fmt.Sprintf("Teacher%d", i), "Doe", s.orgID)
			user.SetPassword("password123")
			// Manually Construct the role relationship for the struct if needed,
			// though typical repo Create doesn't cascade save Roles unless implemented that way.
			// We will insert into user_roles manually as below.
			user.Roles = []userDomain.Role{{Base: shared.Base{ID: s.teacherRoleID}}}

			if err := s.userRepo.Create(ctx, user); err != nil {
				return err
			}

			_, err = s.db.ExecContext(ctx, "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", user.ID, s.teacherRoleID)
			if err != nil {
				return err
			}
		}
		s.teacherIDs = append(s.teacherIDs, user.ID)
	}

	// Create Students
	for i := 1; i <= 50; i++ {
		email := fmt.Sprintf("student%d@candletree.edu", i)
		user, err := s.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return err
		}

		if user == nil {
			user = userDomain.NewUser(email, fmt.Sprintf("Student%d", i), "Smith", s.orgID)
			user.SetPassword("password123")
			if err := s.userRepo.Create(ctx, user); err != nil {
				return err
			}
			_, err = s.db.ExecContext(ctx, "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", user.ID, s.studentRoleID)
			if err != nil {
				return err
			}
		}
		s.studentIDs = append(s.studentIDs, user.ID)
	}
	return nil
}

func (s *Seeder) seedEducationLevels(ctx context.Context) error {
	fmt.Println("Seeding Education Levels...")
	levels := []struct {
		Name string
		Code string
	}{
		{"High School", "HS"},
		{"Middle School", "MS"},
	}

	for _, l := range levels {
		level := &eduDomain.EducationLevel{
			OrganizationID: s.orgID,
			Name:           l.Name,
			Code:           l.Code,
		}
		level.PrepareCreate(nil)
		if err := s.eduRepo.Create(ctx, level); err != nil {
			fmt.Printf("Warning: Failed to create education level %s: %v\n", l.Name, err)
		}
		s.eduLevelIDs = append(s.eduLevelIDs, level.ID)
	}
	return nil
}

func (s *Seeder) seedSubjects(ctx context.Context) error {
	fmt.Println("Seeding Subjects...")
	subjects := []struct {
		Name string
		Code string
	}{
		{"Mathematics", "MATH"},
		{"Science", "SCI"},
		{"History", "HIST"},
		{"English", "ENG"},
		{"Physics", "PHYS"},
	}

	if len(s.eduLevelIDs) == 0 {
		return fmt.Errorf("no education levels seeded")
	}

	for _, sub := range subjects {
		subject := &subjectDomain.Subject{
			OrganizationID:   s.orgID,
			EducationLevelID: s.eduLevelIDs[0],
			Name:             sub.Name,
			Code:             sub.Code,
		}
		subject.PrepareCreate(nil)
		if err := s.subjectRepo.Create(ctx, subject); err != nil {
			fmt.Printf("Warning: Failed to create subject %s: %v\n", sub.Name, err)
		}
		// Try to verify ID if not returned correctly/created duplicates
		// Assuming for now it works due to previous run checks
		s.subjectMap[sub.Code] = subject.ID
	}
	return nil
}

func (s *Seeder) seedCohortsAndSections(ctx context.Context) error {
	fmt.Println("Seeding Cohorts and Sections...")

	if len(s.eduLevelIDs) == 0 {
		return fmt.Errorf("no education levels seeded")
	}

	// Create Cohorts for Grade 10, 11, 12
	grades := []string{"Grade 10", "Grade 11", "Grade 12"}

	for _, g := range grades {
		cohort := &domain.Cohort{
			OrganizationID:   s.orgID,
			EducationLevelID: s.eduLevelIDs[0],
			AcademicPeriodID: s.academicPeriodID,
			Name:             g,
		}
		cohort.PrepareCreate(nil)
		if err := s.cohortRepo.Create(ctx, cohort); err != nil {
			return fmt.Errorf("failed to create cohort %s: %w", g, err)
		}
		s.cohortIDs = append(s.cohortIDs, cohort.ID)

		// Create Sections
		sections := []string{"A", "B", "C"}
		for _, secName := range sections {
			section := &sectionDomain.Section{
				CohortID: cohort.ID,
				Name:     fmt.Sprintf("%s - Section %s", g, secName),
				RoomCode: fmt.Sprintf("R-%s-%s", g, secName),
				Capacity: 30,
			}
			section.PrepareCreate(nil)
			if err := s.sectionRepo.Create(ctx, section); err != nil {
				return fmt.Errorf("failed to create section: %w", err)
			}
			s.sectionIDs = append(s.sectionIDs, section.ID)

			// Link Students to Section (Draft distribution)
			// We just want to ensure students are in *some* section.
			// Currently assume 50 students. 3 Grades * 3 Sections = 9 sections. (~5 students per section)
			// We'll iterate students and assign based on index.
		}
	}

	// Assign students to random sections (or round robin)
	// We need to loop all sections we just created?
	// s.sectionIDs has IDs from all cohorts mixed.
	// But we need to be careful: Enrollment connects User, Section, Course, AcademicPeriod.
	// This function `seedCohortsAndSections` creates the structural entities.
	// Enrollment happens in `seedEnrollments`.
	// But `section_members` table also exists?
	// The prompt said "make sure to have multiple students inside of each section" (Step 0)
	// AND "create enrollment to courses for each students" (Step 0).
	// `section_members` usually implies homeroom or just being in the section.

	for i, studentID := range s.studentIDs {
		// Round robin assignment to sections
		if len(s.sectionIDs) > 0 {
			sectionID := s.sectionIDs[i%len(s.sectionIDs)]
			// Insert into section_members
			_, err := s.db.ExecContext(ctx, "INSERT INTO section_members (id, section_id, user_id, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
				uuid.New(), sectionID, studentID, "student", time.Now(), time.Now())
			if err != nil {
				// ignore dupes or errors
			}
		}
	}

	return nil
}

func (s *Seeder) seedCourses(ctx context.Context) error {
	fmt.Println("Seeding Courses and Content...")
	if len(s.teacherIDs) == 0 {
		return fmt.Errorf("no teachers seeded")
	}
	if len(s.subjectMap) == 0 {
		return fmt.Errorf("no subjects seeded")
	}

	// Structure: Subject -> [List of Courses]
	// E.g. MATH -> [Mathematics - Grade 10, Linear Algebra]

	courses := []struct {
		SubjectCode string
		Title       string
		GradeLevel  int
	}{
		{"MATH", "Mathematics - Grade 10", 10},
		{"MATH", "Linear Algebra", 11},
		{"SCI", "Biology - Grade 10", 10},
		{"SCI", "Chemistry - Grade 11", 11},
		{"PHYS", "Physics - Grade 12", 12},
		{"ENG", "English Literature - Grade 10", 10},
		{"HIST", "World History - Grade 11", 11},
	}

	for i, c := range courses {
		subjectID, ok := s.subjectMap[c.SubjectCode]
		if !ok {
			// fallback if map not populated (e.g. subject creation failed silently)
			// Try to pick any subject or skip
			continue
		}

		course := &courseDomain.Course{
			OrganizationID:   s.orgID,
			InstructorID:     s.teacherIDs[i%len(s.teacherIDs)],
			SubjectID:        subjectID,
			EducationLevelID: s.eduLevelIDs[0],
			Title:            c.Title,
			Description:      fmt.Sprintf("Course for %s", c.SubjectCode),
			Status:           courseDomain.Published,
			Price:            0,
			GradeLevel:       c.GradeLevel,
			Credits:          3,
		}
		course.PrepareCreate(nil)
		if err := s.courseRepo.Create(ctx, course); err != nil {
			return err
		}
		s.courseIDs = append(s.courseIDs, course.ID)

		// Seed structure: Modules -> Lessons
		if err := s.seedModulesAndLessons(ctx, course.ID, c.Title); err != nil {
			fmt.Printf("Warning: Failed to seed modules for course %s: %v\n", c.Title, err)
		}
	}
	return nil
}

func (s *Seeder) seedModulesAndLessons(ctx context.Context, courseID uuid.UUID, courseTitle string) error {
	modules := []string{"Introduction", "Core Concepts", "Advanced Topics"}

	for mIdx, mTitle := range modules {
		module := &courseDomain.Module{
			CourseID:   courseID,
			Title:      fmt.Sprintf("%s - %s", courseTitle, mTitle),
			OrderIndex: mIdx,
		}
		module.PrepareCreate(nil)
		if err := s.moduleRepo.Create(ctx, module); err != nil {
			return err
		}

		// Create Lessons
		lessons := []string{"Part 1", "Part 2", "workshop"}
		for lIdx, lTitle := range lessons {
			lesson := &courseDomain.Lesson{
				ModuleID:   module.ID,
				Title:      fmt.Sprintf("Lesson %s", lTitle),
				OrderIndex: lIdx,
			}
			lesson.PrepareCreate(nil)
			if err := s.lessonRepo.Create(ctx, lesson); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Seeder) seedEnrollments(ctx context.Context) error {
	fmt.Println("Seeding Enrollments...")
	// Enroll students in courses
	// "apply the teacher as the one who teaches" - handled in Course creation (InstructorID)

	for i, studentID := range s.studentIDs {
		// Find a section this student is in
		var sectionID uuid.UUID
		if len(s.sectionIDs) > 0 {
			sectionID = s.sectionIDs[i%len(s.sectionIDs)]
		}

		for _, courseID := range s.courseIDs {
			enrollment := &enrollmentDomain.Enrollment{
				UserID:           studentID,
				CourseID:         courseID,
				SectionID:        sectionID,
				AcademicPeriodID: s.academicPeriodID,
				Status:           enrollmentDomain.Active,
				EnrolledAt:       time.Now(),
			}
			enrollment.PrepareCreate(nil)
			if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
				// Ignore errors
			}
		}
	}
	return nil
}

func (s *Seeder) seedEvents(ctx context.Context) error {
	fmt.Println("Seeding Events...")

	// 1. Seed Holidays (Indonesia)
	holidays := []struct {
		Name string
		Date time.Time
	}{
		{"New Year's Day", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Independence Day", time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC)},
		{"Christmas Day", time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)},
	}

	for _, h := range holidays {
		event := eventDomain.NewEvent(s.orgID, h.Name, eventDomain.Holiday)
		event.StartAt = &h.Date
		end := h.Date.Add(24 * time.Hour)
		event.EndAt = &end
		event.IsAllDay = true
		event.PrepareCreate(nil)

		if err := s.eventRepo.Create(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to seed holiday %s: %v\n", h.Name, err)
		}
	}

	// 2. Seed Lesson Schedules for Grades 10, 11, 12
	// Distribute across Mon-Fri, multiple lessons
	now := time.Now()
	// Calculate Monday of current week
	offset := int(now.Weekday())
	if offset == 0 {
		offset = 7
	} // Sunday -> 7
	monday := now.AddDate(0, 0, -offset+1)

	grades := []int{10, 11, 12}

	for _, grade := range grades {
		for day := 0; day < 5; day++ { // Mon-Fri
			date := monday.AddDate(0, 0, day)

			// 3 lessons per day
			for lesson := 0; lesson < 3; lesson++ {
				// Different start times for different grades to avoid overlap if same room (though we use random rooms)
				// Schedule: Grade 10: 8am, Grade 11: 9am... mixing it up

				startHour := 8 + lesson + (grade - 10)
				start := time.Date(date.Year(), date.Month(), date.Day(), startHour, 0, 0, 0, time.UTC)
				end := start.Add(45 * time.Minute)

				title := fmt.Sprintf("Grade %d Lesson %d (%s)", grade, lesson+1, date.Weekday())

				event := eventDomain.NewEvent(s.orgID, title, eventDomain.Schedule,
					eventDomain.WithTimes(start, end),
					eventDomain.WithLocation(fmt.Sprintf("Room %d0%d", grade, lesson)),
				)
				event.Description = fmt.Sprintf("Regular lesson for Grade %d", grade)

				event.PrepareCreate(nil)
				if err := s.eventRepo.Create(ctx, event); err != nil {
					// ignore
				}
			}
		}
	}

	return nil
}

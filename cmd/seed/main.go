package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
	cohortDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	cohortRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/repository/postgres"
	contentDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/content/domain"
	contentRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/content/repository/postgres"
	courseDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	courseRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/repository/postgres"
	eduDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/domain"
	eduRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/repository/postgres"
	enrollmentDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/domain"
	enrollmentRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/repository/postgres"
	eventDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/domain"
	eventRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/repository/postgres"
	orgDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/domain"
	orgRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/repository/postgres"
	programDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/domain"
	programRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/repository/postgres"
	progressDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/progress_tracker/domain"
	progressRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/progress_tracker/repository/postgres"
	sectionDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	sectionRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/repository/postgres"
	subjectDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	subjectRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/repository/postgres"
	userDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	userRepo "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
	"github.com/google/uuid"
)

type Seeder struct {
	// Repositories
	orgRepo            orgDomain.OrganizationRepository
	academicPeriodRepo orgDomain.AcademicPeriodRepository
	userRepo           userDomain.UserRepository
	roleRepo           userDomain.RoleRepository
	cohortRepo         cohortDomain.CohortRepository
	cohortMemberRepo   cohortDomain.CohortMemberRepository
	sectionRepo        sectionDomain.SectionRepository
	sectionMemberRepo  sectionDomain.SectionMemberRepository
	courseRepo         courseDomain.CourseRepository
	moduleRepo         courseDomain.ModuleRepository
	lessonRepo         courseDomain.LessonRepository
	enrollmentRepo     enrollmentDomain.EnrollmentRepository
	eventRepo          eventDomain.EventRepository
	subjectRepo        subjectDomain.SubjectRepository
	eduRepo            eduDomain.EducationLevelRepository
	programRepo        programDomain.ProgramRepository
	programCourseRepo  programDomain.ProgramCourseRepository
	progressRepo       progressDomain.ProgressTrackerRepository
	contentRepo        contentDomain.ContentRepository

	// Context data
	orgID            uuid.UUID
	adminRoleID      uuid.UUID
	teacherRoleID    uuid.UUID
	studentRoleID    uuid.UUID
	adminIDs         []uuid.UUID
	teacherIDs       []uuid.UUID
	studentIDs       []uuid.UUID
	eduLevelIDs      []uuid.UUID
	subjectMap       map[string]uuid.UUID
	cohortIDs        []uuid.UUID
	sectionIDs       []uuid.UUID
	courseIDs        []uuid.UUID
	lessonIDs        []uuid.UUID
	enrollmentIDs    []uuid.UUID
	programIDs       []uuid.UUID
	academicPeriodID uuid.UUID
	contentIDs       []uuid.UUID
}

func main() {
	v := app.NewViper()
	logger := app.NewLogger(v)
	db := app.NewDatabase(v, logger)
	defer db.Close()

	seeder := &Seeder{
		orgRepo:            orgRepo.NewOrganizationRepo(db),
		academicPeriodRepo: orgRepo.NewAcademicPeriodRepository(db),
		userRepo:           userRepo.NewUserRepo(db),
		roleRepo:           userRepo.NewRoleRepository(db),
		cohortRepo:         cohortRepo.NewCohortRepository(db),
		cohortMemberRepo:   cohortRepo.NewCohortMemberRepository(db),
		sectionRepo:        sectionRepo.NewSectionRepository(db),
		sectionMemberRepo:  sectionRepo.NewSectionMemberRepository(db),
		courseRepo:         courseRepo.NewCourseRepository(db),
		moduleRepo:         courseRepo.NewModuleRepository(db),
		lessonRepo:         courseRepo.NewLessonRepository(db),
		enrollmentRepo:     enrollmentRepo.NewEnrollmentRepository(db),
		eventRepo:          eventRepo.NewEventRepository(db),
		subjectRepo:        subjectRepo.NewSubjectRepository(db),
		eduRepo:            eduRepo.NewEducationLevelRepository(db),
		programRepo:        programRepo.NewProgramRepository(db),
		programCourseRepo:  programRepo.NewProgramCourseRepository(db),
		progressRepo:       progressRepo.NewProgressTrackerRepository(db),
		contentRepo:        contentRepo.NewContentRepository(db),
		subjectMap:         make(map[string]uuid.UUID),
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
	if err := s.seedPrograms(ctx); err != nil {
		return err
	}
	if err := s.seedCourses(ctx); err != nil {
		return err
	}
	if err := s.seedEnrollments(ctx); err != nil {
		return err
	}
	if err := s.seedProgressTrackers(ctx); err != nil {
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
	roles := []string{"Admin", "Teacher", "Student"}

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
		case "Admin":
			s.adminRoleID = role.ID
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

	existing, err := s.academicPeriodRepo.GetActiveByOrganizationID(ctx, s.orgID)
	if err != nil {
		return err
	}
	if existing != nil {
		s.academicPeriodID = existing.ID
		return nil
	}

	period := &orgDomain.AcademicPeriod{
		Name:      fmt.Sprintf("Academic Year %d/%d", time.Now().Year(), time.Now().Year()+1),
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0),
		IsActive:  true,
	}

	if err := s.academicPeriodRepo.Create(ctx, period, s.orgID); err != nil {
		return fmt.Errorf("failed to create academic period: %w", err)
	}
	s.academicPeriodID = period.ID
	return nil
}

func (s *Seeder) seedUsers(ctx context.Context) error {
	fmt.Println("Seeding Users...")

	// Create Admin
	adminEmail := "admin@candletree.edu"
	admin, err := s.userRepo.GetByEmail(ctx, adminEmail)
	if err != nil {
		return err
	}
	if admin == nil {
		admin = userDomain.NewUser(adminEmail, "Admin", "User", s.orgID)
		admin.SetPassword("password123")
		if err := s.userRepo.Create(ctx, admin); err != nil {
			return err
		}
		if err := s.roleRepo.AssignRoleToUser(ctx, admin.ID, s.adminRoleID); err != nil {
			return err
		}
	}
	s.adminIDs = append(s.adminIDs, admin.ID)

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
			if err := s.userRepo.Create(ctx, user); err != nil {
				return err
			}
			if err := s.roleRepo.AssignRoleToUser(ctx, user.ID, s.teacherRoleID); err != nil {
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
			if err := s.roleRepo.AssignRoleToUser(ctx, user.ID, s.studentRoleID); err != nil {
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
		{"Grade School", "GS"},
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
		{"Chemistry", "CHEM"},
		{"Biology", "BIO"},
		{"Geography", "GEO"},
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
		s.subjectMap[sub.Code] = subject.ID
	}
	return nil
}

func (s *Seeder) seedCohortsAndSections(ctx context.Context) error {
	fmt.Println("Seeding Cohorts and Sections...")

	if len(s.eduLevelIDs) == 0 {
		return fmt.Errorf("no education levels seeded")
	}

	grades := []string{"Grade 10", "Grade 11", "Grade 12"}

	sectionIdx := 0
	for cohortIdx, g := range grades {
		cohort := &cohortDomain.Cohort{
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

		// Add students to cohort (distribute evenly)
		studentsPerCohort := len(s.studentIDs) / len(grades)
		startIdx := cohortIdx * studentsPerCohort
		endIdx := startIdx + studentsPerCohort
		if cohortIdx == len(grades)-1 {
			endIdx = len(s.studentIDs)
		}

		for i := startIdx; i < endIdx; i++ {
			member := &cohortDomain.CohortMember{
				CohortID: cohort.ID,
				UserID:   s.studentIDs[i],
			}
			if err := s.cohortMemberRepo.Create(ctx, member); err != nil {
				fmt.Printf("Warning: Failed to add student to cohort: %v\n", err)
			}
		}

		// Create Sections with different role types
		sectionRoles := []sectionDomain.SectionRoleType{
			sectionDomain.Student,
			sectionDomain.Teacher,
			sectionDomain.Assistant,
			sectionDomain.Monitor,
		}
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

			// Assign a teacher to the section
			if sectionIdx < len(s.teacherIDs) {
				teacherMember := &sectionDomain.SectionMember{
					SectionID: section.ID,
					UserID:    s.teacherIDs[sectionIdx%len(s.teacherIDs)],
					RoleType:  sectionDomain.Teacher,
				}
				if err := s.sectionMemberRepo.Create(ctx, teacherMember); err != nil {
					fmt.Printf("Warning: Failed to add teacher to section: %v\n", err)
				}
			}

			// Assign students with varying roles
			for i := startIdx; i < endIdx; i++ {
				roleType := sectionRoles[0] // Default to student
				if i == startIdx {
					roleType = sectionDomain.Monitor // First student is monitor
				} else if i == startIdx+1 && startIdx+1 < endIdx {
					roleType = sectionDomain.Assistant // Second is assistant
				}

				member := &sectionDomain.SectionMember{
					SectionID: section.ID,
					UserID:    s.studentIDs[i],
					RoleType:  roleType,
				}
				if err := s.sectionMemberRepo.Create(ctx, member); err != nil {
					// Ignore duplicates
				}
			}
			sectionIdx++
		}
	}

	return nil
}

func (s *Seeder) seedPrograms(ctx context.Context) error {
	fmt.Println("Seeding Programs...")

	programs := []struct {
		Name        string
		Description string
	}{
		{"Science Track", "Program for students pursuing science and technology"},
		{"Social Studies Track", "Program for students interested in humanities and social sciences"},
		{"General Education", "Comprehensive education program covering all subjects"},
	}

	for _, p := range programs {
		program := &programDomain.Program{
			OrganizationID: s.orgID,
			Name:           p.Name,
			Description:    p.Description,
		}
		if err := s.programRepo.Create(ctx, program); err != nil {
			return fmt.Errorf("failed to create program %s: %w", p.Name, err)
		}
		s.programIDs = append(s.programIDs, program.ID)
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

	// Courses with different statuses
	courses := []struct {
		SubjectCode string
		Title       string
		GradeLevel  int
		Status      courseDomain.CourseStatus
	}{
		{"MATH", "Mathematics - Grade 10", 10, courseDomain.Published},
		{"MATH", "Linear Algebra", 11, courseDomain.Published},
		{"MATH", "Calculus", 12, courseDomain.Draft},
		{"SCI", "Biology - Grade 10", 10, courseDomain.Published},
		{"SCI", "Chemistry - Grade 11", 11, courseDomain.Published},
		{"PHYS", "Physics - Grade 12", 12, courseDomain.Published},
		{"ENG", "English Literature - Grade 10", 10, courseDomain.Published},
		{"ENG", "Advanced Writing", 11, courseDomain.Archived},
		{"HIST", "World History - Grade 11", 11, courseDomain.Published},
		{"GEO", "Geography - Grade 10", 10, courseDomain.Draft},
	}

	for i, c := range courses {
		subjectID, ok := s.subjectMap[c.SubjectCode]
		if !ok {
			continue
		}

		course := &courseDomain.Course{
			OrganizationID:   s.orgID,
			InstructorID:     s.teacherIDs[i%len(s.teacherIDs)],
			SubjectID:        subjectID,
			EducationLevelID: s.eduLevelIDs[0],
			Title:            c.Title,
			Description:      fmt.Sprintf("Comprehensive course for %s", c.SubjectCode),
			Status:           c.Status,
			Price:            0,
			GradeLevel:       c.GradeLevel,
			Credits:          3,
		}
		course.PrepareCreate(nil)
		if err := s.courseRepo.Create(ctx, course); err != nil {
			return err
		}
		s.courseIDs = append(s.courseIDs, course.ID)

		// Seed modules and lessons
		lessonIDs, err := s.seedModulesAndLessons(ctx, course.ID, c.Title)
		if err != nil {
			fmt.Printf("Warning: Failed to seed modules for course %s: %v\n", c.Title, err)
		}
		s.lessonIDs = append(s.lessonIDs, lessonIDs...)

		// Link courses to programs
		if len(s.programIDs) > 0 {
			programIdx := i % len(s.programIDs)
			pc := &programDomain.ProgramCourse{
				ProgramID:  s.programIDs[programIdx],
				CourseID:   course.ID,
				OrderIndex: i,
			}
			if err := s.programCourseRepo.Create(ctx, pc); err != nil {
				fmt.Printf("Warning: Failed to link course to program: %v\n", err)
			}
		}
	}
	return nil
}

func (s *Seeder) seedModulesAndLessons(ctx context.Context, courseID uuid.UUID, courseTitle string) ([]uuid.UUID, error) {
	modules := []string{"Introduction", "Core Concepts", "Advanced Topics", "Practical Applications"}
	var lessonIDs []uuid.UUID

	contentTypes := []contentDomain.ContentType{
		contentDomain.Video,
		contentDomain.Document,
		contentDomain.Quiz,
	}

	for mIdx, mTitle := range modules {
		module := &courseDomain.Module{
			CourseID:   courseID,
			Title:      fmt.Sprintf("%s - %s", courseTitle, mTitle),
			OrderIndex: mIdx,
		}
		module.PrepareCreate(nil)
		if err := s.moduleRepo.Create(ctx, module); err != nil {
			return nil, err
		}

		lessons := []string{"Part 1", "Part 2", "Workshop", "Quiz"}
		for lIdx, lTitle := range lessons {
			lesson := &courseDomain.Lesson{
				ModuleID:   module.ID,
				Title:      fmt.Sprintf("Lesson %s", lTitle),
				OrderIndex: lIdx,
			}
			lesson.PrepareCreate(nil)
			if err := s.lessonRepo.Create(ctx, lesson); err != nil {
				return nil, err
			}
			lessonIDs = append(lessonIDs, lesson.ID)

			// Create content for each lesson (one of each type per lesson cycle)
			contentType := contentTypes[lIdx%len(contentTypes)]
			content := &contentDomain.Content{
				LessonID: lesson.ID,
				Type:     contentType,
			}
			content.PrepareCreate(nil)
			if err := s.contentRepo.Create(ctx, content); err != nil {
				fmt.Printf("Warning: Failed to create content for lesson %s: %v\n", lTitle, err)
			} else {
				s.contentIDs = append(s.contentIDs, content.ID)
			}
		}
	}
	return lessonIDs, nil

}

func (s *Seeder) seedEnrollments(ctx context.Context) error {
	fmt.Println("Seeding Enrollments...")

	// Enroll students with different statuses
	statuses := []enrollmentDomain.EnrollmentStatus{
		enrollmentDomain.Active,
		enrollmentDomain.Completed,
		enrollmentDomain.Dropped,
	}

	for i, studentID := range s.studentIDs {
		var sectionID uuid.UUID
		if len(s.sectionIDs) > 0 {
			sectionID = s.sectionIDs[i%len(s.sectionIDs)]
		}

		for j, courseID := range s.courseIDs {
			// Vary enrollment status
			status := statuses[0] // Most are active
			if j == len(s.courseIDs)-1 && i%10 == 0 {
				status = statuses[1] // Some completed
			} else if j == len(s.courseIDs)-1 && i%15 == 0 {
				status = statuses[2] // Some dropped
			}

			enrollment := &enrollmentDomain.Enrollment{
				UserID:           studentID,
				CourseID:         courseID,
				SectionID:        sectionID,
				AcademicPeriodID: s.academicPeriodID,
				Status:           status,
				EnrolledAt:       time.Now(),
			}
			enrollment.PrepareCreate(nil)
			if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
				// Ignore errors (duplicates)
				continue
			}
			s.enrollmentIDs = append(s.enrollmentIDs, enrollment.ID)
		}
	}
	return nil
}

func (s *Seeder) seedProgressTrackers(ctx context.Context) error {
	fmt.Println("Seeding Progress Trackers...")

	if len(s.enrollmentIDs) == 0 {
		fmt.Println("No enrollments to track progress (enrollments may already exist from previous run)")
		return nil
	}
	if len(s.contentIDs) == 0 {
		fmt.Println("No content to track progress (content may already exist from previous run)")
		return nil
	}

	fmt.Printf("Creating progress trackers for %d enrollments and %d content items...\n", len(s.enrollmentIDs), len(s.contentIDs))

	// Create progress for some enrollments
	created := 0
	for i := 0; i < len(s.enrollmentIDs) && i < 100; i++ {
		for j := 0; j < len(s.contentIDs) && j < 5; j++ {
			tracker := &progressDomain.ProgressTracker{
				EnrollmentID: s.enrollmentIDs[i],
				ContentID:    s.contentIDs[j],
				IsCompleted:  i%3 == 0, // Some completed
				UpdatedAt:    time.Now(),
			}
			if err := s.progressRepo.Create(ctx, tracker); err != nil {
				if i == 0 && j == 0 {
					fmt.Printf("Warning: First progress tracker creation failed: %v\n", err)
				}
				// Ignore subsequent errors (likely duplicates)
			} else {
				created++
			}
		}
	}
	fmt.Printf("Created %d progress tracker records\n", created)
	return nil
}

func (s *Seeder) seedEvents(ctx context.Context) error {
	fmt.Println("Seeding Events...")

	// 1. Holidays (EventType: Holiday) - Indonesian National Holidays 2025 & 2026
	holidays := []struct {
		Name string
		Date time.Time
	}{
		// 2025 Indonesian National Holidays
		{"New Year's Day 2025", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Isra Mi'raj 2025", time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)},
		{"Chinese New Year 2025", time.Date(2025, 1, 29, 0, 0, 0, 0, time.UTC)},
		{"Nyepi (Balinese Day of Silence) 2025", time.Date(2025, 3, 29, 0, 0, 0, 0, time.UTC)},
		{"Eid al-Fitr 2025 (Day 1)", time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)},
		{"Eid al-Fitr 2025 (Day 2)", time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)},
		{"Good Friday 2025", time.Date(2025, 4, 18, 0, 0, 0, 0, time.UTC)},
		{"Labour Day 2025", time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC)},
		{"Waisak Day (Buddha's Birthday) 2025", time.Date(2025, 5, 12, 0, 0, 0, 0, time.UTC)},
		{"Ascension of Jesus Christ 2025", time.Date(2025, 5, 29, 0, 0, 0, 0, time.UTC)},
		{"Pancasila Day 2025", time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)},
		{"Eid al-Adha 2025", time.Date(2025, 6, 7, 0, 0, 0, 0, time.UTC)},
		{"Islamic New Year 1447 H", time.Date(2025, 6, 27, 0, 0, 0, 0, time.UTC)},
		{"Independence Day 2025", time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC)},
		{"Prophet Muhammad's Birthday 2025", time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC)},
		{"Christmas Day 2025", time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)},
		// 2026 Indonesian National Holidays
		{"New Year's Day 2026", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Isra Mi'raj 2026", time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC)},
		{"Chinese New Year 2026", time.Date(2026, 2, 17, 0, 0, 0, 0, time.UTC)},
		{"Nyepi (Balinese Day of Silence) 2026", time.Date(2026, 3, 19, 0, 0, 0, 0, time.UTC)},
		{"Eid al-Fitr 2026 (Day 1)", time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)},
		{"Eid al-Fitr 2026 (Day 2)", time.Date(2026, 3, 21, 0, 0, 0, 0, time.UTC)},
		{"Good Friday 2026", time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)},
		{"Labour Day 2026", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)},
		{"Waisak Day (Buddha's Birthday) 2026", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)},
		{"Ascension of Jesus Christ 2026", time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC)},
		{"Eid al-Adha 2026", time.Date(2026, 5, 27, 0, 0, 0, 0, time.UTC)},
		{"Pancasila Day 2026", time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)},
		{"Islamic New Year 1448 H", time.Date(2026, 6, 17, 0, 0, 0, 0, time.UTC)},
		{"Independence Day 2026", time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC)},
		{"Prophet Muhammad's Birthday 2026", time.Date(2026, 8, 26, 0, 0, 0, 0, time.UTC)},
		{"Christmas Day 2026", time.Date(2026, 12, 25, 0, 0, 0, 0, time.UTC)},
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

	// 2. Deadlines (EventType: Deadline)
	deadlines := []struct {
		Name string
		Date time.Time
	}{
		{"Midterm Exam Registration", time.Now().AddDate(0, 1, 0)},
		{"Final Project Submission", time.Now().AddDate(0, 3, 0)},
		{"Course Enrollment Deadline", time.Now().AddDate(0, 0, 14)},
	}

	for _, d := range deadlines {
		event := eventDomain.NewEvent(s.orgID, d.Name, eventDomain.Deadline)
		event.StartAt = &d.Date
		event.PrepareCreate(nil)

		if err := s.eventRepo.Create(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to seed deadline %s: %v\n", d.Name, err)
		}
	}

	// 3. Sessions (EventType: Session) - for cohort
	if len(s.cohortIDs) > 0 {
		sessionStart := time.Now().AddDate(0, 0, 7)
		sessionEnd := sessionStart.Add(2 * time.Hour)
		event := eventDomain.NewEvent(s.orgID, "Cohort Orientation Session", eventDomain.Session,
			eventDomain.WithTimes(sessionStart, sessionEnd),
			eventDomain.ForCohort(s.cohortIDs[0]),
		)
		event.PrepareCreate(nil)
		if err := s.eventRepo.Create(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to seed session: %v\n", err)
		}
	}

	// 4. Vanilla Events (EventType: Vanilla)
	vanillaStart := time.Now().AddDate(0, 0, 3)
	vanillaEnd := vanillaStart.Add(1 * time.Hour)
	vanillaEvent := eventDomain.NewEvent(s.orgID, "School Assembly", eventDomain.Vanilla,
		eventDomain.WithTimes(vanillaStart, vanillaEnd),
		eventDomain.WithLocation("Main Auditorium"),
	)
	vanillaEvent.PrepareCreate(nil)
	if err := s.eventRepo.Create(ctx, vanillaEvent); err != nil {
		fmt.Printf("Warning: Failed to seed vanilla event: %v\n", err)
	}

	// 5. Meetings (EventType: Meeting) - for section
	if len(s.sectionIDs) > 0 {
		meetingStart := time.Now().AddDate(0, 0, 5)
		meetingEnd := meetingStart.Add(1 * time.Hour)
		event := eventDomain.NewEvent(s.orgID, "Parent-Teacher Meeting", eventDomain.Meeting,
			eventDomain.WithTimes(meetingStart, meetingEnd),
			eventDomain.ForSection(s.sectionIDs[0]),
			eventDomain.WithLocation("Conference Room A"),
		)
		event.PrepareCreate(nil)
		if err := s.eventRepo.Create(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to seed meeting: %v\n", err)
		}
	}

	// 6. Schedules (EventType: Schedule) - lesson schedules
	now := time.Now()
	offset := int(now.Weekday())
	if offset == 0 {
		offset = 7
	}
	monday := now.AddDate(0, 0, -offset+1)

	grades := []int{10, 11, 12}
	for _, grade := range grades {
		for day := 0; day < 5; day++ {
			date := monday.AddDate(0, 0, day)
			for lesson := 0; lesson < 3; lesson++ {
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

	// 7. Announcements (EventType: Announcement)
	announcements := []struct {
		Title string
		Body  string
		Scope eventDomain.EventScope
	}{
		{"Welcome to New Semester", "We are excited to welcome all students back for the new semester!", eventDomain.ScopeGlobal},
		{"Science Fair Coming Up", "Prepare your projects for the annual science fair.", eventDomain.ScopeGlobal},
	}

	for _, a := range announcements {
		event := eventDomain.NewAnnouncement(s.orgID, a.Title, a.Body, a.Scope)
		event.PrepareCreate(nil)
		if err := s.eventRepo.Create(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to seed announcement %s: %v\n", a.Title, err)
		}
	}

	// 8. Personal events
	if len(s.studentIDs) > 0 {
		personalStart := time.Now().AddDate(0, 0, 10)
		personalEnd := personalStart.Add(1 * time.Hour)
		event := eventDomain.NewEvent(s.orgID, "Study Group Session", eventDomain.Session,
			eventDomain.WithTimes(personalStart, personalEnd),
			eventDomain.ForUser(s.studentIDs[0]),
		)
		event.PrepareCreate(nil)
		if err := s.eventRepo.Create(ctx, event); err != nil {
			fmt.Printf("Warning: Failed to seed personal event: %v\n", err)
		}
	}

	return nil
}

package main

import (
	"context"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
	cohr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/repository/postgres"
	coh "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/seed"
	cr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/repository/postgres"
	lr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/repository/postgres"
	mr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/repository/postgres"
	crs "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/seed"
	lessons "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/seed"
	modules "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/seed"
	elr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/repository/postgres"
	el "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/seed"
	er "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/repository/postgres"
	enroll "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/enrollment/seed"
	eventr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/repository/postgres"
	event "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/event/seed"
	or "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/repository/postgres"
	o "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/seed"
	prog "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/repository/postgres"
	pr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/seed"
	sectionDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/domain"
	secr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/repository/postgres"
	sec "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/section/seed"
	subjectDomain "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	subj "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/repository/postgres"
	sub "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/seed"
	ur "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
	u "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/seed"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	"github.com/google/uuid"
)

func main() {
	v := app.NewViper()
	logger := app.NewLogger(v)
	db := app.NewDatabase(v, logger)
	defer db.Close()

	roleRepo := ur.NewRoleRepository(db)
	roleSeeder := u.NewRoleSeeder(roleRepo)

	organizationRepo := or.NewOrganizationRepo(db)
	acadPeriodRepo := or.NewAcademicPeriodRepository(db)
	organizationSeeder := o.NewOrganizationSeeder(organizationRepo, acadPeriodRepo)

	educationLevelRepo := elr.NewEducationLevelRepository(db)
	educationLevelSeeder := el.NewEducationLevelRepository(educationLevelRepo)

	programRepo := prog.NewProgramRepository(db)
	programSeeder := pr.NewProgramSeeder(programRepo)

	subjectRepo := subj.NewSubjectRepository(db)
	subjectSeeder := sub.NewSubjectSeeder(subjectRepo)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Info("Starting role seeding...")
	_, err := roleSeeder.SeedRoles(ctx)
	if err != nil {
		logger.Info("Role seeding failed: ", err.Error())
	} else {
		logger.Info("Role seeding complete...")
	}

	logger.Info("Starting organization and academic period seeding...")
	seededOrg, academicPeriod, err := organizationSeeder.SeedOrganizations(ctx)
	if err != nil {
		logger.Info("Organization and academic period seeding failed: ", err.Error())
		return
	} else {
		logger.Info("Organization and academic period seeding complete...")
	}

	ctx = context.WithValue(ctx, auth.OrgIDKey, seededOrg.ID)

	logger.Info("Starting education level seeding...")
	seededEduLevels, err := educationLevelSeeder.SeedEducationLevels(ctx)
	if err != nil {
		logger.Info("Education Level seeding failed: ", err.Error())
	} else {
		logger.Info("Education Level seeding complete...")
	}

	logger.Info("Starting program seeding...")
	seededPrograms, err := programSeeder.SeedPrograms(ctx)
	if err != nil {
		logger.Info("Program seeding failed: ", err.Error())
	} else {
		logger.Info("Program seeding complete...")
	}

	// User Seeder
	userRepo := ur.NewUserRepo(db)
	userSeeder := u.NewUserSeeder(userRepo, roleRepo)
	logger.Info("Starting user seeding...")
	seededUsers, err := userSeeder.SeedUsers(ctx, seededOrg.ID)
	if err != nil {
		logger.Info("User seeding failed: ", err.Error())
	} else {
		logger.Info("User seeding complete...")
	}

	// Subject Seeder
	logger.Info("Starting subject seeding...")
	var highSchoolLevelID uuid.UUID
	for _, level := range seededEduLevels {
		if level.Code == "HIGH" {
			highSchoolLevelID = level.ID
			break
		}
	}

	var seededSubjects []*subjectDomain.Subject

	if highSchoolLevelID != uuid.Nil {
		seededSubjects, err = subjectSeeder.SeedSubjects(ctx, highSchoolLevelID)
		if err != nil {
			logger.Info("Subject seeding failed: ", err.Error())
		} else {
			logger.Info("Subject seeding complete...")
		}
	} else {
		logger.Info("High School education level not found, skipping subject seeding")
	}

	// Cohort & Section Seeder
	cohortRepo := cohr.NewCohortRepository(db)
	cohortSeeder := coh.NewCohortSeeder(cohortRepo)

	sectionRepo := secr.NewSectionRepository(db)
	sectionSeeder := sec.NewSectionSeeder(sectionRepo)

	logger.Info("Starting cohort seeding...")
	seededCohorts, err := cohortSeeder.SeedCohorts(ctx, academicPeriod.ID, highSchoolLevelID)
	if err != nil {
		logger.Info("Cohort seeding failed: ", err.Error())
	} else {
		logger.Info("Cohort seeding complete...")
	}

	logger.Info("Starting section seeding...")
	seededSections, err := sectionSeeder.SeedSections(ctx, seededCohorts)
	if err != nil {
		logger.Info("Section seeding failed: ", err.Error())
	} else {
		logger.Info("Section seeding complete...")
	}

	// Course, Module, Lesson Seeder
	courseRepo := cr.NewCourseRepository(db)
	courseSeeder := crs.NewCourseSeeder(courseRepo)

	moduleRepo := mr.NewModuleRepository(db)
	moduleSeeder := modules.NewModuleSeeder(moduleRepo)

	lessonRepo := lr.NewLessonRepository(db)
	lessonSeeder := lessons.NewLessonSeeder(lessonRepo)

	logger.Info("Starting course seeding...")
	teacher := seededUsers["teacher@candletree.com"]
	var teacherID uuid.UUID
	if teacher != nil {
		teacherID = teacher.ID
	}

	seededCourses, err := courseSeeder.SeedCourses(ctx, seededSubjects, teacherID, highSchoolLevelID)
	if err != nil {
		logger.Info("Course seeding failed: ", err.Error())
	} else {
		logger.Info("Course seeding complete...")
	}

	logger.Info("Starting module seeding...")
	seededModules, err := moduleSeeder.SeedModules(ctx, seededCourses)
	if err != nil {
		logger.Info("Module seeding failed: ", err.Error())
	} else {
		logger.Info("Module seeding complete...")
	}

	logger.Info("Starting lesson seeding...")
	seededLessons, err := lessonSeeder.SeedLessons(ctx, seededModules)
	if err != nil {
		logger.Info("Lesson seeding failed: ", err.Error())
	} else {
		logger.Info("Lesson seeding complete...")
	}

	// Event Seeder (Holidays + Lesson Schedules)
	eventRepo := eventr.NewEventRepository(db)
	eventSeeder := event.NewEventSeeder(eventRepo)

	logger.Info("Starting Indonesia holiday seeding...")
	googleAPIKey := v.GetString("GOOGLE_API_KEY")
	if googleAPIKey == "" {
		logger.Info("Skipping holiday seeding: GOOGLE_API_KEY not found in configuration")
	} else {
		_, err = eventSeeder.SeedIndonesiaHolidays(ctx, 2026, googleAPIKey)
		if err != nil {
			logger.Info("Indonesia holiday seeding failed: ", err.Error())
		} else {
			logger.Info("Indonesia holiday seeding complete...")
		}
	}

	// Build module -> section mapping based on course grade level
	moduleToCourse := make(map[uuid.UUID]int) // moduleID -> gradeLevel
	for _, module := range seededModules {
		for _, course := range seededCourses {
			if module.CourseID == course.ID {
				moduleToCourse[module.ID] = course.GradeLevel
				break
			}
		}
	}

	sectionsByGrade := make(map[int]*sectionDomain.Section)
	for _, section := range seededSections {
		switch section.Name {
		case "10-A":
			sectionsByGrade[10] = section
		case "11-A":
			sectionsByGrade[11] = section
		}
	}

	sectionsByModule := make(map[uuid.UUID]*sectionDomain.Section)
	for moduleID, gradeLevel := range moduleToCourse {
		if section, ok := sectionsByGrade[gradeLevel]; ok {
			sectionsByModule[moduleID] = section
		}
	}

	logger.Info("Starting lesson schedule seeding...")
	_, err = eventSeeder.SeedLessonSchedules(ctx, seededLessons, sectionsByModule)
	if err != nil {
		logger.Info("Lesson schedule seeding failed: ", err.Error())
	} else {
		logger.Info("Lesson schedule seeding complete...")
	}

	logger.Info("Starting school events seeding...")
	_, err = eventSeeder.SeedSchoolEvents(ctx)
	if err != nil {
		logger.Info("School events seeding failed: ", err.Error())
	} else {
		logger.Info("School events seeding complete...")
	}

	logger.Info("Starting announcements seeding...")
	_, err = eventSeeder.SeedAnnouncements(ctx)
	if err != nil {
		logger.Info("Announcements seeding failed: ", err.Error())
	} else {
		logger.Info("Announcements seeding complete...")
	}

	// Enrollment Seeder
	enrollmentRepo := er.NewEnrollmentRepository(db)
	enrollmentSeeder := enroll.NewEnrollmentSeeder(enrollmentRepo)

	logger.Info("Starting enrollment seeding...")
	_, err = enrollmentSeeder.SeedEnrollments(ctx, seededUsers, seededCourses, seededSections, academicPeriod.ID)
	if err != nil {
		logger.Info("Enrollment seeding failed: ", err.Error())
	} else {
		logger.Info("Enrollment seeding complete...")
	}

	// Cohort Member Seeder
	cohortMemberRepo := cohr.NewCohortMemberRepository(db)
	cohortMemberSeeder := coh.NewCohortMemberSeeder(cohortMemberRepo)

	logger.Info("Starting cohort member seeding...")
	_, err = cohortMemberSeeder.SeedCohortMembers(ctx, seededUsers, seededCohorts)
	if err != nil {
		logger.Info("Cohort member seeding failed: ", err.Error())
	} else {
		logger.Info("Cohort member seeding complete...")
	}

	// Section Member Seeder
	sectionMemberRepo := secr.NewSectionMemberRepository(db)
	sectionMemberSeeder := sec.NewSectionMemberSeeder(sectionMemberRepo)

	logger.Info("Starting section member seeding...")
	_, err = sectionMemberSeeder.SeedSectionMembers(ctx, seededUsers, seededSections)
	if err != nil {
		logger.Info("Section member seeding failed: ", err.Error())
	} else {
		logger.Info("Section member seeding complete...")
	}

	// Program Course Seeder
	programCourseRepo := prog.NewProgramCourseRepository(db)
	programCourseSeeder := pr.NewProgramCourseSeeder(programCourseRepo)

	logger.Info("Starting program course seeding...")
	_, err = programCourseSeeder.SeedProgramCourses(ctx, seededPrograms, seededCourses)
	if err != nil {
		logger.Info("Program course seeding failed: ", err.Error())
	} else {
		logger.Info("Program course seeding complete...")
	}
}

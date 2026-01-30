package main

import (
	"context"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
	elr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/repository/postgres"
	el "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/seed"
	or "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/repository/postgres"
	o "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/seed"
	prog "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/repository/postgres"
	pr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/seed"
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
	seededOrg, err := organizationSeeder.SeedOrganizations(ctx)
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
	_, err = programSeeder.SeedPrograms(ctx)
	if err != nil {
		logger.Info("Program seeding failed: ", err.Error())
	} else {
		logger.Info("Program seeding complete...")
	}

	logger.Info("Starting subject seeding...")
	var highSchoolLevelID uuid.UUID
	for _, level := range seededEduLevels {
		if level.Code == "HIGH" {
			highSchoolLevelID = level.ID
			break
		}
	}

	if highSchoolLevelID != uuid.Nil {
		_, err = subjectSeeder.SeedSubjects(ctx, highSchoolLevelID)
		if err != nil {
			logger.Info("Subject seeding failed: ", err.Error())
		} else {
			logger.Info("Subject seeding complete...")
		}
	} else {
		logger.Info("High School education level not found, skipping subject seeding")
	}
}

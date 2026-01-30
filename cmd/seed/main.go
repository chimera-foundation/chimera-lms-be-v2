package main

import (
	"context"
	"time"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
	or "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/repository/postgres"
	o "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/seed"
	ur "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
	u "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/seed"
	elr "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/repository/postgres"
	el "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/seed"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
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

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    logger.Info("Starting role seeding...")
    err := roleSeeder.SeedRoles(ctx)
    if err != nil {
        logger.Info("Role seeding failed: ", err.Error())
    } else {
        logger.Info("Role seeding complete...")
    }

    logger.Info("Starting organization and academic period seeding...")
    err = organizationSeeder.SeedOrganizations(ctx)
    if err != nil {
        logger.Info("Organization and academic period seeding failed: ", err.Error())
    } else {
        logger.Info("Organization and academic period seeding complete...")
    }

    logger.Info("Fetching Organization...")
    organization, err := organizationRepo.GetBySlug(ctx, "cts")
    if err != nil {
        logger.Info("Organization not found")
    } else {
        logger.Info("Organization found: ", organization.ID)
    }

    ctx = context.WithValue(ctx, auth.OrgIDKey, organization.ID)

    logger.Info("Starting education level seeding...")
    err = educationLevelSeeder.SeedEducationLevels(ctx)
    if err != nil {
        logger.Info("Education Level seeding failed: ", err.Error())
    } else {
        logger.Info("Education Level seeding complete...")
    }
}
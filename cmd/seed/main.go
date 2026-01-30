package main

import (
    "context"
    "time"

    "github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
    ur "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
    or "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/repository/postgres"
    u "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/seed"
    o "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/organization/seed"
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
}
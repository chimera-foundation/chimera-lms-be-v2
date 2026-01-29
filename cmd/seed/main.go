package main

import (
    "context"
    "time"

    "github.com/chimera-foundation/chimera-lms-be-v2/internal/app"
    "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/repository/postgres"
    "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/seed"
)

func main() {
    v := app.NewViper()
    logger := app.NewLogger(v)
    db := app.NewDatabase(v, logger)
    defer db.Close()

    roleRepo := postgres.NewRoleRepository(db)
    roleSeeder := seed.NewRoleSeeder(roleRepo)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    logger.Info("Starting role seeding...")
    roleSeeder.SeedRoles(ctx)
    logger.Info("Role seeding complete...")
}
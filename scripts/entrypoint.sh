#!/bin/sh
set -e

# Use double quotes for safety
if [ "$SEED_DB" = "true" ]; then
    echo "🌱 SEED_DB is true. Running seeder..."
    # Ensure the path matches where the Dockerfile puts it
    ./chimera-seeder || echo "⚠️ Seeder finished with errors or already seeded."
else
    echo "⏩ Seeding skipped (SEED_DB=$SEED_DB)"
fi

echo "🚀 Starting API..."
# 'exec' replaces the shell with the API process
exec ./chimera-api
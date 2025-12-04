#!/bin/bash

# Script to initialize SQLite3 database for secure-notes-server
# Run this from the project root: ./init_database.sh

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}SQLite3 Database Initialization${NC}"
echo -e "${BLUE}========================================${NC}"

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}Error: sqlite3 is not installed${NC}"
    echo "Please install it first:"
    echo "  Ubuntu/Debian: sudo apt-get install sqlite3"
    echo "  macOS: brew install sqlite3"
    exit 1
fi

# Database path
DB_DIR="./server/database"
DB_PATH="${DB_DIR}/secure_notes.db"
MIGRATIONS_DIR="./server/migrations"

# Create database directory if not exists
mkdir -p "$DB_DIR"

# Remove old database if exists (for clean initialization)
if [ -f "$DB_PATH" ]; then
    echo -e "${BLUE}Removing old database...${NC}"
    rm "$DB_PATH"
fi

echo -e "${GREEN}Creating new database at: ${DB_PATH}${NC}"

# Run schema migration
echo -e "${BLUE}Running schema migration (001_init_schema.sql)...${NC}"
if [ -f "${MIGRATIONS_DIR}/001_init_schema.sql" ]; then
    sqlite3 "$DB_PATH" < "${MIGRATIONS_DIR}/001_init_schema.sql"
    echo -e "${GREEN}✓ Schema created successfully${NC}"
else
    echo -e "${RED}Error: ${MIGRATIONS_DIR}/001_init_schema.sql not found${NC}"
    exit 1
fi

# Run seed data (optional)
echo -e "${BLUE}Running seed data (002_seed_data.sql)...${NC}"
if [ -f "${MIGRATIONS_DIR}/002_seed_data.sql" ]; then
    sqlite3 "$DB_PATH" < "${MIGRATIONS_DIR}/002_seed_data.sql"
    echo -e "${GREEN}✓ Seed data inserted successfully${NC}"
else
    echo -e "${BLUE}⚠ No seed data file found (skipped)${NC}"
fi

# Verify database
echo -e "${BLUE}Verifying database tables...${NC}"
TABLE_COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM sqlite_master WHERE type='table';")
echo -e "${GREEN}✓ Database created with ${TABLE_COUNT} tables${NC}"

# Show table list
echo -e "${BLUE}Tables in database:${NC}"
sqlite3 "$DB_PATH" "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Database initialization completed!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "Database location: ${DB_PATH}"
echo -e "Database size: $(du -h "$DB_PATH" | cut -f1)"
echo ""
echo -e "${BLUE}To connect to the database:${NC}"
echo -e "  sqlite3 ${DB_PATH}"
echo ""
echo -e "${BLUE}To view tables:${NC}"
echo -e "  sqlite3 ${DB_PATH} '.tables'"
echo ""

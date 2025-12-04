#!/bin/bash

set -e

DB_PATH="./database/database.db"
SCHEMA_FILE="./database/schema.sql"
SEED_FILE="./database/seed.sql"

# ==== Xử lý Arguments ====
DROP=true
IMPORT_SCHEMA=true
IMPORT_SEED=true

for arg in "$@"; do
    case $arg in
        --no-drop) DROP=false ;;
        --no-schema) IMPORT_SCHEMA=false ;;
        --no-seed) IMPORT_SEED=false ;;
    esac
done


# ========================
# 1. XÓA DATABASE
# ========================
if $DROP; then
    rm -f "$DB_PATH"
    echo "1. Database initialized (database.db created)"
else
    echo "1. Skipping database removal ( --no-drop )"
fi


# ========================
# 2. IMPORT SCHEMA
# ========================
if $IMPORT_SCHEMA; then
    if [ ! -f "$SCHEMA_FILE" ]; then
        echo "! Schema file not found: $SCHEMA_FILE"
        exit 1
    fi

    if sqlite3 "$DB_PATH" < "$SCHEMA_FILE"; then
        echo "2. Schema imported successfully"
    else
        echo "! Failed to import schema"
        exit 1
    fi
else
    echo "2.Skipping schema import ( --no-schema )"
fi


# ========================
# 3. IMPORT SEED DATA
# ========================
if $IMPORT_SEED; then
    if [ ! -f "$SEED_FILE" ]; then
        echo "! Seed file not found: $SEED_FILE"
    else
        if sqlite3 "$DB_PATH" < "$SEED_FILE"; then
            echo "3. Seed data imported successfully"
        else
            echo "! Seed import failed"
        fi
    fi
else
    echo "3. Skipping seed import ( --no-seed )"
fi


# ========================
# 4. DONE
# ========================
echo "###Build completed!"

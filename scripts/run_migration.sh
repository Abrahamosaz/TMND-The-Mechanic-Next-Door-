#!/usr/bin/sh

# Check if the migration name argument is provided
if [ -z "$1" ]; then
  echo "Error: Migration name is required."
  exit 1
fi

# Get the migration name from the first command-line argument
migration_name=$1

atlas migrate diff "$migration_name" --env gorm 
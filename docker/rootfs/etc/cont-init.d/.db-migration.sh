#!/usr/bin/with-contenv bash

source $GFLY_HOME/.bashrc
source $GFLY_HOME/app/.env
echo -e "${BLUE}---------- DB Migration - ${DB_DRIVER} ----------${NC}"

if [[ "$DB_DRIVER" == "postgresql" ]]
then
  MIGRATION_FOLDER="${GFLY_HOME}/app/database/migrations/postgresql"
  DATABASE_URL="postgres://user:secret@db:5432/gfly?sslmode=disable"
fi

if [[ "$DB_DRIVER" == "mysql" ]]
then
  MIGRATION_FOLDER="${GFLY_HOME}/app/database/migrations/mysql"
  DATABASE_URL="mysql://user:secret@tcp(db:3306)/gfly"
fi

migrate -path $MIGRATION_FOLDER -database "$DATABASE_URL" up

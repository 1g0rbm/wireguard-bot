#!/bin/bash

export MIGRATION_DSN="host=pg port=5432 dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable"

echo | ls -la
echo | pwd
echo ${MIGRATION_DIR}

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v
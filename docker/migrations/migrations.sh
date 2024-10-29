#!/bin/bash

echo | ls -la
echo | pwd
echo ${MIGRATION_DIR}

sleep 2 && goose -dir "/app/migrations" postgres "${DATABASE_URL}" up -v
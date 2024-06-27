#!/bin/bash

DB_USER="mr_workers"
DB_PASSWORD="password"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="mr_workers"
SSL_MODE="disable"

~/go/bin/goose -dir db/migrations postgres "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$SSL_MODE" up
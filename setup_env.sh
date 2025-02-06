#!/bin/bash

echo "Quickly setting up .env file."

# Prompt for user input
read -p "PostgreSQL username: " POSTGRES_USER
read -p "PostgreSQL password: " POSTGRES_PASSWORD
read -p "PostgreSQL database name: " POSTGRES_DB
read -p "Grafana admin username: " GF_SECURITY_ADMIN_USER
read -p "Grafana admin password: " GF_SECURITY_ADMIN_PASSWORD
read -p "Database username: " DB_USER
read -p "Database password: " DB_PASS
read -p "Database name: " DB_NAME

# Define the environment variables
cat <<EOF > .env
POSTGRES_USER=$POSTGRES_USER
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
POSTGRES_DB=$POSTGRES_DB

GF_SECURITY_ADMIN_USER=$GF_SECURITY_ADMIN_USER
GF_SECURITY_ADMIN_PASSWORD=$GF_SECURITY_ADMIN_PASSWORD

DB_HOST=postgres
DB_PORT=5432
DB_USER=$DB_USER
DB_PASS=$DB_PASS
DB_NAME=$DB_NAME
EOF

echo ".env file created successfully."

echo "To edit, edit the .env file or run this again."
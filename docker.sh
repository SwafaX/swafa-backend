#!/bin/bash

POSTGRES_USER=todo_user
POSTGRES_PASSWORD=todo_password
POSTGRES_DB=tododb
POSTGRES_PORT=5435

docker pull postgres:15-alpine

docker run -d \
  --name postgres_tododb \
  -e POSTGRES_USER=$POSTGRES_USER \
  -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
  -e POSTGRES_DB=$POSTGRES_DB \
  -p $POSTGRES_PORT:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  postgres:15-alpine

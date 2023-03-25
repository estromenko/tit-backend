#!/bin/sh

set -x

export PG_USER="${PG_USER:="postgres"}"

export PG_PASSWORD="${PG_PASSWORD:="secret"}"

export PG_HOST="${PG_HOST:="127.0.0.1"}"

export PG_PORT="${PG_PORT:="5432"}"

export PG_NAME="${PG_NAME:="tutorintech"}"

migrate \
  -source file:///migrations \
  -database "postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_NAME}?sslmode=disable" \
  up

exec tit-backend

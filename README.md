# tit-backend

## Dependencies

### Mandatory

- Golang version 18(Installation instruction [link](https://go.dev/doc/install))

## Getting started

### Development

- You need to install dev dependencies

```bash
$ make install-dev-requirements
```

- Run PostgreSQL database

```bash
$ docker run -d --rm --name tutorintech_postgres -v /srv/_tutorintech_postgres:/var/lib/postgresql/data -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:15-alpine
$ docker run --rm -it --link tutorintech_postgres:postgres -e PGPASSWORD=secret postgres:15-alpine createdb -h postgres -U postgres tutorintech
```

- Install tool for migrations and apply them:

```
$ make install-migrate
$ migrate -source file://./migrations -database postgres://postgres:secret@localhost:5432/tutorintech\?sslmode=disable up
```

- Build required images:
```bash
$ docker compose -f docker/dashboard/docker-compose.yml build
```

- Run application. You can configure app by passing env variables directly or create .env 
file in project root.

```bash
$ make dev
```

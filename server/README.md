# Backend

This directory contains backend for Expense Management System.

## Prerequisites

Before running the project locally, make sure you have the following installed:

- [Go](https://go.dev/) (version 1.24.3 or later)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [golang-migrate](https://github.com/golang-migrate/migrate) – follow installation instructions: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Optional tools for development and debugging:

- **PostgreSQL client** (`psql`) or GUI tools (e.g., [DBeaver](https://dbeaver.io/)) for inspecting database, running queries, and checking migrations
- **Redis client** (`redis-cli`) – for inspecting keys and debugging

## Development Setup

The following commands should be run from the `server/` directory.

### Start Docker Services

This command starts all the necessary services (PostgreSQL, Redis, Kafka) in Docker.

```bash
make docker-up
```

To stop all services:

```bash
make docker-down
```

> Make sure all containers are running before running database migration and starting the app.
>
> Kafka UI is available at http://localhost:8080/

### Run Database Migrations

Before starting the application, you need to apply the database schema.

To apply all migrations:

```bash
make migrate-up
```

To rollback migrations:

```bash
make migrate-down
```

To create a new migration file:

```bash
make migrate-create name=<migration_name>
```

To run database seeding:

```bash
go run dev/seeder/main.go
```

### Run the Application

First, copy the sample environment file and adjust the configuration values as needed.

```bash
cp env.sample .env
```

The application consists of 2 main processes that need to be run separately.

To run the API server:

```bash
make run-api
```

> The API will be available at http://localhost:8500/
>
> API documentation (Swagger) is available at http://localhost:8500/swagger/

To run the Kafka consumer worker:

```bash
make run-consumer
```

### Testing

To run unit tests:

```bash
make test
```

### Mock Payment API

By default, the mock payment API is running in Docker and can be accessed via http://localhost:9500/v1/payments.

To run it manually:

```bash
go run dev/payment-api/main.go
```

> This local mock was created because the public Postman mock API provided intermittently returns a 403 Forbidden error.

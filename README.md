# CAT SOCIAL NETWORK

This is a mini e-commerce application where some people can register as cashier a.k.a staff to maintain products and anyone can register as customer to buy anything available from the store. This project is developed using [Golang](https://golang.org), [PostgreSQL](https://www.postgresql.org), and [Docker](https://www.docker.com). We use [Gofiber](https://gofiber.io) as the main web framework for this project, [pgx](https://github.com/jackc/pgx) as the main PostgreSQL driver and [golang-migrate](https://github.com/golang-migrate/migrate) for the migration.

To run this project, you need to have Docker installed on your machine. If you don't have Docker installed, you can download it [here](https://www.docker.com/products/docker-desktop). For Mac users, we recommend using [Orbstack](https://orbstack.dev/) to run your Docker containers since we have a better experience using it.

For the ease of local development, we utilize [Docker Compose](https://docs.docker.com/compose/gettingstarted/) and [Air](https://github.com/cosmtrek/air) for hot-reloading. Docker Compose is installed by default if you are using either Docker Desktop or Orbstack on your machine. Please follow through the [Air's documentation](https://github.com/cosmtrek/air?tab=readme-ov-file#installation) to install Air on your machine and [Docker Compose installation manual](https://docs.docker.com/compose/install/) if you need to install it manually.

You can run the project using the following command:

```bash
docker compose up
```

Please make sure that the required environment variables was set before running the project. You can set the variables on the `.env` (`.env.dev` for local development) file.

You will also need to make sure that golang-migrate is installed on your machine. Please refer to the [documentation](https://github.com/golang-migrate/migrate) for more details.

At this stage, you will need to run the migration manually to create the tables on the database. **You won't need to do this manually in the next stage.** Please use the following command to run the migration:

```bash
migrate -database "postgresql://[USER]:[PASSWORD]@[HOST]:[PORT]/[DB_NAME]?connect_timeout=10&application_name=[APP_NAME]&sslmode=disable" -path db/migrate/primary/ up
```

To add a new migration, you can run the following command:

```bash
migrate create -ext sql -dir ./db/migrate/primary/ -tz "Asia/Jakarta" [MIGRATION_NAME]
```

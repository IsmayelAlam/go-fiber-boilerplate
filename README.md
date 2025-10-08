# Go Fiber Feature Based Backend Boilerplate

A clean, production-ready **Go backend boilerplate** built with the [Fiber](https://gofiber.io/) framework, **PostgreSQL**, and **Swagger documentation**.  
This template follows a **feature-based modular architecture** for scalable, maintainable backend development.

---

##  Key Features

- **High-performance HTTP framework** — Powered by Fiber
- **Feature-based modular design** — Scales cleanly as your app grows
- **PostgreSQL integration** — via `pgx` driver
- **Docker-ready setup** — with `.env` auto-loading and `docker-compose`
- **JWT authentication middleware**
- **Swagger API docs** — Auto-generated OpenAPI 3.0 spec
- **Unit testing support** `(todo)`
- **Graceful shutdown & centralized logging**
- **Config management** using `flags` and structured configs

---

## Project Structure

```
├── cmd/
│ └── api/
│   └── main.go # Application entrypoint

├── config/ # Config loading and environment management
├── docs/ # Swagger documentation files
│
├── internal/
│ ├── database/ # DB initialization & migrations
│ ├── middleware/ # Global and feature-based middlewares
│ ├── modules/ # all modules
│   ├── auth/ # Auth feature (JWT, login, register)
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── routes.go
│   ├── user/ # User module
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── routes.go
│   └── ...
│
├── .env
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Getting Started

### Prerequisites

- **Go** v1.21+
- **Docker** & **Docker Compose**
- **Make** (optional convenience)

### Clone and Setup

```bash
git clone https://github.com/ismayelalam/go-fiber-boilerplate.git
cd go-fiber-boilerplate
```

### Install Dependencies

```
go mod tidy
```

### Configure Environment

Copy `.env.example` → `.env`:

```
cp .env.example .env
```

Run with Docker
`Build and start the application with PostgreSQL:`

```
docker-compose up --build
```

### This will spin up:

`app` → Go Fiber backend
`db` → PostgreSQL database

Visit the API at: http://localhost:8080/api/v1

### Local Development (without Docker)

If you prefer running locally:

```
air
```

Ensure your .env file points to your local DB and have [air](https://github.com/air-verse/air) installed.

### License

This project is licensed under the `MIT` License.

### Contributing

- Fork the repo

- Create a feature branch (feature/awesome-feature)

- Commit your changes

- Open a PR 

### Author

[**Ismayel Alam**](https://ismayelalam.com/)

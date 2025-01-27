Finance API

This project is a REST API built in Go using the Gin framework, designed for managing financial operations. It allows users to perform transactions, transfer funds between users, and view recent transaction history. The application interacts with a PostgreSQL database and supports database migrations via Goose.

Key features:

User balance management.
Money transfers between users.
Viewing the last 10 transactions for a user.
Clean architecture with handler, service, and repository layers.
Dockerized for easy deployment with Docker Compose.
Technologies used:

Go
Gin
PostgreSQL
Goose (for migrations)
Docker & Docker Compose

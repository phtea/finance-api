# Finance API

This project is a REST API built in Go using the Gin framework, designed for managing financial operations. It allows users to perform transactions, transfer funds between users, and view recent transaction history. The application interacts with a PostgreSQL database and supports database migrations via Goose. Service layer is covered with unit-tests using Mocks.

## Key Features

- User balance management
- Money transfers between users
- Viewing the last 10 transactions for a user
- Clean architecture with handler, service, and repository layers
- Dockerized for easy deployment with Docker Compose

## Technologies Used

- **Go**: The main programming language
- **Gin**: Web framework for building the REST API
- **PostgreSQL**: Relational database for storing user and transaction data
- **Goose**: Tool for database migrations
- **Docker & Docker Compose**: For containerization and easy deployment

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/phtea/finance-api.git
   ```
2. Navigate to the project directory:
   ```bash
   cd finance-api
   ```
3. Build and run the application with Docker Compose:
   ```bash
   make run
   ```
4. Apply database migrations:
   ```bash
   make migrate-up
   ```
Endpoints
- **POST /users**: Create a new user
- **POST /users/balance**: Add balance to a user
- **POST /users/transfer**: Transfer balance between users
- **GET /users/:id**: Get user details by ID
- **GET /transactions**: Get the last 10 transactions

## License
This project is licensed under the MIT License.

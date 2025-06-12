# Healthcare Portal Backend

A RESTful API backend for the Healthcare Portal system built with Go, Gin, and PostgreSQL.

## Features

- JWT-based authentication
- Role-based access control (Receptionist & Doctor)
- Patient management with CRUD operations
- Repository pattern implementation
- Swagger API documentation
- Comprehensive error handling
- Database migrations

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/healthcare-portal.git
cd healthcare-portal/backend
```

2. Install dependencies:
```BASH
go mod download
```
3. Set up environment variables:
```BASH
cp .env.example .env
# Edit .env with your configuration
```
4. Run database migrations:
```BASH
go run cmd/server/main.go
```

5. (Optional) Seed the database:
```BASH
go run cmd/seeder/main.go
```
## Running the Application
### Development
```BASH
go run cmd/server/main.go
```
### Production
```BASH
go build -o healthcare-portal cmd/server/main.go
./healthcare-portal
```

## API Documentation
Swagger documentation is available at:
```BASH
http://localhost:8080/swagger/index.html
```
### Running Tests
```BASH
go test ./...
```
For coverage:
```BASH
go test -cover ./...
```
## API Endpoints
### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
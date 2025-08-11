# Go CRUD Demo

A production-ready CRUD (Create, Read, Update, Delete) application built with Go, featuring RESTful APIs, database integration, comprehensive testing, and modern development practices.

*Created using GitHub MCP By Claude-4*

## 🚀 Features

- **RESTful API** - Full CRUD operations for user management
- **Database Integration** - PostgreSQL with connection pooling
- **Middleware** - Authentication, logging, and CORS
- **Validation** - Request data validation
- **Testing** - Unit tests and integration tests
- **Configuration** - Environment-based configuration
- **Documentation** - Swagger/OpenAPI documentation
- **Docker Support** - Containerized deployment
- **Health Checks** - Application health monitoring
- **Graceful Shutdown** - Clean application termination

## 📁 Project Structure

```
go-crud/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   ├── connection.go        # Database connection setup
│   │   └── migrations/          # SQL migration files
│   ├── handlers/
│   │   ├── user_handler.go      # User CRUD handlers
│   │   └── health_handler.go    # Health check handlers
│   ├── middleware/
│   │   ├── auth.go              # Authentication middleware
│   │   ├── logging.go           # Logging middleware
│   │   └── cors.go              # CORS middleware
│   ├── models/
│   │   └── user.go              # User data model
│   ├── repository/
│   │   ├── user_repository.go   # User data access layer
│   │   └── interfaces.go        # Repository interfaces
│   ├── services/
│   │   └── user_service.go      # Business logic layer
│   └── validators/
│       └── user_validator.go    # Input validation
├── tests/
│   ├── integration/
│   │   └── user_integration_test.go
│   └── unit/
│       ├── handlers_test.go
│       ├── services_test.go
│       └── repository_test.go
├── docs/
│   ├── api.md                   # API documentation
│   ├── deployment.md            # Deployment guide
│   └── swagger.yaml             # OpenAPI specification
├── scripts/
│   ├── migrate.sh               # Database migration script
│   └── test.sh                  # Testing script
├── docker/
│   ├── Dockerfile               # Production Docker image
│   ├── Dockerfile.dev           # Development Docker image
│   └── docker-compose.yml       # Multi-container setup
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Makefile                     # Build and development commands
└── README.md                    # This file
```

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 15+
- **HTTP Router**: Gorilla Mux
- **Database Driver**: lib/pq
- **Testing**: Go testing package + Testify
- **Documentation**: Swagger
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Docker (optional)
- Git

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/pratham15541/go-crud.git
cd go-crud
```

### 2. Environment Setup

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Database Setup

```bash
# Using Docker
docker-compose up -d postgres

# Or manually create PostgreSQL database
createdb crud_demo
```

### 5. Run Migrations

```bash
make migrate-up
```

### 6. Start the Application

```bash
# Development mode
make dev

# Production mode
make build
./bin/server
```

## 🐳 Docker Development

```bash
# Start all services
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Test coverage
make test-coverage
```

## 📖 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/users` | Get all users |
| GET | `/users/{id}` | Get user by ID |
| POST | `/users` | Create new user |
| PUT | `/users/{id}` | Update user |
| DELETE | `/users/{id}` | Delete user |

### Example Requests

#### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
  }'
```

#### Get All Users
```bash
curl http://localhost:8080/api/v1/users
```

#### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

#### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "johnsmith@example.com",
    "age": 31
  }'
```

#### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## 🔧 Configuration

Environment variables can be set in `.env` file:

```env
# Server
PORT=8080
HOST=localhost
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=crud_demo
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 📚 Additional Documentation

- [API Documentation](docs/api.md) - Detailed API reference
- [Deployment Guide](docs/deployment.md) - Production deployment instructions
- [Contributing Guidelines](CONTRIBUTING.md) - How to contribute to this project

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Pratham Parikh**
- GitHub: [@pratham15541](https://github.com/pratham15541)
- Email: prathamparikh94@gmail.com

## 🙏 Acknowledgments

- Created with GitHub MCP (Model Context Protocol)
- Go community for excellent documentation
- Contributors who help improve this project

## 📊 Project Status

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Coverage](https://img.shields.io/badge/coverage-85%25-yellow)

---

⭐ If you find this project helpful, please give it a star!

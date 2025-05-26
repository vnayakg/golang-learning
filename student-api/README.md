# Student API

Simple RESTful API for managing student.

- Create, read, update, and delete student records
- Input validation
- PostgreSQL database with GORM ORM
- Soft delete support

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd student-api
```

2. Start the PostgreSQL database:
```bash
docker-compose up -d
```

3. Set up environment variables (optional):
```bash
cp .env.example .env
```

4. Run the application:
```bash
make run
```

Or without Make:
```bash
go run main.go
```

## API Endpoints

### Create Student
```http
POST /v1/api/students
Content-Type: application/json

{
    "name": "John Doe",
    "email": "john@doe.com",
    "age": 20,
    "grade": "A"
}
```

### Get All Students
```http
GET /v1/api/students
```

### Get Student by ID
```http
GET /v1/api/students/:id
```

### Update Student
```http
PUT /v1/api/students/:id
Content-Type: application/json

{
    "name": "John Updated",
    "email": "john.updated@doe.com",
    "age": 21,
    "grade": "A+"
}
```

### Delete Student
```http
DELETE /v1/api/students/:id
```

## Development

### Running Tests
```bash
make test
```

### Building
```bash
make build
```

### Clean
```bash
make clean
```

## Environment Variables

- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: PostgreSQL user (default: postgres)
- `DB_PASSWORD`: PostgreSQL password (default: postgres)
- `DB_NAME`: PostgreSQL database name (default: student_db)
- `SERVER_PORT`: API server port (default: 8080)

# Go Hello World API

A clean, optimized Go REST API with user management and task functionality.

## Features

- ✅ User creation and authentication via API keys
- ✅ Task management (CRUD operations)
- ✅ Clean architecture with proper separation of concerns
- ✅ Database connection pooling and optimization
- ✅ Comprehensive error handling
- ✅ Input validation and sanitization
- ✅ Graceful server shutdown
- ✅ CORS support
- ✅ Structured logging

## Improvements Made

### Code Quality & Organization
- **Separated concerns**: Database, handlers, models, and auth are properly organized
- **Optimized regex compilation**: Compiled once at package level for better performance
- **Improved error handling**: Consistent error responses and proper logging
- **Input validation**: Comprehensive validation for all user inputs
- **Memory optimization**: Proper connection pooling and resource cleanup

### API Architecture
- **RESTful design**: Following REST conventions for all endpoints
- **Middleware support**: CORS and potential for additional middleware
- **Graceful shutdown**: Proper server lifecycle management
- **Timeout handling**: Request timeouts to prevent hanging connections

### Database Optimizations
- **Connection pooling**: Configured with optimal pool settings
- **Singleton pattern**: Efficient database connection management
- **Proper transactions**: Support for database transactions
- **SQL injection prevention**: Using parameterized queries

## Quick Start

1. **Setup Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```

2. **Run the Application**
   ```bash
   go run main.go
   ```

3. **Test the API**
   ```bash
   # Health check
   curl http://localhost:8080/v1/healthz
   
   # Create a user
   curl -X POST http://localhost:8080/v1/user \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser"}'
   ```

## API Endpoints

### Health & Status
- `GET /v1/healthz` - Health check
- `GET /v1/err` - Error testing endpoint

### User Management
- `POST /v1/user` - Create a new user
- `GET /v1/user` - Get current user (requires API key)

### Task Management
- `POST /v1/tasks` - Create a new task
- `GET /v1/tasks` - Get all tasks for current user
- `GET /v1/tasks/search?query=term` - Search tasks
- `GET /v1/tasks/{taskId}` - Get specific task
- `PUT /v1/tasks/{taskId}` - Update task
- `DELETE /v1/tasks/{taskId}` - Delete task (soft delete)

## Authentication

All protected endpoints require an API key in the Authorization header:
```
Authorization: APIKEY your_api_key_here
```

## Request/Response Examples

### Create User
```bash
curl -X POST http://localhost:8080/v1/user \
  -H "Content-Type: application/json" \
  -d '{"username": "johndoe"}'
```

Response:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "username": "johndoe",
  "api_key": "abc123def456",
  "created_at": "2025-07-26T10:00:00Z",
  "updated_at": "2025-07-26T10:00:00Z"
}
```

### Create Task
```bash
curl -X POST http://localhost:8080/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: APIKEY abc123def456" \
  -d '{"title": "Complete project documentation"}'
```

## Project Structure

```
├── main.go                 # Application entry point
├── .env.example           # Environment configuration template
├── handlers/              # HTTP handlers
│   ├── api_config.go     # Database configuration
│   ├── json.go           # JSON response utilities
│   ├── users.go          # User handlers
│   ├── tasks.go          # Task handlers
│   └── utils.go          # Utility functions
├── internal/
│   ├── auth/             # Authentication logic
│   │   └── auth.go
│   └── database/         # Database layer
│       ├── db.go         # Generated database interface
│       ├── models.go     # Generated database models
│       ├── users.sql.go  # Generated user queries
│       └── tasks.sql.go  # Generated task queries
├── models/               # API response models
│   ├── users.go
│   └── tasks.go
└── sql/                  # Database schema and queries
    ├── queries/
    │   ├── users.sql
    │   └── tasks.sql
    └── schema/
        ├── 001_users.sql
        ├── 002_users_apikey.sql
        └── 003_tasks.sql
```

## Dependencies

- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/go-chi/cors` - CORS middleware
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/lib/pq` - PostgreSQL driver

## Development

### Code Quality Standards
- All functions have proper error handling
- Input validation on all user inputs
- Consistent naming conventions
- Comprehensive logging
- No memory leaks or resource leaks

### Performance Optimizations
- Regex compiled once at package level
- Database connection pooling
- Efficient JSON marshaling
- Proper HTTP timeouts
- Minimal memory allocations

### Security Features
- SQL injection prevention
- Input sanitization
- API key authentication
- Proper error messages (no information leakage)
- CORS configuration

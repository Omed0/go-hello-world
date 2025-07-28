# Go Task Management API

A clean, well-structured REST API built with Go for managing tasks with user authentication. This project demonstrates Go best practices, clean architecture, and is perfect for beginners learning Go web development.

## 🚀 Features

- **User Management**: Create users with automatic API key generation
- **Task CRUD Operations**: Create, read, update, delete tasks
- **Task Completion**: Mark tasks as complete/incomplete
- **Search & Filter**: Search tasks by title and description
- **Authentication**: API key-based authentication
- **Soft Delete**: Tasks are soft-deleted (can be restored)
- **Clean Architecture**: Proper separation of concerns
- **Middleware**: Authentication, logging, and panic recovery
- **Graceful Shutdown**: Proper server lifecycle management
- **Connection Pooling**: Optimized database connections

## 📋 Requirements

- Go 1.19+
- PostgreSQL database
- Git

## 🛠️ Installation & Setup

### 1. Clone the Repository
```bash
git clone <your-repo-url>
cd # 🚀 GO Hello World - Complete Task Management System

A modern, production-ready task management system built with Go, featuring a REST API, beautiful CLI interface, and comprehensive user management with RBAC (Role-Based Access Control).

## ✨ Features

### 🔐 Authentication & Authorization
- **Secure Authentication**: API key-based authentication with bcrypt password hashing
- **Role-Based Access Control (RBAC)**: Support for User, Moderator, Admin, and Owner roles
- **Organization Management**: Multi-tenant organization support
- **User Profiles**: Complete user management with age, gender, and organization fields

### 📋 Task Management
- **Full CRUD Operations**: Create, Read, Update, Delete tasks
- **Task Status Tracking**: Mark tasks as finished/unfinished
- **Search & Filter**: Advanced task search capabilities
- **Soft Delete**: Safe task deletion with recovery options

### 🎨 Beautiful CLI Interface
- **Server-Connected**: CLI connects to your API server for real-time data
- **Interactive TUI**: Built with Charm.sh for beautiful terminal experience
- **Full Feature Parity**: Access all server features through the CLI
- **Authentication**: Login with your server credentials
- **Real-time Updates**: Live task management synchronized with server
- **Intuitive Navigation**: Easy-to-use keyboard shortcuts

### 🏗️ Clean Architecture
- **Middleware System**: Authentication, logging, recovery, and RBAC middleware
- **Database Layer**: PostgreSQL with SQLC for type-safe queries
- **Connection Pooling**: Optimized database connections
- **Graceful Shutdown**: Proper cleanup and signal handling
- **Environment Configuration**: Flexible configuration with .env support

## 📁 Project Structure

```
go-hello-world/
├── cmd/
│   └── cli/                    # Beautiful CLI application
│       └── main.go
├── handlers/                   # HTTP request handlers
│   ├── api_config.go          # API configuration
│   ├── json.go                # JSON utilities
│   ├── organizations.go       # Organization management
│   ├── tasks.go              # Task management
│   ├── users.go              # User management
│   └── utils.go              # Utility handlers
├── internal/
│   ├── auth/                  # Authentication system
│   │   └── auth.go
│   ├── config/                # Configuration management
│   │   └── config.go
│   ├── database/              # Database layer (SQLC generated)
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── organizations.sql.go
│   │   ├── tasks.sql.go
│   │   └── users.sql.go
│   └── middleware/            # HTTP middleware
│       ├── auth.go
│       ├── cors.go
│       ├── logging.go
│       ├── rbac.go
│       └── recovery.go
├── models/                    # API response models
│   ├── organizations.go
│   ├── tasks.go
│   └── users.go
├── sql/
│   ├── queries/              # SQL queries for SQLC
│   │   ├── organizations.sql
│   │   ├── tasks.sql
│   │   └── users.sql
│   └── schema/               # Database migrations
│       ├── 001_users.sql
│       ├── 002_users_apikey.sql
│       ├── 003_tasks.sql
│       └── 004_tasks_fields.sql
├── .env.example              # Environment variables template
├── go.mod                    # Go modules
├── main.go                   # Server entry point
├── requirements.txt          # Feature requirements
└── sqlc.yaml                # SQLC configuration
```

## 🚀 Quick Start

### Prerequisites
- Go 1.19 or higher
- PostgreSQL database
- Git

### 1. Clone & Setup
```bash
git clone <repository-url>
cd go-hello-world
go mod download
```

### 2. Database Setup
```bash
# Create PostgreSQL database
createdb helloworlddb

# Run migrations (in order)
psql -d helloworlddb -f sql/schema/001_users.sql
psql -d helloworlddb -f sql/schema/002_users_apikey.sql
psql -d helloworlddb -f sql/schema/003_tasks.sql
psql -d helloworlddb -f sql/schema/004_tasks_fields.sql
```

### 3. Environment Configuration
```bash
cp .env.example .env
# Edit .env with your database credentials
```

Example `.env`:
```env
PORT=8080
DB_URL=postgres://username:password@localhost/helloworlddb?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key
```

### 4. Build & Run

#### API Server
```bash
# Build the server
go build -o bin/hello-world .

# Run the server
./bin/hello-world
```

#### CLI Application
```bash
# Build the CLI
go build -o bin/task-cli ./cmd/cli

# Run the CLI
./bin/task-cli
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/v1
```

### Authentication
All protected endpoints require an API key in the Authorization header:
```
Authorization: ApiKey your-api-key-here
```

### Endpoints

#### 🔐 Authentication
- `POST /user` - Create new user account
- `POST /login` - Login and get API key

#### 👤 User Management
- `GET /user` - Get current user profile
- `PUT /user` - Update user profile

#### 🏢 Organization Management
- `POST /organizations` - Create organization (requires authentication)
- `GET /organizations/{orgId}` - Get organization details
- `PUT /organizations/{orgId}` - Update organization (admin/owner only)
- `DELETE /organizations/{orgId}` - Delete organization (owner only)
- `GET /organizations/{orgId}/users` - List organization users

#### 📋 Task Management
- `POST /tasks` - Create new task
- `GET /tasks` - List all tasks
- `GET /tasks/search?q=query` - Search tasks
- `GET /tasks/{taskId}` - Get specific task
- `PUT /tasks/{taskId}` - Update task
- `DELETE /tasks/{taskId}` - Delete task (soft delete)

#### 🛠️ Utilities
- `GET /healthz` - Health check
- `GET /err` - Error endpoint (for testing)

### Example Requests

#### Create User
```bash
curl -X POST http://localhost:8080/v1/user 
  -H "Content-Type: application/json" 
  -d '{
    "username": "johndoe",
    "password": "securepassword123",
    "age": 25,
    "gender": "male"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/v1/login 
  -H "Content-Type: application/json" 
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

#### Create Task
```bash
curl -X POST http://localhost:8080/v1/tasks 
  -H "Content-Type: application/json" 
  -H "Authorization: ApiKey your-api-key" 
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive README and API docs"
  }'
```

## 🎮 CLI Usage

The CLI provides an interactive terminal interface that connects to your API server:

### Getting Started
1. **Start the API Server**: First run the server with `./bin/hello-world`
2. **Run the CLI**: Execute `./bin/task-cli` 
3. **Login**: Use your server credentials to authenticate
4. **Enjoy**: Full access to all server features in a beautiful TUI

### CLI Features
- 🔐 **Server Authentication**: Login with username/password
- 📋 **Task Management**: Create, edit, view, and delete tasks
- 👤 **User Profile**: View your user information and role
- 🏢 **Organization**: Access organization features (if applicable)
- 🔄 **Real-time Sync**: All changes sync with the server immediately

### Navigation
- **Login**: Enter credentials and press Enter
- **Main Menu**: Use number keys or shortcuts to navigate
- **Task List**: Arrow keys to navigate, various shortcuts for actions
- **Create/Edit**: Tab between fields, Enter to save, Esc to cancel

### Commands
- **Main Menu**:
  - `1` or `t` - Task Management
  - `2` or `p` - User Profile  
  - `3` or `o` - Organization
  - `l` - Logout
  - `q` - Quit

- **Task Management**:
  - `n` - Create new task
  - `Enter` - View task details
  - `t` - Toggle task status (finished/unfinished)
  - `d` - Delete task
  - `r` - Refresh task list
  - `b` - Back to main menu

- **Task Details**:
  - `e` - Edit task
  - `t` - Toggle status
  - `b` - Back to task list

### Configuration
The CLI stores configuration in `cli-config.json`:
```json
{
  "server_url": "http://localhost:8080/v1",
  "api_key": "your-api-key-after-login",
  "username": "your-username"
}
```

## 🔧 Development

### Code Generation
This project uses SQLC for type-safe database queries:

```bash
# Generate database code
sqlc generate
```

### Database Migrations
To create new migrations:

1. Add SQL file to `sql/schema/`
2. Update queries in `sql/queries/`
3. Run `sqlc generate`

### Testing
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./handlers
go test ./internal/auth
```

### Build for Production
```bash
# Build optimized binary
go build -ldflags="-w -s" -o bin/hello-world .

# Build CLI
go build -ldflags="-w -s" -o bin/task-cli ./cmd/cli
```

## 🏗️ Architecture

### Clean Architecture Principles
- **Separation of Concerns**: Clear separation between handlers, business logic, and data layer
- **Dependency Injection**: Configurable dependencies through interfaces
- **Middleware Pattern**: Composable request processing pipeline
- **Repository Pattern**: Abstract database operations

### Security Features
- **Password Hashing**: Bcrypt with salt for secure password storage
- **API Key Authentication**: Secure API key generation and validation
- **Role-Based Access Control**: Hierarchical permission system
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries with SQLC

### Performance Optimizations
- **Connection Pooling**: Efficient database connection management
- **Prepared Statements**: Pre-compiled SQL queries
- **Graceful Shutdown**: Proper cleanup of resources
- **Optimized JSON**: Efficient JSON marshaling/unmarshaling

## 🌟 Technologies Used

### Backend
- **Go 1.19+**: Modern, performant language
- **Chi Router**: Lightweight, fast HTTP router
- **PostgreSQL**: Reliable, feature-rich database
- **SQLC**: Type-safe SQL query generation
- **Bcrypt**: Secure password hashing

### CLI
- **Bubble Tea**: Framework for building terminal apps
- **Lip Gloss**: Styling and layout for terminal UIs
- **Bubbles**: Common TUI components

### Development Tools
- **Go Modules**: Dependency management
- **SQLC**: SQL to Go code generation
- **Vendor**: Dependency vendoring for reproducible builds

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and conventions
- Add tests for new features
- Update documentation for API changes
- Use meaningful commit messages
- Keep functions small and focused

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔗 Links

- [Go Documentation](https://golang.org/doc/)
- [Chi Router](https://github.com/go-chi/chi)
- [SQLC](https://sqlc.dev/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [PostgreSQL](https://postgresql.org/)

---

**Built with ❤️ using Go and modern development practices**
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Set Up Environment Variables
```bash
cp .env.example .env
# Edit .env with your database configuration
```

Required environment variables:
- `PORT`: Server port (default: 8080)
- `DB_URL`: PostgreSQL connection string

### 4. Database Setup
Make sure PostgreSQL is running and create a database. The connection string format:
```
postgres://username:password@localhost:5432/database_name?sslmode=disable
```

### 5. Run the Application
```bash
go run main.go
```

The server will start on the configured port (default: 8080).

## 📚 API Documentation

### Base URL
```
http://localhost:8080/v1
```

### Authentication
Most endpoints require authentication via API key in the Authorization header:
```
Authorization: APIKEY your_api_key_here
```

### Endpoints

#### Health Check
```http
GET /v1/healthz
```
**Response:** Returns server status

#### Create User
```http
POST /v1/user
Content-Type: application/json

{
  "username": "john_doe"
}
```
**Response:** Returns user details with API key

#### Get Current User
```http
GET /v1/user
Authorization: APIKEY your_api_key
```
**Response:** Returns authenticated user details

#### Create Task
```http
POST /v1/tasks
Authorization: APIKEY your_api_key
Content-Type: application/json

{
  "title": "Learn Go programming",
  "description": "Complete Go tutorial and build a project"
}
```

#### Get All Tasks
```http
GET /v1/tasks
Authorization: APIKEY your_api_key
```

#### Get Specific Task
```http
GET /v1/tasks/{taskId}
Authorization: APIKEY your_api_key
```

#### Update Task
```http
PUT /v1/tasks/{taskId}
Authorization: APIKEY your_api_key
Content-Type: application/json

{
  "title": "Updated title",
  "description": "Updated description",
  "is_completed": true
}
```

#### Toggle Task Completion
```http
PATCH /v1/tasks/{taskId}/complete
Authorization: APIKEY your_api_key
Content-Type: application/json

{
  "is_completed": true
}
```

#### Delete Task
```http
DELETE /v1/tasks/{taskId}
Authorization: APIKEY your_api_key
```

#### Search Tasks
```http
GET /v1/tasks/search?query=go&limit=10
Authorization: APIKEY your_api_key
```

## 🏗️ Project Structure

```
├── main.go                     # Application entry point
├── .env.example               # Environment configuration template
├── handlers/                  # HTTP handlers (controllers)
│   ├── api_config.go         # Database configuration
│   ├── json.go               # JSON response utilities
│   ├── users.go              # User handlers
│   ├── tasks.go              # Task handlers
│   └── utils.go              # Utility functions
├── internal/
│   ├── auth/                 # Authentication logic
│   │   └── auth.go
│   ├── config/               # Configuration management
│   │   └── config.go
│   ├── middleware/           # HTTP middleware
│   │   ├── auth.go          # Authentication middleware
│   │   └── logging.go       # Logging and recovery middleware
│   └── database/             # Database layer (generated by sqlc)
│       ├── db.go            # Database interface
│       ├── models.go        # Database models
│       ├── users.sql.go     # User queries
│       └── tasks.sql.go     # Task queries
├── models/                   # API response models
│   ├── users.go
│   └── tasks.go
├── sql/                      # Database schema and queries
│   ├── queries/
│   │   ├── users.sql
│   │   └── tasks.sql
│   └── schema/
│       ├── 001_users.sql
│       ├── 002_users_apikey.sql
│       ├── 003_tasks.sql
│       └── 004_tasks_fields.sql
└── vendor/                   # Go modules dependencies
```

## 🎯 Design Patterns & Best Practices

### 1. Clean Architecture
- **Separation of Concerns**: Each layer has a single responsibility
- **Dependency Injection**: Database dependencies are injected
- **Interface Segregation**: Clean interfaces for database operations

### 2. Error Handling
- **Consistent Error Responses**: Standardized JSON error format
- **Proper HTTP Status Codes**: RESTful status code usage
- **Error Logging**: Comprehensive error logging with context

### 3. Security
- **API Key Authentication**: Secure authentication mechanism
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries
- **Authorization**: Users can only access their own resources

### 4. Performance
- **Connection Pooling**: Optimized database connection management
- **Middleware Pipeline**: Efficient request processing
- **Graceful Shutdown**: Proper resource cleanup

### 5. Code Organization
- **Package Structure**: Logical package organization
- **Constants**: Error messages and configuration as constants
- **Regular Expressions**: Compiled once for performance
- **Context Usage**: Request context for cancellation and timeouts

## 🔧 Configuration

The application supports configuration through environment variables:

### Required
- `DB_URL`: PostgreSQL connection string

### Optional (with defaults)
- `PORT`: Server port (default: 8080)
- `SERVER_TIMEOUT`: Graceful shutdown timeout (default: 30s)
- `READ_TIMEOUT`: Request read timeout (default: 10s)
- `WRITE_TIMEOUT`: Response write timeout (default: 15s)
- `IDLE_TIMEOUT`: Connection idle timeout (default: 60s)
- `DB_MAX_OPEN_CONNS`: Max open DB connections (default: 25)
- `DB_MAX_IDLE_CONNS`: Max idle DB connections (default: 25)
- `DB_CONN_MAX_LIFETIME`: Max connection lifetime (default: 5m)

## 🧪 Testing the API

### Using curl

1. **Create a user:**
```bash
curl -X POST http://localhost:8080/v1/user \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser"}'
```

2. **Create a task:**
```bash
curl -X POST http://localhost:8080/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: APIKEY your_api_key_here" \
  -d '{"title": "My first task", "description": "Learning Go API development"}'
```

3. **Get all tasks:**
```bash
curl -X GET http://localhost:8080/v1/tasks \
  -H "Authorization: APIKEY your_api_key_here"
```

## 🚀 Next Steps for Learning

This project is designed to be beginner-friendly while demonstrating production-ready patterns. Here are suggested next steps:

### Beginner Level
1. **Understand the Request Flow**: Follow a request from `main.go` through middleware to handlers
2. **Study the Database Layer**: Examine how SQLC generates type-safe database code
3. **Explore Middleware**: Understand how authentication and logging middleware work
4. **Practice API Testing**: Use the provided curl examples to test all endpoints

### Intermediate Level
1. **Add Unit Tests**: Write tests for handlers and business logic
2. **Implement Pagination**: Add offset/limit pagination to task listing
3. **Add Input Validation**: Implement more sophisticated validation rules
4. **Metrics & Monitoring**: Add Prometheus metrics and health checks

### Advanced Level
1. **Add RBAC**: Implement role-based access control
2. **Caching Layer**: Add Redis for performance optimization
3. **Rate Limiting**: Implement API rate limiting
4. **Database Migrations**: Add automatic database migration system
5. **Containerization**: Add Docker and Docker Compose setup

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 💡 Learning Resources

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Web Examples](https://gowebexamples.com/)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [SQLC Documentation](https://docs.sqlc.dev/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

---

Happy coding! 🎉 This project demonstrates clean Go practices and is perfect for learning modern web API development with Go.

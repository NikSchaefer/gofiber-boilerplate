# 🚀 Go Fiber Boilerplate

> A modern, production-ready Go REST API boilerplate built with [Fiber](https://gofiber.io/), [Ent ORM](https://entgo.io/), and PostgreSQL.

## ✨ Features

- **Authentication & Authorization** - JWT-based sessions with OAuth support
- **Email & SMS Integration** - Resend and Twilio integration
- **Analytics** - PostHog integration for user analytics
- **Database** - PostgreSQL with Ent ORM for type-safe queries
- **Docker Support** - Multi-stage Docker builds
- **Security** - CORS, security headers, input validation
- **OTP Authentication** - One-time password support
- **Production Ready** - Graceful shutdown, proper error handling

## 🏗️ Architecture

```
gofiber-boilerplate/
├── config/          # Configuration management
├── ent/            # Ent ORM generated code
├── internal/       # Private application code
│   ├── database/   # Database connection & setup
│   ├── handlers/   # HTTP request handlers
│   ├── middleware/ # Custom middleware
│   ├── router/     # Route definitions
│   └── services/   # Business logic
├── model/          # Data models
├── pkg/            # Public packages
│   ├── analytics/  # Analytics integration
│   ├── notifications/ # Email/SMS services
│   ├── utils/      # Utility functions
│   └── validator/  # Input validation
└── seeds/          # Database seeders
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24.2 or higher
- PostgreSQL 12 or higher
- Docker (optional)

### 1. Clone the Repository

```bash
git clone https://github.com/NikSchaefer/go-fiber
cd go-fiber
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Set Up Environment Variables

Create a `.env` file in the root directory:

```bash
# Copy the example environment file
cp .env.example .env
```

Configure your environment variables:

```env
# Database Configuration
DATABASE_URL="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"

# Server Configuration
PORT=8000
STAGE=development
ALLOWED_ORIGINS="http://localhost:3000,http://localhost:3001"

# External Services (Optional for development)
POSTHOG_KEY=your_posthog_key_here
RESEND_KEY=your_resend_key_here
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token

# OAuth Configuration (Optional for development)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Application Configuration
APP_DOMAIN=localhost:8000
TWILIO_PHONE_NUMBER=+1234567890
```

### 4. Set Up Database

#### Option A: Using Docker (Recommended)

```bash
# Start PostgreSQL container
docker run --name postgres-db \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=postgres \
  -p 5432:5432 \
  -d postgres:alpine

# Wait a few seconds for the database to start
```

#### Option B: Local PostgreSQL

Make sure PostgreSQL is running and create a database:

```sql
CREATE DATABASE postgres;
```

### 5. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8000`

## 🐳 Docker Deployment

### Build and Run with Docker

```bash
# Build the Docker image
docker build -t go-fiber-app .

# Run the container
docker run -p 8000:8000 \
  --env-file .env \
  --name go-fiber-container \
  go-fiber-app
```

## 📚 API Documentation

### Authentication Endpoints

#### User Registration

```http
POST /auth/signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### Password Login

```http
POST /auth/login/password
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

#### OTP Login Request

```http
POST /auth/login/otp/request
Content-Type: application/json

{
  "email": "john@example.com"
}
```

#### OTP Verification

```http
POST /auth/login/otp/verify
Content-Type: application/json

{
  "email": "john@example.com",
  "otp": "123456"
}
```

#### Logout

```http
DELETE /auth/logout
Cookie: session=<session_token>
```

### User Management

#### Get Current User

```http
GET /users/me
Cookie: session=<session_token>
```

#### Update User Profile

```http
PATCH /users/profile
Cookie: session=<session_token>
Content-Type: application/json

{
  "bio": "Software Developer",
  "location": "San Francisco"
}
```

#### Change Password

```http
POST /auth/password/change
Cookie: session=<session_token>
Content-Type: application/json

{
  "currentPassword": "oldpassword",
  "newPassword": "newpassword123"
}
```

### OAuth Integration

#### Google OAuth

```http
POST /auth/oauth/google
Content-Type: application/json

{
  "redirectUri": "http://localhost:3000/callback"
}
```

## 🔧 Configuration

### Environment Variables

| Variable               | Description                  | Default               | Required |
| ---------------------- | ---------------------------- | --------------------- | -------- |
| `DATABASE_URL`         | PostgreSQL connection string | -                     | ✅       |
| `PORT`                 | Server port                  | `8000`                | ❌       |
| `STAGE`                | Environment stage            | `development`         | ❌       |
| `ALLOWED_ORIGINS`      | CORS allowed origins         | `localhost:3000,3001` | ❌       |
| `POSTHOG_KEY`          | PostHog analytics key        | -                     | ❌       |
| `RESEND_KEY`           | Resend email API key         | -                     | ❌       |
| `TWILIO_ACCOUNT_SID`   | Twilio account SID           | -                     | ❌       |
| `TWILIO_AUTH_TOKEN`    | Twilio auth token            | -                     | ❌       |
| `GOOGLE_CLIENT_ID`     | Google OAuth client ID       | -                     | ❌       |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret   | -                     | ❌       |

### Database Schema

The application uses Ent ORM with the following entities:

- **User** - User accounts and profiles
- **Session** - User sessions and authentication
- **OTP** - One-time passwords for authentication
- **Account** - OAuth account connections
- **Profile** - User profile information

## 🛠️ Development

### Project Structure

```
internal/
├── database/       # Database connection and setup
├── handlers/       # HTTP request handlers
│   ├── auth/       # Authentication handlers
│   └── users/      # User management handlers
├── middleware/     # Custom middleware
│   ├── auth.go     # Authentication middleware
│   ├── security.go # Security headers
│   └── json.go     # JSON parsing middleware
├── router/         # Route definitions
└── services/       # Business logic layer
```

### Adding New Endpoints

1. **Create a handler** in `internal/handlers/`
2. **Add business logic** in `internal/services/`
3. **Define routes** in `internal/router/router.go`
4. **Add validation** in `pkg/validator/`

### Database Migrations

The application uses Ent ORM for database management:

```bash
# Generate Ent code after schema changes
go generate ./ent

# Run migrations (automatic in development)
go run main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/handlers/auth
```

## 🔒 Security Features

- **CORS Protection** - Configurable allowed origins
- **Security Headers** - XSS protection, content type options
- **Input Validation** - Request validation using validator
- **Session Management** - Secure session handling
- **Password Hashing** - bcrypt password hashing
- **Rate Limiting** - Built-in rate limiting (configurable)

## 📊 Monitoring & Analytics

### Health Check

```http
GET /
```

Returns a simple health check response.

### Analytics Integration

The application includes PostHog integration for user analytics:

```go
// Track user events
analytics.Track("user_signed_up", map[string]interface{}{
    "user_id": user.ID,
    "email": user.Email,
})
```

## 🚀 Deployment

### Production Checklist

- [ ] Set `STAGE=production` in environment
- [ ] Configure `ALLOWED_ORIGINS` with your domain
- [ ] Set up SSL/TLS certificates
- [ ] Configure database connection pooling
- [ ] Set up monitoring and logging
- [ ] Configure backup strategy
- [ ] Set up CI/CD pipeline

### Environment-Specific Configurations

#### Development

```env
STAGE=development
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

#### Production

```env
STAGE=production
ALLOWED_ORIGINS=https://yourdomain.com
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Fast HTTP framework
- [Ent](https://entgo.io/) - Type-safe ORM
- [PostgreSQL](https://www.postgresql.org/) - Reliable database
- [Resend](https://resend.com/) - Email delivery
- [Twilio](https://www.twilio.com/) - SMS services
- [PostHog](https://posthog.com/) - Product analytics

## 📞 Support

If you have any questions or need help:

- Create an [issue](https://github.com/NikSchaefer/go-fiber/issues)
- Check the [documentation](https://gofiber.io/)
- Join our [Discord community](https://gofiber.io/discord)

---

**Made with ❤️**

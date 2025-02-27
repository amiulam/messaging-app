# Messaging Application

A modern messaging application built with Go and Fiber framework, featuring various advanced capabilities.

## Technologies Used

- **Go** (version 1.22.2)
- **Fiber** - Fast and flexible web framework
- **MongoDB** - Primary database
- **MySQL** - Relational database
- **WebSocket** - For real-time communication
- **JWT** - For authentication
- **ELK Stack** - For logging and monitoring
- **Docker** - For containerization
- **GitHub Actions** - For CI/CD pipeline

## Key Features

- JWT authentication system
- Real-time messaging using WebSocket
- MongoDB and MySQL database integration
- System logging with ELK Stack
- Input validation using validator
- HTML template rendering
- Environment variables system
- Docker containerization
- Automated CI/CD pipeline

## System Requirements

- Go 1.22 or higher
- MongoDB
- MySQL
- Docker (optional)
- ELK Stack (optional)

## Getting Started

1. Clone this repository
2. Copy `.env.example` to `.env` and adjust the configuration
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

## Environment Configuration

Copy the `.env.example` file to `.env` and adjust the following values:

- `APP_HOST`: Application host (default: localhost)
- `APP_PORT`: Application port (default: 4000)
- [Add other configurations according to .env]

## Project Structure

```
.
├── app/            # Application logic
├── bootstrap/      # Application initialization
├── elk_stack/     # ELK Stack configuration
├── pkg/           # Reusable packages
├── views/         # HTML templates
├── logs/          # Log files
├── main.go        # Application entry point
└── Dockerfile     # Docker configuration
```

## Docker Usage

To run the application using Docker:

```bash
docker build -t messaging-app .
docker run -p 4000:4000 messaging-app
```

## CI/CD Pipeline

The pipeline configuration can be found in `.github/workflows/ci-cd.yml`.

Required secrets for CI/CD:

- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub password
- `PORT`: Application port


# Gateway README

## Introduction
This gateway serves as the entry point to the FinMan application ecosystem. It provides a unified interface for accessing various services through gRPC endpoints and exposes Swagger documentation for easy API reference.

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Build the Docker images:
   ```bash
   docker-compose build
   ```

3. Start the services:
   ```bash
   docker-compose up
   ```

4. Access the gateway at:
   - Gateway endpoint: `http://{gateway-ip}:{gateway-port}`
   - Swagger documentation: `http://{gateway-ip}:{gateway-port}/openapi/`

## Services
The gateway orchestrates access to the following services:

- **Authentication Service**: Handles user authentication and authorization.
- **User Service**: Manages user information and interactions.
- **Role Service**: Manages roles and permissions for users.

## API Documentation
Explore the APIs using Swagger UI:

- Swagger UI: `http://{gateway-ip}:{gateway-port}/openapi/`

## Configuration
Ensure your environment variables are correctly set in the `.env` file for each service:

```dotenv
JWT_SECRET=eDM!":jmx2/QoHBlY'.O8e4?Uy,",9
JWT_EXPIRE_MINUTE=20
PORT=8080
IP=0.0.0.0
USER_SERVICE_ADDR=finman-user-service:8081
```

## Troubleshooting
- If services fail to connect, ensure Docker containers are running and ports are accessible.
- Check network configurations (`docker network ls`) to ensure services are on the same network.

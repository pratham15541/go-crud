# Deployment Guide

This guide covers different deployment strategies for the Go CRUD application.

## Table of Contents

- [Local Development](#local-development)
- [Docker Deployment](#docker-deployment)
- [Production Deployment](#production-deployment)
- [Cloud Deployment](#cloud-deployment)
- [Monitoring and Logging](#monitoring-and-logging)
- [Security Considerations](#security-considerations)

## Local Development

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Git

### Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/pratham15541/go-crud.git
   cd go-crud
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Setup environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start PostgreSQL:**
   ```bash
   # Using Docker
   docker run -d \
     --name postgres \
     -e POSTGRES_PASSWORD=password \
     -e POSTGRES_DB=crud_demo \
     -p 5432:5432 \
     postgres:15-alpine
   ```

5. **Run migrations:**
   ```bash
   chmod +x scripts/migrate.sh
   ./scripts/migrate.sh up
   ```

6. **Start the application:**
   ```bash
   make dev
   # or
   go run cmd/server/main.go
   ```

## Docker Deployment

### Using Docker Compose (Recommended)

1. **Start all services:**
   ```bash
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   docker-compose logs -f app
   ```

3. **Stop services:**
   ```bash
   docker-compose down
   ```

### Manual Docker Build

1. **Build the image:**
   ```bash
   docker build -t go-crud:latest .
   ```

2. **Run the container:**
   ```bash
   docker run -d \
     --name go-crud-app \
     -p 8080:8080 \
     --env-file .env \
     go-crud:latest
   ```

## Production Deployment

### Server Setup

1. **System Requirements:**
   - CPU: 2+ cores
   - RAM: 4GB+ recommended
   - Storage: 20GB+ available space
   - OS: Ubuntu 20.04+ or similar

2. **Install dependencies:**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sh get-docker.sh
   sudo usermod -aG docker $USER
   
   # Install Docker Compose
   sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   ```

3. **Setup application:**
   ```bash
   # Create application directory
   sudo mkdir -p /opt/go-crud
   cd /opt/go-crud
   
   # Clone repository
   git clone https://github.com/pratham15541/go-crud.git .
   
   # Setup environment
   cp .env.example .env
   # Edit .env for production settings
   ```

4. **Production environment variables:**
   ```bash
   # Server Configuration
   HOST=0.0.0.0
   PORT=8080
   GIN_MODE=release
   
   # Database Configuration
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=<strong-password>
   DB_NAME=crud_demo
   DB_SSLMODE=require
   
   # JWT Configuration
   JWT_SECRET=<strong-jwt-secret-key>
   JWT_EXPIRATION=24h
   
   # Logging
   LOG_LEVEL=info
   LOG_FORMAT=json
   ```

5. **Start production services:**
   ```bash
   docker-compose -f docker-compose.yml up -d
   ```

### Reverse Proxy (Nginx)

1. **Install Nginx:**
   ```bash
   sudo apt install nginx
   ```

2. **Configure Nginx:**
   ```nginx
   # /etc/nginx/sites-available/go-crud
   server {
       listen 80;
       server_name your-domain.com;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```

3. **Enable site:**
   ```bash
   sudo ln -s /etc/nginx/sites-available/go-crud /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```

### SSL with Let's Encrypt

1. **Install Certbot:**
   ```bash
   sudo apt install certbot python3-certbot-nginx
   ```

2. **Obtain SSL certificate:**
   ```bash
   sudo certbot --nginx -d your-domain.com
   ```

## Cloud Deployment

### AWS Deployment

#### Using AWS ECS

1. **Build and push Docker image:**
   ```bash
   # Build image
   docker build -t go-crud:latest .
   
   # Tag for ECR
   docker tag go-crud:latest 123456789012.dkr.ecr.us-east-1.amazonaws.com/go-crud:latest
   
   # Push to ECR
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 123456789012.dkr.ecr.us-east-1.amazonaws.com
   docker push 123456789012.dkr.ecr.us-east-1.amazonaws.com/go-crud:latest
   ```

2. **Create ECS task definition:**
   ```json
   {
     "family": "go-crud",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "256",
     "memory": "512",
     "executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
     "containerDefinitions": [
       {
         "name": "go-crud",
         "image": "123456789012.dkr.ecr.us-east-1.amazonaws.com/go-crud:latest",
         "portMappings": [
           {
             "containerPort": 8080,
             "protocol": "tcp"
           }
         ],
         "environment": [
           {
             "name": "DB_HOST",
             "value": "your-rds-endpoint"
           }
         ],
         "logConfiguration": {
           "logDriver": "awslogs",
           "options": {
             "awslogs-group": "/ecs/go-crud",
             "awslogs-region": "us-east-1",
             "awslogs-stream-prefix": "ecs"
           }
         }
       }
     ]
   }
   ```

#### Using AWS RDS for Database

1. **Create RDS PostgreSQL instance:**
   ```bash
   aws rds create-db-instance \
     --db-instance-identifier go-crud-db \
     --db-instance-class db.t3.micro \
     --engine postgres \
     --master-username postgres \
     --master-user-password <password> \
     --allocated-storage 20 \
     --vpc-security-group-ids sg-xxxxxxxxx
   ```

### Google Cloud Platform

#### Using Cloud Run

1. **Build and deploy:**
   ```bash
   # Build and push to Google Container Registry
   gcloud builds submit --tag gcr.io/PROJECT-ID/go-crud
   
   # Deploy to Cloud Run
   gcloud run deploy go-crud \
     --image gcr.io/PROJECT-ID/go-crud \
     --platform managed \
     --region us-central1 \
     --allow-unauthenticated
   ```

### Heroku

1. **Create Heroku app:**
   ```bash
   heroku create go-crud-app
   ```

2. **Add PostgreSQL addon:**
   ```bash
   heroku addons:create heroku-postgresql:hobby-dev
   ```

3. **Deploy:**
   ```bash
   git push heroku main
   ```

## Monitoring and Logging

### Health Checks

The application provides a health check endpoint at `/api/v1/health`. Configure your load balancer or orchestrator to use this endpoint.

### Logging

1. **Application logs** are written to stdout in JSON format.
2. **Access logs** are handled by the logging middleware.
3. **Error logs** include stack traces and context.

### Metrics

Implement monitoring using:
- **Prometheus** for metrics collection
- **Grafana** for visualization
- **AlertManager** for alerting

## Security Considerations

### Environment Variables

- Use strong passwords and JWT secrets
- Never commit secrets to version control
- Use secret management services in production

### Database Security

- Enable SSL/TLS for database connections
- Use connection pooling
- Implement proper backup strategies

### Network Security

- Use HTTPS in production
- Implement rate limiting
- Configure firewalls properly

### Application Security

- Keep dependencies updated
- Implement input validation
- Use security headers
- Enable CORS properly

## Troubleshooting

### Common Issues

1. **Database connection failed:**
   - Check database credentials
   - Verify network connectivity
   - Check firewall settings

2. **Application won't start:**
   - Check environment variables
   - Verify port availability
   - Check application logs

3. **High memory usage:**
   - Check for memory leaks
   - Tune database connection pool
   - Monitor goroutines

### Debugging

1. **Enable debug logging:**
   ```bash
   export LOG_LEVEL=debug
   ```

2. **Check application health:**
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

3. **View application logs:**
   ```bash
   docker-compose logs -f app
   ```
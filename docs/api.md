# API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Common Response Format

### Success Response
```json
{
  "message": "Operation successful",
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "error": "Error Type",
  "message": "Detailed error message",
  "code": 400
}
```

## Endpoints

### Health Check

#### GET /health
Check the health status of the application and its dependencies.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-08-11T05:34:07Z",
  "version": "1.0.0",
  "uptime": "1h30m45s",
  "checks": {
    "database": {
      "status": "healthy",
      "error": ""
    },
    "memory": {
      "status": "healthy"
    }
  }
}
```

### Users

#### POST /users
Create a new user.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30
}
```

**Response (201 Created):**
```json
{
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2025-08-11T05:34:07Z",
    "updated_at": "2025-08-11T05:34:07Z"
  }
}
```

#### GET /users
Retrieve all users with pagination.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Number of users per page (default: 10, max: 100)

**Example:**
```
GET /users?page=1&limit=10
```

**Response (200 OK):**
```json
{
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "age": 30,
        "created_at": "2025-08-11T05:34:07Z",
        "updated_at": "2025-08-11T05:34:07Z"
      }
    ],
    "pagination": {
      "total": 1,
      "page": 1,
      "limit": 10
    }
  }
}
```

#### GET /users/{id}
Retrieve a specific user by ID.

**Path Parameters:**
- `id`: User ID (integer)

**Response (200 OK):**
```json
{
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30,
    "created_at": "2025-08-11T05:34:07Z",
    "updated_at": "2025-08-11T05:34:07Z"
  }
}
```

#### PUT /users/{id}
Update an existing user.

**Path Parameters:**
- `id`: User ID (integer)

**Request Body:**
```json
{
  "name": "John Smith",
  "email": "johnsmith@example.com",
  "age": 31
}
```

**Note:** All fields are optional. Only provided fields will be updated.

**Response (200 OK):**
```json
{
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "name": "John Smith",
    "email": "johnsmith@example.com",
    "age": 31,
    "created_at": "2025-08-11T05:34:07Z",
    "updated_at": "2025-08-11T05:35:00Z"
  }
}
```

#### DELETE /users/{id}
Delete a user.

**Path Parameters:**
- `id`: User ID (integer)

**Response (200 OK):**
```json
{
  "message": "User deleted successfully",
  "data": null
}
```

## HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `422 Unprocessable Entity` - Validation error
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

## Rate Limiting

The API implements rate limiting to prevent abuse:
- 100 requests per minute per IP address
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`: Request limit
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset time (Unix timestamp)

## CORS

Cross-Origin Resource Sharing (CORS) is enabled for:
- Origins: `http://localhost:3000`, `http://localhost:8080`
- Methods: `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`
- Headers: `Content-Type`, `Authorization`

## Validation Rules

### User Validation
- **Name**: Required, 2-100 characters
- **Email**: Required, valid email format, unique
- **Age**: Required, integer between 1-150

## Error Codes

| Code | Description |
|------|-------------|
| 1001 | Invalid JSON payload |
| 1002 | Validation error |
| 1003 | User not found |
| 1004 | Email already exists |
| 1005 | Database connection error |
| 1006 | Authentication failed |
| 1007 | Access denied |

## Examples

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "age": 28
  }'
```

### Get All Users
```bash
curl http://localhost:8080/api/v1/users?page=1&limit=5
```

### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "age": 29
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```
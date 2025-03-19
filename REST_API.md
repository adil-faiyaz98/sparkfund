# Money-Pulse REST API Documentation

This document outlines the REST API endpoints provided by the Money-Pulse service.

## Base URL

```
https://[your-api-gateway-id].execute-api.[your-region].amazonaws.com/[stage]
```

## Endpoints

### GET /health

Health check endpoint to verify the API is operational.

**Response:**

```json
{
  "status": "healthy"
}
```

### [Add your other endpoints here]

## Authentication

[Document your authentication method here]

## Error Handling

All errors return appropriate HTTP status codes with a JSON body containing:

```json
{
  "error": "Error description"
}
```

## Deployment

This API is deployed as an AWS Lambda function behind API Gateway, providing a serverless REST architecture.

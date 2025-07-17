# API Routes Documentation

This document describes all available API routes, their parameters, request bodies, and responses.

## Base URL
All routes are prefixed with `/v1` unless otherwise specified.

## Authentication
Most routes require authentication using JWT tokens passed in the `Authorization` header:
```
Authorization: Bearer <token>
```

## Status Routes (Public)

### GET `/status`
**Description**: Simple health check endpoint  
**Authentication**: None  
**Parameters**: None  
**Response**:
```json
{
  "status": "ok",
  "timestamp": "2023-12-01T14:30:22Z"
}
```

### GET `/health`
**Description**: Detailed health check with service status  
**Authentication**: None  
**Parameters**: None  
**Response**:
```json
{
  "status": "healthy",
  "timestamp": "2023-12-01T14:30:22Z",
  "services": {
    "database": "healthy",
    "cache": "healthy",
    "application": "healthy"
  }
}
```

## User Authentication Routes

### POST `/v1/users/login`
**Description**: User login  
**Authentication**: None  
**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
**Response**:
```json
{
  "token": "jwt-token-here"
}
```

### POST `/v1/users/logout`
**Description**: User logout  
**Authentication**: Required (sales, driver)  
**Headers**: `Authorization: Bearer <token>`  
**Response**:
```json
{
  "message": "Successfully logged out"
}
```

### GET `/v1/users/me`
**Description**: Get current user profile  
**Authentication**: Required (sales, driver)  
**Response**:
```json
{
  "user_id": "user-uuid",
  "email": "user@example.com",
  "role": "driver",
  "profile": {
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890",
    "last_connection": "2023-12-01T14:25:00Z"
  }
}
```

### PUT `/v1/users/reset-password`
**Description**: Reset user password (requires current password)  
**Authentication**: Required (sales, driver)  
**Request Body**:
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}
```
**Response**:
```json
{
  "message": "Password reset successfully"
}
```

## User Management Routes

### POST `/v1/users`
**Description**: Register new user  
**Authentication**: Required (any role)  
**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+1234567890",
  "role": "driver"
}
```
**Response**:
```json
{
  "message": "User registered successfully",
  "user": {
    "user_id": "user-uuid",
    "email": "user@example.com",
    "role": "driver"
  }
}
```

### GET `/v1/users`
**Description**: Get all users with pagination  
**Authentication**: Required (sales)  
**Query Parameters**:
- `limit` (optional): Number of results per page (default: 100)
- `offset` (optional): Number of results to skip (default: 0)

**Response**:
```json
{
  "users": [
    {
      "user_id": "user-uuid",
      "email": "user@example.com",
      "role": "driver",
      "first_name": "John",
      "last_name": "Doe",
      "phone_number": "+1234567890",
      "last_connection": "2023-12-01T14:25:00Z",
      "created_at": "2023-11-01T10:00:00Z",
      "updated_at": "2023-12-01T14:25:00Z"
    }
  ],
  "total": 1,
  "limit": 100,
  "offset": 0
}
```

### GET `/v1/users/:id`
**Description**: Get specific user by ID  
**Authentication**: Required (any role)  
**Path Parameters**:
- `id`: User UUID

**Response**:
```json
{
  "user": {
    "user_id": "user-uuid",
    "email": "user@example.com",
    "role": "driver",
    "created_at": "2023-11-01T10:00:00Z",
    "updated_at": "2023-12-01T14:25:00Z"
  },
  "profile": {
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890",
    "last_connection": "2023-12-01T14:25:00Z"
  }
}
```

### GET `/v1/users/:id/orders`
**Description**: Get orders assigned to a specific user (driver only)  
**Authentication**: Required (driver role)  
**Path Parameters**:
- `id`: User UUID (must match authenticated user's ID)

**Response**:
```json
[
  {
    "order_id": "order-uuid",
    "order_number": "123456",
    "email": "customer@example.com",
    "address": "123 Main St, City, State",
    "status": "pending",
    "assigned_to": "driver-uuid",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T14:25:00Z"
  }
]
```

### DELETE `/v1/users/:id`
**Description**: Delete user  
**Authentication**: Required (any role)  
**Path Parameters**:
- `id`: User UUID

**Response**:
```json
{
  "message": "User deleted successfully"
}
```

## User Location Routes (Driver Only)

### PUT `/v1/users/location`
**Description**: Update driver location  
**Authentication**: Required (driver)  
**Request Body**:
```json
{
  "latitude": 40.7128,
  "longitude": -74.0060
}
```
**Response**:
```json
{
  "message": "Location updated successfully"
}
```

### GET `/v1/users/location`
**Description**: Get current driver location  
**Authentication**: Required (driver)  
**Response**:
```json
{
  "latitude": 40.7128,
  "longitude": -74.0060,
  "updated_at": "2023-12-01T14:30:22Z"
}
```

### GET `/v1/users/drivers/active`
**Description**: Get all active drivers with locations (last 10 minutes)  
**Authentication**: Required (sales)  
**Response**:
```json
{
  "drivers": [
    {
      "user_id": "driver-uuid",
      "email": "driver@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "phone_number": "+1234567890",
      "latitude": 40.7128,
      "longitude": -74.0060,
      "last_connection": "2023-12-01T14:25:00Z",
      "location_updated_at": "2023-12-01T14:20:00Z"
    }
  ],
  "count": 1
}
```

### GET `/v1/users/drivers`
**Description**: Get all drivers ordered by last connection  
**Authentication**: Required (sales)  
**Response**:
```json
{
  "drivers": [
    {
      "user_id": "driver-uuid",
      "email": "driver@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "phone_number": "+1234567890",
      "role": "driver",
      "last_connection": "2023-12-01T14:25:00Z",
      "created_at": "2023-11-01T10:00:00Z",
      "updated_at": "2023-12-01T14:25:00Z"
    }
  ],
  "count": 1
}
```

## Order Management Routes

### POST `/v1/orders`
**Description**: Create new order  
**Authentication**: Required (sales)  
**Request Body**:
```json
{
  "email": "customer@example.com",
  "address": "123 Main St, City, State",
  "order_number": "123456"
}
```
**Notes**: 
- `order_number` must be 1-6 characters long and contain only numbers

**Response**:
```json
{
  "order_id": "order-uuid",
  "order_number": "123456",
  "created_by": "user-uuid",
  "created_by_username": "sales@example.com",
  "assigned_to": null,
  "assigned_to_username": null,
  "delivery_address": "123 Main St, City, State",
  "status": "pending",
  "created_at": "2023-12-01T14:30:22Z",
  "updated_at": "2023-12-01T14:30:22Z"
}
```

**Error Responses**:
- **400 Bad Request**: Invalid order_number format
```json
{
  "error": "order_number must contain only numbers"
}
```
- **409 Conflict**: Duplicate order number
```json
{
  "error": "Ya existe una orden con este número"
}
```

### GET `/v1/orders`
**Description**: Get all orders with optional pagination  
**Authentication**: Required (sales, admin)  
**Query Parameters**:
- `limit` (optional): Number of results per page
- `offset` (optional): Number of results to skip

**Response**:
```json
{
  "orders": [
    {
      "order_id": "order-uuid",
      "order_number": "123456",
      "created_by": "user-uuid",
      "created_by_username": "sales@example.com",
      "assigned_to": "driver-uuid",
      "assigned_to_username": "driver@example.com",
      "delivery_address": "123 Main St, City, State",
      "status": "pending",
      "created_at": "2023-12-01T14:30:22Z",
      "updated_at": "2023-12-01T14:30:22Z"
    }
  ],
  "total": 1
}
```

### GET `/v1/orders/:id`
**Description**: Get specific order by ID  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID

**Response**:
```json
{
  "order_id": "order-uuid",
  "order_number": "123456",
  "created_by": "user-uuid",
  "created_by_username": "sales@example.com",
  "assigned_to": "driver-uuid",
  "assigned_to_username": "driver@example.com",
  "delivery_address": "123 Main St, City, State",
  "status": "pending",
  "created_at": "2023-12-01T14:30:22Z",
  "updated_at": "2023-12-01T14:30:22Z"
}
```

### PUT `/v1/orders/:id`
**Description**: Update order  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID

**Request Body**:
```json
{
  "assigned_to": "driver-uuid",
  "address": "456 New St, City, State",
  "status": "in_progress"
}
```

**Notes**:
- `assigned_to` can be set to `null` to unassign the order from any driver
- All fields are optional - only provided fields will be updated

**Response**:
```json
{
  "order_id": "order-uuid",
  "order_number": "123456",
  "created_by": "user-uuid",
  "created_by_username": "sales@example.com",
  "assigned_to": "driver-uuid",
  "assigned_to_username": "driver@example.com",
  "delivery_address": "456 New St, City, State",
  "status": "in_progress",
  "created_at": "2023-12-01T14:30:22Z",
  "updated_at": "2023-12-01T14:35:22Z"
}
```

### PATCH `/v1/orders/:id/status`
**Description**: Update only the status of an order  
**Authentication**: Required (sales, driver)  
**Path Parameters**:
- `id`: Order UUID

**Request Body**:
```json
{
  "status": "completed"
}
```

**Valid Status Values**:
- `pending`
- `in_progress`
- `completed`
- `cancelled`

**Response**:
```json
{
  "order_id": "order-uuid",
  "order_number": "123456",
  "created_by": "user-uuid",
  "created_by_username": "sales@example.com",
  "assigned_to": "driver-uuid",
  "assigned_to_username": "driver@example.com",
  "delivery_address": "123 Main St, City, State",
  "status": "completed",
  "created_at": "2023-12-01T14:30:22Z",
  "updated_at": "2023-12-01T14:40:22Z"
}
```

### DELETE `/v1/orders/:id`
**Description**: Delete order  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID

**Response**:
```json
{
  "message": "Order deleted successfully"
}
```

## Order Items Routes

### POST `/v1/orders/:id/items`
**Description**: Add item to order  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID

**Request Body**:
```json
{
  "product_name": "Product Name",
  "quantity": 5
}
```
**Response**:
```json
{
  "item_id": "item-uuid",
  "order_id": "order-uuid",
  "product_name": "Product Name",
  "quantity": 5,
  "created_at": "2023-12-01T14:30:22Z",
  "updated_at": "2023-12-01T14:30:22Z"
}
```

### POST `/v1/orders/:id/items/bulk`
**Description**: Add multiple items to order  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID

**Request Body**:
```json
{
  "items": [
    {
      "product_name": "Product 1",
      "quantity": 3
    },
    {
      "product_name": "Product 2",
      "quantity": 2
    }
  ]
}
```
**Response**:
```json
{
  "message": "Items added successfully",
  "items": [
    {
      "item_id": "item-uuid-1",
      "order_id": "order-uuid",
      "product_name": "Product 1",
      "quantity": 3,
      "created_at": "2023-12-01T14:30:22Z",
      "updated_at": "2023-12-01T14:30:22Z"
    },
    {
      "item_id": "item-uuid-2",
      "order_id": "order-uuid",
      "product_name": "Product 2",
      "quantity": 2,
      "created_at": "2023-12-01T14:30:22Z",
      "updated_at": "2023-12-01T14:30:22Z"
    }
  ]
}
```

### GET `/v1/orders/:id/items`
**Description**: Get all items for an order  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID

**Response**:
```json
{
  "items": [
    {
      "item_id": "item-uuid",
      "order_id": "order-uuid",
      "product_name": "Product Name",
      "quantity": 5,
      "created_at": "2023-12-01T14:30:22Z",
      "updated_at": "2023-12-01T14:30:22Z"
    }
  ],
  "total": 1
}
```

### PUT `/v1/orders/:id/items/:item_id`
**Description**: Update order item  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID
- `item_id`: Item UUID

**Request Body**:
```json
{
  "product_name": "Updated Product Name",
  "quantity": 10
}
```
**Response**:
```json
{
  "item_id": "item-uuid",
  "order_id": "order-uuid",
  "product_name": "Updated Product Name",
  "quantity": 10,
  "created_at": "2023-12-01T14:30:22Z",
  "updated_at": "2023-12-01T14:35:22Z"
}
```

### DELETE `/v1/orders/:id/items/:item_id`
**Description**: Delete order item  
**Authentication**: Required (sales)  
**Path Parameters**:
- `id`: Order UUID
- `item_id`: Item UUID

**Response**:
```json
{
  "message": "Item deleted successfully"
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

Common validation errors:
- `order_number must contain only numbers` - Order number contains non-numeric characters
- Various field validation errors based on binding rules

### 401 Unauthorized
```json
{
  "error": "User not authenticated"
}
```

### 403 Forbidden
```json
{
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

## Role-Based Access Control

- **Public**: No authentication required
- **Any Role**: Any authenticated user can access
- **Sales**: Users with 'sales' role
- **Driver**: Users with 'driver' role
- **Admin**: Users with 'admin' role

## Data Types

- **UUID**: Standard UUID format (e.g., `123e4567-e89b-12d3-a456-426614174000`)
- **Timestamp**: ISO 8601 format (e.g., `2023-12-01T14:30:22Z`)
- **Order Number**: 1-6 digit numeric string (e.g., `"123456"`, `"1"`, `"000123"`)
- **Coordinates**: Decimal degrees (e.g., `40.7128` for latitude, `-74.0060` for longitude)
- **Order Status**: `pending`, `in_progress`, `completed`, `cancelled`
- **User Role**: `admin`, `driver`, `sales`

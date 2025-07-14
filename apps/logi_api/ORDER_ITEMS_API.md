# Order Items API Documentation

## Overview

The Order Items functionality has been fully implemented with comprehensive CRUD operations and bulk operations support. This includes individual item management and bulk adding of multiple items to an order.

## API Endpoints

### Order Item Management

#### 1. Add Single Item to Order
```
POST /v1/orders/{id}/items
```
**Authentication**: Required (sales role)

**Request Body**:
```json
{
  "product_name": "Product Name",
  "quantity": 5
}
```

**Response** (201 Created):
```json
{
  "item_id": "uuid",
  "order_id": "uuid",
  "product_name": "Product Name",
  "quantity": 5,
  "created_at": "2025-07-13T22:30:00Z",
  "updated_at": "2025-07-13T22:30:00Z"
}
```

#### 2. Add Multiple Items to Order (Bulk)
```
POST /v1/orders/{id}/items/bulk
```
**Authentication**: Required (sales role)

**Request Body**:
```json
{
  "items": [
    {
      "product_name": "Product A",
      "quantity": 3
    },
    {
      "product_name": "Product B",
      "quantity": 7
    }
  ]
}
```

**Response** (201 Created):
```json
{
  "items": [
    {
      "item_id": "uuid1",
      "order_id": "uuid",
      "product_name": "Product A",
      "quantity": 3,
      "created_at": "2025-07-13T22:30:00Z",
      "updated_at": "2025-07-13T22:30:00Z"
    },
    {
      "item_id": "uuid2",
      "order_id": "uuid",
      "product_name": "Product B",
      "quantity": 7,
      "created_at": "2025-07-13T22:30:00Z",
      "updated_at": "2025-07-13T22:30:00Z"
    }
  ],
  "total": 2
}
```

#### 3. Get All Items for an Order
```
GET /v1/orders/{id}/items
```
**Authentication**: Required (sales role)

**Response** (200 OK):
```json
{
  "items": [
    {
      "item_id": "uuid",
      "order_id": "uuid",
      "product_name": "Product Name",
      "quantity": 5,
      "created_at": "2025-07-13T22:30:00Z",
      "updated_at": "2025-07-13T22:30:00Z"
    }
  ],
  "total": 1
}
```

#### 4. Update Order Item
```
PUT /v1/orders/{id}/items/{item_id}
```
**Authentication**: Required (sales role)

**Request Body** (all fields optional):
```json
{
  "product_name": "Updated Product Name",
  "quantity": 10
}
```

**Response** (200 OK):
```json
{
  "item_id": "uuid",
  "order_id": "uuid", 
  "product_name": "Updated Product Name",
  "quantity": 10,
  "created_at": "2025-07-13T22:30:00Z",
  "updated_at": "2025-07-13T22:35:00Z"
}
```

#### 5. Delete Order Item
```
DELETE /v1/orders/{id}/items/{item_id}
```
**Authentication**: Required (sales role)

**Response** (200 OK):
```json
{
  "message": "Order item deleted successfully"
}
```

### Enhanced Order Endpoint

#### Get Order with Items
```
GET /v1/orders/{id}
```
**Authentication**: Required (sales role)

**Response** (200 OK):
```json
{
  "order_id": "uuid",
  "created_by": "user_uuid",
  "assigned_to": null,
  "delivery_address": "123 Main St",
  "status": "pending",
  "created_at": "2025-07-13T22:30:00Z",
  "updated_at": "2025-07-13T22:30:00Z",
  "items": [
    {
      "item_id": "uuid",
      "order_id": "uuid",
      "product_name": "Product Name",
      "quantity": 5,
      "created_at": "2025-07-13T22:30:00Z",
      "updated_at": "2025-07-13T22:30:00Z"
    }
  ]
}
```

## Implementation Details

### Architecture

The implementation follows a clean architecture pattern with clear separation of concerns:

1. **Models** (`order_model.go`):
   - `Order` struct with `Items []OrderItem` field
   - `OrderItem` struct representing individual items

2. **Repository Interface** (`order_repository.go`):
   - Extended with order item CRUD operations
   - Bulk operations support
   - Transaction support for multiple items

3. **Service Layer** (`order_service.go`):
   - Business logic validation
   - Order existence checks
   - Quantity validation
   - UUID generation

4. **Database Layer** (`database/order_repository.go`):
   - Efficient database queries
   - Transaction support for bulk operations
   - Proper error handling

5. **HTTP Handlers** (`handlers/order_handler.go`):
   - Request validation
   - Response formatting
   - Error handling

### Database Schema

The implementation uses the existing `order_items` table:

```sql
CREATE TABLE IF NOT EXISTS order_items (
    item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES orders(order_id) ON DELETE CASCADE,
    product_name TEXT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### Key Features

1. **Bulk Operations**: The `/bulk` endpoint allows adding multiple items in a single transaction
2. **Validation**: Comprehensive validation at all layers
3. **Error Handling**: Proper error messages and HTTP status codes
4. **Transaction Safety**: Bulk operations use database transactions
5. **Performance**: Efficient queries and minimal database calls
6. **Security**: All endpoints protected with authentication middleware

### Validation Rules

- **Product Name**: Required, non-empty string
- **Quantity**: Required, must be greater than 0
- **Order**: Must exist before adding items
- **Item**: Must exist before updating/deleting

### Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message describing what went wrong"
}
```

Common error scenarios:
- 400 Bad Request: Invalid input data
- 404 Not Found: Order or item not found
- 500 Internal Server Error: Database or system errors

## Testing

To test the functionality:

1. **Create an order** using the existing order creation endpoint
2. **Add items** using either single or bulk endpoints
3. **Retrieve order** with items using the enhanced GET endpoint
4. **Update items** as needed
5. **Delete items** when no longer needed

Example workflow:
```bash
# 1. Create order
curl -X POST /v1/orders \
  -H "Authorization: Bearer <token>" \
  -d '{"email": "user@example.com", "address": "123 Main St"}'

# 2. Add items (bulk)
curl -X POST /v1/orders/{id}/items/bulk \
  -H "Authorization: Bearer <token>" \
  -d '{"items": [{"product_name": "Product A", "quantity": 3}]}'

# 3. Get order with items
curl -X GET /v1/orders/{id} \
  -H "Authorization: Bearer <token>"
```

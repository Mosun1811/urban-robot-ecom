# FutureMarket API

A production-style e-commerce backend built in Go, featuring authentication, product catalog, shopping cart, orders, reviews, and deployment to Render with PostgreSQL.

## Tech Stack

- **Language:** Go 1.22
- **Frameworks/Libraries:** net/http, gorilla/mux, GORM
- **Database:** PostgreSQL (Render managed DB)
- **Auth:** JWT with token blacklist for logout
- **Deployment:** Docker + Render Web Service
- **Other:** gomod, environment-based config

## Features

### Auth & Users
- User registration with:
  - Email format validation
  - Name validation (length & characters)
  - Strong password policy (length + upper/lower/number/special)
- JWT login (24h expiry) and logout via token blacklist.
- Admin seeding (`admin@futuremarket.com / AdminPass123!`).

### Product Catalog
- List products with pagination & filters:
  - `GET /api/v1/products?page=&limit=&min_price=&max_price=&category=`
- Get product details:
  - `GET /api/v1/products/{id}`
- Admin product management:
  - `POST /api/v1/admin/products`
  - `PATCH /api/v1/admin/products/{id}`
- Stock tracking via separate `stocks` table.

### Cart & Orders
- Authenticated cart endpoints:
  - `GET /api/v1/cart`
  - `POST /api/v1/cart` (add item)
  - `PATCH /api/v1/cart/{product_id}` (update qty)
  - `DELETE /api/v1/cart/{product_id}` (remove item)
- Stock validation against real `Stock` records.
- Checkout:
  - `POST /api/v1/checkout`
  - Creates `orders` + `order_items` and decrements stock inside a DB transaction.
- Order history:
  - `GET /api/v1/orders`
  - `GET /api/v1/orders/paginated?page=&limit=`

### Reviews & Ratings
- Public, paginated reviews:
  - `GET /api/v1/products/{id}/reviews?page=&limit=`
- Authenticated create/update:
  - `POST /api/v1/products/{id}/reviews`
- One review per user per product (update instead of duplicate).
- Denormalised rating fields stored on `products`:
  - `average_rating`, `review_count`.

## Local Development

### 1. Clone

```bash
git clone git@github.com:Mariana-Tech-Academy/urban-robot.git
cd urban-robot/urban-robot-1



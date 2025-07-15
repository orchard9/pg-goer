# PostgreSQL Database Documentation

Generated on: 2025-07-14 21:50:16

## Table of Contents

- [Database Summary](#database-summary)
- [Database Relationships](#database-relationships)
- [Tables](#tables)
  - [categories](#categories)
  - [order_items](#order-items)
  - [orders](#orders)
  - [products](#products)
  - [users](#users)

## Database Summary

**Total Tables:** 5
**Total Rows:** 47

## Database Relationships

```mermaid
erDiagram
    categories ||--o{ categories : "parent_id"
    orders ||--o{ order_items : "order_id"
    users ||--o{ orders : "user_id"
    categories ||--o{ products : "category_id"

    categories {
        integer id PK
        varchar name UK
        text description
        integer parent_id
    }
    order_items {
        integer id PK
        integer order_id
        varchar product_name
        integer quantity
        decimal unit_price
        decimal total_price
    }
    orders {
        integer id PK
        integer user_id
        decimal total
        timestamp order_date
        varchar status
        text shipping_address
        text notes
    }
    products {
        integer id PK
        varchar name
        text description
        decimal price
        integer category_id
        boolean in_stock
        timestamp created_at
    }
    users {
        integer id PK
        varchar email UK
        varchar first_name
        varchar last_name
        timestamp created_at
        varchar status
        boolean is_verified
    }
```

## Tables

## categories

<a id="categories"></a>

Row Count: 6

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY | nextval('categories_id_seq'::regclass) |
| name | character varying(100) | NO | UNIQUE |  |
| description | text | YES |  |  |
| parent_id | integer | YES |  |  |

---

## order_items

<a id="order-items"></a>

Row Count: 13

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY | nextval('order_items_id_seq'::regclass) |
| order_id | integer | NO |  |  |
| product_name | character varying(255) | NO |  |  |
| quantity | integer | NO |  | 1 |
| unit_price | numeric | NO |  |  |
| total_price | numeric | YES |  |  |

---

## orders

<a id="orders"></a>

Row Count: 10

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY | nextval('orders_id_seq'::regclass) |
| user_id | integer | NO |  |  |
| total | numeric | NO |  | 0.00 |
| order_date | timestamp without time zone | YES |  | CURRENT_TIMESTAMP |
| status | character varying(20) | YES |  | 'pending'::character varying |
| shipping_address | text | YES |  |  |
| notes | text | YES |  |  |

---

## products

<a id="products"></a>

Row Count: 8

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY | nextval('products_id_seq'::regclass) |
| name | character varying(255) | NO |  |  |
| description | text | YES |  |  |
| price | numeric | NO |  |  |
| category_id | integer | YES |  |  |
| in_stock | boolean | YES |  | true |
| created_at | timestamp without time zone | YES |  | CURRENT_TIMESTAMP |

---

## users

<a id="users"></a>

Row Count: 10

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY | nextval('users_id_seq'::regclass) |
| email | character varying(255) | NO | UNIQUE |  |
| first_name | character varying(100) | NO |  |  |
| last_name | character varying(100) | NO |  |  |
| created_at | timestamp without time zone | YES |  | CURRENT_TIMESTAMP |
| status | character varying(20) | YES |  | 'active'::character varying |
| is_verified | boolean | YES |  | false |

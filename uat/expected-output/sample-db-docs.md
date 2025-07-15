# PostgreSQL Database Documentation

Generated on: 2024-01-01 12:00:00

## Table of Contents

- [Database Summary](#database-summary)
- [Database Relationships](#database-relationships)
- [Tables](#tables)
  - [users](#users)
  - [orders](#orders)
  - [order_items](#order-items)
  - [categories](#categories)
  - [products](#products)

## Database Summary

**Total Tables:** 5
**Total Rows:** 41

## Database Relationships

```mermaid
erDiagram
    users ||--o{ orders : ""
    orders ||--o{ order_items : ""
    categories ||--o{ categories : ""
    categories ||--o{ products : ""
```

## Tables

## users

<a id="users"></a>

Row Count: 10

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY |  |
| email | varchar(255) | NO | UNIQUE |  |
| first_name | varchar(100) | NO |  |  |
| last_name | varchar(100) | NO |  |  |
| created_at | timestamp without time zone | YES |  | CURRENT_TIMESTAMP |
| status | varchar(20) | YES |  | 'active'::character varying |
| is_verified | boolean | YES |  | false |

---

## orders

<a id="orders"></a>

Row Count: 10

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY |  |
| user_id | integer | NO |  |  |
| total | numeric(10,2) | NO |  | 0.00 |
| order_date | timestamp without time zone | YES |  | CURRENT_TIMESTAMP |
| status | varchar(20) | YES |  | 'pending'::character varying |
| shipping_address | text | YES |  |  |
| notes | text | YES |  |  |

---

## order_items

<a id="order-items"></a>

Row Count: 13

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY |  |
| order_id | integer | NO |  |  |
| product_name | varchar(255) | NO |  |  |
| quantity | integer | NO |  | 1 |
| unit_price | numeric(10,2) | NO |  |  |
| total_price | numeric(10,2) | YES |  |  |

---

## categories

<a id="categories"></a>

Row Count: 6

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY |  |
| name | varchar(100) | NO | UNIQUE |  |
| description | text | YES |  |  |
| parent_id | integer | YES |  |  |

---

## products

<a id="products"></a>

Row Count: 8

### Columns

| Column | Type | Nullable | Constraints | Default |
|--------|------|----------|-------------|---------|
| id | integer | NO | PRIMARY KEY |  |
| name | varchar(255) | NO |  |  |
| description | text | YES |  |  |
| price | numeric(10,2) | NO |  |  |
| category_id | integer | YES |  |  |
| in_stock | boolean | YES |  | true |
| created_at | timestamp without time zone | YES |  | CURRENT_TIMESTAMP |
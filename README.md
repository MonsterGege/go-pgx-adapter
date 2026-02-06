# Go PGX Adapter

A lightweight PostgreSQL database adapter for Go applications using the pgx driver with connection pooling and transaction support.

## Features

- Connection pooling with configurable parameters
- Thread-safe database initialization
- Transaction management
- Context-aware operations
- SSL support
- Automatic connection lifecycle management

## Installation

```bash
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
```

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "log"
    
    "your-module/database"
)

func main() {
    ctx := context.Background()
    
    db, err := database.NewDatabase(
        ctx,
        "localhost",    // host
        "username",     // user
        "password",     // password
        "dbname",       // database name
        "5432",         // port
        "disable",      // ssl mode
    )
    if err != nil {
        log.Fatal(err)
    }
    defer db.CloseConnection(ctx)
    
    // Use the database connection
    pool := db.GetConnection(ctx)
    // ... your database operations
}
```

### Using Transactions

```go
func performTransaction(ctx context.Context, db *database.Database) error {
    tx, err := db.UsingTransactions(ctx)
    if err != nil {
        return err
    }
    
    // Execute queries within transaction
    err = db.ExecuteTransaction(ctx, tx, "INSERT INTO users (name) VALUES ($1)", "John Doe")
    if err != nil {
        db.Rollback(ctx, tx)
        return err
    }
    
    // Query within transaction
    rows, err := db.QueryTransaction(ctx, tx, "SELECT id FROM users WHERE name = $1", "John Doe")
    if err != nil {
        db.Rollback(ctx, tx)
        return err
    }
    defer rows.Close()
    
    // Commit transaction
    return db.Commit(ctx, tx)
}
```

## Configuration

The adapter uses the following connection pool settings by default:

- `pool_max_conns`: 30
- `pool_min_conns`: 5
- `pool_max_conn_lifetime`: 1 hour
- `pool_max_conn_idle_time`: 30 minutes

## API Reference

### Database

#### `NewDatabase(ctx, host, user, password, dbname, port, sslmode) (*Database, error)`
Creates a new database instance with connection pool.

#### `GetConnection(ctx) *pgxpool.Pool`
Returns the connection pool for direct database operations.

#### `CloseConnection(ctx)`
Closes the database connection pool.

### Transactions

#### `UsingTransactions(ctx) (pgx.Tx, error)`
Begins a new database transaction.

#### `QueryTransaction(ctx, tx, query, args...) (pgx.Rows, error)`
Executes a query within a transaction.

#### `ExecuteTransaction(ctx, tx, query, args...) error`
Executes a statement within a transaction.

#### `Commit(ctx, tx) error`
Commits the transaction.

#### `Rollback(ctx, tx) error`
Rolls back the transaction.

## SSL Modes

Supported SSL modes:
- `disable` - No SSL
- `require` - Require SSL
- `verify-ca` - Verify certificate authority
- `verify-full` - Full certificate verification

## License

This project is open source. Please check the license file for details.

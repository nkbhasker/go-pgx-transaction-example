## Overview

This repository contains an example of how to use PGX, a PostgreSQL driver and toolkit for Golang, to manage database transactions. The example demonstrates how to start a transaction, execute queries, and commit or rollback the transaction based on the success or failure of the operations.

## Prerequisites

- Go 1.16 or later
- PostgreSQL 9.5 or later
- PGX library for Go

## Installation

1. Clone the repository:
  ```sh
  git clone https://github.com/yourusername/go-pgx-transaction.git
  cd go-pgx-transaction
  ```

2. Install the dependencies:
  ```sh
  go mod tidy
  ```

## Usage

1. Update the database connection settings in `main.go`:
  ```go
  connConfig, err := pgxpool.ParseConfig("postgresql://username:password@localhost:5432/database_name")
  ```

2. Run the example:
  ```sh
  go run main.go
  ```

## Example

The example code in `main.go` shows how to:

- Establish a connection to the PostgreSQL database using PGX.
- Start a transaction.
- Execute multiple queries within the transaction.
- Commit the transaction if all queries succeed.
- Rollback the transaction if any query fails.


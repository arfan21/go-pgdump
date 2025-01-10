# Go PostgreSQL Dump

A utility to dump PostgreSQL database schema (tables, functions, and procedures) to individual .sql files.  

## Requirements

- Go (1.18+)
- PostgreSQL
- pg_dump CLI installed

## Installation

1. Clone this repository.
2. Ensure you have Go installed.
3. Run `go mod download` to retrieve dependencies.

## Configuration

Provide a .env file in the project root directory:

```properties
DATABASE_URL="postgres://username:password@hostname:5432/dbname?sslmode=disable"
DUMP_DIR=./db_dump
```

• DATABASE_URL: Connection string to your PostgreSQL database.  
• DUMP_DIR: Output directory for dump files.

## Usage

1. Build and run the project:

   ```bash
   go build -o github.com/arfan21/go-pgdump
   ./github.com/arfan21/go-pgdump
   ```

2. Check the DUMP_DIR for dumped .sql files.

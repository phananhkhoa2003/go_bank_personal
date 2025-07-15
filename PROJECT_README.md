# Simple Bank Project

A simple banking application built with Go, featuring RESTful APIs for account management, transfers, and user authentication.

## Environment Setup

### 1. Configuration
Copy the environment template and configure your settings:
```bash
cp app.env.example app.env
```

Edit `app.env` with your actual database credentials and settings:
```env
DB_SOURCE=postgresql://username:password@hostname:port/database_name
DB_DRIVER=postgres
SERVER_ADDRESS=0.0.0.0:8080
ACCESS_TOKEN_DURATION=15m
TOKEN_SYMMETRIC_KEY=your_32_character_secret_key_here
```

### 2. Database Setup
```bash
# Start PostgreSQL container
make postgres

# Create database
make createdb

# Run migrations
make migrateup
```

### 3. Running Tests
```bash
# Run all tests
make test

# Run benchmark tests
make benchmark

# Run and save benchmark results
make benchmark-save
```

### 4. Available Make Commands
```bash
make benchmark-help  # Show all benchmark commands
make help            # Show all available commands
```

## Project Structure
- `/api` - HTTP handlers and routing
- `/db` - Database queries, migrations, and tests
- `/token` - JWT and PASETO token management
- `/util` - Utility functions and configuration
- `/benchmark_results` - Benchmark test results (excluded from git)

## Security Notes
- `app.env` contains sensitive data and is excluded from version control
- Use `app.env.example` as a template for new environments
- Benchmark results and temporary files are automatically ignored by git

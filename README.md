# Sustainability Tracker Activity Service

This service manages CRUD operations for sustainability activities and implements Redis caching. It is built using Golang with Gin, PostgreSQL, and Redis.

## Setup

1. **Environment Variables**

   Set the following environment variables:
   - `DATABASE_URL` (e.g., `postgres://user:password@localhost:5432/sustainability?sslmode=disable`)
   - `REDIS_ADDR` (e.g., `redis:6379`)
   - `REDIS_PASSWORD` (if any)
   - `PORT` (optional, default is 8081)

2. **Database Migration**

   To run the database migrations (which create the `activities` table), use:

   ```sh
   make migrate
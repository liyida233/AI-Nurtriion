# AI-Enhanced Fitness and Nutrition Tracking System

Full-stack FYP project for an AI-assisted fitness and nutrition decision-support platform.

## Stack

- Frontend: React, TypeScript, Vite
- Backend: Go, Gin, GORM
- Database: MySQL
- Cache: Redis
- AI: mock provider first, designed for later LLM/API integration

## Project Structure

```text
backend/    Go REST API, domain models, analytics, mock AI recommendation
frontend/   React dashboard and module forms
docs/       Proposal and UML diagrams
tools/      Local documentation tooling
```

## Start Services

```powershell
docker compose up -d mysql redis
```

## Run Backend

```powershell
cd backend
go mod tidy
go run ./cmd/api
```

Backend API: `http://localhost:8080/api`

The Docker MySQL service is exposed on host port `13306` to avoid conflicts with a locally installed MySQL server and Windows reserved port ranges.

## Run Frontend

PowerShell may block `npm.ps1`, so use `npm.cmd`:

```powershell
cd frontend
npm.cmd install
npm.cmd run dev
```

Frontend: `http://localhost:5173`

## Current Implementation

- User registration and login with JWT
- Auth session support with Redis-backed logout and `/auth/me`
- Profile creation/update
- Exercise reference data CRUD
- Food reference data CRUD
- Workout session CRUD with training volume indicators
- Meal log CRUD with macro calculation
- Body record CRUD
- Goal CRUD with feasibility checks, generated milestones, status updates, and milestone completion
- Dashboard weekly/monthly summary with Redis cache support
- Persistable analytics snapshots
- Analytics indicators for progressive overload, muscle-group balance, macro ratios, nutrition gaps, meal quality, calorie status, and 7-day weight moving average
- Mock AI recommendation flow with structured prompt context, safety validation, and feedback CRUD
- Optional OpenAI-compatible AI provider skeleton
- Weekly/monthly progress report generation, detail, listing, and deletion
- HTML report export and download endpoint
- Admin role guard for reference data and user role management
- AI recommendation rate limiting with Redis

## AI Provider

The default provider is `mock`, which keeps local development stable and free. To enable an OpenAI-compatible API later, set:

```text
AI_PROVIDER=openai_compatible
AI_API_KEY=your-api-key
AI_BASE_URL=https://api.openai.com/v1
AI_MODEL=gpt-4o-mini
AI_RATE_LIMIT_PER_HOUR=10
```

## API Coverage

```text
POST   /api/auth/register
POST   /api/auth/login
GET    /api/auth/me
POST   /api/auth/logout

GET    /api/profile
PUT    /api/profile

GET    /api/workouts/exercises
POST   /api/workouts/exercises
GET    /api/workouts/exercises/:id
PUT    /api/workouts/exercises/:id
DELETE /api/workouts/exercises/:id
GET    /api/workouts
POST   /api/workouts
GET    /api/workouts/:id
PUT    /api/workouts/:id
DELETE /api/workouts/:id

GET    /api/nutrition/foods
POST   /api/nutrition/foods
GET    /api/nutrition/foods/:id
PUT    /api/nutrition/foods/:id
DELETE /api/nutrition/foods/:id
GET    /api/nutrition/meals
POST   /api/nutrition/meals
GET    /api/nutrition/meals/:id
PUT    /api/nutrition/meals/:id
DELETE /api/nutrition/meals/:id

GET    /api/body-records
POST   /api/body-records
GET    /api/body-records/:id
PUT    /api/body-records/:id
DELETE /api/body-records/:id

GET    /api/goals
POST   /api/goals
GET    /api/goals/:id
PUT    /api/goals/:id
PATCH  /api/goals/:id/status
PATCH  /api/goals/:id/milestones/:milestoneId
DELETE /api/goals/:id

GET    /api/dashboard/summary?period=weekly|monthly&persist=true|false
GET    /api/dashboard/snapshots

GET    /api/recommendations
POST   /api/recommendations/generate
GET    /api/recommendations/:id
DELETE /api/recommendations/:id
GET    /api/recommendations/:id/feedback
POST   /api/recommendations/:id/feedback

GET    /api/reports
POST   /api/reports/generate
GET    /api/reports/:id
GET    /api/reports/:id/download
DELETE /api/reports/:id

GET    /api/admin/users
GET    /api/admin/users/:id
PATCH  /api/admin/users/:id/role
POST   /api/admin/workouts/exercises
PUT    /api/admin/workouts/exercises/:id
DELETE /api/admin/workouts/exercises/:id
POST   /api/admin/nutrition/foods
PUT    /api/admin/nutrition/foods/:id
DELETE /api/admin/nutrition/foods/:id
```

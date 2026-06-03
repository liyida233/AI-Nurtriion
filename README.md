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
- Profile creation/update
- Exercise and food seed data
- Workout logging with training volume indicator
- Nutrition logging with macro calculation
- Body record logging
- Goal creation with feasibility checks and generated milestones
- Dashboard summary with Redis cache support
- Mock AI recommendation flow with structured prompt context and feedback endpoint

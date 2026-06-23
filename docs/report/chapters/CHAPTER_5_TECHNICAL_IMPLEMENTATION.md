# CHAPTER 5: TECHNICAL IMPLEMENTATION

## 5.1 Introduction

This chapter describes the technical implementation of the AI-Enhanced Fitness and Nutrition Tracking System. While Chapter 4 explained the requirements and design, this chapter presents how the design was converted into an operational full-stack application.

The implemented system consists of a React and TypeScript frontend, a Go and Gin REST API, a MySQL relational database, Redis supporting infrastructure, and local HTML report storage. The source code is maintained in one repository with separate frontend and backend directories. At the time of writing, the backend contains 55 Go source files, 54 route declarations, and 9 automated tests. The frontend is implemented as a single-page application using shared API and interface utilities.

Implementation evidence in this chapter uses an anonymised local account named `FYP Demo User`. The demonstration data is synthetic and does not contain real personal health information.

## 5.2 Project Source-Code Structure

The project uses a monorepository so that frontend, backend, documentation, UML source files, and development tools can be revised together.

```text
AI Nurtirion/
|-- backend/
|   |-- cmd/api/
|   |-- internal/
|       |-- cache/
|       |-- config/
|       |-- database/
|       |-- httpctx/
|       |-- middleware/
|       |-- models/
|       |-- modules/
|       |-- server/
|-- frontend/
|   |-- src/
|-- docs/
|   |-- report/
|   |-- UML/
|-- tools/
|-- docker-compose.yml
|-- README.md
```

The backend follows feature-based organisation. Each main domain is represented by a module under `backend/internal/modules`. This allows the code associated with one business capability to remain together.

**Table 5.1: Backend File Responsibilities**

| File Type | Responsibility |
|---|---|
| `controller.go` | Registers routes, reads HTTP input, obtains authenticated identity, and returns HTTP responses |
| `dto.go` | Defines request structures and binding validation |
| `service.go` | Implements business use cases and coordinates repositories, analytics, Redis, and providers |
| `repository.go` | Performs GORM database queries and ownership filtering |
| `analyzer.go` | Implements calculations and classifications that can be tested independently |
| `provider.go` | Defines recommendation-provider behaviour and external-provider integration |
| `generator.go` | Generates deterministic report narrative and HTML |

The frontend currently uses four TypeScript or TSX source files. `main.tsx` creates the React root, `App.tsx` contains the application shell and module panels, `api.ts` centralises HTTP communication, and `vite-env.d.ts` provides Vite environment typing. Styling is maintained in `styles.css`.

Although `App.tsx` is functional, it has grown substantially as all module panels were developed. A future maintainability improvement is to divide it into route-level pages, reusable forms, tables, domain types, and smaller components.

## 5.3 Application Start-Up and Infrastructure

### 5.3.1 Local Infrastructure

MySQL and Redis are defined in `docker-compose.yml`. The MySQL container uses MySQL 8.4 and exposes container port 3306 through host port 13306. The alternative host port avoids collision with an existing local MySQL installation and Windows port restrictions. Redis 7.4 is exposed through port 6379.

Persistent MySQL data is stored in a named Docker volume. Restarting the container therefore does not remove existing application records.

The backend can be started after the infrastructure becomes healthy:

```powershell
docker compose up -d mysql redis
cd backend
go run ./cmd/api
```

The frontend is started separately:

```powershell
cd frontend
npm.cmd install
npm.cmd run dev
```

The backend port is controlled by the `APP_PORT` environment variable and defaults to 8080. The frontend API location is controlled by `VITE_API_BASE_URL`. During the latest report-evidence session, port 8180 was used because the Windows environment had reserved port 8080.

### 5.3.2 Backend Start-Up Sequence

The backend entry point performs the following operations:

1. Load environment-based configuration.
2. Connect to MySQL.
3. run GORM AutoMigrate for the current entities.
4. Insert seed exercise and food references when they are absent.
5. Connect to Redis.
6. Construct the Gin router and register routes.
7. Start the HTTP server using the configured application port.

MySQL is a required dependency because the application data is relational and persistent. Redis is treated as a supporting dependency. If Redis cannot be reached, the connector returns `nil` and the backend continues without caching, logout-aware session storage, or recommendation rate limiting. This fallback is useful for development, although a production environment should treat the loss of security-related Redis functionality more strictly.

## 5.4 Frontend Implementation

### 5.4.1 React Application Shell

The frontend is implemented as a single React application. The root component stores:

- The JWT access token.
- The authenticated user.
- The active navigation module.
- The current dashboard summary.
- The selected analytical date range.
- A shared user-facing message.

The token and basic user object are persisted in browser local storage. When the application starts with a stored token, it requests `/auth/me` to refresh the authenticated identity and role.

The navigation list is generated from a configuration array. Each item contains a key, label, and Lucide icon. The administrator item is filtered out unless the current user has the `admin` role. Selecting a navigation item switches the displayed module panel without a full-page reload.

The application shell provides a fixed navigation area, current user and role, current module heading, refresh action, and shared notification area. This ensures that every module uses the same overall interaction pattern.

### 5.4.2 Central API Client

HTTP communication is centralised in `frontend/src/api.ts`. The API base URL is read from `VITE_API_BASE_URL` and falls back to `http://localhost:8080/api`.

The generic `api<T>` function performs four tasks:

1. Combines the base URL with the requested path.
2. Adds JSON headers and the bearer token when supplied.
3. Parses the standard backend response.
4. Returns either typed data or a readable error.

```typescript
const response = await fetch(`${API_BASE}${path}`, {
  ...options,
  headers: {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options.headers
  }
});
```

A separate download function returns a `Blob` instead of attempting to parse JSON. The Reports panel creates a temporary browser URL and starts the HTML file download.

### 5.4.3 Authentication Interface

**Figure 5.1: Login Interface**

*[Insert `docs/report/screenshots/05_login.png` here.]*

The authentication screen supports login and registration using one panel. Registration adds the name field, while both modes require email and password. The password input enforces a minimum length of eight characters on registration.

After successful authentication, the token and user object are stored and the application shell is displayed. Authentication errors such as an invalid password or duplicate email are shown inside the form rather than only appearing in the browser console.

### 5.4.4 Dashboard Interface

**Figure 5.2: Implemented Dashboard**

*[Insert `docs/report/screenshots/05_dashboard.png` here. Use a landscape page or full-width figure.]*

The Dashboard panel requests a `DashboardSummary` from the backend. It supports daily, weekly, monthly, and custom ranges. Custom mode displays start and end date controls and adds them to the query string.

The interface displays:

- Workout sessions, volume, consistency, and progressive-load status.
- Calories, meal quality, and calorie classification.
- Weight trend and seven-day average.
- Goal adherence.
- Protein, carbohydrate, and fat in grams.
- Protein, carbohydrate, and fat as percentages.
- Muscle-group distribution.
- Nutrition gaps and training warnings.

Gram values were added alongside percentages because a percentage alone does not show the total quantity consumed. The selected range is shown under the period control so that users can interpret the values in the correct time context.

### 5.4.5 Daily Workout Journal

**Figure 5.3: Workout Journal and Exercise Search**

*[Insert `docs/report/screenshots/05_workout_journal.png` here. Use a landscape page.]*

The Workout panel combines a selected date, daily totals, a searchable exercise selector, the workout form, workout history, and exercise reference table.

The exercise search compares the keyword against the exercise name, category, muscle group, and equipment. The filtered list is used by both the selector and reference table.

Daily totals are calculated on the frontend from the sessions that match the selected date. The totals include session count, duration, training load, and number of exercise entries. This immediate calculation improves the daily logging experience, while the backend remains responsible for date-range and cross-session analytics.

Each session can be edited or deleted. Editing loads the existing values into the form. The backend replaces the session entries in a database transaction to prevent a partially updated workout from being stored.

### 5.4.6 Daily Nutrition Journal

**Figure 5.4: Nutrition Journal and Food Search**

*[Insert `docs/report/screenshots/05_nutrition_journal.png` here. Use a landscape page.]*

The Nutrition panel follows the same daily-journal pattern as Workout. A selected date controls the displayed meal history and daily totals. The food search filters by food name and serving size.

When the user selects a food and quantity, the backend retrieves the food reference and calculates total calories, protein, carbohydrate, and fat. The calculated values are stored with the meal log.

The daily interface sums the stored meal totals to show the selected day's calories and macronutrients. This avoids mixing a daily logging workflow with a long historical list and allows old records to be edited by changing the date.

### 5.4.7 Body, Goal, Recommendation, and Report Interfaces

**Figure 5.5: Body-Progress Timeline**

*[Insert `docs/report/screenshots/05_body_progress.png` here.]*

The Body panel supports dated weight records with optional notes. Records are shown as a timeline with edit and delete actions. The dashboard uses these records for latest weight, trend, and calendar-based seven-day average calculations.

**Figure 5.6: Goal and Custom Milestone Interface**

*[Insert `docs/report/screenshots/05_goals_milestones.png` here. Use a landscape page.]*

The Goal panel allows the user to enter a goal type, target metric, target value, deadline, and priority. Three milestone groups allow custom titles, target values, and due dates. Existing goals show status controls and clickable milestone completion actions.

**Figure 5.7: Recommendation Interface**

*[Insert `docs/report/screenshots/05_recommendations.png` here. Use a landscape page.]*

The Recommendation panel provides weekly, workout, and meal recommendation actions. Generated content is displayed with its type and creation time. Users can mark a recommendation as useful or not useful. The feedback operation stores the rating through a separate endpoint.

**Figure 5.8: Report Generation and History**

*[Insert `docs/report/screenshots/05_reports.png` here. Use a landscape page.]*

The Reports panel uses the same four period types as the dashboard. Generated reports are listed with their deterministic narrative. Each report can be downloaded as HTML or deleted.

### 5.4.8 Administration Interface

**Figure 5.9: Administration Interface**

*[Insert `docs/report/screenshots/05_admin.png` here. Use only the anonymised demonstration screenshot.]*

The Admin panel is available only when the current authenticated user has the administrator role. It contains:

- A user table with role-change actions.
- A form for adding or editing food references.
- A form for adding or editing exercise references.
- Food and exercise library tables with edit and delete actions.

The frontend role check controls navigation visibility, but it is not treated as the security boundary. The backend administrator middleware independently validates the role for every protected administration request.

## 5.5 Backend API Implementation

### 5.5.1 Gin Routing and Middleware

The Gin server registers all application routes below `/api`. Public routes include health, registration, and login. The remaining application routes are placed in a protected group using authentication middleware.

Administrator routes are nested inside the protected group and apply a second role middleware. This means an administrator request must first pass JWT and session validation and then pass the database-backed role check.

The server also applies CORS configuration for the known local frontend origins. Allowed methods include GET, POST, PUT, PATCH, DELETE, and OPTIONS. The `Authorization` and `Content-Type` headers are explicitly permitted.

### 5.5.2 Controller, Service, and Repository Flow

Most modules follow the same execution flow:

```text
HTTP request
  -> controller and DTO binding
  -> authenticated user extraction
  -> service business rule
  -> repository database operation
  -> cache invalidation when required
  -> standard JSON response
```

Controllers are intentionally thin. They select the HTTP status, bind input, obtain path and query values, and call the service. Services own use-case logic. Repositories use GORM and accept `context.Context` so that database work follows the request lifecycle.

The standard success response is:

```json
{
  "data": {}
}
```

The standard failure response is:

```json
{
  "error": "readable error message"
}
```

This response format allows the frontend API client to handle all modules consistently.

### 5.5.3 Ownership-Filtered Persistence

Private resource retrieval includes both the resource identifier and authenticated user identifier:

```go
Where("id = ? AND user_id = ?", id, userID)
```

This pattern is used for workouts, meals, body records, goals, recommendations, and reports. A valid token is therefore insufficient to access another user's private resource.

Nested operations verify the parent first. For example, milestone completion verifies that the goal belongs to the user before retrieving and updating the milestone.

Transactions are used when several related database operations must succeed together. Updating a workout saves the session, removes existing entries, and inserts replacement entries inside one GORM transaction. Deleting a workout removes entries and the owned session in one transaction.

## 5.6 Authentication and Authorisation Implementation

### 5.6.1 Password and Token Handling

Registration validates a name, email, and password. Passwords are hashed using bcrypt's default cost. Only the resulting hash is stored in MySQL, and the model's JSON tag prevents it from appearing in API responses.

Login retrieves the user by email and compares the submitted password with the hash. A successful operation creates a JWT containing:

- `sub`: the user UUID.
- `iat`: token issue time.
- `exp`: expiry time.

The token is signed using the configured JWT secret and currently expires after 24 hours.

### 5.6.2 Redis Session Validation

When a token is issued, Redis stores:

```text
session:{userId} = active
```

The Redis key lifetime matches the token expiry. Protected middleware validates the token signature and subject and then checks the active-session key. Logout removes this key.

This design makes logout effective before JWT expiry. A purely stateless JWT cannot normally be revoked without additional state.

### 5.6.3 Role-Based Access Control

Administrator middleware reads the current user role from MySQL. A non-administrator receives `403 Forbidden`. The role is checked against current database state instead of relying only on the role value that may have been stored in frontend local storage.

The implementation currently supports `user` and `admin`. Role updates are validated by the admin service before persistence.

## 5.7 Database Implementation

### 5.7.1 GORM Models and Migration

The relational model is implemented as Go structures with GORM tags. UUID strings use `CHAR(36)` primary keys. User ownership and foreign-reference fields are indexed.

The startup migration includes:

- User and UserProfile.
- Exercise, WorkoutSession, and WorkoutEntry.
- FoodItem and MealLog.
- BodyRecord.
- Goal and GoalMilestone.
- AnalyticsSnapshot.
- AIRecommendation and RecommendationFeedback.
- ProgressReport.

The user email field is unique. The profile user identifier is also unique so that one user cannot have multiple profile rows.

### 5.7.2 Seed Data

The database seeding function checks for each item by name before insertion. It currently provides:

- Squat.
- Bench Press.
- Lat Pulldown.
- Running.
- Chicken Breast.
- Cooked Rice.
- Egg.
- Banana.

The seed is idempotent at the item-name level, so restarting the backend does not insert the same records repeatedly. These references allow the application to demonstrate search, logging, analysis, and administration immediately after initial setup.

### 5.7.3 UUID Selection

UUID values are generated in the service or authentication layer before records are stored. The same identifier format is used across modules.

This is useful for public APIs because identifiers are not simple sequential values. However, UUIDs are not considered an access-control mechanism. Repository ownership filtering remains necessary.

## 5.8 Analytics Implementation

### 5.8.1 Date-Range Resolution

The analytics controller accepts `period`, `startDate`, `endDate`, and optional `persist`. `ResolveRange` converts these values into start and end timestamps.

- Daily uses the current local calendar day.
- Weekly uses the current day and previous six days.
- Monthly uses a rolling one-month interval ending on the current date.
- Custom requires two valid dates and rejects an end date before the start date.

The same range resolver is reused by report generation. This ensures that a weekly dashboard and weekly report interpret the range consistently.

### 5.8.2 Workout Indicators

For each workout entry, the analytics service calculates:

```text
sets x repetitions x weight
```

Progressive-overload analysis uses a fallback load of one when weight is absent. This means bodyweight or unloaded exercise sessions can still be compared through sets and repetitions.

Sessions are sorted and divided into earlier and recent groups. A difference greater than 5% is improving, less than -5% is declining, and the remaining range is stable.

Workout consistency scales a target of four sessions per seven days to the selected period and limits the result to 100%.

Muscle-group distribution uses the linked exercise reference. Rules identify missing lower-body or back training and a concentration greater than 60% in one group.

### 5.8.3 Nutrition Indicators

Meal creation uses a small independent function:

```go
return MealIndicators{
    Calories: food.Calories * quantity,
    Protein:  food.Protein * quantity,
    Carbs:    food.Carbohydrates * quantity,
    Fat:      food.Fat * quantity,
}
```

The analytics service sums the stored meal values for the selected range. Macronutrient percentages convert grams to energy using the 4-4-9 factors.

Meal quality begins at 100 and applies documented deductions. When protein expectation is calculated, the implementation uses the number of days that contain meal logs instead of assuming that every selected day has a complete dietary record.

### 5.8.4 Body and Calorie Indicators

BMR is calculated from profile age, gender, height, and weight. Activity factors convert BMR to estimated TDEE. Total range expenditure is estimated by multiplying TDEE by the selected number of days.

Body records are sorted by date. The seven-day moving average establishes a start date six days before the latest record and includes only records within that calendar window.

Weight-trend classification compares the earliest and latest available records. A difference smaller than 0.3 kg is stable.

### 5.8.5 Goal Indicators

The goal analyzer prevents:

- Deadlines that allow fewer than seven days.
- Workout-frequency targets above seven sessions per week.
- Weight-loss targets above the configured weekly threshold.

The current implementation returns `422 Unprocessable Entity` for these business-validation errors.

Goal adherence divides completed milestones by total milestones. When no milestones exist, adherence remains zero instead of attempting to divide by zero.

## 5.9 Redis Implementation

Redis is used for three limited responsibilities.

### 5.9.1 Weekly Dashboard Cache

The weekly dashboard summary is serialised as JSON and stored using:

```text
dashboard:{userId}
```

The expiration time is five minutes. The analytics service only reads this cache for the weekly period and only when the request does not ask to persist a snapshot.

Profile, workout, nutrition, body, and goal changes call a shared invalidation function that deletes the user's dashboard key. This prevents the dashboard from displaying a known stale weekly summary after a relevant write.

### 5.9.2 Active Sessions

Redis stores active session keys as described in Section 5.6.2. The same infrastructure supports logout-aware authentication without adding a separate session table to MySQL.

### 5.9.3 Recommendation Rate Limiting

Recommendation generation increments an hourly counter using the user UUID and hour number:

```text
rate:ai:{userId}:{hour}
```

The first increment assigns a one-hour expiry. Requests above the configured limit are rejected. If Redis is unavailable, recommendation generation continues without the rate limit in the current development configuration.

## 5.10 Recommendation Implementation

### 5.10.1 Provider Interface

The recommendation module defines:

```go
type Provider interface {
    Generate(
        ctx context.Context,
        kind string,
        summary analytics.DashboardSummary,
    ) (string, error)
}
```

The service is independent of the concrete provider. `NewProvider` selects the mock provider unless the environment configures `openai` or `openai_compatible`.

### 5.10.2 Mock Provider

The mock provider produces deterministic output from weekly indicators. Example rules include:

- Suggest increasing protein when the weekly protein total is low.
- Identify a high calorie balance.
- Suggest scheduling three short sessions when workout frequency is low.
- Suggest progressive increases when frequency is sufficient.
- Compare upward weight trend with positive calorie balance.
- Encourage creating a SMART goal when no active goal exists.

This output is not generated by a production LLM. The interface and workflow are complete, but the content remains predictable and free during the current phase.

### 5.10.3 OpenAI-Compatible Provider

The optional provider builds a JSON prompt from the analytical summary and sends a chat-completions request. The provider supports configurable API key, base URL, and model.

The system prompt restricts output to concise general-wellness guidance and instructs the model to avoid medical diagnosis, extreme dieting, and unsafe training advice. The HTTP client applies a 25-second timeout.

### 5.10.4 Safety and Feedback

After provider generation, a safety validator converts the content to lowercase and checks for blocked terms including diagnosis, cure, starvation, and extreme fasting. Rejected output is not saved.

Accepted output is stored with the structured prompt context. Storing the context supports later examination of which indicators produced a recommendation. Feedback is stored as a separate entity, allowing more than one feedback record if required by future evaluation.

The current validator is intentionally simple. A production implementation would require stronger policy rules, structured output validation, prompt-injection protection, model evaluation, and human review for high-risk cases.

## 5.11 Progress-Report Implementation

Report generation first resolves the requested period and builds a dashboard summary without Redis caching. A deterministic narrative describes workout, nutrition, calorie, body, and goal values.

The HTML generator escapes the narrative before inserting it into the document. This prevents summary text from being interpreted as arbitrary HTML. The document contains styled metric cards and a readable summary.

Generated files are stored in:

```text
backend/storage/reports/{reportId}.html
```

The database stores the period, generation time, summary, and protected download URL. The download endpoint retrieves the owned report before returning the file as an attachment.

Deleting a report removes its database metadata and attempts to remove the local HTML file. Production deployment should replace local storage with protected durable object storage.

## 5.12 Testing and Verification

### 5.12.1 Automated Calculation Tests

The current backend automated test suite contains nine test functions. The analytics tests cover:

- BMR calculation.
- Macronutrient ratio calculation.
- Progressive-load classification.
- Progressive-load fallback without external weight.
- Muscle-group warnings.
- Meal-quality scoring.
- Sparse meal-log date-range behaviour.
- Calendar-based seven-day body-weight average.

An additional goal-controller test verifies the classification of feasibility and status errors as validation errors.

### 5.12.2 Build Verification

The following backend command was executed:

```powershell
go test ./...
```

All packages compiled, and all available tests passed.

The frontend was verified using:

```powershell
npm.cmd run build
```

TypeScript compilation and Vite production bundling completed successfully. The latest build transformed 1,580 modules and produced the HTML, CSS, and JavaScript assets without a build error.

### 5.12.3 Integrated Workflow Verification

An anonymised demonstration account was used to verify:

- Registration and login.
- Profile update.
- Exercise and food retrieval.
- Workout creation and daily display.
- Meal creation and daily nutrition totals.
- Body-record creation and timeline.
- Goal creation with custom milestones.
- Milestone completion and adherence update.
- Weekly dashboard aggregation.
- Recommendation generation.
- Weekly report generation.
- Administrator navigation and protected data management.

The browser frontend communicated with the backend through the configured API base URL. MySQL and Redis containers were healthy during the workflow.

**Table 5.2: Latest Verification Results**

| Verification | Result | Evidence |
|---|---|---|
| MySQL container | Passed | Container healthy on local port 13306 |
| Redis container | Passed | Container listening on local port 6379 |
| Backend health endpoint | Passed | `/api/health` returned service status `ok` |
| Backend compilation and tests | Passed | `go test ./...` |
| Frontend production build | Passed | TypeScript and Vite build completed |
| Authentication workflow | Passed | Demo user login and protected navigation |
| Workout and nutrition workflow | Passed | Daily records and totals displayed |
| Dashboard integration | Passed | Weekly indicators generated from stored records |
| Goal and milestone workflow | Passed | Goal displayed with completed milestone |
| Recommendation workflow | Passed | Mock recommendation stored and displayed |
| Report workflow | Passed | Weekly HTML report created and listed |
| Administrator visibility | Passed | Admin module shown for administrator role |

### 5.12.4 Remaining Test Gaps

The current test suite is strongest around analytical calculations. It does not yet provide complete automated coverage for:

- Controller and HTTP integration.
- Repository queries against a test database.
- React component behaviour.
- Browser end-to-end workflows.
- Concurrent updates and transactions.
- Security penetration testing.
- Load and response-time testing.
- External AI-provider quality and failure handling.

These gaps are reported explicitly because a successful build is not equivalent to complete system validation.

## 5.13 Implementation Challenges and Solutions

**Table 5.3: Main Implementation Challenges**

| Challenge | Observed Problem | Implemented Solution |
|---|---|---|
| Fragmented time periods | Weekly and monthly views did not support daily editing or arbitrary ranges | Added daily and custom ranges and date-based journals |
| Sparse body records | The previous average used recent records regardless of calendar distance | Changed the moving average to a seven-calendar-day window |
| Unloaded exercises | Progressive load remained insufficient when no weight was entered | Added sets-and-repetitions fallback |
| Long-range meal quality | Missing days reduced the score even when they were not logged | Scaled protein expectation using actual meal-log days |
| Fixed milestones | Generated milestones could not represent the user's own checkpoints | Added custom milestone input with default generation fallback |
| Stale analytics | Writes could leave a cached weekly dashboard unchanged | Added shared cache invalidation to relevant services |
| JWT logout | A token could remain valid until expiry | Added Redis active-session validation and logout deletion |
| External AI dependency | API cost and availability could block the demonstration | Added a mock provider behind a replaceable interface |
| Report portability | PDF generation added unnecessary rendering complexity | Implemented downloadable HTML for the current phase |
| Local Windows port conflict | Port 8080 was reserved during the report session | Used configurable `APP_PORT=8180` and matching frontend environment |

## 5.14 Current Implementation Limitations

The implemented system is a complete local prototype but has several limitations.

First, the frontend is concentrated in one large `App.tsx` file. This should be divided before a larger team or feature set is introduced.

Second, the nutrition and exercise reference datasets are small. A larger source would require licensing, data-cleaning, search optimisation, and source attribution.

Third, the recommendation system uses the mock provider by default. A real model has not yet been evaluated for relevance, consistency, cost, latency, or safety.

Fourth, the safety validator uses blocked phrases and cannot identify every unsafe implication. It is a development safeguard rather than a clinical safety system.

Fifth, Redis fallback reduces security and performance features when Redis is unavailable. Production configuration should make the required behaviour explicit.

Sixth, reports are HTML files stored on the backend filesystem. They are not currently PDF documents and are not stored in durable cloud storage.

Seventh, there is no wearable-device integration, automatic food import, notification system, native mobile application, or production deployment.

Finally, formal independent usability, performance, and security evaluations remain incomplete. The current evidence demonstrates implementation and internal workflow correctness, not clinical effectiveness or large-scale readiness.

## 5.15 Chapter Summary

This chapter described the implementation of the AI-Enhanced Fitness and Nutrition Tracking System as a working React, Go, MySQL, and Redis application.

The frontend provides authentication, persistent module navigation, date-based workout and nutrition journals, dashboard analytics, body tracking, custom goals and milestones, recommendations, reports, and administration. A shared typed API client connects these interfaces to the backend.

The backend uses modular controllers, DTOs, services, repositories, and analyzers. JWT and Redis support authenticated sessions, while database-backed role middleware protects administration. GORM implements the relational entities and ownership-filtered queries. Redis supports the weekly dashboard cache and recommendation rate limiting.

Deterministic analytics calculate workout, nutrition, calorie, body, and goal indicators. The recommendation module applies a provider abstraction and safety check. The report module generates protected downloadable HTML files without depending on the AI provider.

Backend tests and the frontend production build passed during the latest verification. An anonymised demonstration account was used to verify the integrated workflows and capture the implementation evidence presented in this chapter. Remaining limitations identify the work required for production hardening and the next FYP phase.

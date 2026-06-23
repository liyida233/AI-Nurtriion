# CHAPTER 4: SYSTEM ANALYSIS AND DESIGN

## 4.1 Introduction

This chapter presents the analysis and design of the AI-Enhanced Fitness and Nutrition Tracking System. The analysis converts the problems and objectives identified in Chapter 1 into functional and non-functional requirements. The design then explains how those requirements are represented through actors, modules, data entities, interfaces, application workflows, security controls, and analytical components.

The system is designed as a modular full-stack web application. A React and TypeScript client provides task-based interfaces, while a Go REST API implements authentication, validation, business rules, persistence, analytics, recommendations, reports, and administration. MySQL stores relational application data. Redis supports active sessions, selected dashboard caching, cache invalidation, and recommendation rate limiting. The recommendation subsystem uses a provider abstraction with a deterministic mock provider as the current default and an optional OpenAI-compatible provider for future use.

The design follows four main principles:

1. User records from different health domains should be connected through a single authenticated identity.
2. Numerical indicators should be calculated by transparent application logic.
3. Generated recommendations should consume structured indicators rather than unrestricted raw data.
4. Each module should have a clear responsibility while sharing common authentication, persistence, caching, and response conventions.

## 4.2 Requirement Analysis

### 4.2.1 Stakeholders and Actors

The principal stakeholder is the general user who wants to record fitness and nutrition behaviour and review progress. The user is expected to be an adult using the system for general-wellness support rather than diagnosis or treatment.

The administrator is a secondary actor. The administrator maintains shared reference information and user roles. Administrative access is separated because normal users should not be able to modify exercise and food records that are shared across the system.

The configured recommendation provider is an external supporting actor from the use-case perspective. In the current configuration, this actor is represented by an internal deterministic mock provider. The same interface can optionally call an OpenAI-compatible external API. The provider does not directly access the database; it receives only the structured summary prepared by the backend.

The project supervisor and potential collaborator are project stakeholders rather than operational system actors. They influence requirements and evaluation, but they do not require a dedicated runtime role in the implemented application.

### 4.2.2 Functional Requirements

Functional requirements describe the operations that the system must perform. The requirements were derived from the project objectives, literature review, existing-system comparison, supervisor feedback, and iterative use of the prototype.

**Table 4.1: Functional Requirements**

| ID | Functional Requirement | Priority |
|---|---|---|
| FR-01 | The system shall allow a new user to register using a name, unique email address, and password. | Must |
| FR-02 | The system shall authenticate registered users and provide an authenticated session. | Must |
| FR-03 | The system shall allow an authenticated user to view and update age, gender, height, weight, activity level, and primary goal. | Must |
| FR-04 | The system shall allow a user to search exercise reference records. | Must |
| FR-05 | The system shall allow a user to create, view, update, and delete workout sessions and exercise entries for a selected date. | Must |
| FR-06 | The system shall calculate workout volume, consistency, progressive-overload status, muscle-group distribution, and basic warnings. | Must |
| FR-07 | The system shall allow a user to search food reference records. | Must |
| FR-08 | The system shall allow a user to create, view, update, and delete meal records for a selected date. | Must |
| FR-09 | The system shall calculate calorie totals, protein, carbohydrate, fat, macronutrient ratios, nutrition gaps, and meal-quality score. | Must |
| FR-10 | The system shall allow a user to create, view, update, and delete date-based body-weight records. | Must |
| FR-11 | The system shall calculate latest weight, weight trend, and seven-day body-weight moving average. | Must |
| FR-12 | The system shall allow a user to create, update, delete, pause, cancel, or complete a goal. | Must |
| FR-13 | The system shall validate basic goal feasibility and support custom milestones with target values and due dates. | Must |
| FR-14 | The system shall allow a user to mark milestones as completed and calculate milestone adherence. | Must |
| FR-15 | The dashboard shall support daily, weekly, monthly, and custom date ranges. | Must |
| FR-16 | The dashboard shall combine workout, nutrition, calorie, body, and goal indicators for the selected range. | Must |
| FR-17 | The system shall generate structured general-wellness recommendations from calculated weekly indicators. | Must |
| FR-18 | The system shall validate recommendation safety and allow the user to submit a rating, suitability value, and comment. | Must |
| FR-19 | The system shall generate, list, download, and delete daily, weekly, monthly, and custom-range progress reports. | Must |
| FR-20 | The system shall generate the current progress report as an HTML document. | Must |
| FR-21 | The administrator shall be able to view users and change a user's role. | Must |
| FR-22 | The administrator shall be able to create, update, and delete shared exercise references. | Must |
| FR-23 | The administrator shall be able to create, update, and delete shared food references. | Must |
| FR-24 | The system shall isolate each user's private operational records, recommendations, feedback, and reports. | Must |
| FR-25 | The system shall support logout and reject protected requests when the active session is no longer valid. | Must |

The Must priority indicates that the requirement belongs to the implemented current scope. Functions such as wearable-device import, image-based food recognition, clinical assessment, native mobile applications, PDF generation, and production cloud deployment are not included as current functional requirements.

### 4.2.3 Non-Functional Requirements

Non-functional requirements define quality expectations and operational constraints.

**Table 4.2: Non-Functional Requirements**

| ID | Category | Requirement |
|---|---|---|
| NFR-01 | Usability | Main user tasks shall be reachable through persistent authenticated navigation and use consistent form and action patterns. |
| NFR-02 | Usability | Workout and nutrition records shall be organised by a user-selected date to support daily entry and historical editing. |
| NFR-03 | Performance | Normal API operations should respond interactively in the local environment, and the weekly dashboard summary may be cached for five minutes. |
| NFR-04 | Reliability | The system shall invalidate the dashboard cache when relevant profile, workout, nutrition, body, or goal data changes. |
| NFR-05 | Security | Passwords shall be hashed and never returned through the API. |
| NFR-06 | Security | Protected endpoints shall require a valid JWT and an active Redis-backed session when Redis is available. |
| NFR-07 | Security | Administrative endpoints shall require the administrator role. |
| NFR-08 | Privacy | Database queries for private resources shall be filtered using the authenticated user's identifier. |
| NFR-09 | Maintainability | Backend business domains shall be separated into modules with controller, service, repository, DTO, and analyzer responsibilities where applicable. |
| NFR-10 | Maintainability | The recommendation provider shall be replaceable through a common provider interface. |
| NFR-11 | Portability | MySQL and Redis shall be runnable through Docker Compose in the local development environment. |
| NFR-12 | Explainability | Numerical indicators and rule-based classifications shall be implemented through documented deterministic logic. |
| NFR-13 | Safety | Recommendations shall be limited to general-wellness guidance and rejected when they contain blocked unsafe or medical claims. |
| NFR-14 | Compatibility | The web client shall communicate with the backend using JSON-based HTTP APIs and support modern desktop browsers. |
| NFR-15 | Testability | High-risk analytical calculations shall be implemented as independently testable functions. |

### 4.2.4 Assumptions and Constraints

The system design assumes that users enter records honestly and select appropriate food or exercise references. The analytical output depends on the completeness and accuracy of this data. Missing records cannot be interpreted as confirmed absence of eating or exercise.

The food and exercise datasets are curated demonstration datasets. They are sufficient for system evaluation but are not comprehensive commercial datasets. Nutrition values are estimates and may differ according to preparation method, brand, serving measurement, and data source.

The current system operates locally and is not designed as a production clinical service. It does not currently provide HTTPS termination, production secret management, automated backups, disaster recovery, or legal compliance certification. These constraints define the boundary of the current prototype.

## 4.3 Use-Case Analysis

**Figure 4.1: System Use-Case Diagram**

*[Insert `docs/UML/rendered/png/01_use_case_diagram.png` here.]*

Figure 4.1 groups the user-facing functions into account and tracking, decision support, and administration.

Account and tracking functions provide the operational data required by the rest of the system. The user registers and logs in, maintains a profile, and manages workout, nutrition, body, goal, and milestone records.

Decision-support functions transform those records into a combined dashboard, recommendation, feedback record, and downloadable progress report. The `Analyse User Data` use case is included by the dashboard, recommendation, and report functions because all three rely on the same deterministic analytics model. This reduces duplicated calculation logic and improves consistency.

Administration functions are restricted to users with the administrator role. An administrator can manage roles and shared food and exercise references. The administrator may also use ordinary user functions, but administrator-only operations are protected separately.

The configured provider participates only in generating recommendation text. It does not calculate the dashboard, create reports, manage goals, or directly retrieve user data. This separation is an important safety and architectural decision.

## 4.4 System Architecture Design

### 4.4.1 Architectural Style

The system uses a client-server architecture with a single-page web client and a modular REST API. It is deployed as a modular monolith rather than as multiple microservices.

A modular monolith was selected because the project is developed and operated by one student, the modules share a relational data model, and the current workload does not justify independent service deployment. The structure still maintains domain boundaries, so selected modules could later be separated if operational requirements change.

The architecture separates four major areas:

1. **Client layer:** Presents forms, journals, analytics, recommendations, reports, and administrative views.
2. **API and domain layer:** Handles routing, authentication, validation, business rules, analytics, report generation, and provider coordination.
3. **Persistence and infrastructure layer:** Stores relational records, cache/session values, and generated HTML files.
4. **Recommendation-provider layer:** Produces recommendation text through the default mock provider or optional external provider.

### 4.4.2 Component Design

**Figure 4.2: System Component Diagram**

*[Insert `docs/UML/rendered/png/02_component_diagram.png` here. A landscape page is recommended.]*

The React client sends JSON requests to the Gin router. CORS configuration allows the known local frontend origins. JWT and role middleware protect the relevant routes before requests reach domain modules.

The primary domain modules are authentication and profile, workout, nutrition, body, goal, and administration. Analytics is represented as a shared service because dashboard, recommendation, and report functions require the same calculated indicators.

GORM repositories provide structured access to MySQL. This separates persistence queries from HTTP handling and most business rules. Redis support is centralised conceptually but used for three distinct purposes:

- Authentication sessions and logout-aware validation.
- A five-minute cache for the weekly dashboard summary.
- Hourly per-user recommendation rate limiting.

The workout, nutrition, body, goal, and profile services invalidate the dashboard cache after relevant changes. Daily, monthly, and custom dashboard ranges are calculated directly because the current implementation caches only the default weekly summary.

The report service uses the analytics service and writes deterministic HTML. It does not call the recommendation provider. The recommendation service uses a weekly analytical summary, calls the configured provider, validates the returned text, and stores both the prompt context and accepted content.

### 4.4.3 Backend Module Structure

The backend applies feature-based modular organisation. Each major business domain has its own package. Within a module:

- The **controller** registers routes, reads path/query/body data, obtains the authenticated user identifier, and converts errors into HTTP responses.
- The **DTO** defines request payloads and input validation constraints.
- The **service** implements use-case and business logic.
- The **repository** performs database operations and ownership-filtered queries.
- The **analyzer** contains calculation or classification logic when the module requires independent analysis.

Not every small module requires every file type. For example, simple profile handling does not need a large standalone analyzer. The structure is applied according to responsibility rather than creating empty layers mechanically.

**Table 4.3: Backend Module Responsibilities**

| Module | Main Responsibility | Main Stored Entities |
|---|---|---|
| Auth/Profile | Registration, login, session, logout, identity, and profile | User, UserProfile |
| Workout | Exercise references and workout-session records | Exercise, WorkoutSession, WorkoutEntry |
| Nutrition | Food references, meal records, and nutrient totals | FoodItem, MealLog |
| Body | Date-based weight records | BodyRecord |
| Goal | Feasibility checks, goal lifecycle, and milestones | Goal, GoalMilestone |
| Analytics | Date ranges, aggregation, classifications, snapshots, and dashboard cache | AnalyticsSnapshot |
| Recommendation | Provider coordination, safety, rate limiting, history, and feedback | AIRecommendation, RecommendationFeedback |
| Report | Requested-period analytics, narrative construction, HTML generation, and download | ProgressReport |
| Admin | User-role management and protected reference-data operations | User, Exercise, FoodItem |

## 4.5 Data Design

### 4.5.1 Domain Class Model

**Figure 4.3: Core Domain Class Diagram**

*[Insert `docs/UML/rendered/png/03_class_diagram.png` here. Use a full landscape page.]*

The class model is centred on `User`. A user may have one profile and multiple workout sessions, meal logs, body records, goals, analytical snapshots, recommendations, and reports.

A workout session represents one dated workout. It contains one or more workout entries. Each entry links to an exercise reference and stores sets, repetitions, external weight, and rest time.

A meal log links one user and one food reference. Nutrient totals are stored in the meal record when it is created or updated. Storing these calculated totals makes historical retrieval simple and prevents each display request from repeatedly recalculating the same quantity.

A goal owns its milestones. Composition is used in the model because milestones have no independent purpose without their parent goal. The same concept applies to workout entries and recommendation feedback.

`AnalyticsSnapshot` stores selected summary fields only when persistence is requested. It is not the source of current dashboard truth. The dashboard is normally calculated from operational records so that it reflects the latest data.

### 4.5.2 Relational Database Design

**Figure 4.4: Entity-Relationship Diagram**

*[Insert `docs/UML/rendered/png/08_erd_bonus.png` here. Use a full landscape page.]*

The database contains fourteen main tables. Primary identifiers use UUID strings stored as `CHAR(36)`. UUIDs provide a consistent identifier format across modules and reduce the exposure of predictable sequential identifiers through public APIs. They do not replace authorisation; ownership checks are still required.

The main relationships are:

- One user to zero or one profile.
- One user to many workout sessions.
- One workout session to many workout entries.
- One exercise to many workout entries.
- One user to many meal logs.
- One food item to many meal logs.
- One user to many body records.
- One user to many goals.
- One goal to many milestones.
- One user to many analytics snapshots.
- One user to many recommendations.
- One recommendation to many feedback records.
- One user to many progress reports.

User identifiers, dates, period types, session identifiers, exercise identifiers, food identifiers, goal identifiers, and recommendation identifiers are indexed where they are commonly used for filtering or relationships. The user email address has a unique index.

### 4.5.3 Data Ownership and Integrity

Private entities contain a `user_id` either directly or through an owned parent. Repository queries combine the requested resource identifier with the authenticated user's identifier. This prevents a user from retrieving or modifying another user's records simply by guessing a UUID.

Shared exercise and food references do not belong to an individual user. General users can read and search them, while only administrators can create, update, or delete them.

Workout entries, goal milestones, and recommendation feedback depend on parent records. The service verifies ownership of the parent before allowing a nested record to be changed.

GORM AutoMigrate creates and updates the schema in the development environment. The application also seeds selected exercise and food reference records when they are absent. For a production system, versioned migration files and controlled deployment procedures would be preferable to unrestricted automatic migration.

## 4.6 API Design

The frontend and backend communicate through REST-style JSON APIs under `/api`. Public authentication operations are separated from protected application operations. The backend uses conventional HTTP methods:

- `GET` retrieves a resource or list.
- `POST` creates a resource or starts a generation process.
- `PUT` replaces editable resource data.
- `PATCH` changes a limited state such as a role, goal status, or milestone completion.
- `DELETE` removes an owned resource.

**Table 4.4: Main API Groups**

| API Group | Representative Endpoints | Purpose |
|---|---|---|
| Authentication | `POST /auth/register`, `POST /auth/login`, `GET /auth/me`, `POST /auth/logout` | Account and session operations |
| Profile | `GET /profile`, `PUT /profile` | User health profile |
| Workout | `/workouts`, `/workouts/:id`, `/workouts/exercises` | Daily workout journal and references |
| Nutrition | `/nutrition/meals`, `/nutrition/meals/:id`, `/nutrition/foods` | Daily meal journal and food references |
| Body | `/body-records`, `/body-records/:id` | Body-weight timeline |
| Goal | `/goals`, `/goals/:id/status`, `/goals/:id/milestones/:milestoneId` | Goal lifecycle and milestone completion |
| Dashboard | `GET /dashboard/summary`, `GET /dashboard/snapshots` | Date-range analytics |
| Recommendation | `/recommendations`, `/recommendations/generate`, `/recommendations/:id/feedback` | Recommendation workflow and feedback |
| Report | `/reports`, `/reports/generate`, `/reports/:id/download` | Report generation, history, and HTML download |
| Administration | `/admin/users`, `/admin/workouts/exercises`, `/admin/nutrition/foods` | Role and shared reference management |

API responses use a consistent JSON response helper. Validation errors return a client-error status, missing resources return `404 Not Found`, unauthenticated requests return `401 Unauthorized`, and administrator violations return `403 Forbidden`.

Date-range requests use the values `daily`, `weekly`, `monthly`, or `custom`. A custom request must include `startDate` and `endDate` in `YYYY-MM-DD` form, and the end date must not precede the start date.

## 4.7 Analytics Design

### 4.7.1 Shared Analytics Model

The dashboard, recommendation, and report modules use the same `DashboardSummary` analytical model. This design prevents each output channel from calculating its own version of calories, workout consistency, body trends, or goal adherence.

The analytics service performs the following general steps:

1. Resolve the requested date range.
2. Retrieve user-owned operational records.
3. Calculate workout indicators.
4. Calculate nutrition indicators.
5. Calculate BMR, TDEE, calorie balance, and body indicators.
6. Calculate goal and milestone indicators.
7. Optionally persist a snapshot.
8. Cache the result when it is the weekly summary.

### 4.7.2 Dashboard Activity Design

**Figure 4.5: Dashboard Analytics Activity Diagram**

*[Insert `docs/UML/rendered/png/16_dashboard_activity_diagram.png` here.]*

Figure 4.5 shows the dashboard decision flow. A weekly request first checks Redis. Other period types proceed directly to calculation. The service retrieves the user's profile, workouts, meals, body records, goals, and milestones and then calculates each group of indicators.

When the request includes `persist=true`, selected fields are saved as an analytics snapshot. This allows historical analytical summaries to be retained without forcing every dashboard request to create a database record.

The cached weekly summary has a five-minute lifetime. Any change to relevant operational data invalidates the cache. This balances repeated dashboard performance with consistency.

### 4.7.3 Workout Analysis Design

Workout load is calculated from sets, repetitions, and external load. If no external load is recorded, a value of one is used so that changes in sets and repetitions can still be compared.

Progressive-overload status requires at least two sessions. Sessions are ordered by date and divided into earlier and more recent groups. A change above 5% is classified as improving, below -5% as declining, and otherwise as stable. When an earlier baseline cannot be calculated, the system returns insufficient data.

Muscle-group distribution counts workout entries using the muscle group stored in the exercise reference. Warnings identify absent lower-body or back records and a distribution in which one group represents more than 60% of entries.

### 4.7.4 Nutrition and Calorie Design

Meal totals are derived from reference nutrients multiplied by the selected quantity. Protein and carbohydrate contribute four kcal per gram, while fat contributes nine kcal per gram when calculating the macro split.

Meal-quality scoring starts at 100 and applies deductions for low protein, high fat ratio, low carbohydrate ratio, and surplus classification. Protein expectation is scaled using days that actually contain meal records. This prevents a long custom range with sparse logs from being interpreted as many confirmed low-protein days.

BMR uses the Mifflin-St Jeor equation. Estimated TDEE multiplies BMR by a profile activity factor. The calorie-balance value compares total recorded calorie intake with estimated TDEE multiplied by the number of selected days. The value is an estimate and is not presented as a clinical measurement.

### 4.7.5 Body and Goal Design

The seven-day weight average uses the calendar window ending on the latest available record. This design excludes old records even when fewer than seven measurements exist.

Goal adherence is the percentage of completed milestones. Feasibility rules require at least seven days before the deadline, limit workout-frequency targets to seven sessions per week, and reject weight-loss targets that exceed the configured safe weekly threshold.

## 4.8 Key Interaction Designs

### 4.8.1 Goal and Milestone Workflow

**Figure 4.6: Goal Creation and Milestone Management Sequence**

*[Insert `docs/UML/rendered/png/15_goal_management_sequence.png` here.]*

The goal workflow begins with parsing and validating the deadline. The feasibility analyzer then evaluates the goal type, metric, target, and available time. An infeasible goal is rejected with a message that allows the user to revise the input.

When a goal is feasible, the service accepts custom milestones. If the request contains no milestones, the service generates three proportional default milestones. The goal and milestones are saved together, and the weekly dashboard cache is invalidated. Milestone completion is handled through a separate protected patch operation.

### 4.8.2 Recommendation Workflow

**Figure 4.7: Recommendation Generation and Feedback Sequence**

*[Insert `docs/UML/rendered/png/05_ai_recommendation_sequence.png` here. A landscape page is recommended.]*

Recommendation generation begins with a Redis rate-limit check. The service then builds a fresh weekly summary from MySQL. It stores the provider type, request type, and summary as structured prompt context.

The current mock provider applies deterministic rules to generate general-wellness guidance. If an OpenAI-compatible provider is configured, the same interface builds a constrained prompt and sends it to the external endpoint.

The content is validated before storage. Rejected content is not shown or saved. Accepted content is stored with its context. The user may later submit feedback, but only after ownership of the recommendation is verified.

### 4.8.3 Progress-Report Workflow

**Figure 4.8: Progress Report Generation and Download Sequence**

*[Insert `docs/UML/rendered/png/06_report_generation_sequence.png` here. A landscape page is recommended.]*

Report generation accepts daily, weekly, monthly, or custom periods. It uses the shared analytics service to calculate the requested range. A deterministic narrative and HTML document are then generated. The file is stored locally, while metadata and the protected download URL are stored in MySQL.

When downloading, the service verifies that the report belongs to the authenticated user before returning the HTML file. The design intentionally does not use the AI provider for report generation. This makes report content predictable and avoids unnecessary external processing.

## 4.9 Interface Design

### 4.9.1 Navigation Structure

**Figure 4.9: Interface Navigation Structure**

*[Insert `docs/UML/rendered/png/20_interface_navigation.png` here.]*

Registration and login are public interfaces. Successful authentication opens the application shell, which provides persistent navigation to the main modules. Administrator options are displayed only when the current authenticated user has the administrator role.

Workout and nutrition use dedicated search controls for reference selection. Dashboard data is derived from workout, nutrition, body, and goal records, but users can navigate directly to the source module when records need to be changed.

### 4.9.2 Screen Design

The interface uses a consistent application shell with a navigation sidebar, page heading, action area, message feedback, and module content. The design avoids requiring the user to move between unrelated applications or understand database relationships.

**Table 4.5: Interface Design by Screen**

| Screen | Main Design Elements | Main User Task |
|---|---|---|
| Login/Registration | Compact credential form, validation feedback, switch between account actions | Enter the protected system |
| Dashboard | Date-range selector, custom dates, metric cards, macro grams and ratios, distributions, warnings | Review combined progress |
| Profile | Personal measurements, activity level, and primary goal form | Maintain calculation inputs |
| Workout | Date selector, exercise search, entry form, daily history, edit and delete actions | Record and review daily training |
| Nutrition | Date selector, food search, quantity and meal fields, daily nutrient totals, history actions | Record and review daily meals |
| Body | Weight-entry form and chronological timeline | Record weight and inspect progress |
| Goals | Goal form, feasibility feedback, three editable milestones, status and completion actions | Define and monitor targets |
| Recommendations | Type action, generated content history, feedback controls, delete action | Obtain and evaluate guidance |
| Reports | Period selector, custom dates, generation action, report history, download and delete actions | Retain a structured summary |
| Administration | User-role table and CRUD forms for shared food and exercise data | Maintain controlled system data |

The daily journal pattern used by workout and nutrition screens was selected after prototype observation. Users normally think about today's workout or meals first, and historical editing becomes a date-selection action. Displaying every historical record in the same operational form would increase visual density and make daily totals less clear.

### 4.9.3 Dashboard Interface Mock-Up Evidence

**Figure 4.10: Implemented Dashboard Interface Design**

*[Insert a current dashboard screenshot here. Show the date-range selector, workout indicators, calorie/macronutrient values, meal-quality score, body indicators, and goal adherence. Use an anonymised demonstration account.]*

The dashboard places the selected time range before analytical output so that the interpretation context is visible. Macronutrients are shown using both gram totals and percentages. Long-range selection changes aggregated intake and other range-dependent values, while indicators that use fixed logic retain their documented interpretation.

### 4.9.4 Daily Journal Interface Mock-Up Evidence

**Figure 4.11: Implemented Daily Journal Design**

*[Insert either the Workout or Nutrition screen here. The screenshot should include the selected date, search control, entry form, daily summary, and history table.]*

The daily journal design combines record creation and historical review in one interface. Search narrows the shared reference dataset. Editing uses the selected record to populate the form, while deletion uses a separate action. This keeps the main workflow efficient without hiding management functions.

## 4.10 Security and Privacy Design

### 4.10.1 Authentication

Passwords are hashed using bcrypt before storage. Successful login produces a signed JWT containing the user identifier and expiry information. The JWT is sent in the `Authorization` header for protected requests.

Redis stores an active-session key for the user. The authentication middleware validates both the JWT and, when Redis is available, the session key. Logout deletes the session key so that the protected session cannot continue only because the JWT has not yet expired.

### 4.10.2 Authorisation

The general authentication middleware protects profile, workout, nutrition, body, goal, dashboard, recommendation, and report routes. A second middleware retrieves the user role and requires `admin` for administrative routes.

Resource-level authorisation is handled through ownership-filtered queries. A valid authenticated user is not automatically allowed to access every resource identifier.

### 4.10.3 Input and Output Protection

Gin request binding validates required fields. Services validate dates, status values, roles, goal constraints, and ownership. GORM parameterised query construction reduces direct SQL-injection risk.

The API does not return password hashes. Recommendation content is checked before storage. CORS allows only the configured local frontend origins in the current environment.

The current design should still be strengthened before production deployment. Required improvements include HTTPS, secure cookie or token-storage review, stronger password policy, login rate limiting, CSRF analysis, security headers, dependency scanning, audit logging, and a formal penetration test.

### 4.10.4 Privacy Boundaries

The current mock provider keeps recommendation processing inside the application and avoids transmitting personal summaries to an external model. If an external provider is enabled later, the system should minimise the context, obtain appropriate user consent, review provider retention terms, and avoid sending direct identifiers where possible.

Generated reports are stored as local files. The download route verifies ownership, but production storage should additionally use protected object storage, retention controls, backup rules, and secure deletion procedures.

## 4.11 Requirement Traceability

Traceability demonstrates that project objectives are connected to design components and planned verification.

**Table 4.6: Requirement Traceability Matrix**

| Objective Area | Related Requirements | Design Components | Verification Evidence |
|---|---|---|---|
| Authentication and profile | FR-01 to FR-03, FR-25 | Auth/Profile module, JWT middleware, Redis session | Registration, login, profile, logout workflow |
| Workout tracking and analysis | FR-04 to FR-06 | Workout module, Exercise and Workout entities, analytics service | CRUD checks and progressive-load unit tests |
| Nutrition tracking and analysis | FR-07 to FR-09 | Nutrition module, FoodItem and MealLog, analytics service | Meal workflow and macro/meal-quality unit tests |
| Body and calorie progress | FR-10 to FR-11 | Body module, profile inputs, analytics service | Body workflow, BMR and moving-average unit tests |
| Goal management | FR-12 to FR-14 | Goal service, feasibility analyzer, GoalMilestone | Goal sequence and milestone workflow |
| Dashboard | FR-15 to FR-16 | Shared analytics, range resolver, Redis weekly cache | Daily/weekly/monthly/custom manual checks |
| Recommendation | FR-17 to FR-18 | Recommendation service, provider, safety validator, Redis limiter | Mock generation and feedback workflow |
| Reporting | FR-19 to FR-20 | Report service, HTML builder, protected file download | Generate, list, download, and delete workflow |
| Administration | FR-21 to FR-23 | Admin middleware and protected CRUD routes | User-role and reference-management checks |
| Privacy and isolation | FR-24, NFR-05 to NFR-08 | Password hashing, middleware, ownership repositories | Unauthenticated and unauthorised access checks |

The matrix also shows remaining evaluation gaps. Automated tests currently provide strong evidence for selected calculations but less evidence for HTTP integration, frontend components, security, and performance. These areas are identified for expansion rather than being presented as fully validated.

## 4.12 Design Decisions and Alternatives

### 4.12.1 Monorepository

The frontend, backend, documentation, UML sources, and local tooling are stored in one repository. This simplifies coordinated changes to API contracts, local setup, diagrams, and documentation. Separate repositories would provide stronger independent release boundaries but would add unnecessary coordination overhead for the current individual project.

### 4.12.2 Modular Monolith Instead of Microservices

Microservices were not selected because the project does not require independent scaling, separate teams, or complex deployment. A modular monolith provides simpler local operation and transactions while still demonstrating domain separation.

### 4.12.3 MySQL Instead of a Document Database

The application contains clear relationships among users, sessions, entries, reference records, goals, milestones, recommendations, feedback, and reports. A relational database supports these relationships and ownership queries naturally.

### 4.12.4 Redis as Supporting Infrastructure

Redis is not used as the primary database. It stores temporary operational state: active sessions, the weekly summary cache, and recommendation counters. The system's durable source of truth remains MySQL.

### 4.12.5 Deterministic Analytics Before AI

Core calculations are not delegated to an LLM. This improves repeatability, testability, and explanation. The provider receives the resulting indicators and produces natural-language guidance. This design also allows the current phase to use a mock provider without removing the recommendation workflow.

### 4.12.6 HTML Report Instead of PDF

HTML was selected for the current implementation because it can be generated without a heavy rendering dependency, opened in common browsers, styled clearly, and downloaded through the existing web architecture. PDF export remains a future enhancement when fixed-layout printing becomes necessary.

## 4.13 Chapter Summary

This chapter translated the project objectives into twenty-five functional and fifteen non-functional requirements. The use-case analysis identified general users, administrators, and the configured recommendation provider as the principal operational actors.

The system was designed as a React client and modular Go REST API supported by MySQL, Redis, and local HTML report storage. The component, class, and entity-relationship designs established clear module and data responsibilities. Shared analytics provide consistent workout, nutrition, calorie, body, and goal indicators to the dashboard, recommendation, and report modules.

The sequence and activity designs described dashboard aggregation, goal feasibility and milestones, controlled recommendation generation, feedback, and progress-report generation. Interface design was organised around date-based journals and persistent module navigation. Authentication, role checks, ownership filtering, safety validation, and privacy boundaries were included as explicit design concerns.

The following chapter describes how these designs were implemented in the frontend, backend, database, cache, analytics, recommendation, and reporting components and presents evidence from the running application.


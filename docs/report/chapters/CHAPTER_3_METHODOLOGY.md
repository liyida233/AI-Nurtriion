# CHAPTER 3: METHODOLOGY

## 3.1 Introduction

This chapter describes the methodology used to analyse, design, develop, and evaluate the AI-Enhanced Fitness and Nutrition Tracking System. The project is both a software-engineering project and an applied personal health decision-support project. Therefore, its methodology must address two related concerns. First, the system must be developed through a controlled process that converts requirements into a working full-stack application. Second, the calculations, recommendations, and user-facing outputs must be evaluated for correctness, clarity, privacy, and safety.

An iterative and incremental development methodology was selected. The system was divided into functional modules, and each module was developed through repeated cycles of requirement refinement, design, implementation, integration, and testing. This approach allowed the project scope to evolve in response to supervisor feedback and observations from the working prototype. It was particularly suitable because several requirements, including custom date ranges, daily journal behaviour, custom goal milestones, progressive-overload calculation, and nutrition presentation, became clearer after the relevant interfaces were implemented and used.

The chapter explains the selected development methodology, requirement-gathering methods, development phases, technologies and tools, data-analysis approach, testing strategy, quality controls, ethical considerations, and limitations of the evaluation.

## 3.2 Research and Development Approach

### 3.2.1 Applied System-Development Approach

The project follows an applied system-development approach. Its primary output is a functional web-based platform rather than a purely theoretical model. Research activities were used to identify the problem, study existing systems, select suitable analytical methods, and establish safety boundaries. Software-engineering activities were then used to convert these findings into system requirements, architecture, source code, database structures, interfaces, and testable outputs.

The project does not attempt to conduct clinical research or prove that the platform improves medical outcomes. Instead, it evaluates whether an integrated software system can reliably record fitness and nutrition data, calculate transparent general-wellness indicators, present understandable progress information, and support a controlled recommendation workflow.

### 3.2.2 Iterative and Incremental Development

The development process was influenced by the Agile principle of delivering working software and responding to change (Beck et al., 2001). However, the project did not apply a complete formal Scrum process because it was developed as an individual final-year project without a permanent Scrum team, product owner, or fixed sprint ceremonies. The more accurate description is an iterative and incremental methodology.

In an incremental process, the system is built as a sequence of usable parts. In an iterative process, an existing part is reviewed and improved when testing or feedback identifies a limitation. The project applied both ideas. Authentication and profile management established the system foundation. Workout, nutrition, body, and goal modules were then implemented as independent increments. Analytics, recommendations, reports, administration, and frontend integration were added after the core records were available. Existing features were subsequently refined instead of being treated as permanently complete after their first implementation.

Examples of iterative refinement included:

- Replacing fixed weekly and monthly dashboard behaviour with daily, weekly, monthly, and custom date ranges.
- Changing workout and nutrition pages into date-based daily journals.
- Adding food and exercise search rather than requiring users to browse unfiltered lists.
- Changing goal milestones from fixed generated values to user-defined milestones with generated defaults as a fallback.
- Correcting the seven-day body-weight average to use a seven-calendar-day window rather than seven historical records.
- Allowing progressive-overload analysis to use sets and repetitions when external weight is unavailable.
- Revising meal-quality analysis so that days without meal logs do not automatically reduce the score across a long selected range.

These changes demonstrate why a single-pass sequential process would have been less appropriate. Several usability and analytical issues became visible only when real data was entered into the integrated interface.

**Figure 3.1: Iterative and Incremental Development Workflow**

*[Insert the rendered development methodology diagram here. Source: `docs/UML/19_development_methodology.puml`.]*

Figure 3.1 shows the repeated cycle used in the project. Requirements were first identified and prioritised. The relevant architecture, database entities, API contracts, and interface flow were then designed. A small functional increment was implemented and checked. If its behaviour did not meet the requirement, the design or implementation was revised. Accepted increments were integrated into the full system before the next requirement was selected.

## 3.3 Requirement-Gathering Methods

Requirements were gathered using a combination of document analysis, literature review, existing-system analysis, supervisor consultation, prototype observation, and developer testing. Using several sources reduced dependence on a single opinion and allowed requirements to be checked against both academic expectations and practical system behaviour.

### 3.3.1 Academic and Project Document Analysis

The project brief, proposal requirements, assessment criteria, submission requirements, and sample FYP report were reviewed to identify the expected academic and technical outputs. This analysis established that the project needed measurable objectives, sufficient technical complexity, complete system analysis and design, evidence of implementation, and appropriate documentation.

The document analysis also influenced the report structure and project evidence. UML diagrams, an entity-relationship model, interface evidence, implementation details, testing results, limitations, and up-to-date references were treated as required project artefacts rather than optional supporting material.

### 3.3.2 Literature Review

Recent journal articles, systematic reviews, standards, and official guidance were examined to identify suitable theories and methods. The literature review covered digital self-monitoring, behaviour-change techniques, resistance-training analysis, nutrition tracking, energy estimation, body-weight trends, digital goal setting, and AI-assisted recommendations.

The literature directly influenced several requirements. Self-monitoring and feedback supported the use of daily journals and dashboard summaries. Goal-setting research supported measurable goals and milestones. Resistance-training literature informed training-volume and progressive-overload indicators. Nutrition methods informed calorie and macronutrient calculations. WHO guidance on AI for health informed the decision to restrict the system to general-wellness recommendations and to separate deterministic calculations from generated text.

### 3.3.3 Existing-System Analysis

MyFitnessPal, Cronometer, Strong, and Fitbod were reviewed as representative existing products. The comparison identified common functions, specialised strengths, and opportunities for the proposed system. Nutrition platforms demonstrated effective food logging and nutrient summaries, while workout platforms demonstrated exercise histories and training-oriented feedback.

The proposed platform was not designed to reproduce the full food database of a mature nutrition product or the adaptive programme generation of a specialised training product. Instead, existing-system analysis supported the decision to integrate several domains and make the project's analytical rules visible for academic evaluation.

### 3.3.4 Supervisor Consultation and Monitoring Feedback

Supervisor communication was used to refine the project direction and presentation of the work. The recorded project activities included:

- A discussion on 10 April 2026 concerning how to identify a collaborator and obtain a collaboration letter.
- An online meeting on 15 April 2026 concerning the proposal, final-project expectations, assessment criteria, and recommendations for improving the project.
- Submission of the monitoring-session presentation on 16 May 2026 to obtain feedback.

These activities helped clarify the expected module complexity, the distinction between supporting functions and core project contributions, and the need to explain analysis and design rather than presenting the project as a collection of CRUD pages.

The signed project logbook is retained as evidence of supervisor consultation. A collaboration letter may be included as supplementary stakeholder evidence if obtained, but it is not treated as evidence of completed stakeholder evaluation unless actual requirements or feedback were collected from that stakeholder.

### 3.3.5 Prototype Observation and Self-Evaluation

The working prototype was repeatedly used with controlled sample data. This method was important because it revealed problems that were not obvious from static requirements. For example, entering body weights with a large gap between dates exposed an error in the interpretation of a seven-day average. Entering multiple workout sessions exposed the need to define the minimum usable data and fallback calculation for progressive overload. Selecting a long dashboard date range exposed inconsistent behaviour between macro summaries and meal-quality scoring.

Prototype observation therefore served as both a requirement-refinement method and an informal usability review. It did not replace formal usability testing with independent participants, but it provided practical evidence for correcting workflow and calculation issues before such testing.

**Table 3.1: Requirement-Gathering Sources**

| Source | Method | Main Information Obtained | Resulting Project Evidence |
|---|---|---|---|
| Assessment and project documents | Document analysis | Expected report sections, technical depth, deliverables, and evaluation dimensions | Report plan, objectives, UML and implementation evidence |
| Academic literature and standards | Secondary research | Relevant theories, formulas, safety concerns, and evaluation principles | Literature review and analytical requirements |
| Existing products | Feature comparison | Common tracking features, specialised strengths, and system gaps | Product-comparison table and integrated scope |
| Supervisor discussions | Consultation and feedback | Scope refinement, collaborator process, proposal expectations, and project advice | Logbook entries and revised proposal |
| Working prototype | Observation with sample records | Workflow problems, ambiguous indicators, calculation defects, and missing search functions | Iterative interface and analytics improvements |
| Future sample users or stakeholder | Planned task-based evaluation | Ease of use, clarity of analytics, and usefulness of reports | Usability findings and final-phase refinements |

## 3.4 Development Phases

### 3.4.1 Phase 1: Project Initiation and Scope Refinement

The initial phase defined the project topic, target users, general problem, and expected contribution. The system was positioned as a personal health decision-support platform rather than a medical system. The scope was divided into supporting functions and core analytical modules.

Authentication, profile management, database access, and shared reference data were considered supporting foundations. Workout analysis, nutrition analysis, body and calorie analysis, adaptive goal monitoring, recommendations, dashboard analytics, and reporting were treated as the main project contributions.

Scope boundaries were also identified during this phase. Clinical diagnosis, direct wearable integration, image-based food recognition, native mobile applications, production deployment, and real-time use of a production LLM were excluded from the current minimum scope.

### 3.4.2 Phase 2: System Analysis and Design

The second phase translated the project scope into functional and non-functional requirements. Actors, use cases, data entities, relationships, module boundaries, API responsibilities, and user flows were analysed.

UML diagrams were used to represent the system from several perspectives:

- A use-case diagram described user and administrator functions.
- A component diagram described the frontend, backend, MySQL, Redis, and recommendation provider.
- A class diagram and ERD described the application data.
- Activity diagrams represented registration, dashboard analysis, and general data-recording workflows.
- Sequence diagrams represented authentication, profile, workout, nutrition, body, goal, recommendation, report, and administration interactions.

Interface planning focused on direct operational workflows. The dashboard served as the combined analytical view, while workout and nutrition pages used date selection to support daily recording and historical editing. Detailed design artefacts are presented in Chapter 4.

### 3.4.3 Phase 3: Development Environment and Data Layer

The project was organised as a monorepository containing separate `frontend`, `backend`, `docs`, and `tools` directories. MySQL and Redis were configured through Docker Compose to provide a reproducible local infrastructure.

The relational model was implemented using GORM. UUID values were selected as identifiers across the major entities. Database migrations establish tables for users, profiles, exercise references, workout sessions and entries, food references, meal logs, body records, goals, milestones, analytical snapshots, recommendations, feedback, and reports.

Seed data was used to provide initial exercise and food references. This allowed search, logging, analytics, and administration workflows to be demonstrated without requiring the user to manually create every shared reference item.

### 3.4.4 Phase 4: Backend Development

The backend was implemented using Go, Gin, and GORM. Development proceeded by module so that each domain owned its HTTP handling, data-transfer objects, business logic, persistence logic, and analytical logic where required.

The backend increments were implemented in the following general order:

1. Configuration, database connection, Redis connection, routing, and common response handling.
2. User registration, login, JWT generation, session validation, logout, and profile management.
3. Exercise and food reference data.
4. Workout, meal, body, and goal CRUD operations.
5. Workout, nutrition, calorie, body-progress, and goal-adherence analytics.
6. Dashboard aggregation and Redis caching.
7. Recommendation context generation, provider abstraction, safety validation, rate limiting, and feedback.
8. Progress-report generation and HTML download.
9. Administrator role protection and reference-data management.

The backend was tested after each major increment through compilation, direct API use, database inspection, and targeted unit tests.

### 3.4.5 Phase 5: Frontend Development

The frontend was implemented using React, TypeScript, and Vite. Shared API utilities and authentication state were established before module pages were connected.

The interface was developed around user tasks rather than backend table names. For example, workout and nutrition interfaces provide a selected-date journal, searchable references, data-entry forms, daily totals, and history actions. The dashboard provides date-range selection and combined indicators. Goal pages combine target information and milestone progress. Recommendation and report pages expose the corresponding backend workflows without requiring users to understand the underlying provider or storage design.

TypeScript compilation and the Vite production build were used to detect type and bundling problems. Browser-based use of the complete application was used to inspect responsive layout, navigation, form behaviour, API errors, empty states, and the display of returned data.

### 3.4.6 Phase 6: Integration and Iterative Refinement

Frontend-backend integration connected the React application to RESTful endpoints under `/api`. MySQL stored persistent records, while Redis supported active-session checks, logout behaviour, selected dashboard caching, cache invalidation, and recommendation rate limiting.

Integration was performed module by module. The typical workflow was:

1. Create or select a reference item.
2. Submit a record through the frontend.
3. Validate the API response.
4. Inspect the record in the daily journal.
5. Confirm persistence in MySQL.
6. Confirm that dashboard analytics changed as expected.
7. Edit or delete the record and confirm recalculation and cache invalidation.

Issues found during integration were returned to the relevant requirement, service, analyzer, or interface. This feedback loop produced the refinements described in Section 3.2.2.

### 3.4.7 Phase 7: Evaluation and Documentation

The final current phase consolidates implementation evidence and evaluates whether the project objectives have been achieved. Source code, UML diagrams, test results, interface screenshots, database evidence, and limitations are being prepared for the report and presentation.

Formal usability testing and production-scale performance testing remain areas for further evaluation. They should be conducted using defined tasks, participant consent, anonymised observations, and measurable criteria before making claims about user acceptance or long-term behavioural impact.

**Table 3.2: Development Phases and Outputs**

| Phase | Main Activities | Main Outputs |
|---|---|---|
| Initiation | Define problem, users, contribution, scope, and boundaries | Revised project scope and objectives |
| Requirements | Review documents, literature, products, and feedback | Functional and non-functional requirements |
| Analysis and design | Model actors, processes, components, data, and interfaces | UML diagrams, ERD, API and interface plan |
| Backend implementation | Build modular REST APIs, persistence, analytics, recommendations, and reports | Go backend and database schema |
| Frontend implementation | Build task-based pages and API integration | React and TypeScript web application |
| Integration and refinement | Run complete workflows and correct discovered issues | Integrated locally runnable system |
| Evaluation and reporting | Run tests, gather evidence, document limitations, and prepare presentation | Test evidence, screenshots, report, and demo |

## 3.5 Development Technologies and Tools

Technology selection considered development productivity, type safety, modularity, performance, local reproducibility, and suitability for demonstrating full-stack engineering.

**Table 3.3: Development Technologies**

| Layer or Purpose | Technology | Reason for Selection |
|---|---|---|
| Frontend | React 19 | Component-based construction of interactive module interfaces |
| Frontend language | TypeScript 5 | Static type checking for API data and interface logic |
| Frontend tooling | Vite 6 | Fast development server and production build process |
| Backend | Go 1.24 | Strong typing, simple deployment model, and suitability for modular HTTP services |
| Web framework | Gin | Lightweight routing, middleware, validation, and JSON API support |
| ORM | GORM | Structured relational mapping, queries, associations, and migration support |
| Primary database | MySQL 8.4 | Relational consistency for connected user, workout, nutrition, body, goal, and report data |
| Cache and session support | Redis 7.4 | Fast session checks, cache storage, invalidation, and rate limiting |
| Local infrastructure | Docker Compose | Repeatable MySQL and Redis configuration |
| Authentication | JWT and bcrypt | Stateless token identity with secure password hashing and Redis-backed session checks |
| Diagramming | PlantUML | Version-controlled and reproducible UML diagrams |
| Version control | Git | Source history and controlled project changes |

The technology choices are discussed as implementation decisions rather than claims that they are the only suitable technologies. Alternative frameworks could implement the same architecture. The selected stack was appropriate because it supported the project's modular requirements and could be operated locally without paid infrastructure.

## 3.6 Data Processing and Analytical Method

### 3.6.1 Structured Data Collection

The platform collects data entered directly by the authenticated user. The main operational data consists of:

- Profile data such as age, sex, height, weight, activity level, and fitness objective.
- Workout data such as date, exercise, sets, repetitions, load, duration, rest time, and notes.
- Nutrition data such as date, food, quantity, meal type, meal time, and nutrient values.
- Body data such as record date, weight, and notes.
- Goal data such as metric, target value, deadline, priority, status, and milestones.
- Recommendation feedback such as rating and optional comments.

No wearable or clinical records are automatically imported. Reference foods and exercises are selected from system-managed data. This reduces ambiguity in the current prototype, although it also limits the breadth of available foods and exercises.

### 3.6.2 Deterministic Analytics

Core numerical outputs are calculated through deterministic application logic. The same input therefore produces the same analytical result. The main methods include:

- Workout load based on sets, repetitions, and external load, with an unloaded fallback.
- Workout consistency based on recorded sessions and the selected date range.
- Progressive-overload classification based on earlier and recent usable sessions.
- Muscle-group distribution based on exercise reference categories.
- Daily calories and nutrient totals based on food values and quantity.
- Macronutrient energy ratios using 4 kcal per gram for protein and carbohydrate and 9 kcal per gram for fat.
- BMR using the Mifflin-St Jeor equation.
- Estimated TDEE using BMR and a profile activity factor.
- Calorie balance based on recorded intake and estimated expenditure.
- Seven-day body-weight average based on calendar dates.
- Goal adherence based on completed and total milestones.
- Meal-quality and warning rules based on transparent thresholds.

These calculations are general-wellness estimates. They are not treated as diagnoses, clinical assessments, or guarantees of future results.

### 3.6.3 Hybrid Recommendation Method

The recommendation workflow separates controlled analysis from natural-language generation. The backend first calculates structured indicators. It then prepares a limited context containing relevant profile, workout, nutrition, body, and goal summaries. A provider converts this context into a readable recommendation.

The current provider is a deterministic mock implementation. This makes the flow stable during development and prevents personal information from being transmitted to an external AI service. The provider interface allows a future OpenAI-compatible model to be configured without replacing the surrounding recommendation workflow.

Generated content is passed through a basic safety validator before it is stored or shown. Redis limits recommendation generation frequency. User ratings and comments are stored separately to support later evaluation. This hybrid method keeps numerical calculations under application control and treats generated text as an explanation layer.

## 3.7 Testing and Evaluation Strategy

Testing was organised according to functional correctness, analytical correctness, integration, interface behaviour, security controls, and build quality. ISO/IEC 25010:2023 identifies a software product quality model that can be used when specifying and evaluating quality properties. The project places particular emphasis on functional suitability, usability, reliability, performance efficiency, security, and maintainability (International Organization for Standardization [ISO], 2023).

### 3.7.1 Unit and Calculation Testing

Targeted Go unit tests were implemented for analytical functions that have a high risk of producing misleading output. The current automated tests cover:

- Mifflin-St Jeor BMR estimation.
- Protein, carbohydrate, and fat energy ratios.
- Progressive-overload classification with external load.
- Progressive-overload fallback when weight is unavailable.
- Muscle-group warning rules.
- Meal-quality scoring.
- Meal-quality behaviour across sparse long-range records.
- Seven-calendar-day body-weight averaging.

Controlled inputs and known expected outputs are used. For example, the seven-day moving-average test includes records separated by two months and verifies that the older record is excluded. This test represents a defect discovered during prototype use and prevents the same behaviour from returning unnoticed.

### 3.7.2 Functional and API Testing

Functional tests verify that each major user operation can be completed and that invalid or unauthorised actions are rejected. The main test areas are:

- Registration, login, authenticated identity, and logout.
- Profile creation and update.
- Creation, retrieval, update, and deletion of workout, meal, body, and goal records.
- Exercise and food search.
- Goal status and milestone completion.
- Daily, weekly, monthly, and custom dashboard queries.
- Recommendation generation, history, feedback, and deletion.
- Report generation, retrieval, download, and deletion.
- Administrator role checks and reference-data management.

API responses are checked for status, response structure, ownership isolation, validation errors, and corresponding database changes.

### 3.7.3 Integration Testing

Integration testing verifies communication among the React frontend, Go backend, MySQL database, and Redis. It also verifies that changing operational data affects downstream analytics. Important integration cases include:

- A new meal changes daily nutrient totals and dashboard summaries.
- A changed workout updates training-load indicators.
- A new body record changes latest weight and the seven-day average.
- Completing a milestone changes goal adherence.
- Editing or deleting relevant data invalidates the cached dashboard summary.
- Logging out removes the active Redis session and prevents reuse of the protected session.

### 3.7.4 Build and Static Verification

The backend is checked using `go test ./...`, which compiles all packages and runs available unit tests. The frontend is checked using `npm run build`, which performs TypeScript compilation and creates the production Vite bundle.

At the time of this chapter's latest verification, the backend test command passed and the frontend production build completed successfully. These results demonstrate compilation and current unit-test success, but they do not replace complete security, usability, load, or production-environment testing.

### 3.7.5 Interface and Usability Evaluation

Current interface evaluation is based on repeated task execution and visual inspection using realistic sample records. The evaluated tasks include logging a workout, logging a meal, changing dates, finding a reference item, reviewing analytics, completing a milestone, generating a recommendation, and downloading a report.

A more formal usability evaluation is proposed for the next project phase. Participants should be asked to perform defined tasks, after which completion success, observed errors, completion time, and qualitative feedback can be recorded. The evaluation should focus on:

- Ease of daily workout and meal entry.
- Clarity of date selection and historical editing.
- Understanding of calories, macronutrients, progressive load, meal quality, and body trends.
- Usefulness of goals, recommendations, and progress reports.
- Confusing labels, unnecessary steps, and missing information.

Until this evaluation is performed, the report should describe the interface as implemented and internally reviewed rather than claiming proven usability.

**Table 3.4: Testing Matrix**

| Test Type | Example Method | Expected Evidence | Current Status |
|---|---|---|---|
| Unit/calculation | Execute analyzer tests with known inputs | Passing automated tests and expected values | Completed for selected high-risk calculations |
| Backend compilation | Run `go test ./...` | All Go packages compile | Completed |
| Frontend build | Run `npm run build` | TypeScript and Vite build succeeds | Completed |
| Functional API | Execute CRUD and role workflows | Correct responses, validation, and persistence | Performed during module development |
| Integration | Use complete frontend-to-database workflows | Consistent records and recalculated analytics | Performed during local development |
| Security checks | Test authentication, ownership, roles, and logout | Rejected unauthorised access | Partially completed |
| Usability | Participant task evaluation | Completion measures and structured feedback | Planned |
| Performance | Repeated dashboard and report requests | Response-time and cache comparison | Planned |
| AI quality | Review relevance, safety, and context use | Ratings and reviewer observations | Mock flow completed; real-provider study planned |

## 3.8 Quality and Risk Control

### 3.8.1 Functional and Data Quality

Request validation is applied at the API boundary, while business rules are applied in services and analyzers. Database relationships and user ownership are used to prevent records from being mixed between users. Reference data provides consistent exercise and food metadata.

High-risk calculations are kept in named functions so that they can be tested independently. Analytics are not embedded only in frontend display code because this could produce different results on different pages. Central backend calculation also makes the rules easier to document and revise.

### 3.8.2 Security Control

The security method includes bcrypt password hashing, JWT validation, Redis-backed active-session checks, authenticated route middleware, resource ownership filtering, and administrator role middleware. The OWASP Application Security Verification Standard provides a reference for testing web application technical security controls and secure-development requirements (OWASP Foundation, n.d.).

The current project does not claim full ASVS compliance or a professional penetration test. The standard is used as a guide for important control categories. Further work should include systematic tests for token misuse, injection, cross-site scripting, insecure configuration, brute-force protection, and dependency vulnerabilities.

### 3.8.3 Performance and Availability

MySQL indexes and date-based queries support operational retrieval. Redis caches the selected dashboard summary for a limited period and is invalidated when relevant user data changes. Redis is also used for session checks and recommendation rate limiting.

The backend can continue without Redis in selected development scenarios, but caching, active-session enforcement, and rate limiting are reduced when Redis is unavailable. This graceful degradation supports local development while making the operational dependency visible.

### 3.8.4 Project Risks

**Table 3.5: Project Risks and Mitigation**

| Risk | Potential Effect | Mitigation |
|---|---|---|
| Project scope becomes too large | Core modules remain incomplete | Prioritise integrated core workflows and defer wearables, mobile apps, PDF, and cloud deployment |
| Incorrect analytical formula | Misleading user feedback | Use documented formulas, transparent rules, controlled inputs, and unit tests |
| Sparse or inconsistent user data | Unstable trends or insufficient indicators | Show insufficient-data states and use calendar-aware logic |
| External AI cost or unavailability | Recommendation feature cannot be demonstrated | Use provider abstraction and a stable mock provider |
| Unsafe generated advice | User may interpret output as medical guidance | Restrict to general wellness, validate output, and display limitations |
| Exposure of personal data | Privacy and ethical concern | Minimise collected data, authenticate access, enforce ownership, and avoid external AI transfer in the current phase |
| Cache becomes stale | Dashboard does not reflect recent records | Invalidate dashboard cache after relevant create, update, and delete operations |
| Lack of independent user evaluation | Limited evidence of usability | Report the limitation and conduct structured testing in the next phase |

## 3.9 Ethical, Privacy, and Safety Considerations

The system processes personal health-related information and must therefore be developed conservatively even though it is not a clinical system.

### 3.9.1 Informed and Minimal Data Collection

Only information required for the implemented functions should be collected. Users enter their own profile and tracking records. Future usability participants should be informed about the purpose of the evaluation, the tasks performed, the data recorded, and their right to stop participating.

The evaluation should avoid requesting actual sensitive medical information. Synthetic or non-sensitive sample data can be used when the purpose is to test interface behaviour rather than health outcomes.

### 3.9.2 Privacy and Access Control

User records are associated with an authenticated UUID and filtered by ownership. Passwords are stored as hashes rather than plaintext. Administrative functions are separated from ordinary user routes.

The current local prototype is not presented as a production health-data service. Production deployment would require stronger secret management, HTTPS, secure backup, retention policies, monitoring, incident response, and compliance review under the relevant jurisdiction.

### 3.9.3 Recommendation Safety

Recommendations are restricted to general fitness and nutrition support. The system must not diagnose conditions, prescribe treatment, claim to cure disease, or encourage starvation, extreme fasting, or unsafe exercise. Numerical indicators are estimates and are accompanied by appropriate interpretation boundaries.

The current mock provider prevents uncontrolled model output and external data transmission. A future real-provider integration should include explicit consent, data minimisation, prompt review, output monitoring, provider terms and privacy review, and stronger safety evaluation.

### 3.9.4 Academic Integrity

Literature, standards, formulas, and external product information must be cited. Generated assistance and software tools should be used in accordance with university rules. The student remains responsible for understanding, verifying, presenting, and defending the system design and implementation during the viva.

## 3.10 Methodological Limitations

The selected methodology has several limitations.

First, the project is developed by one student. Requirement interpretation, implementation, and much of the current testing are therefore performed by the same person, which can introduce confirmation bias.

Second, prototype observation has identified several useful improvements, but it is not equivalent to an independent usability study. User acceptance and long-term engagement cannot yet be concluded.

Third, the food and exercise datasets are curated demonstration datasets rather than comprehensive commercial or clinical databases. Analytical results depend on the accuracy and completeness of the selected reference values and user records.

Fourth, the current recommendation provider is a mock implementation. The architecture and safety flow can be evaluated, but the quality, variability, latency, cost, and risk of a production LLM require separate testing.

Fifth, current automated test coverage concentrates on selected analytical functions. Broader service, repository, HTTP integration, frontend component, security, performance, and end-to-end automation should be added in the next development phase.

These limitations do not invalidate the system as a functional prototype. They define the boundaries within which the current results should be interpreted and provide a clear plan for later evaluation.

## 3.11 Chapter Summary

This chapter presented the methodology used to develop the AI-Enhanced Fitness and Nutrition Tracking System. An applied, iterative, and incremental approach was selected because the project required a working system and because several requirements became clearer through repeated prototype use.

Requirements were obtained from academic and assessment documents, literature, existing products, supervisor consultation, prototype observation, and developer testing. Development proceeded through initiation, requirements analysis, system design, backend implementation, frontend implementation, integration, refinement, evaluation, and documentation.

The project combines deterministic analytics with a controlled recommendation-provider workflow. Testing includes selected calculation unit tests, backend compilation, frontend production building, functional API checks, and integrated workflow checks. Formal usability, performance, security, and real-AI evaluation remain future activities. Ethical controls emphasise data minimisation, access control, general-wellness boundaries, recommendation safety, and honest reporting of limitations.


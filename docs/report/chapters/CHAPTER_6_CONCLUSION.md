# CHAPTER 6: CONCLUSION AND FUTURE WORK

## 6.1 Introduction

This chapter concludes the development of the AI-Enhanced Fitness and Nutrition Tracking System for the current project phase. It summarises the completed work, evaluates the achievement of the project objectives, discusses the project's contribution, identifies limitations, and proposes future improvements.

The project began from the observation that personal workout, nutrition, body, and goal information is often fragmented or presented as separate records. Raw data alone may not explain whether training load is improving, nutrition is balanced, body weight is changing meaningfully, or personal milestones are being achieved. The project therefore aimed to implement an integrated web-based decision-support platform that converts user records into understandable general-wellness indicators, recommendations, and progress reports.

The resulting prototype is a locally runnable full-stack system rather than only a design proposal. It includes a React and TypeScript frontend, a modular Go REST API, MySQL persistence, Redis session and cache support, deterministic health-related analytics, recommendation-provider abstraction, HTML report generation, and role-based administration.

## 6.2 Summary of Project Outcomes

The completed system supports an end-to-end workflow from account creation to progress review.

An ordinary user can register, log in, configure a profile, record workouts, record meals, record body weight, create goals with milestones, view date-range analytics, request recommendations, submit recommendation feedback, and generate downloadable reports. An administrator can manage user roles and shared exercise and food references.

Workout data is transformed into training volume, consistency, progressive-load status, muscle-group distribution, and warnings. Nutrition data is transformed into calorie and macronutrient totals, macro ratios, nutrition gaps, and a rule-based meal-quality score. Profile, nutrition, and body data are used to estimate BMR, TDEE, calorie balance, body-weight trend, and a seven-day moving average. Goal milestones are used to calculate adherence.

The dashboard supports daily, weekly, monthly, and custom ranges. Workout and nutrition pages use date-based daily journals, allowing users to switch dates when reviewing or editing historical records. Reports support the same range types and can be downloaded as HTML.

The recommendation workflow currently uses a deterministic mock provider. This provider consumes a structured weekly analytical summary and produces general-wellness feedback. The backend also contains an optional OpenAI-compatible provider implementation, but it has not been enabled or formally evaluated in the current phase.

The latest verification confirmed that the backend compiles and its current automated tests pass. The frontend also passes TypeScript compilation and the Vite production build. An anonymised demonstration account was used to verify the integrated workflows with MySQL and Redis running locally.

## 6.3 Evaluation of Project Objectives

The project defined ten objectives in Chapter 1. Table 6.1 evaluates each objective against the implemented evidence.

**Table 6.1: Achievement of Project Objectives**

| No. | Project Objective | Achievement | Evidence and Qualification |
|---|---|---|---|
| 1 | Design a modular full-stack application integrating profiles, workouts, nutrition, body progress, goals, analytics, recommendations, and reports | Achieved | React frontend, modular Go backend, MySQL entities, Redis support, REST APIs, UML, and integrated module navigation were implemented |
| 2 | Implement registration, login, session validation, profile management, and role-based access control | Achieved | bcrypt password hashing, JWT authentication, Redis active sessions, logout, profile operations, ownership filtering, and administrator middleware were implemented |
| 3 | Implement daily workout recording and workout-performance analysis | Achieved | Daily workout journal, exercise search, CRUD operations, training volume, consistency, progressive-load status, muscle distribution, and warnings were implemented |
| 4 | Implement daily meal recording and nutrition analysis | Achieved | Daily nutrition journal, food search, CRUD operations, calories, macro grams, macro ratios, nutrition gaps, and meal-quality score were implemented |
| 5 | Implement body-progress and calorie-balance analysis | Achieved at prototype level | BMR, estimated TDEE, calorie balance, weight trend, latest weight, and calendar-based seven-day average were implemented as general-wellness estimates |
| 6 | Implement goal management with milestones and feasibility checks | Achieved | Goal lifecycle, priorities, deadlines, custom milestones, generated fallback milestones, feasibility validation, milestone completion, and adherence were implemented |
| 7 | Develop daily, weekly, monthly, and custom dashboard analytics | Achieved | Shared date-range resolver, dashboard controls, aggregated metrics, optional snapshots, and weekly Redis caching were implemented |
| 8 | Develop a structured recommendation workflow with safety, feedback, and provider abstraction | Achieved for workflow; real-AI evaluation pending | Structured context, mock provider, optional OpenAI-compatible provider, rate limiting, basic safety validation, history, and feedback were implemented; no production LLM study was completed |
| 9 | Generate downloadable daily, weekly, monthly, and custom progress reports | Achieved | Deterministic narrative generation, HTML document generation, history, protected download, and deletion were implemented |
| 10 | Provide administration of roles and shared reference data | Achieved | Database-backed role checks, user-role updates, and food and exercise reference CRUD operations were implemented |

The objectives were therefore achieved within the defined prototype scope. The qualifications are important. The system does not claim clinical accuracy, production security certification, proven user behaviour change, or validated production-AI quality.

## 6.4 Contribution of the Project

### 6.4.1 Integrated Personal Health Workflow

The first contribution is the integration of several related health-tracking domains in one application. Workout, nutrition, body, and goal data share one authenticated user model and can be interpreted together through a combined dashboard and report.

This integration allows the user to move beyond isolated logs. For example, calorie intake can be interpreted together with estimated expenditure and body trend. Workout history can be interpreted through consistency, progressive load, and muscle-group distribution. Goal progress can be represented through measurable milestones and adherence.

The project does not attempt to exceed mature commercial applications in database size or specialised coaching depth. Its contribution is a transparent end-to-end prototype showing how these domains can be connected.

### 6.4.2 Transparent Decision-Support Logic

The second contribution is the separation of deterministic analytics from natural-language recommendation.

Core calculations are implemented as application functions rather than delegated to a language model. The same inputs therefore produce repeatable outputs that can be documented and tested. These outputs include:

- Training volume and progressive-load classification.
- Workout consistency and muscle-group distribution.
- Meal nutrient totals and macro ratios.
- BMR, estimated TDEE, and calorie balance.
- Seven-day body-weight average and weight trend.
- Goal milestone adherence.

This approach improves explainability and makes defects easier to identify. Several calculation issues were corrected during development, including sparse body-record averaging, unloaded workout comparison, and long-range meal-quality scoring.

### 6.4.3 Controlled Recommendation Architecture

The third contribution is a recommendation architecture that does not require unrestricted AI generation.

The backend first calculates structured indicators. A provider then converts the summary into readable feedback. The provider can be replaced without changing the surrounding recommendation service, history, feedback, rate limiting, or safety workflow.

The mock provider allows the workflow to be demonstrated without external cost, API availability, or personal-data transfer. The optional OpenAI-compatible provider establishes a technical foundation for the next phase without falsely presenting an untested model as a completed result.

### 6.4.4 Modular Full-Stack Engineering

The project also demonstrates practical full-stack software engineering:

- Feature-based backend modules.
- Controller, service, repository, DTO, and analyzer responsibilities.
- Relational modelling and UUID identifiers.
- REST-style API design.
- Password hashing and JWT authentication.
- Redis-backed active sessions and cache support.
- Role-based authorisation and resource ownership.
- Frontend-backend integration.
- Report generation and protected download.
- Automated calculation tests and production build verification.

The codebase and documentation provide a foundation that can be extended during the next FYP phase.

## 6.5 Lessons Learned

### 6.5.1 Working Interfaces Improve Requirements

Several requirements became clearer only after the application was used with realistic records. A static requirement such as "show workout progress" did not initially define how sparse records, unloaded exercises, or different time ranges should behave.

The iterative process showed that interface implementation is also a form of requirement discovery. Changing workout and nutrition into daily journals, adding search, supporting custom date ranges, and allowing custom milestones all resulted from observing the working system rather than only reviewing an initial specification.

### 6.5.2 Analytical Labels Require Precise Definitions

Terms such as progressive load, meal quality, consistency, and seven-day average can appear self-explanatory but may hide different interpretations.

The project demonstrated the importance of documenting:

- The records included in a calculation.
- The time window used.
- The minimum required data.
- The threshold for each classification.
- The behaviour when data is missing.
- The limitation of the resulting score.

Without these definitions, a dashboard may appear sophisticated while producing confusing or unfair feedback.

### 6.5.3 Caching Requires Invalidation

Adding Redis caching was straightforward compared with ensuring that the cache remained consistent. Workout, meal, body, goal, and profile changes can all affect the weekly summary.

A shared cache-invalidation function was therefore added and called by the relevant write operations. This illustrates that caching should be introduced for a specific need and accompanied by an explicit invalidation strategy.

### 6.5.4 AI Integration Is More Than an API Call

The project also clarified that a useful AI feature requires more than sending a prompt. A responsible workflow needs:

- Controlled input context.
- Deterministic calculations outside the model.
- Provider configuration.
- Timeouts and error handling.
- Rate limiting.
- Safety validation.
- Stored output and context.
- User feedback.
- Evaluation criteria.

Implementing the surrounding architecture before activating a paid model reduced dependency risk and established clearer boundaries for later evaluation.

### 6.5.5 Documentation Must Follow the Implemented System

During development, several UML diagrams became outdated as the implementation changed. Retaining them would have caused contradictions, such as suggesting that reports call an LLM or that every record automatically creates an analytical snapshot.

The diagrams were therefore audited, regenerated, and consolidated. This reinforces that design documentation should be maintained as part of development rather than treated as a one-time proposal artefact.

## 6.6 Current Limitations

### 6.6.1 Limited Independent Evaluation

The strongest current evidence is implementation evidence, automated calculation tests, build verification, and internal workflow testing. A formal usability study with independent participants has not yet been completed.

The report therefore cannot conclude that the system is easy for all target users, improves long-term engagement, or changes health behaviour. These questions require participant tasks, consent, structured measurements, and analysis.

### 6.6.2 Small Reference Datasets

The food and exercise libraries are demonstration datasets. They are sufficient for testing the workflow but do not provide the coverage expected from a commercial nutrition or exercise platform.

A larger dataset introduces additional issues, including licensing, provenance, serving units, duplicate items, regional foods, data quality, indexing, and search relevance.

### 6.6.3 Rule-Based Analytical Simplification

The current indicators are intentionally understandable but simplified.

Training volume does not include perceived exertion, movement velocity, technique, range of motion, or exercise-specific strength differences. Meal quality does not evaluate micronutrients, fibre, food processing, allergies, culture, or clinical requirements. Calorie expenditure is estimated rather than measured. Goal adherence treats milestones equally.

The outputs should therefore be interpreted as decision-support prompts rather than precise physiological conclusions.

### 6.6.4 Mock AI Provider

The default recommendation provider is deterministic and rule-based. It demonstrates the complete provider workflow but not the variability or language capability of a production LLM.

The optional external provider has not been formally evaluated for:

- Recommendation relevance.
- Factual accuracy.
- Safety.
- Consistency across repeated requests.
- Latency.
- Cost.
- Privacy implications.
- Failure and retry behaviour.

### 6.6.5 Basic Safety Validation

The recommendation safety validator currently checks a small list of blocked phrases. Unsafe advice may be expressed without these exact words, while safe educational text might contain a blocked term in another context.

This mechanism is a basic defensive layer and should not be described as comprehensive medical safety validation.

### 6.6.6 Frontend Maintainability

The frontend is currently concentrated in a large `App.tsx` file. The application remains understandable at the current size, but future development should separate domain types, API hooks, routes, pages, forms, and shared components.

The current application also relies on local storage for the bearer token. A production security review should compare this approach with secure HTTP-only cookie sessions and appropriate CSRF controls.

### 6.6.7 Deployment and Operational Limitations

The system currently runs locally. Production requirements such as HTTPS, secret management, environment separation, database backups, durable report storage, observability, audit logging, incident response, and automated deployment have not been completed.

Redis currently supports graceful degradation. This is convenient for development, but the loss of session enforcement and rate limiting may not be acceptable in production.

## 6.7 Future Work

Future work should be prioritised in stages instead of adding every possible technology immediately.

### 6.7.1 Phase 1: Code Quality and Automated Testing

The first priority should be strengthening the existing application:

1. Split `App.tsx` into domain pages and reusable components.
2. Add backend service and HTTP integration tests.
3. Add repository tests using an isolated test database.
4. Add React component tests.
5. Add browser end-to-end tests for the main user workflows.
6. Add security checks for ownership, role access, token expiry, and invalid input.
7. Add continuous integration to run tests and builds on every change.

This work provides more value than introducing additional infrastructure before the current behaviour is comprehensively protected.

### 6.7.2 Phase 2: Formal User and Stakeholder Evaluation

A structured usability evaluation should be conducted with suitable participants. Example tasks include:

- Register and complete the profile.
- Log a workout and meal.
- Edit a historical record.
- Explain the meaning of progressive load and meal quality.
- Create a goal and complete a milestone.
- Generate a recommendation and report.

Evaluation measures can include task completion, completion time, observed errors, perceived clarity, usefulness, and qualitative comments. The findings should be linked to specific interface changes.

Stakeholder collaboration should also be documented through actual requirements or feedback, not only through a collaboration letter.

### 6.7.3 Phase 3: Real AI Integration and Evaluation

After the deterministic analytics and test coverage are stable, the OpenAI-compatible provider can be enabled in a controlled environment.

The evaluation should compare mock and model-generated recommendations using:

- Relevance to the supplied indicators.
- Correct use of goals and trends.
- Absence of unsupported medical claims.
- Actionability.
- Consistency.
- Response time and cost.
- User ratings.

Structured JSON output should be considered instead of unrestricted text. The backend can validate the structure before converting it into user-facing content.

Personal identifiers should be removed from model context, and the user should be informed when external processing occurs.

### 6.7.4 Phase 4: Data and Search Improvements

The food and exercise libraries can be expanded using appropriately licensed and cited data sources. Search can then be improved through:

- Normalised names and aliases.
- Categories and filters.
- Serving-unit conversion.
- Pagination.
- Full-text indexing.
- Frequently used and recent items.

A conventional database search or search engine should be evaluated before introducing a vector database. Vector retrieval may become useful for semantic food queries or knowledge-grounded recommendations, but it is not required for the current structured reference data.

### 6.7.5 Phase 5: Asynchronous AI and Reporting

A message queue may be introduced when recommendation or report jobs become slow, expensive, or require retries. Suitable use cases include:

- Long-running LLM requests.
- Batch weekly report generation.
- Notification delivery.
- File conversion.
- Retry and failure recovery.

The current synchronous mock workflow does not require a queue. Introducing one now would add operational complexity without solving a demonstrated problem.

### 6.7.6 Phase 6: Production Deployment

Production readiness should include:

- Containerised backend and frontend deployment.
- Managed MySQL and Redis.
- HTTPS.
- Secure secret storage.
- Database migration versioning.
- Automated backups.
- Durable report storage.
- Monitoring and structured logging.
- Login rate limiting.
- Security headers and dependency scanning.
- Privacy and data-retention policy.

PDF export, responsive mobile refinement, wearable integration, notifications, and a native mobile client can be considered after the web application is stable and evaluated.

## 6.8 Final Conclusion

The project successfully developed a web-based fitness and nutrition decision-support prototype that integrates personal records and transforms them into structured progress information.

The system addresses the identified problems by placing workout, nutrition, body, and goal data in one platform; applying transparent analytical rules; presenting multiple date ranges; generating contextual recommendations; and producing downloadable progress reports.

Its main strength is not any single formula or interface. The value lies in the complete and explainable workflow:

```text
User record
  -> validated and owned data
  -> deterministic analysis
  -> dashboard, recommendation, and report
  -> user review and feedback
```

The project demonstrates that substantial functionality can be implemented before depending on a production LLM. This decision produced a stable prototype that can be demonstrated, tested, and defended while retaining a clear path to future AI integration.

The current system should not be interpreted as a clinical product or production-ready service. Nevertheless, it provides a strong technical and academic foundation for the next project phase. Future work can now focus on independent evaluation, broader automated testing, real-AI assessment, improved reference data, and secure deployment rather than rebuilding the application foundation.


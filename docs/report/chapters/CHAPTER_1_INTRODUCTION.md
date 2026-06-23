# CHAPTER 1: INTRODUCTION

## 1.1 Background of Study

Physical activity and balanced nutrition are important components of personal health and general well-being. Regular physical activity can contribute to the prevention and management of non-communicable diseases, while also supporting mental health, sleep quality, cognitive health, and healthy body composition. Nevertheless, the World Health Organization reported that approximately 31% of adults worldwide did not meet the recommended physical activity level in 2022 (World Health Organization [WHO], 2024b). This situation highlights the continuing need for practical approaches that encourage individuals to monitor and improve their daily health behaviour.

The widespread use of smartphones, web applications, wearable devices, and digital health platforms has made personal health data easier to record. Users can track information such as workout frequency, exercise duration, food intake, calorie intake, body weight, and personal goals. Digital tracking tools can support behaviour-change techniques such as self-monitoring, feedback, and goal setting. A systematic review by Ferguson et al. (2022) found that wearable activity trackers can improve physical activity and related health outcomes across clinical and non-clinical populations. Similarly, Tong et al. (2024) reported that mobile technologies have potential to promote physical activity and reduce sedentary behaviour, although the effectiveness of these systems also depends on user engagement and implementation quality.

Although digital tracking tools are widely available, recording data alone does not necessarily help users understand their progress. Workout data may be stored separately from nutrition data, while body weight and personal goals may be recorded using other applications or manual notes. This fragmentation makes it difficult for users to obtain a unified view of their behaviour. Even when an application stores multiple types of data, it may present the information mainly as independent logs without explaining the relationships among workout consistency, training load, calorie intake, macronutrient distribution, body-weight trends, and goal achievement.

A personal health decision-support platform should therefore do more than store information. It should transform user records into indicators that are understandable and actionable. For example, workout records can be analysed to estimate training load, consistency, progressive overload, and muscle-group distribution. Meal records can be used to calculate daily calories, protein, carbohydrate, fat, macronutrient ratios, and a basic meal-quality score. Profile and body records can support estimations such as Basal Metabolic Rate (BMR), Total Daily Energy Expenditure (TDEE), calorie balance, and a seven-day moving average of body weight. Goal and milestone records can be used to estimate adherence and provide structured progress feedback.

This project develops an AI-Enhanced Fitness and Nutrition Tracking System as a web-based personal health decision-support platform. It integrates workout, nutrition, body progress, and goal data into one system. The platform provides daily data-recording workflows, date-range analytics, structured recommendations, and downloadable progress reports. The current AI workflow uses a mock recommendation provider for stable and safe development, while the backend includes a provider abstraction that can support an OpenAI-compatible application programming interface (API) in a future phase.

The project is intended for general fitness and nutrition support rather than medical diagnosis or treatment. Its calculations and recommendations are designed to help users interpret their own records and make informed general-wellness decisions. Users with medical conditions or specialised dietary requirements should seek guidance from qualified healthcare professionals.

## 1.2 Problem Statement

An effective personal health tracking platform should allow a user to record relevant health behaviour, understand progress over time, and receive feedback that supports the user's goals. Ideally, workout, nutrition, body progress, and goal data should be integrated so that the system can analyse relationships among these records and present them in a clear form.

However, several problems remain in existing personal tracking practices.

### 1.2.1 Fragmented Fitness and Nutrition Data

Users may use different applications or manual methods to record exercise, meals, body weight, and personal goals. As a result, related information is separated across multiple sources. A user may be able to see a list of completed workouts or consumed meals, but may not be able to review how these records collectively affect calorie balance, body-weight trends, or goal achievement. The absence of an integrated data source reduces the usefulness of historical records and makes long-term progress more difficult to evaluate.

### 1.2.2 Limited Interpretation of Recorded Data

Many tracking workflows focus on data entry and historical lists. Raw records such as sets, repetitions, food quantities, calorie values, and body weights may be available, but users may not know how to interpret them. For example, users may not know whether their workout volume is improving, whether their training is concentrated on a limited number of muscle groups, whether their macronutrient intake is balanced, or whether their recent body-weight change represents a meaningful trend.

### 1.2.3 Generic and Weakly Contextualised Recommendations

General workout and nutrition advice may not consider the user's profile, historical records, current consistency, dietary patterns, or active goals. Advice that is not linked to actual user data may be less relevant and less actionable. A more useful recommendation workflow should first calculate structured indicators and then use these indicators to generate concise feedback related to the user's current situation.

### 1.2.4 Limited Progress Reporting

Users may need to review their behaviour over different time periods rather than only viewing a single daily log. Without daily, weekly, monthly, and custom-range summaries, it is difficult to compare short-term activity with longer-term progress. In addition, the absence of an exportable progress report limits the user's ability to retain, review, or share a structured summary of recorded information.

To address these problems, this project proposes and implements a unified web-based platform that records health-related data and converts the records into analytical indicators. The system combines daily workout and nutrition journals, body-progress tracking, adaptive goal and milestone management, date-range dashboard analytics, structured recommendation generation, and progress-report generation. The resulting platform provides a complete workflow from data entry to interpretation and feedback.

## 1.3 Project Aim

The aim of this project is to design and develop a web-based personal health decision-support platform that integrates fitness and nutrition tracking, analyses user progress, and generates structured recommendations and progress reports based on user profile information, historical records, and personal goals.

## 1.4 Project Objectives

The objectives of this project are:

1. To design a modular full-stack web application that integrates user profiles, workouts, nutrition records, body-progress records, goals, analytics, recommendations, and reports.
2. To implement secure user registration, login, session validation, profile management, and role-based access control.
3. To implement daily workout recording and workout-performance analysis, including training-load calculation, workout consistency, progressive-overload status, and muscle-group distribution.
4. To implement daily meal recording and nutrition analysis, including calorie and macronutrient calculation, macronutrient ratios, nutrition-gap detection, and meal-quality scoring.
5. To implement body-progress and calorie-balance analysis using BMR, estimated TDEE, calorie intake, body-weight trends, and a seven-day body-weight moving average.
6. To implement goal management with target values, deadlines, priorities, custom milestones, feasibility checks, status management, and milestone-adherence analysis.
7. To develop dashboard analytics that support daily, weekly, monthly, and custom date ranges.
8. To develop a structured recommendation workflow that uses calculated user indicators, safety validation, user feedback, and an AI-provider abstraction.
9. To generate downloadable daily, weekly, monthly, and custom-range progress reports.
10. To provide administrative functions for managing user roles, exercise reference data, and food reference data.

## 1.5 Project Scope

### 1.5.1 Target Users

The system has two main user roles:

**General user.** The primary user is an adult who wants to record and review general fitness and nutrition information. A general user can register an account, configure a health profile, record daily workouts and meals, record body weight, manage goals and milestones, view dashboard analytics, request recommendations, submit recommendation feedback, and generate progress reports.

**Administrator.** The administrator maintains selected system reference data and user access. An administrator can manage user roles, exercise reference records, and food reference records. These functions support data consistency and prevent normal users from directly modifying shared reference data.

### 1.5.2 Functional Scope

The functional scope consists of the following modules:

| Module | Main Functions |
|---|---|
| Authentication and profile | Registration, login, logout, JWT authentication, session validation, profile creation and update |
| Workout tracking | Daily workout journal, exercise search, sets, repetitions, weight, rest time, duration, notes, editing and deletion |
| Workout analysis | Training load, consistency, progressive-overload status, muscle-group distribution, and training warnings |
| Nutrition tracking | Daily meal journal, food search, quantity, meal type, meal time, editing and deletion |
| Nutrition analysis | Daily calories and macronutrients, macronutrient ratios, nutrition gaps, and meal-quality score |
| Body progress | Date-based body-weight records, body timeline, latest weight, weight trend, and seven-day moving average |
| Goal management | Goal creation, target metric, target value, deadline, priority, status, custom milestones, and milestone completion |
| Dashboard analytics | Daily, weekly, monthly, and custom-range aggregation of workout, nutrition, body, and goal indicators |
| Recommendation | Mock recommendation generation, structured context, safety validation, recommendation history, and user feedback |
| Reporting | Daily, weekly, monthly, and custom-range report generation, history, deletion, and HTML download |
| Administration | User-role management and CRUD management of shared food and exercise reference data |

### 1.5.3 Technical Scope

The frontend is implemented as a React and TypeScript web application using Vite. It communicates with the backend through RESTful HTTP APIs. The backend is implemented in Go using the Gin web framework and GORM for database access. MySQL stores relational application data, while Redis supports logout-aware session validation, dashboard caching, and recommendation rate limiting. Docker Compose is used to run MySQL and Redis in the local development environment.

The system currently operates as a local web application. It has been designed so that the frontend, backend, database, cache, and AI provider can later be deployed as separate services. UUID values are used as primary identifiers to provide consistent identifiers across modules and to avoid exposing predictable sequential resource identifiers through public APIs.

### 1.5.4 Scope Boundaries

The following items are outside the current minimum project scope:

- Medical diagnosis, clinical treatment, or emergency-health guidance
- Replacement of advice from doctors, dietitians, or other qualified professionals
- Direct integration with wearable devices or smart scales
- Automatic import from commercial fitness or nutrition platforms
- Image-based food recognition
- A comprehensive clinical-grade food-composition database
- Real-time integration with a production large language model in the current phase
- Native Android or iOS applications
- Production cloud deployment
- PDF report export

The current report output is a downloadable HTML document. The recommendation workflow uses a mock provider to demonstrate the complete application flow without depending on an external API key. Real AI integration, PDF export, mobile applications, and cloud deployment are identified as future enhancements.

## 1.6 Significance of Project

The project is significant from user, technical, and academic perspectives.

From the user's perspective, the platform centralises multiple types of personal health data. Instead of reviewing isolated workout, meal, weight, and goal records, the user can review a combined dashboard and progress report. Daily journals simplify data entry, while date-range analytics allow users to inspect both short-term behaviour and longer-term progress. The calculated indicators also reduce the effort required to manually interpret raw records.

From a decision-support perspective, the system demonstrates how deterministic analytics and natural-language recommendations can be combined. Calculations such as training load, macronutrient ratios, BMR, TDEE, calorie balance, seven-day moving average, and goal adherence are handled by controlled application logic. These indicators can then be supplied to a recommendation provider as structured context. This approach is more transparent than sending unrestricted raw user data directly to a language model.

From a software-engineering perspective, the project demonstrates modular full-stack development. The backend separates major business domains into independent modules and applies controller, service, repository, data-transfer object, and analysis responsibilities where appropriate. The system also demonstrates RESTful API design, relational database modelling, JWT authentication, role-based authorisation, Redis-backed session and cache support, report generation, and frontend-backend integration.

From an academic perspective, the project provides a practical environment for evaluating personal health analytics, usability, recommendation relevance, and the limitations of AI-assisted general-wellness systems. It also provides a foundation for future work involving real AI integration, improved visualisation, larger food datasets, automated testing, and deployment.

## 1.7 Project Schedule

The project followed an iterative and incremental development process. Requirements, architecture, backend modules, frontend interfaces, and integration tests were developed in stages. Table 1.1 summarises the major project activities.

**Table 1.1: Project Schedule**

| Phase | Main Activities | Expected Output | Status |
|---|---|---|---|
| Project initiation | Topic refinement, collaborator identification, collaboration-letter preparation, and proposal planning | Confirmed project direction and initial scope | Completed |
| Requirement analysis | Review project expectations, identify target users, define core and supporting modules, and refine supervisor feedback | Revised proposal and module requirements | Completed |
| Literature and existing-system review | Review digital health tracking, workout and nutrition systems, analytics methods, and AI-assisted recommendation approaches | Literature-review evidence and identified gaps | In progress |
| System analysis and design | Define architecture, database model, use cases, module responsibilities, API boundaries, and interface design | UML diagrams, ERD, architecture, and requirements | Completed with ongoing refinement |
| Backend development | Implement authentication, profile, workout, nutrition, body, goal, analytics, recommendation, report, and admin modules | Functional Go REST API | Completed for the current phase |
| Frontend development | Implement authentication and module interfaces, daily journals, dashboard, recommendation, report, and admin pages | Functional React web interface | Completed for the current phase |
| Integration and testing | Connect frontend and backend, verify MySQL and Redis, perform unit tests, API smoke tests, and workflow testing | Integrated and locally testable system | Completed with further testing planned |
| Monitoring and viva preparation | Prepare monitoring slides, gather supervisor feedback, prepare viva slides, record presentations, and demonstrate the system | Monitoring and viva deliverables | Completed |
| Report preparation | Consolidate literature, requirements, UML diagrams, implementation evidence, testing results, and limitations | FYP 1 report | In progress |

The detailed Gantt chart will be inserted as **Figure 1.1: Project Gantt Chart** after the final academic-week dates are confirmed.

## 1.8 Structure of Report

This report is organised into six chapters.

**Chapter 1: Introduction** introduces the project background, problem statement, aim, objectives, scope, significance, and project schedule.

**Chapter 2: Literature Review** reviews recent academic research related to digital health tracking, fitness and nutrition monitoring, behaviour change, health analytics, goal management, and AI-assisted recommendations. It also compares relevant existing systems and identifies the research and system gaps addressed by this project.

**Chapter 3: Methodology** explains the iterative and incremental development approach, requirement-gathering activities, supervisor and stakeholder collaboration, development phases, testing approach, and ethical considerations.

**Chapter 4: System Analysis and Design** presents the functional and non-functional requirements, system modules, use cases, architecture, class model, database design, activity diagrams, sequence diagrams, security design, analytics design, and interface design.

**Chapter 5: Technical Implementation** describes the implemented frontend, backend, database, authentication, analytics, recommendation, reporting, administration, integration, and testing components. Real system screenshots and selected implementation evidence are presented in this chapter.

**Chapter 6: Conclusion** summarises the project outcomes, evaluates the achievement of the objectives, identifies current limitations, and proposes future improvements.

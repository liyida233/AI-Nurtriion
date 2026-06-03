# Final Year Project Proposal

# AI-Enhanced Fitness and Nutrition Tracking System

**Subtitle:** A Web-Based Personal Health Decision-Support Platform  
**Student Name:** Li AO  
**Programme:** Final Year Project  
**Project Type:** Full-stack Web Application with AI-Assisted Recommendation  
**Proposed Backend:** Go, MySQL, Redis  
**Revision Focus:** Core modules enhanced beyond basic CRUD and dashboard functions  

**Revised Version**  
This version keeps user authentication as a supporting foundation and strengthens Modules 2–7 with analytical, adaptive, and AI-assisted decision-support functions.

---

## 1. Project Overview

**Project Title:** AI-Enhanced Fitness and Nutrition Tracking System  
**Subtitle:** A Web-Based Personal Health Decision-Support Platform

This project proposes the design and development of an AI-enhanced fitness and nutrition tracking system that helps users manage workout activities, dietary intake, calorie balance, fitness goals, and long-term progress through a unified web-based platform.

The revised scope moves the project beyond a basic tracking website by introducing rule-based analysis, progress prediction, adaptive goal monitoring, and AI-assisted personalized recommendations.

The system will be developed as a modular full-stack web application. The user login and profile setup component will remain necessary for system operation, but it will be treated as a supporting foundation rather than a major functional contribution. The core contribution of the project will be concentrated in Modules 2–7, which focus on workout performance analysis, dietary pattern analysis, calorie and body progress analysis, AI recommendation, adaptive goal management, and automated reporting.

---

## 2. Problem Statement

Many existing fitness applications allow users to record workouts or meals, but the recorded data is often presented as simple logs without deeper interpretation.

Users may struggle to understand whether their workout frequency, calorie balance, macronutrient intake, and body weight trends are aligned with their goals.

Generic workout or meal suggestions are often not personalized using the user's historical records, consistency, progress, and goal status.

A more useful system should convert daily fitness and nutrition data into actionable insights, adaptive feedback, and structured progress reports.

---

## 3. Aim and Objectives

### Aim

The aim of this project is to develop a web-based personal health decision-support platform that records fitness and nutrition data, analyzes user progress, and provides AI-assisted recommendations based on user goals and historical behaviour.

### Objectives

1. To design a modular full-stack web application for managing workout, nutrition, calorie, goal, and progress data.
2. To implement workout performance analysis, including training volume calculation, workout consistency tracking, and progressive overload detection.
3. To implement nutrition analysis, including calorie and macronutrient calculation, dietary pattern detection, and nutrition gap identification.
4. To calculate calorie balance and body progress trends using BMR, estimated TDEE, exercise expenditure, and time-series progress data.
5. To integrate an AI-assisted recommendation module that generates personalized workout plans, meal suggestions, and progress feedback using structured user data.
6. To develop adaptive goal monitoring that evaluates goal feasibility, adherence, milestone completion, and potential goal adjustment.
7. To generate dashboards and automated reports that summarize weekly and monthly user progress in a clear and actionable format.

---

## 4. Revised Module Structure

**Key revision:** The authentication-related module is retained as a supporting module, while Modules 2–7 are strengthened as the main FYP contribution. The emphasis is shifted from simple create-read-update-delete operations to analysis, evaluation, recommendation, prediction, adaptation, and reporting.

| No. | Module | Status | Main Technical Depth |
|---:|---|---|---|
| 1 | System Foundation: Authentication and Profile Setup | Supporting | Secure access, profile setup, health profile, role-based access |
| 2 | Workout Tracking and Performance Analysis | Core | Training volume, consistency score, progressive overload, muscle group balance |
| 3 | Nutrition Logging and Dietary Pattern Analysis | Core | Food database, macro analysis, meal quality score, nutrition gap detection |
| 4 | Calorie Balance and Body Progress Analysis | Core | BMR, TDEE, calorie balance, moving average, progress prediction |
| 5 | AI-Based Personalized Recommendation | Core | Structured prompt generation, rule-based preprocessing, LLM feedback, feedback loop |
| 6 | Adaptive Goal Setting and Achievement Monitoring | Core | SMART goals, feasibility check, adherence score, dynamic adjustment |
| 7 | Progress Analytics, Dashboard, and Automated Reporting | Core | Comparative analytics, visual dashboard, weekly/monthly reports, export function |

---

## 5. Detailed Module Requirements

### 5.1 Supporting Module 1: System Foundation — Authentication and Profile Setup

This module supports the operation of the application but is not treated as a core FYP contribution. It provides secure access control and collects basic user data required by the analytical modules.

#### Main Functions

- User registration and login
- Password hashing and session/JWT-based authentication
- User profile management
- Health profile setup such as age, gender, height, weight, activity level, and fitness goals
- Role-based access for normal users and administrators

#### Technical Complexity

- Authentication flow
- User-profile schema design
- Access-control middleware
- Validation of profile data required by analytics modules

---

### 5.2 Core Module 2: Workout Tracking and Performance Analysis Module

This module allows users to record workout activities and converts workout logs into measurable performance indicators. It is designed to support both activity recording and training progress interpretation.

#### Main Functions

- Record exercise type, duration, sets, repetitions, weights, rest time, and workout category
- Maintain an exercise database with category, target muscle group, equipment, and intensity level
- Calculate training volume using formulas such as `sets × reps × weight`
- Track weekly workout frequency and planned-versus-completed workout adherence
- Detect progressive overload by comparing historical changes in weight, repetitions, sets, or total volume
- Identify simple workout imbalance patterns, such as repeated upper-body training with low lower-body activity

#### Technical Complexity

- Time-series workout history analysis
- Rule-based progressive overload detection
- Workout consistency scoring
- Muscle-group distribution analysis

---

### 5.3 Core Module 3: Nutrition Logging and Dietary Pattern Analysis Module

This module enables users to log meals and evaluates their dietary patterns against personal calorie and macronutrient targets. It is not limited to food entry; it includes nutrient computation and dietary quality analysis.

#### Main Functions

- Log meals such as breakfast, lunch, dinner, and snacks
- Record food name, quantity, serving size, meal time, and nutrition values
- Maintain a food database containing calories, protein, carbohydrate, fat, and optional sugar or sodium values
- Calculate daily total calories and macronutrient distribution
- Compare actual intake against targets based on user goals
- Detect nutrition gaps such as low protein intake, repeated calorie excess, or unbalanced meal patterns
- Support reusable meal templates for frequently consumed meals

#### Technical Complexity

- Nutrition data modelling
- Macro-ratio calculation
- Target-versus-actual comparison
- Meal quality scoring
- Pattern detection across daily and weekly logs

---

### 5.4 Core Module 4: Calorie Balance and Body Progress Analysis Module

This module analyzes the relationship between calorie intake, calorie expenditure, and body progress. It provides formula-based calculations and trend analysis to help users understand whether their current routine aligns with their fitness goals.

#### Main Functions

- Estimate Basal Metabolic Rate based on profile data
- Estimate Total Daily Energy Expenditure using activity level and exercise records
- Calculate calorie balance as intake minus expenditure
- Classify days as calorie deficit, maintenance, or surplus
- Track body weight and body progress records over time
- Generate 7-day moving average for body weight trends
- Predict short-term progress direction based on recent calorie balance and weight trend
- Detect inconsistent logging or abnormal progress patterns

#### Technical Complexity

- Formula-based calorie computation
- BMR and TDEE estimation
- Time-series smoothing
- Trend classification
- Basic progress prediction and anomaly detection

---

### 5.5 Core Module 5: AI-Based Personalized Recommendation Module

This module is the intelligent decision-support component of the system. The AI component will not simply generate generic fitness advice. Instead, the system will first process structured user data and then pass concise, controlled context to an LLM/API-based recommendation engine.

#### Main Functions

- Generate personalized workout plans based on goal, availability, workout history, and weak muscle groups
- Suggest meal ideas based on calorie targets, macronutrient gaps, and dietary history
- Generate weekly natural-language feedback from workout, nutrition, calorie, and goal data
- Use rule-based preprocessing to summarize user state before sending it to the AI model
- Apply prompt templates to ensure recommendations are structured, explainable, and within safe general wellness boundaries
- Collect user feedback such as useful, too difficult, too easy, or not suitable, and use it to refine future recommendations

#### Technical Complexity

- Structured prompt generation
- Hybrid rule-based and LLM-based recommendation pipeline
- Explainable recommendation output
- Feedback-loop design
- Safety filtering to avoid medical diagnosis or unsafe diet advice

---

### 5.6 Core Module 6: Adaptive Goal Setting and Achievement Monitoring Module

This module helps users define, evaluate, and adjust fitness or nutrition goals. It extends basic goal management by adding feasibility checking, adherence scoring, milestone tracking, and adaptive suggestions.

#### Main Functions

- Set SMART goals with target value, metric, deadline, priority, and progress status
- Support goals such as workout frequency, calorie target, protein target, weight trend, or consistency target
- Evaluate basic feasibility of goals, such as overly aggressive short-term weight changes or unrealistic workout schedules
- Calculate goal adherence score using workout completion, nutrition target achievement, and progress trend
- Break long-term goals into weekly or monthly milestones
- Suggest goal adjustment when repeated non-completion is detected

#### Technical Complexity

- Goal modelling
- Feasibility rules
- Adherence scoring
- Milestone tracking
- Adaptive goal adjustment logic

---

### 5.7 Core Module 7: Progress Analytics, Dashboard, and Automated Reporting Module

This module presents analyzed information in a structured dashboard and generates progress reports. It focuses on turning stored data and calculated indicators into readable summaries for users.

#### Main Functions

- Daily, weekly, and monthly dashboard views
- Visualize workout frequency, training volume, calorie intake, calorie expenditure, calorie balance, weight trend, and goal completion
- Compare current week against previous week and current month against previous month
- Show planned-versus-actual performance for workouts and nutrition targets
- Generate automated weekly or monthly progress reports
- Include AI-generated summaries and recommendations in reports
- Export reports as PDF for review or record keeping

#### Technical Complexity

- Aggregation queries
- Comparative analytics
- Data visualization
- Automated report generation
- AI-generated report summaries
- PDF export workflow

---

## 6. AI Recommendation Workflow

The AI module will follow a controlled pipeline rather than relying on unrestricted text generation. The system will first compute structured indicators from workout, nutrition, calorie, and goal records. These indicators will then be summarized into a prompt template used by the AI model to generate recommendations and explanations.

| Step | Input / Process | Output |
|---:|---|---|
| 1. Data Collection | Workout logs, meal logs, profile, body records, goals | Structured user dataset |
| 2. Rule-Based Analysis | Calculate volume, macro ratio, calorie balance, adherence, trend | Analytical indicators |
| 3. Context Builder | Convert indicators into compact JSON-like context | Prompt-ready user state |
| 4. AI Generation | LLM/API generates plan, feedback, and explanation | Workout plan, meal suggestion, summary |
| 5. Safety and Validation | Check against general wellness boundaries and formatting rules | Filtered recommendation |
| 6. Feedback Loop | User rates recommendation suitability | Improved future recommendation context |

---

## 7. Proposed System Architecture

**Architecture style:** The system will follow a modular full-stack architecture with a web-based frontend, Go backend services, relational database storage, Redis-based caching or session support, and an external or local AI service integration layer.

| Layer | Responsibility |
|---|---|
| Frontend Dashboard | User interface for workout logging, nutrition logging, goal setup, analytics dashboard, and report viewing. |
| Backend API Layer | REST API implemented in Go for authentication, module logic, validation, and data processing. |
| Business Logic Layer | Workout analysis, nutrition analysis, calorie computation, goal scoring, recommendation preparation, and report generation. |
| Database Layer | MySQL tables for users, profiles, workouts, exercises, meals, food items, goals, body records, reports, and recommendation feedback. |
| Cache / Performance Layer | Redis for session storage, frequently accessed dashboard summaries, or rate-limiting support. |
| AI Integration Layer | LLM/API service for structured recommendation generation, progress explanation, and report summary generation. |

---

## 8. Preliminary Data Model

The database design will support both operational records and analytical processing. The following entities are expected:

| Entity | Description |
|---|---|
| User | Account information, role, authentication metadata |
| UserProfile | Age, gender, height, weight, activity level, primary fitness goal |
| Exercise | Exercise name, category, muscle group, equipment, intensity level |
| WorkoutLog | Workout date, exercise, sets, reps, weight, duration, notes |
| FoodItem | Food name, serving size, calories, protein, carbohydrates, fat, optional sugar or sodium |
| MealLog | Meal type, food item, quantity, meal time, calculated nutrition values |
| BodyRecord | Date, body weight, optional body measurement notes |
| Goal | Goal type, target metric, target value, deadline, status, priority |
| AnalyticsSnapshot | Weekly or monthly computed indicators for dashboard and reports |
| RecommendationFeedback | AI recommendation ID, user rating, suitability feedback, adjustment notes |

---

## 9. Proposed Technologies

| Component | Proposed Technology |
|---|---|
| Backend | Go with a web framework such as Gin or Fiber |
| Database | MySQL for relational data storage |
| Cache / Session / Rate Support | Redis |
| Frontend | Web-based dashboard using a modern frontend framework or server-rendered frontend |
| AI Integration | LLM/API-based recommendation and natural-language summary generation |
| Reporting | Backend-generated PDF reports or frontend-exported progress reports |
| Visualization | Charts for weekly/monthly trends, goal completion, workout volume, and calorie balance |

---

## 10. Development Methodology

1. **Requirement refinement:** Confirm core functional scope, supporting scope, and success criteria for each module.
2. **System design:** Define architecture, database schema, API endpoints, and module boundaries.
3. **Implementation phase 1:** Build authentication, profile setup, and core CRUD operations for workout, nutrition, and goals.
4. **Implementation phase 2:** Add analytical logic for workout performance, nutrition patterns, calorie balance, progress trends, and goal adherence.
5. **Implementation phase 3:** Integrate AI recommendation workflow using structured context and prompt templates.
6. **Implementation phase 4:** Build dashboard, reporting, export, and administrative views where applicable.
7. **Testing and evaluation:** Test correctness of calculations, API responses, user flows, dashboard accuracy, and recommendation relevance.

---

## 11. Evaluation Plan

The project can be evaluated using functional, analytical, usability, and recommendation-quality criteria.

- **Functional testing:** Verify that users can log workouts, meals, body records, and goals correctly.
- **Calculation testing:** Validate workout volume, macro totals, calorie balance, BMR/TDEE estimates, trend calculations, and adherence scores using controlled test cases.
- **Dashboard testing:** Compare displayed summaries against database records and expected analytical values.
- **AI output review:** Evaluate whether AI recommendations correctly use structured user context and provide relevant explanations.
- **Usability testing:** Collect feedback from sample users on clarity of dashboards, report usefulness, and ease of logging.
- **Performance testing:** Test common dashboard and reporting queries, with optional Redis caching for frequently accessed summaries.

---

## 12. Expected Deliverables

- A functional web application with user profile setup, workout logging, nutrition logging, calorie analysis, goal management, AI recommendation, dashboard, and reporting features.
- A structured backend implemented in Go with modular service boundaries.
- A MySQL database schema supporting operational data and analytical summaries.
- Redis support for selected caching, session, or performance-related use cases.
- AI-assisted personalized recommendation workflow with structured prompts and feedback collection.
- Dashboard visualizations and exportable progress reports.
- Project documentation including system design, implementation details, testing results, and limitations.

---

## 13. Scope Boundaries and Safety Considerations

- The system provides general fitness and nutrition support; it is not a medical diagnosis or treatment system.
- Recommendations should avoid unsafe or extreme dieting, excessive exercise, or medical claims.
- The AI model will be used as a recommendation and explanation assistant, while calculations and user data processing remain controlled by the application logic.
- The project focuses on personal tracking and decision support; integration with wearable devices is optional and outside the minimum scope.
- Food nutrition values may depend on the available food database and should be treated as estimates rather than exact medical measurements.

---

## 14. Main Project Contribution

The main contribution of this project is not merely to build a fitness tracking website. The contribution is to design a modular personal health decision-support platform that transforms user logs into structured insights.

The strengthened modules provide analytical depth through workout performance analysis, dietary pattern evaluation, calorie balance computation, body progress trend analysis, adaptive goal monitoring, and AI-assisted personalized recommendations.

The revised project therefore demonstrates technical depth in full-stack development, database design, data analysis, rule-based reasoning, AI integration, reporting, and user-centred dashboard design.

---

## 15. Conclusion

This revised proposal strengthens the original project scope by clearly separating the supporting authentication module from the core FYP modules. Modules 2–7 are enhanced with analytical, adaptive, and AI-assisted functions so that the system is positioned as an AI-assisted personal fitness and nutrition decision-support platform rather than a basic logging application.

The final system is expected to provide practical user value while demonstrating sufficient technical complexity for an individual final year project.

---

## Appendix A: Response to Supervisor Feedback

| Feedback Issue | Revision Applied | Resulting Improvement |
|---|---|---|
| Login module is not counted as a valid core module. | Renamed it as a supporting system foundation and removed it from the main contribution. | The proposal now focuses evaluation on Modules 2–7. |
| Modules 2–7 need more complexity. | Added analysis, scoring, prediction, adaptation, and AI-assisted decision-support features. | The modules now show stronger technical depth beyond CRUD. |
| AI feature may look too generic. | Defined a structured AI workflow with rule-based preprocessing, prompt templates, safety filtering, and user feedback loop. | The AI component becomes a controlled system feature rather than simple API usage. |
| Dashboard may be too simple. | Expanded dashboard into comparative analytics and automated report generation. | The reporting module now involves aggregation, visualization, and exportable reports. |

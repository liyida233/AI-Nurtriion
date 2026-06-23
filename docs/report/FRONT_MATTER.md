# TITLE PAGE

**FACULTY OF COMPUTER SCIENCE AND INFORMATION TECHNOLOGY**  
**UNIVERSITI MALAYA**

**WIA3002/WIB3002: ACADEMIC PROJECT I**

# AI-ENHANCED FITNESS AND NUTRITION TRACKING SYSTEM

**Prepared by**

**LI AO**

Matric Number: **[INSERT MATRIC NUMBER]**

Supervisor: **[INSERT SUPERVISOR NAME]**

Session: **2025/2026**

Submission Date: **[INSERT FINAL SUBMISSION DATE]**

---

# ACKNOWLEDGEMENT

I would like to express my sincere appreciation to my supervisor, **[insert supervisor name]**, for the guidance, feedback, and encouragement provided throughout this project. The discussions concerning the proposal, project requirements, assessment expectations, and system direction helped me refine both the technical scope and presentation of the work.

I am also grateful to the Faculty of Computer Science and Information Technology, Universiti Malaya, for providing the academic environment and resources required to complete this project. I would like to thank the individuals who provided comments on the project idea, monitoring presentation, and prototype. Their observations contributed to the improvement of the system requirements and user workflows.

Finally, I would like to thank my family and friends for their patience and support during the development, testing, documentation, and presentation of this project.

---

# ABSTRACT

Personal fitness information is often divided among separate workout, nutrition, body-weight, and goal-tracking tools, while recorded data may be presented without sufficient interpretation. This project designed and developed an AI-Enhanced Fitness and Nutrition Tracking System as a web-based personal health decision-support platform. The system integrates daily workout and meal journals, body-progress records, goals and milestones, date-range analytics, recommendations, reports, and administrative reference-data management.

An iterative and incremental methodology was used to refine requirements through literature review, existing-system analysis, supervisor feedback, prototype observation, implementation, and testing. The frontend was developed using React and TypeScript, while the backend used Go, Gin, and GORM. MySQL stores relational application data, and Redis supports active sessions, weekly dashboard caching, cache invalidation, and recommendation rate limiting.

The implemented analytics calculate workout volume, consistency, progressive-load status, muscle-group distribution, calorie and macronutrient totals, meal-quality indicators, BMR, estimated TDEE, calorie balance, body-weight trends, a seven-day moving average, and milestone adherence. These deterministic indicators are shared by the dashboard, recommendation, and reporting workflows. The current recommendation workflow uses a mock provider with safety validation and feedback collection, while an optional OpenAI-compatible provider has been prepared for future evaluation. Progress reports are generated as downloadable HTML documents.

Backend tests and the frontend production build passed during the latest verification. The result is a functional and explainable prototype that demonstrates integrated personal health tracking and structured decision support. Future work includes formal usability testing, broader automated testing, real-AI evaluation, expanded reference datasets, and secure production deployment.

**Keywords:** fitness tracking, nutrition tracking, personal health analytics, decision support, web application, artificial intelligence

---

# TABLE OF CONTENTS

Use **Insert > Table of contents** in Google Docs after applying:

- Heading 1 to chapter titles
- Heading 2 to numbered main sections
- Heading 3 to numbered subsections

The main document order should be:

1. Acknowledgement
2. Abstract
3. Table of Contents
4. List of Figures
5. List of Tables
6. Chapter 1: Introduction
7. Chapter 2: Literature Review
8. Chapter 3: Methodology
9. Chapter 4: System Analysis and Design
10. Chapter 5: Technical Implementation
11. Chapter 6: Conclusion and Future Work
12. References
13. Appendices

---

# LIST OF FIGURES

| Figure | Title | Page |
|---|---|---|
| Figure 3.1 | Iterative and Incremental Development Workflow | ___ |
| Figure 4.1 | System Use-Case Diagram | ___ |
| Figure 4.2 | System Component Diagram | ___ |
| Figure 4.3 | Core Domain Class Diagram | ___ |
| Figure 4.4 | Entity-Relationship Diagram | ___ |
| Figure 4.5 | Dashboard Analytics Activity Diagram | ___ |
| Figure 4.6 | Goal Creation and Milestone Management Sequence | ___ |
| Figure 4.7 | Recommendation Generation and Feedback Sequence | ___ |
| Figure 4.8 | Progress Report Generation and Download Sequence | ___ |
| Figure 4.9 | Interface Navigation Structure | ___ |
| Figure 4.10 | Implemented Dashboard Interface Design | ___ |
| Figure 4.11 | Implemented Daily Journal Design | ___ |
| Figure 5.1 | Login Interface | ___ |
| Figure 5.2 | Implemented Dashboard | ___ |
| Figure 5.3 | Workout Journal and Exercise Search | ___ |
| Figure 5.4 | Nutrition Journal and Food Search | ___ |
| Figure 5.5 | Body-Progress Timeline | ___ |
| Figure 5.6 | Goal and Custom Milestone Interface | ___ |
| Figure 5.7 | Recommendation Interface | ___ |
| Figure 5.8 | Report Generation and History | ___ |
| Figure 5.9 | Administration Interface | ___ |

Do not add Figure 1.1 unless a final Gantt chart is inserted into Chapter 1.

---

# LIST OF TABLES

| Table | Title | Page |
|---|---|---|
| Table 1.1 | Project Schedule | ___ |
| Table 2.1 | Comparison of Existing Systems and the Proposed System | ___ |
| Table 3.1 | Requirement-Gathering Sources | ___ |
| Table 3.2 | Development Phases and Outputs | ___ |
| Table 3.3 | Development Technologies | ___ |
| Table 3.4 | Testing Matrix | ___ |
| Table 3.5 | Project Risks and Mitigation | ___ |
| Table 4.1 | Functional Requirements | ___ |
| Table 4.2 | Non-Functional Requirements | ___ |
| Table 4.3 | Backend Module Responsibilities | ___ |
| Table 4.4 | Main API Groups | ___ |
| Table 4.5 | Interface Design by Screen | ___ |
| Table 4.6 | Requirement Traceability Matrix | ___ |
| Table 5.1 | Backend File Responsibilities | ___ |
| Table 5.2 | Latest Verification Results | ___ |
| Table 5.3 | Main Implementation Challenges | ___ |
| Table 6.1 | Achievement of Project Objectives | ___ |


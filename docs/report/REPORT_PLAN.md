# FYP 1 Report Working Plan

## Report Structure

1. Front Matter
   - Title Page
   - Acknowledgement
   - Abstract (maximum 300 words)
   - Table of Contents
   - List of Figures
   - List of Tables
2. Chapter 1: Introduction
3. Chapter 2: Literature Review
4. Chapter 3: Methodology
5. Chapter 4: System Analysis and Design
6. Chapter 5: Technical Implementation
7. Chapter 6: Conclusion
8. References
9. Appendices

## UML Asset Inventory

### Suitable after a small consistency review

- `01_use_case_diagram`
- `03_class_diagram`
- `05_ai_recommendation_sequence`
- `06_report_generation_sequence`
- `07_goal_state_diagram`
- `09_user_login_sequence`
- `10_user_registration_activity`
- `11_profile_update_sequence`
- `12_workout_logging_sequence`
- `13_nutrition_logging_sequence`
- `14_body_record_sequence`
- `17_admin_reference_data_sequence`
- `18_generic_data_recording_activity`

### Must be updated before insertion

- `01_use_case_diagram`: fix `UC141`; remove unsupported admin analytics use case or align it with implementation.
- `02_component_diagram`: remove direct Report-to-LLM flow; align Redis usage and current mock/provider architecture.
- `03_class_diagram`: add `GoalMilestone` and current nutrition/report fields.
- `08_erd_bonus`: add `GoalMilestone` and align entity fields with the implemented schema.
- `15_goal_management_sequence`: show user-defined milestones with default milestone generation as fallback.
- `16_dashboard_activity_diagram`: support daily, weekly, monthly, and custom date ranges; remove unsupported previous-period comparison.

## Recommended Figures by Chapter

### Chapter 3: Methodology

- Iterative and incremental development workflow
- Project schedule or Gantt chart

### Chapter 4: System Analysis and Design

- Use case diagram
- Component/system architecture diagram
- Class diagram
- ERD
- Dashboard activity diagram
- AI recommendation sequence diagram
- Report generation sequence diagram
- Goal management sequence diagram

The remaining module sequence diagrams may be placed in the appendix to avoid making Chapter 4 unnecessarily long.

### Chapter 5: Technical Implementation

- Login and registration interface
- Dashboard with date-range selection
- Daily workout journal and exercise search
- Daily nutrition journal and food search
- Body timeline
- Goal and custom milestone interface
- Recommendation interface
- Report generation and download interface
- Admin management interface
- Selected backend folder structure or API evidence

## Writing Order

1. Chapter 1: Introduction
2. Chapter 3: Methodology
3. Chapter 4: Requirements and architecture
4. Chapter 2: Literature Review
5. Chapter 5: Technical Implementation
6. Abstract and conclusion

## Evidence Rules

- Use the current implemented system as the source of truth.
- Update diagrams that no longer match the implementation.
- Use screenshots from the running application, not old proposal mock-ups.
- Add a figure number, caption, and in-text explanation for every figure.
- Do not include large code dumps; use short selected snippets only when they explain a design decision.
- Use recent academic references where possible and format references in APA style.

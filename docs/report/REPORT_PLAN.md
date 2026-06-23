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

The UML directory has been consolidated. All remaining diagrams match the
current implementation and have corresponding PNG and SVG renders.

| Diagram | Main Report Use |
|---|---|
| `01_use_case_diagram` | Chapter 4 system use cases |
| `02_component_diagram` | Chapter 4 system architecture |
| `03_class_diagram` | Chapter 4 domain model |
| `05_ai_recommendation_sequence` | Chapter 4 recommendation workflow |
| `06_report_generation_sequence` | Chapter 4 report workflow |
| `08_erd_bonus` | Chapter 4 relational database design |
| `15_goal_management_sequence` | Chapter 4 goal and milestone workflow |
| `16_dashboard_activity_diagram` | Chapter 4 dashboard analytics workflow |
| `19_development_methodology` | Chapter 3 development methodology |
| `20_interface_navigation` | Chapter 4 interface navigation |

Obsolete diagrams were removed after their workflows were replaced or found
to conflict with the current implementation. The removed diagrams should not
be used in the report or appendix.

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

Additional module screenshots and API evidence may be placed in the appendix.

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

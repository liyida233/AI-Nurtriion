# CHAPTER 2: LITERATURE REVIEW

## 2.1 Introduction

This chapter reviews literature related to digital health self-monitoring, fitness tracking, nutrition tracking, body-progress analysis, digital goal setting, and AI-assisted health recommendations. It also reviews four existing products: MyFitnessPal, Cronometer, Strong, and Fitbod. The purpose of the review is to identify useful theories, methods, system features, and limitations that inform the design of the AI-Enhanced Fitness and Nutrition Tracking System.

The review focuses mainly on research and systems published or available between 2021 and 2026. Older sources are included only when they describe foundational methods that remain directly relevant, such as the Mifflin-St Jeor resting-energy equation and established resistance-training progression principles.

## 2.2 Digital Health Self-Monitoring

Digital self-monitoring refers to the use of mobile applications, websites, wearable devices, or other digital tools to record and review personal behaviour or health-related data. Examples include exercise frequency, step count, food intake, body weight, sleep, and adherence to personal goals. Self-monitoring can improve awareness by making behaviour visible and measurable. However, the value of self-monitoring depends on whether the recorded information is understandable and whether the system provides useful feedback.

Tong et al. (2022) investigated the use of mobile applications and fitness trackers for supporting healthy behaviour. Among 552 adults, users of mobile applications or fitness trackers had almost twice the odds of meeting aerobic physical-activity guidelines compared with non-users. The authors noted that differences in adoption among user groups and the changing context of use indicate that technologies should adapt to different user needs. Although the cross-sectional study cannot establish a direct causal relationship, it demonstrates the practical relevance of digital tracking tools for supporting health behaviour.

Ferguson et al. (2022) conducted a systematic review of systematic reviews and meta-analyses on wearable activity trackers. The study concluded that wearable activity trackers can improve physical activity and related physiological and psychosocial outcomes across clinical and non-clinical populations. These findings support the use of self-monitoring technology as a low-cost method of encouraging users to become more aware of their activity.

Self-monitoring alone is not always sufficient. Users may record large amounts of information but still struggle to interpret what the information means. Krukowski et al. (2024) reviewed the impact of feedback on self-monitoring of dietary intake, physical activity, and weight. They identified feedback as an important part of behavioural interventions because it can facilitate goal setting and engagement. However, the evidence on the most effective method of feedback generation and presentation was mixed. This suggests that a digital platform should avoid treating every calculated result as a definitive health judgement. Feedback should be understandable, contextualised, and presented as general decision support.

The proposed system applies this principle by separating data recording from analysis and feedback. Workout, nutrition, body, and goal records are stored as structured operational data. The analytics module then generates indicators such as training load, macronutrient ratios, calorie balance, body-weight trends, and goal adherence. These indicators are presented through a dashboard, recommendations, and progress reports rather than leaving the user with raw records only.

## 2.3 Behaviour Change Techniques in Digital Health

Digital health applications frequently use behaviour change techniques to encourage continued participation. Common techniques include self-monitoring, goal setting, feedback, prompts, rewards, and progress visualisation.

Zhu et al. (2024) reviewed 41 studies of digital behaviour-change interventions for habit formation. The most frequently applied techniques were self-monitoring, goal setting, and prompts or cues. Common implementation strategies included automatic monitoring, descriptive feedback, self-set goals, general guidelines, time-based cues, and virtual rewards. Their proposed framework distinguished between target-mediated strategies, such as personalisation, and technology-mediated interaction strategies.

This finding is relevant to the proposed system for three reasons. First, daily workout and nutrition journals provide explicit self-monitoring. Second, goals and user-defined milestones allow users to convert broad intentions into measurable targets. Third, dashboard indicators and reports provide descriptive feedback. The current system does not implement persuasive notifications, social competition, or virtual rewards. These features are intentionally outside the current scope because the project focuses on integrated analysis and decision support.

The literature also indicates that the existence of a feature does not guarantee behaviour change. Salas-Groves et al. (2023) reviewed nutrition applications for people with chronic diseases and found that fewer than half of the included studies explicitly based their applications on behavioural theory or its constructs. The review found positive health outcomes but also identified declining application usage in some studies and limited evidence of sustained behaviour change. Therefore, an application should minimise unnecessary complexity and provide clear value each time the user records information.

The proposed platform uses a direct workflow: select a date, record data, review the daily journal, and then view aggregate analytics. This workflow reduces navigation between unrelated pages and allows users to edit historical records by switching to the relevant date.

## 2.4 Fitness and Workout Tracking

### 2.4.1 Workout Data Recording

A workout-tracking system commonly records exercise name, session date, sets, repetitions, weight, duration, rest time, and notes. These fields provide both an historical training log and the data required for later analysis.

The basic training-volume formula used by many resistance-training systems is:

`Training volume = sets x repetitions x external load`

This metric provides a simple representation of the total external work performed. It is useful for comparing similar exercises or sessions, although it does not fully represent exercise technique, range of motion, fatigue, movement speed, or perceived effort. For bodyweight or unloaded exercises, the proposed system uses sets multiplied by repetitions as a fallback training-load indicator because an external load may not be recorded.

The system also stores exercise reference data, including category, primary muscle group, equipment, and intensity level. The reference data supports exercise search and enables basic muscle-group distribution analysis. This is useful because total volume alone may hide an imbalanced training routine.

### 2.4.2 Progressive Overload

Progressive overload refers to gradually increasing training demands so that adaptation can continue. Training demand may be increased through external load, repetitions, sets, training frequency, exercise difficulty, or other variables. The IUSCA position stand by Schoenfeld et al. (2021) reviewed resistance-training recommendations and emphasised the manipulation of training variables for supporting hypertrophy.

The proposed system implements a simplified progressive-overload status. Workout sessions within the selected period are ordered by date and separated into earlier and more recent groups. The calculated load of the recent group is compared with the earlier group:

- An increase greater than the defined threshold is classified as `improving`.
- A decrease greater than the threshold is classified as `declining`.
- A small difference is classified as `stable`.
- Fewer than two usable sessions are classified as `insufficient data`.

This approach is intentionally simple and explainable. It is not intended to replace a coach or a complete periodisation model. A limitation is that changes in exercise selection may reduce the validity of a direct volume comparison. Future versions could compare each exercise independently and include rating of perceived exertion, repetition maximum, or velocity data.

### 2.4.3 Workout Consistency and Muscle-Group Distribution

Consistency can be represented as the number of recorded sessions compared with an expected frequency. The proposed system estimates a target of approximately four sessions per seven-day period and scales this target to the selected date range. The result is capped at 100%.

Muscle-group distribution is calculated by counting workout entries associated with each exercise's muscle group. Rule-based warnings identify absent lower-body or back training and excessive concentration on one muscle group. These warnings are not clinical or prescriptive. They serve as basic prompts that encourage the user to review the balance of recorded activities.

## 2.5 Nutrition Tracking and Dietary Analysis

### 2.5.1 Digital Food Logging

Digital food logging allows users to record food items, serving sizes, quantities, meal types, and meal times. The application can then calculate energy and nutrient totals. Compared with manual calculation, a food database can reduce repeated work and make daily dietary self-monitoring more practical.

Salas-Groves et al. (2023) reviewed 46 nutrition-application studies involving 256,430 participants with chronic diseases. The review suggested that mobile nutrition interventions can improve health outcomes, particularly when tailored to the relevant population. However, the authors also identified considerable variation in application features, behaviour-change techniques, study design, and long-term engagement. This indicates that nutrition applications should be evaluated not only by the quantity of food data available, but also by whether the results are relevant and understandable.

The current project implements a curated reference database rather than attempting to reproduce a commercial food database. Each food record contains serving size, calories, protein, carbohydrates, fat, sugar, and sodium. The user selects a food, specifies quantity, meal type, and time, and the system calculates totals. The daily nutrition page shows only records for the selected date to match the mental model of a food diary.

### 2.5.2 Calorie and Macronutrient Calculation

For a logged food quantity `q`, the system calculates the nutrient total using:

`Nutrient total = nutrient value per serving x q`

Daily energy intake is the sum of calories from all recorded meals. Daily protein, carbohydrate, and fat totals are calculated using the same approach.

Macronutrient ratios are estimated by converting grams into energy:

- Protein: 4 kcal per gram
- Carbohydrate: 4 kcal per gram
- Fat: 9 kcal per gram

The percentage for each macronutrient is calculated as:

`Macro percentage = macro energy / total macro energy x 100`

The ratio provides information about dietary composition, while the gram totals show the quantity consumed. Both representations are necessary because an unchanged percentage may hide a large change in total intake.

### 2.5.3 Nutrition-Gap Detection and Meal-Quality Score

The project uses transparent rule-based analysis rather than claiming clinical dietary assessment. Nutrition-gap rules currently identify:

- Low protein intake relative to the number of days with meal records
- Absence of recent meal records
- A high proportion of energy from fat
- A low carbohydrate proportion when calorie intake exists

Meal quality is represented as a score from 0 to 100. The score begins at 100 and applies deductions for the identified conditions. The protein target is based on days that contain meal records instead of the entire selected date range. This avoids unfairly lowering the score when a user selects a long date range but records meals on only a few days.

The score has clear limitations. It does not evaluate micronutrients, food processing, fibre, dietary restrictions, cultural eating patterns, or clinical nutritional needs. It is therefore presented as a basic log-quality indicator rather than a medical nutrition score.

## 2.6 Calorie Balance and Body-Progress Analysis

### 2.6.1 BMR and TDEE

Basal Metabolic Rate represents the approximate energy required to support basic physiological functions at rest. The proposed system uses the Mifflin-St Jeor equation because it requires variables that are available in the user profile: age, sex, height, and weight. The equation is:

For males:

`BMR = 10W + 6.25H - 5A + 5`

For females:

`BMR = 10W + 6.25H - 5A - 161`

where `W` is body weight in kilograms, `H` is height in centimetres, and `A` is age in years (Mifflin et al., 1990).

Total Daily Energy Expenditure is estimated by multiplying BMR by an activity factor:

`Estimated TDEE = BMR x activity factor`

The activity factors used by the system range from sedentary to active. Calorie balance is then calculated as:

`Calorie balance = recorded calorie intake - estimated TDEE x selected days`

The result is classified as deficit, maintenance, or surplus using rule-based thresholds.

These calculations are estimates. Individual energy expenditure varies according to body composition, occupation, exercise, health status, and measurement error. The system therefore presents these values as general indicators rather than exact physiological measurements.

### 2.6.2 Body-Weight Trend and Moving Average

Day-to-day body weight may change because of hydration, food intake, glycogen, and measurement conditions. A moving average can reduce the visual effect of short-term fluctuations.

The project calculates a seven-day moving average using records within the seven-day period ending on the latest recorded date:

`Seven-day average = sum of weight records in the seven-day window / number of records in the window`

The algorithm uses calendar days rather than simply selecting the most recent seven records. This distinction prevents an old record from being included when the user has recorded weight infrequently.

The basic trend classification compares the earliest and latest available records. Small differences are classified as stable, while larger positive or negative differences are classified as increasing or decreasing. This approach is understandable but does not provide statistical forecasting. Future work could apply regression, outlier detection, or confidence intervals.

## 2.7 Digital Goal Setting and Progress Monitoring

Goal setting is widely used in digital behaviour-change systems. Zhu et al. (2024) identified goal setting as one of the most frequently implemented techniques in digital habit-formation interventions. Goals provide a target against which behaviour and outcomes can be reviewed.

The proposed system represents a goal using:

- Goal type
- Target metric
- Target value
- Deadline
- Priority
- Status
- User-defined milestones

The deadline and target value allow the system to apply basic feasibility checks. For example, excessively aggressive weight-change targets or unrealistic workout-frequency targets can be rejected. The system does not determine medical suitability; it applies simple safety-oriented constraints.

Milestones divide a larger target into smaller checkpoints. Users can define three custom milestones with titles, values, and due dates, and mark them as completed. Goal adherence is calculated as:

`Goal adherence = completed milestones / total milestones x 100`

This design combines self-set goals, progress visibility, and descriptive feedback. A limitation is that every goal currently contributes equally to adherence. Future work could apply priority weights, overdue-milestone penalties, or metric-specific automatic completion.

## 2.8 AI-Assisted Recommendations

### 2.8.1 Potential of Large Language Models

Large language models can generate natural-language explanations and recommendations from structured context. Lai et al. (2025) reviewed LLM applications in exercise recommendations and physical activity. The review found that LLMs showed potential for tailored recommendations, accessibility, engagement, and time saving. However, only 11 studies met the review criteria, indicating that this field remains at an early stage.

The same review emphasised that LLMs should supplement rather than replace professional expertise. Expert validation was considered necessary to reduce risk. This is particularly important because generated advice may be plausible but inaccurate, overly general, or inappropriate for the user's condition.

### 2.8.2 Safety and Governance

WHO guidance on AI for health identifies ethical and governance concerns such as autonomy, transparency, accountability, bias, privacy, and safety (WHO, 2021, 2024). Generative models may produce inconsistent output and may not clearly communicate uncertainty.

The project addresses these limitations through a hybrid workflow:

1. User data is stored in structured application records.
2. Deterministic analytics calculate workout, nutrition, body, and goal indicators.
3. The indicators are converted into a limited prompt context.
4. A provider generates general-wellness recommendations.
5. A safety validator rejects content containing unsafe or medical claims.
6. The user can rate the recommendation.

The current implementation uses a mock provider. This allows the complete recommendation workflow to be demonstrated without sending personal data to an external model or depending on an API key. The provider interface is compatible with a future OpenAI-style API. Rate limiting is applied through Redis.

The mock provider is a limitation because it does not demonstrate the variability or reasoning capability of a real model. However, it also makes the current system predictable, testable, and suitable for evaluating data flow and user interaction before external AI integration.

## 2.9 Review of Existing Systems

### 2.9.1 MyFitnessPal

MyFitnessPal is an all-in-one food, exercise, calorie, macro, weight, and goal-tracking platform. Its official website reports a food database containing more than 20 million entries and support for integrations with fitness devices and health platforms (MyFitnessPal, n.d.). The system provides daily and weekly nutritional breakdowns and uses the Mifflin-St Jeor equation when estimating calorie requirements.

Its main strength is the breadth and convenience of food tracking. It also integrates several types of health data. However, its core identity remains nutrition and calorie tracking. Detailed resistance-training interpretation, custom milestone workflows, and transparent explanation of application-specific analytical rules are not its primary focus.

### 2.9.2 Cronometer

Cronometer focuses on accurate and detailed nutrition tracking. Its official materials state that users can monitor calories and up to 95 nutrients and compounds, connect devices, and track additional biometrics (Cronometer, n.d.). This makes Cronometer particularly suitable for users who require detailed micronutrient information.

Cronometer's strength is nutritional depth and its emphasis on verified food information. In comparison, the current project has a much smaller food reference database and only analyses selected nutrients. However, the project combines nutrition records with workout analysis, custom goal milestones, rule-based training warnings, recommendation feedback, and report generation within an educational full-stack architecture.

### 2.9.3 Strong

Strong is a specialised workout tracker. It supports custom exercises, supersets, RPE, advanced charts, body measurements, timers, workout scheduling, health integrations, and CSV export (Strong, n.d.). Its focused design supports efficient gym logging and detailed training history.

Strong provides greater workout-specific functionality than the proposed system. The gap relevant to this project is that a specialised workout log does not serve as a complete nutrition and calorie decision-support platform. The proposed system sacrifices some specialised workout functionality to combine workouts with meal, body, goal, recommendation, and report data.

### 2.9.4 Fitbod

Fitbod provides personalised strength workouts based on goals, experience, available equipment, progressive overload, and muscle recovery (Fitbod, n.d.). This demonstrates the value of converting historical training data into future workout decisions rather than displaying records only.

Fitbod's strength is adaptive strength-training planning. Its algorithm is proprietary, which limits external examination of the exact decision process. The current project uses simpler and more transparent rules. It also includes nutrition, calorie balance, body tracking, goal milestones, and reports rather than focusing primarily on strength-program generation.

## 2.10 Comparison of Existing Systems

**Table 2.1: Comparison of Existing Systems and the Proposed System**

| Feature | MyFitnessPal | Cronometer | Strong | Fitbod | Proposed System |
|---|---|---|---|---|---|
| User profile | Yes | Yes | Yes | Yes | Yes |
| Food and meal logging | Strong | Strong | No | Limited | Yes |
| Macro tracking | Yes | Yes | No | Limited | Yes |
| Detailed micronutrients | Yes | Strong | No | No | Limited |
| Workout logging | Basic/integrated | Basic/integrated | Strong | Strong | Yes |
| Exercise search/reference | Yes | Limited | Yes | Yes | Yes |
| Training-load analysis | Limited | No | Yes | Yes | Yes |
| Progressive-overload status | Limited | No | Charts/log based | Adaptive | Rule-based |
| Muscle-group distribution | Limited | No | Available | Recovery based | Rule-based |
| Body-weight tracking | Yes | Yes | Yes | Limited | Yes |
| Seven-day weight average | Varies by feature | Trend charts | No | No | Yes |
| Custom goals | Yes | Yes | Workout goals | Training goals | Yes |
| Custom milestones | Limited | Limited | No | No | Yes |
| Daily/weekly/monthly/custom analytics | Partial | Partial | Partial | Partial | Yes |
| AI/personalised recommendation | Increasingly available | Limited | No | Strong | Provider-based prototype |
| Recommendation feedback | Not a main feature | Not a main feature | No | Indirect | Yes |
| Downloadable progress report | Data export/premium dependent | Export available | CSV export | Limited | HTML report |
| Admin reference management | Internal only | Internal only | Internal only | Internal only | Demonstrated module |
| Transparent rules for academic evaluation | Limited | Limited | Limited | Proprietary | Yes |

The comparison is based on publicly described features available at the time of writing. Commercial products change frequently, and some features may require premium subscriptions or specific platforms.

## 2.11 Identified Research and System Gap

The literature indicates that self-monitoring, goal setting, and feedback are important components of digital behaviour-change interventions. Existing commercial systems demonstrate strong individual capabilities: MyFitnessPal provides large-scale food logging, Cronometer provides detailed nutrient analysis, Strong provides specialised workout recording, and Fitbod provides adaptive workout recommendations.

However, the review identifies several gaps relevant to this project.

First, specialised systems often emphasise one primary domain. Nutrition-oriented applications may provide only basic workout interpretation, while workout-oriented applications may provide limited nutrition analysis.

Second, commercial systems may provide recommendations or analytics without exposing the underlying rules. This is appropriate for commercial software but limits their usefulness as an academic demonstration of system analysis, explainable calculations, and modular architecture.

Third, raw records and visualisations do not always provide a complete decision-support workflow. A useful workflow should connect recording, deterministic analysis, goal monitoring, recommendation generation, user feedback, and report generation.

Fourth, AI-assisted recommendations require controlled context and safety boundaries. Literature on LLM exercise recommendations remains limited, and researchers emphasise expert validation, transparency, and safety. A system should therefore not rely entirely on unrestricted text generation.

The proposed system addresses these gaps by integrating workout, nutrition, body, and goal records in a single relational model. It uses transparent rule-based analytics before recommendation generation. It provides daily journals and multiple dashboard date ranges, custom milestones, feedback on recommendations, and downloadable reports. Its purpose is not to exceed commercial systems in food-database size or workout-program sophistication. Its contribution is an explainable, modular, end-to-end prototype that demonstrates how personal health records can be transformed into structured decision support.

## 2.12 Chapter Summary

This chapter reviewed research on digital self-monitoring, behaviour change, workout analysis, nutrition analysis, body progress, goal setting, and AI-assisted recommendations. The literature supports the use of self-monitoring, feedback, and goal setting, but also identifies challenges involving sustained engagement, interpretation, safety, and personalisation.

The comparison of existing products showed that mature applications provide strong features within their specialised areas. Nevertheless, there is an opportunity to demonstrate a transparent and integrated system that combines data recording, deterministic analytics, structured recommendations, and reporting. These findings inform the requirements, architecture, analytics rules, and safety boundaries presented in the following chapters.

## Working References for Chapter 2

Cronometer. (n.d.). *The most accurate nutrition tracking app and calorie counter*. https://cronometer.com/

Ferguson, T., Olds, T., Curtis, R., Blake, H., Crozier, A. J., Dankiw, K., Dumuid, D., Kasai, D., O'Connor, E., Virgara, R., & Maher, C. (2022). Effectiveness of wearable activity trackers to increase physical activity and improve health: A systematic review of systematic reviews and meta-analyses. *The Lancet Digital Health, 4*(8), e615-e626. https://doi.org/10.1016/S2589-7500(22)00111-X

Fitbod. (n.d.). *Less planning. More progress*. https://fitbod.me/

Krukowski, R. A., Denton, A. H., & Konig, L. M. (2024). Impact of feedback generation and presentation on self-monitoring behaviors, dietary intake, physical activity, and weight: A systematic review and meta-analysis. *International Journal of Behavioral Nutrition and Physical Activity, 21*, Article 3. https://doi.org/10.1186/s12966-023-01555-6

Lai, X., Chen, J., Lai, Y., Huang, S., Cai, Y., Sun, Z., Wang, X., Pan, K., Gao, Q., & Huang, C. (2025). Using large language models to enhance exercise recommendations and physical activity in clinical and healthy populations: Scoping review. *JMIR Medical Informatics, 13*, e59309. https://doi.org/10.2196/59309

Mifflin, M. D., St Jeor, S. T., Hill, L. A., Scott, B. J., Daugherty, S. A., & Koh, Y. O. (1990). A new predictive equation for resting energy expenditure in healthy individuals. *The American Journal of Clinical Nutrition, 51*(2), 241-247. https://doi.org/10.1093/ajcn/51.2.241

MyFitnessPal. (n.d.). *Calorie tracker and BMR calculator to reach your goals*. https://www.myfitnesspal.com/

Salas-Groves, E., Galyean, S., Alcorn, M., & Childress, A. (2023). Behavior change effectiveness using nutrition apps in people with chronic diseases: Scoping review. *JMIR mHealth and uHealth, 11*, e41235. https://doi.org/10.2196/41235

Schoenfeld, B., Fisher, J., Grgic, J., Haun, C., Helms, E., Phillips, S., Steele, J., & Vigotsky, A. (2021). Resistance training recommendations to maximize muscle hypertrophy in an athletic population: Position stand of the IUSCA. *International Journal of Strength and Conditioning, 1*(1). https://doi.org/10.47206/ijsc.v1i1.81

Strong. (n.d.). *Strong: Workout tracker and gym log*. https://www.strong.app/

Tong, H. L., Maher, C., Parker, K., Pham, T. D., Neves, A. L., Riordan, B., Chow, C. K., Laranjo, L., & Quiroz, J. C. (2022). The use of mobile apps and fitness trackers to promote healthy behaviors during COVID-19: A cross-sectional survey. *PLOS Digital Health, 1*(8), e0000087. https://doi.org/10.1371/journal.pdig.0000087

World Health Organization. (2021). *Ethics and governance of artificial intelligence for health: WHO guidance*. https://www.who.int/publications/i/item/9789240029200

World Health Organization. (2024). *Ethics and governance of artificial intelligence for health: Guidance on large multi-modal models*. https://www.who.int/publications/i/item/9789240084759

Zhu, Y., Long, Y., Wang, H., Lee, K. P., Zhang, L., & Wang, S. J. (2024). Digital behavior change intervention designs for habit formation: Systematic review. *Journal of Medical Internet Research, 26*, e54375. https://doi.org/10.2196/54375

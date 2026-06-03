package report

import (
	"fmt"
	"html"
	"strings"

	"ai-nutrition/backend/internal/modules/analytics"
)

func BuildNarrativeSummary(summary analytics.DashboardSummary) string {
	lines := []string{
		fmt.Sprintf("%s progress report", title(summary.Period)),
		fmt.Sprintf("Workout: %d sessions with %.0f kg total training volume.", summary.WorkoutSessions, summary.TrainingVolume),
		fmt.Sprintf("Nutrition: %.0f kcal consumed, including %.0f g protein, %.0f g carbs, and %.0f g fat.", summary.CaloriesIn, summary.Protein, summary.Carbohydrates, summary.Fat),
		fmt.Sprintf("Calorie balance: %.0f kcal against estimated expenditure.", summary.CalorieBalance),
		fmt.Sprintf("Body progress: latest weight %.1f kg with %s trend.", summary.LatestWeightKg, summary.WeightTrend),
		fmt.Sprintf("Goals: %d active goals with %.0f%% milestone adherence.", summary.ActiveGoals, summary.GoalAdherence),
	}

	if summary.WorkoutSessions < 3 {
		lines = append(lines, "Insight: workout consistency is below the recommended weekly target for this system.")
	}
	if summary.Protein < 420 {
		lines = append(lines, "Insight: weekly protein intake appears low; future meal suggestions should prioritize protein.")
	}
	if summary.ActiveGoals == 0 {
		lines = append(lines, "Insight: no active goals were found, so progress interpretation is less targeted.")
	}

	return strings.Join(lines, "\n")
}

func title(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func BuildHTML(reportID string, summary analytics.DashboardSummary, narrative string) string {
	escapedNarrative := html.EscapeString(narrative)
	escapedNarrative = strings.ReplaceAll(escapedNarrative, "\n", "<br>")
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Progress Report %s</title>
  <style>
    body { font-family: Arial, sans-serif; color: #18202d; margin: 40px; }
    h1 { color: #1f4e79; }
    .grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; margin: 24px 0; }
    .metric { border: 1px solid #d9e2ec; border-radius: 8px; padding: 14px; }
    .metric span { color: #667085; display: block; font-size: 12px; margin-bottom: 6px; }
    .metric strong { font-size: 20px; }
    .summary { line-height: 1.6; }
  </style>
</head>
<body>
  <h1>%s Progress Report</h1>
  <p>Report ID: %s</p>
  <div class="grid">
    <div class="metric"><span>Workout sessions</span><strong>%d</strong></div>
    <div class="metric"><span>Training volume</span><strong>%.0f kg</strong></div>
    <div class="metric"><span>Calories in</span><strong>%.0f kcal</strong></div>
    <div class="metric"><span>Calorie status</span><strong>%s</strong></div>
    <div class="metric"><span>Meal quality</span><strong>%.0f%%</strong></div>
    <div class="metric"><span>Goal adherence</span><strong>%.0f%%</strong></div>
  </div>
  <div class="summary">%s</div>
</body>
</html>`,
		html.EscapeString(reportID),
		html.EscapeString(title(summary.Period)),
		html.EscapeString(reportID),
		summary.WorkoutSessions,
		summary.TrainingVolume,
		summary.CaloriesIn,
		html.EscapeString(summary.CalorieStatus),
		summary.MealQualityScore,
		summary.GoalAdherence,
		escapedNarrative,
	)
}

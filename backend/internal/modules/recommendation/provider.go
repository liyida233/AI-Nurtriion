package recommendation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/modules/analytics"
)

type Provider interface {
	Generate(ctx context.Context, kind string, summary analytics.DashboardSummary) (string, error)
}

type MockProvider struct{}

func (p MockProvider) Generate(_ context.Context, kind string, summary analytics.DashboardSummary) (string, error) {
	lines := []string{
		"Weekly recommendation:",
		"- Keep this as general wellness guidance, not medical advice.",
	}

	if kind == "meal" || kind == "weekly" {
		if summary.Protein < 420 {
			lines = append(lines, "- Protein intake looks low for the week; add lean protein to one or two meals per day.")
		}
		if summary.CalorieBalance > 2500 {
			lines = append(lines, "- Recent calorie balance is high; consider smaller portions or lower-calorie snacks.")
		}
	}

	if kind == "workout" || kind == "weekly" {
		if summary.WorkoutSessions < 3 {
			lines = append(lines, "- Workout consistency can improve; schedule three short sessions before increasing intensity.")
		} else {
			lines = append(lines, "- Workout frequency is on track; progress by adding small increases in reps, sets, or weight.")
		}
	}

	if summary.WeightTrend == "increasing" && summary.CalorieBalance > 0 {
		lines = append(lines, "- Body trend and calorie balance both point upward; check whether this matches your current goal.")
	}
	if summary.ActiveGoals == 0 {
		lines = append(lines, "- Add one SMART goal so future recommendations can be more targeted.")
	}

	return strings.Join(lines, "\n"), nil
}

type OpenAICompatibleProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func NewProvider(cfg config.Config) Provider {
	if cfg.AIProvider == "openai_compatible" || cfg.AIProvider == "openai" {
		return OpenAICompatibleProvider{
			apiKey:  cfg.AIAPIKey,
			baseURL: strings.TrimRight(cfg.AIBaseURL, "/"),
			model:   cfg.AIModel,
			client:  &http.Client{Timeout: 25 * time.Second},
		}
	}
	return MockProvider{}
}

func (p OpenAICompatibleProvider) Generate(ctx context.Context, kind string, summary analytics.DashboardSummary) (string, error) {
	if p.apiKey == "" {
		return "", errors.New("AI_API_KEY is required for openai_compatible provider")
	}

	prompt, err := buildPrompt(kind, summary)
	if err != nil {
		return "", err
	}

	payload := chatCompletionRequest{
		Model: p.model,
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: "You generate concise, safe, general wellness recommendations. Avoid medical diagnosis, extreme dieting, or unsafe training advice.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.4,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("AI provider returned non-success status")
	}

	var decoded chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if len(decoded.Choices) == 0 {
		return "", errors.New("AI provider returned no choices")
	}
	return decoded.Choices[0].Message.Content, nil
}

func buildPrompt(kind string, summary analytics.DashboardSummary) (string, error) {
	payload, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return "", err
	}
	return "Recommendation type: " + kind + "\nUser analytics JSON:\n" + string(payload) + "\nReturn 4-6 bullet points with brief explanations.", nil
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

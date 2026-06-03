package report

type GenerateRequest struct {
	PeriodType string `json:"periodType" binding:"required"`
}

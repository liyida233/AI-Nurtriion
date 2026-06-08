package report

type GenerateRequest struct {
	PeriodType string `json:"periodType" binding:"required"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

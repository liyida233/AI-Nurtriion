package recommendation

type GenerateRequest struct {
	Type string `json:"type" binding:"required"`
}

type FeedbackRequest struct {
	Rating      string `json:"rating" binding:"required"`
	Suitability string `json:"suitability"`
	Comment     string `json:"comment"`
}

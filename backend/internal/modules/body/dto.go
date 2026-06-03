package body

type RecordRequest struct {
	RecordDate string  `json:"recordDate" binding:"required"`
	WeightKg   float64 `json:"weightKg" binding:"required,min=25,max=350"`
	Note       string  `json:"note"`
}

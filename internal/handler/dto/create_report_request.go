package dto

type CreateReportRequest struct {
	Name      string `json:"name" validate:"required"`
	BirthDate string `json:"birth_date" validate:"required"`
	BirthTime string `json:"birth_time" validate:"required"`
	Location  string `json:"location" validate:"required"`
}

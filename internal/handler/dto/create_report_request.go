package dto

type CreateReportRequest struct {
	Name         string  `json:"name" validate:"required"`
	BirthDate    string  `json:"birth_date" validate:"required"`
	BirthTime    *string `json:"birth_time,omitempty"`
	PlaceOfBirth *string `json:"place_of_birth,omitempty"`
}

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/kshitij-nehete/astro-report/internal/handler/dto"
	"github.com/kshitij-nehete/astro-report/internal/middleware"
	"github.com/kshitij-nehete/astro-report/internal/response"
	"github.com/kshitij-nehete/astro-report/internal/usecase"
)

type ReportHandler struct {
	reportUsecase *usecase.ReportUsecase
	validator     *validator.Validate
}

func NewReportHandler(reportUsecase *usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
		validator:     validator.New(),
	}
}

func (h *ReportHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateReportRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteJSONError(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	report, err := h.reportUsecase.CreateReport(
		r.Context(),
		userID,
		req.Name,
		req.BirthDate,
		req.BirthTime,
		req.Location,
	)

	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"id":        report.ID,
		"status":    report.Status,
		"expiresAt": report.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

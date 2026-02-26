package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/kshitij-nehete/astro-report/internal/domain"
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
		*req.BirthTime,
		*req.PlaceOfBirth,
	)

	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]interface{}{
		"id":        json.NewEncoder(w).Encode(report),
		"status":    report.Status,
		"expiresAt": report.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ReportHandler) GetUserReports(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteJSONError(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	reports, err := h.reportUsecase.GetUserReports(r.Context(), userID)
	if err != nil {
		response.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var summaries []dto.ReportSummaryResponse

	for _, report := range reports {
		summaries = append(summaries, dto.ReportSummaryResponse{
			ID:        report.ID.Hex(),
			Name:      report.Name,
			CreatedAt: report.CreatedAt,
			ExpiresAt: report.ExpiresAt,
			Status:    string(report.Status),
		})
	}

	json.NewEncoder(w).Encode(summaries)
}

func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteJSONError(w, http.StatusUnauthorized, "invalid user context")
		return
	}

	reportID := chi.URLParam(r, "id")

	report, err := h.reportUsecase.GetReportByID(
		r.Context(),
		userID,
		reportID,
	)
	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if report.Status == domain.StatusExpired {
		// You can optionally blank numerology data
		// but we keep it for now and let frontend blur it.
	}

	json.NewEncoder(w).Encode(report)

}

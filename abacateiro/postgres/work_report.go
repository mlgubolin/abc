package postgres

import (
	"application"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
)

// In-memory store
var workReports = make(map[int]application.WorkReport)
var idCounter = 1

type WorkReportService struct {
	db *pgxpool.Pool
}

func NewWorkReportService(db *pgxpool.Pool) *WorkReportService {
	return &WorkReportService{
		db: db,
	}
}

func RegisterWorkReportRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/workreports", createWorkReport)
	router.Get("/workreports", getAllWorkReports)
	router.Get("/workreports/{id}", getWorkReport)
	router.Put("/workreports/{id}", updateWorkReport)
	router.Delete("/workreports/{id}", deleteWorkReport)
}

func createWorkReport(w http.ResponseWriter, r *http.Request) {
	var report application.WorkReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	report.ID = idCounter
	idCounter++
	workReports[report.ID] = report

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

func getAllWorkReports(w http.ResponseWriter, r *http.Request) {
	reports := make([]application.WorkReport, 0, len(workReports))
	for _, report := range workReports {
		reports = append(reports, report)
	}
	json.NewEncoder(w).Encode(reports)
}

func getWorkReport(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || workReports[id].ID == 0 {
		http.Error(w, "WorkReport not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(workReports[id])
}

func updateWorkReport(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || workReports[id].ID == 0 {
		http.Error(w, "WorkReport not found", http.StatusNotFound)
		return
	}

	var updatedReport application.WorkReport
	if err := json.NewDecoder(r.Body).Decode(&updatedReport); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedReport.ID = id
	workReports[id] = updatedReport
	json.NewEncoder(w).Encode(updatedReport)
}

func deleteWorkReport(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || workReports[id].ID == 0 {
		http.Error(w, "WorkReport not found", http.StatusNotFound)
		return
	}
	delete(workReports, id)
	w.WriteHeader(http.StatusNoContent)
}

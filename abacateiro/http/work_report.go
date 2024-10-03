package http

import (
	"application"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// In-memory store
var workReports = make(map[int]application.WorkReport)
var idCounter = 1

func (s *Server) RegisterWorkReportRoutes(router chi.Router) {
	router.Post("/work-reports", s.handleCreateWorkReport)
	router.Get("/work-reports", s.handleGetAllWorkReports)
	router.Get("/work-reports/{id}", s.handleGetWorkReport)
	//router.Put("/work-reports/{id}", s.handleUpdateWorkReport)
	//router.Delete("/work-reports/{id}", s.handleDeleteWorkReport)
}

func (s *Server) handleCreateWorkReport(w http.ResponseWriter, r *http.Request) {
	var report application.WorkReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	report.ID = idCounter
	idCounter++
	workReports[report.ID] = report

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(report)
	if err != nil {
		return
	}
}

func (s *Server) handleGetAllWorkReports(w http.ResponseWriter, r *http.Request) {
	reports := make([]application.WorkReport, 0, len(workReports))
	for _, report := range workReports {
		reports = append(reports, report)
	}
	err := json.NewEncoder(w).Encode(reports)
	if err != nil {
		return
	}
}

func (s *Server) handleGetWorkReport(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || workReports[id].ID == 0 {
		http.Error(w, "WorkReport not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(workReports[id])
	if err != nil {
		return
	}
}

//
//func (s *Server) handleUpdateWorkReport(w http.ResponseWriter, r *http.Request) {
//	id, err := strconv.Atoi(chi.URLParam(r, "id"))
//	if err != nil || workReports[id].ID == 0 {
//		http.Error(w, "WorkReport not found", http.StatusNotFound)
//		return
//	}
//
//	var updatedReport application.WorkReport
//	if err := json.NewDecoder(r.Body).Decode(&updatedReport); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	updatedReport.ID = id
//	workReports[id] = updatedReport
//	err = json.NewEncoder(w).Encode(updatedReport)
//	if err != nil {
//		return
//	}
//}
//
//func (s *Server) handleDeleteWorkReport(w http.ResponseWriter, r *http.Request) {
//	id, err := strconv.Atoi(chi.URLParam(r, "id"))
//	if err != nil || workReports[id].ID == 0 {
//		http.Error(w, "WorkReport not found", http.StatusNotFound)
//		return
//	}
//	delete(workReports, id)
//	w.WriteHeader(http.StatusNoContent)
//}

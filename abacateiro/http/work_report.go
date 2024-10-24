package http

import (
	"application"
	workreport "application/work_report"

	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// func dump(data interface{}) {
// 	jsonData, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao serializar dados:", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(jsonData))
// }

type findWorkReportResponse struct {
	WorkReports []*application.WorkReport `json:"work_reports"`
	Metadata    application.Metadata      `json:"metadata"`
}

type findWorkReportTopicResponse struct {
	WorkReportTopics []*application.WorkReportTopic `json:"work_report_topics"`
	Metadata         application.Metadata           `json:"metadata"`
}

type findWRAdvSearchResponse struct {
	Results  []*application.WRAdvSearchResult `json:"work_report_adv_search_results"`
	Metadata application.Metadata             `json:"metadata"`
}

func (s *Server) RegisterWorkReportRoutes(router chi.Router) {
	router.Get("/work-reports", s.handleWorkReportList)
	router.Post("/work-reports/{file_name}", s.handleCreateWorkReport)
	router.Get("/work-report-topics", s.handleWorkReportTopicList)
	router.Get("/work-report-topics/adv-search", s.handleWorkReportAdvSearch)

}

func (s *Server) handleCreateWorkReport(w http.ResponseWriter, r *http.Request) {

	fileName := chi.URLParam(r, "file_name")

	if filepath.Ext(fileName) != ".docx" {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Formato inválido: apenas docx é aceito"))
		return
	}

	workReport, err := application.GetWorkReportFromFileName(fileName)

	if err != nil {
		s.Error(w, r, err)
		return
	}

	units, _, _ := s.unitService.FindUnits(r.Context(), application.UnitFilter{Name: &workReport.UnitName})
	if len(units) != 1 {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Unidade não reconhecida: %s", workReport.UnitName))
		return
	}

	workReport.UnitID = units[0].ID
	err = r.ParseMultipartForm(200 << 20)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "maximum size allowed: 200MB"))
		return
	}

	mpFile, mpHeader, err := r.FormFile("file")
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}
	defer mpFile.Close()

	zr, err := zip.NewReader(mpFile, mpHeader.Size)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}

	text, topics, err := workreport.ExtractText(zr)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}
	workReport.Text = text

	if _, err := mpFile.Seek(0, io.SeekStart); err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}

	workReport.Data, err = io.ReadAll(mpFile)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "arquivo inválido: %v", err))
		return
	}

	// verifica se o relatório já existe
	wrs, _, err := s.workReportService.FindWorkReports(application.WorkReportFilter{DocName: &workReport.DocName})
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Erro ao verificar relatório de trabalho: %v", err))
		return
	}

	if len(wrs) > 0 {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Relatório de trabalho já existe: %s", workReport.DocName))
		return
	}

	err = s.workReportService.CreateWorkReport(workReport)
	if err != nil {
		s.Error(w, r, application.Errorf(application.ErrInvalid, "Erro ao criar relatório de trabalho: %v", err))
		return
	}

	for _, topic := range topics {
		t := &application.WorkReportTopic{
			Title:        topic.Title,
			Text:         topic.Text,
			WorkReportID: workReport.ID,
		}
		if err := s.workReportService.CreateWorkReportTopic(t); err != nil {
			s.Error(w, r, err)
			return
		}
	}

	go s.workReportService.DeleteDuplicatesWorkReportTopics()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(workReport); err != nil {
		s.Error(w, r, err)
		return
	}
}

func (s *Server) handleWorkReportList(w http.ResponseWriter, r *http.Request) {

	var filter application.WorkReportFilter

	switch r.Header.Get("Content-type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			s.Error(w, r, application.Errorf(application.ERRINVALID, "Invalid JSON body"))
			return
		}
	default:
		filter.Pagination.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		filter.Pagination.PageSize = 20
	}

	filter.LimitPagination()

	workReports, meta, err := s.workReportService.FindWorkReports(filter)
	if err != nil {
		s.Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(findWorkReportResponse{
		WorkReports: workReports,
		Metadata:    meta,
	}); err != nil {
		s.Error(w, r, err)
		return
	}
}

func (s *Server) handleWorkReportTopicList(w http.ResponseWriter, r *http.Request) {

	var filter application.WorkReportTopicFilter

	switch r.Header.Get("Content-type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			s.Error(w, r, application.Errorf(application.ERRINVALID, "Invalid JSON body"))
			return
		}
	default:
		filter.Pagination.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		filter.Pagination.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	}

	filter.LimitPagination()

	search := r.URL.Query().Get("search")

	filter.GlobalSearch = &search

	topics, meta, err := s.workReportService.FindWorkReportTopics(filter)
	if err != nil {
		s.Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(findWorkReportTopicResponse{
		WorkReportTopics: topics,
		Metadata:         meta,
	}); err != nil {
		s.Error(w, r, err)
		return
	}
}

// func dump(data interface{}) {
// 	jsonData, err := json.MarshalIndent(data, "", "  ")
// 	if err != nil {
// 		fmt.Println("Erro ao serializar dados:", err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(string(jsonData))
// }

func (s *Server) handleWorkReportAdvSearch(w http.ResponseWriter, r *http.Request) {

	var filter application.WRAdvSearchFilter

	// Define content type and decode JSON body if applicable using switch
	switch r.Header.Get("Content-type") {
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
			s.Error(w, r, application.Errorf(application.ERRINVALID, "Invalid JSON body"))
			return
		}
	default:
		filter.Pagination.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
		filter.Pagination.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	}

	filter.LimitPagination()
	filter.SortDescending = true

	// Extract query parameters
	search := r.URL.Query().Get("search")
	unitID := r.URL.Query().Get("unit_id")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	years := r.URL.Query()["year[]"]

	if search != "" {
		filter.GlobalSearch = &search
	}

	if unitID != "" {
		unitIDInt, _ := strconv.Atoi(unitID)
		filter.UnitID = &unitIDInt
	}

	if from != "" {
		fromDate, _ := time.Parse("2006-01-02", from)
		filter.From = &fromDate
	}

	if to != "" {
		toDate, _ := time.Parse("2006-01-02", to)
		filter.To = &toDate
	}

	if len(years) > 0 {
		yearsInt, err := getYears(years)
		if err != nil {
			s.Error(w, r, application.Errorf(application.ErrInvalid, "Invalid year format: %v", err))
			return
		}
		filter.Years = yearsInt
	}

	// dump(map[string]interface{}{
	// 	"years":   years,
	// 	"unit_id": unitID,
	// 	"search":  search,
	// 	"from":    from,
	// 	"to":      to,
	// })

	// dump(filter)

	// Call the service to get results
	results, meta, err := s.workReportService.FindWorkReportTopicsAdvSearch(filter)
	if err != nil {
		s.Error(w, r, err)
		return
	}

	// Return the response as JSON
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(findWRAdvSearchResponse{
		Results:  results,
		Metadata: meta,
	}); err != nil {
		s.Error(w, r, err)
		return
	}
}

func getYears(yearsParams interface{}) ([]int, error) {
	var years []int

	// Verifica se é um array de strings e maior que 0
	params, ok := yearsParams.([]string)
	if !ok || len(params) == 0 {
		return nil, fmt.Errorf("year must be a non-empty array of strings")
	}

	for _, yearStr := range params {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return nil, fmt.Errorf("invalid year format: %s", yearStr)
		}
		years = append(years, year)
	}
	return years, nil
}

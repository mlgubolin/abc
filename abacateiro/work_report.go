package application

import (
	"errors"
	"net/http"
	"time"
)

type WorkReport struct {
	From     time.Time `json:"work_report_from"`
	To       time.Time `json:"work_report_to"`
	DocName  string    `json:"work_report_doc_name"`
	UnitName string    `json:"-"`
	Text     string    `json:"-"`
	Data     []byte    `json:"-"`
	ID       int       `json:"work_report_id"`
	UnitID   int       `json:"unit_id"`
}

func (wr *WorkReport) Validate() error {
	if wr.From.After(wr.To) {
		return errors.New("the 'From' date must be before the 'To' date")
	}
	if wr.DocName == "" {
		return errors.New("the 'DocName' field cannot be empty")
	}
	if wr.UnitName == "" {
		return errors.New("the 'DocName' field cannot be empty")
	}
	if wr.ID <= 0 {
		return errors.New("the 'ID' must be a positive integer")
	}
	if wr.UnitID <= 0 {
		return errors.New("the 'UnitID' must be a positive integer")
	}
	return nil
}

type WorkReportService interface {
	CreateWorkReport(http.ResponseWriter, *http.Request)
	GetWorkReport(http.ResponseWriter, *http.Request)
	GetWorkReports(http.ResponseWriter, *http.Request)
	//DeleteWorkReport(http.ResponseWriter, *http.Request)
	//UpdateWorkReport(http.ResponseWriter, *http.Request)
	//DeleteWorkReports(http.ResponseWriter, *http.Request)
}

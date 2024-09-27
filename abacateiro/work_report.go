package application

import (
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

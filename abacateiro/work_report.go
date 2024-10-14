package application

import (
	"context"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	MaxDocNameLen              = 100
	MaxWorkReportTopicTitleLen = 300
)

type WorkReport struct {
	From     time.Time `json:"work_report_from"`
	To       time.Time `json:"work_report_to"`
	Unit     *Unit     `json:"unit"`
	DocName  string    `json:"work_report_docname"`
	UnitName string    `json:"-"` // auxiliar
	Text     string    `json:"-"` // auxiliar - temporário
	Data     []byte    `json:"-"`
	ID       int       `json:"work_report_id"`
	UnitID   int       `json:"unit_id"`
}

type WorkReportTopic struct {
	WorkReport *WorkReport `json:"work_report"`
	Title      string      `json:"work_report_topic_title"`
	Text       string      `json:"work_report_topic_text"`

	ID int `json:"work_report_topic_id"`

	WorkReportID int `json:"work_report_id"`
}

type WRAdvSearchResult struct {
	HighLightTitle string `json:"highlight_title"`
	HighLightText  string `json:"highlight_text"`

	WorkReportTopicTitle string  `json:"-"` // apenas para uso interno em AltSearch
	WorkReportTopicText  string  `json:"work_report_topic_text"`
	WorkReportDocName    string  `json:"work_report_docname"`
	Rank                 float64 `json:"rank"`

	WorkReportTopicID int `json:"work_report_topic_id"`
	WorkReportID      int `json:"work_report_id"`
}

type WorkReportFilter struct {
	ID      *int       `json:"work_report_id"`
	DocName *string    `json:"work_report_docname"`
	UnitID  *int       `json:"unit_id"`
	From    *time.Time `json:"work_report_from"`
	To      *time.Time `json:"work_report_to"`

	GlobalSearch *string `json:"global_search"`

	Pagination
	LoadData bool `json:"-"`
}

type WorkReportTopicFilter struct {
	ID     *int       `json:"work_report_topic_id"`
	Title  *string    `json:"work_report_topic_title"`
	Text   *string    `json:"work_report_topic_text"`
	UnitID *int       `json:"unit_id"`
	From   *time.Time `json:"work_report_from"`
	To     *time.Time `json:"work_report_to"`

	GlobalSearch *string `json:"global_search"`

	Pagination
}

type WRAdvSearchFilter struct {
	UnitID *int       `json:"unit_id"`
	From   *time.Time `json:"work_report_from"`
	To     *time.Time `json:"work_report_to"`
	Years  []int      `json:"years"`

	GlobalSearch *string `json:"global_search"`

	Pagination
}

func GetWorkReportFromFileName(fname string) (wr *WorkReport, err error) {
	re := regexp.MustCompile(`^(P.*|SAR|MAR)_(\d{6})_(\d{6}) - (.*).docx$`)
	match := re.FindStringSubmatch(fname)
	if len(match) != 5 {
		return nil, Errorf(ErrInvalid, "Nome de arquivo inválido: esperado formato UEP_YYMMDD_YYMMDD - Autor, eg. P66_220101_220115 - Victor.docx")
	}

	wr = &WorkReport{}

	// remove extensao docx do arquivo
	// nota: o regex garante que o LastIndexByte retornará um índice válido
	wr.DocName = match[0][:strings.LastIndexByte(match[0], '.')]

	wr.UnitName = match[1]
	if wr.UnitName == "SAR" || wr.UnitName == "MAR" {
		wr.UnitName = "EDISA"
	}
	fromStr := match[2]
	toStr := match[3]
	author := match[4]
	_ = author

	const dateLayout = "060102"
	wr.From, err = time.Parse(dateLayout, fromStr)
	if err != nil {
		return nil, Errorf(ErrInvalid, "Data inválida, formato esperado: YYMMDD (eg. 220615), recebido: %s", fromStr)
	}
	wr.To, err = time.Parse(dateLayout, toStr)
	if err != nil {
		return nil, Errorf(ErrInvalid, "Data inválida: %s", toStr)
	}

	if wr.From.After(wr.To) {
		return nil, Errorf(ErrInvalid, "Data inválida: o dia %s é posterior ao dia %s", fromStr, toStr)
	}

	return
}

func (w *WorkReport) Validate() error {
	if w.DocName == "" {
		return Errorf(ERRINVALID, "DocName required.")
	} else if utf8.RuneCountInString(w.DocName) > MaxDocNameLen {
		return Errorf(ERRINVALID, "DocName too long.")
	} else if w.UnitID < 1 {
		return Errorf(ERRINVALID, "UnitID required.")
	}
	return nil
}

func (w *WorkReportTopic) Validate() error {
	if w.Title == "" {
		return Errorf(ERRINVALID, "Title required.")
	} else if utf8.RuneCountInString(w.Title) > MaxWorkReportTopicTitleLen {
		return Errorf(ERRINVALID, "Title too long.")
	} else if w.Text == "" {
		return Errorf(ERRINVALID, "Text required.")
	}
	return nil
}

type WorkReportUpdate struct {
	DocName *string    `json:"work_report_docname"`
	UnitID  *int       `json:"unit_id"`
	From    *time.Time `json:"work_report_from"`
	To      *time.Time `json:"work_report_to"`
}

type WorkReportService interface {
	CreateWorkReport(ctx context.Context, wr *WorkReport) error
	FindWorkReportByID(ctx context.Context, id int) (*WorkReport, error)
	FindWorkReports(ctx context.Context, filter WorkReportFilter) ([]*WorkReport, Metadata, error)
	UpdateWorkReport(ctx context.Context, id int, upd WorkReportUpdate) (*WorkReport, error)
	DeleteWorkReport(ctx context.Context, id int) error

	CreateWorkReportTopic(ctx context.Context, wr *WorkReportTopic) error
	DeleteDuplicatesWorkReportTopics(ctx context.Context) error
	FindWorkReportTopicByID(ctx context.Context, id int) (*WorkReportTopic, error)
	FindWorkReportTopics(ctx context.Context, filter WorkReportTopicFilter) ([]*WorkReportTopic, Metadata, error)

	FindWorkReportTopicsAdvSearch(ctx context.Context, filter WRAdvSearchFilter) ([]*WRAdvSearchResult, Metadata, error)
}

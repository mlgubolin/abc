package postgres

import (
	"application"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkReportService struct {
	db *pgxpool.Pool
}

func NewWorkReportService(db *pgxpool.Pool) *WorkReportService {
	return &WorkReportService{
		db: db,
	}
}

func (s *WorkReportService) CreateWorkReport(workReport application.WorkReport) (application.WorkReport, error) {

	if err := workReport.Validate(); err != nil {
		return application.WorkReport{}, fmt.Errorf("invalid user: %w", err)
	}
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO work_reports (
			work_report_docname,
			work_report_from,
			work_report_to,
			work_report_text,
			work_report_data,
			unit_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING work_report_id
		`,
		workReport.DocName,
		workReport.From,
		workReport.To,
		workReport.Text,
		workReport.Data,
		workReport.UnitID,
	).Scan(&workReport.ID)

	if err != nil {
		return application.WorkReport{}, fmt.Errorf("failed to create user: %w", err)
	}
	return workReport, nil
}

package postgres

import (
	"application"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
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

func (s *WorkReportService) CreateWorkReport(workReport *application.WorkReport) error {

	if err := workReport.Validate(); err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO work_reports (
			report_name,
			data_from,
			data_to,
			content,
			file_data _data,
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
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
func (s *WorkReportService) FindWorkReportByID(id int) (*application.WorkReport, error) {
	workReport := application.WorkReport{}
	err := s.db.QueryRow(context.Background(), `
		SELECT
			report_name,
			data_from,
			data_to,
			content,
			file_data _data,
			unit_id
		WHERE work_report_id = $1
		`,
		id,
	).Scan(&workReport.DocName,
		&workReport.From,
		&workReport.To,
		&workReport.Text,
		&workReport.Data,
		&workReport.UnitID)

	if err != nil {
		return &application.WorkReport{}, fmt.Errorf("failed to find user: %w", err)
	}
	return &workReport, nil
}

func (s *WorkReportService) UpdateWorkReport(id int, upd application.WorkReportUpdate) (*application.WorkReport, error) {
	var workReport application.WorkReport
	var updateFields []string
	var args []interface{}
	argCount := 1

	if upd.DocName != nil {
		updateFields = append(updateFields, fmt.Sprintf("report_name = $%d", argCount))
		args = append(args, *upd.DocName)
		argCount++
	}
	if upd.From != nil {
		updateFields = append(updateFields, fmt.Sprintf("data_from = $%d", argCount))
		args = append(args, *upd.From)
		argCount++
	}
	if upd.To != nil {
		updateFields = append(updateFields, fmt.Sprintf("data_to = $%d", argCount))
		args = append(args, *upd.To)
		argCount++
	}
	if upd.UnitID != nil {
		updateFields = append(updateFields, fmt.Sprintf("unit_id = $%d", argCount))
		args = append(args, *upd.UnitID)
		argCount++
	}

	if len(updateFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE work_reports
		SET %s
		WHERE work_report_id = $%d
		RETURNING work_report_id, report_name, data_from, data_to, content, file_data, unit_id
	`, strings.Join(updateFields, ", "), argCount)

	args = append(args, id)

	err := s.db.QueryRow(context.Background(), query, args...).Scan(
		&workReport.ID,
		&workReport.DocName,
		&workReport.From,
		&workReport.To,
		&workReport.Text,
		&workReport.Data,
		&workReport.UnitID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("work report not found")
		}
		return nil, fmt.Errorf("failed to update work report: %w", err)
	}

	return &workReport, nil
}

// }
func (s *WorkReportService) FindWorkReports(filter application.WorkReportFilter) ([]*application.WorkReport, application.Metadata, error) {
	var workReports []*application.WorkReport
	var metadata application.Metadata

	query := `
		SELECT wr.work_report_id, wr.report_name, wr.data_from, wr.data_to, wr.content, wr.file_data, wr.unit_id
		FROM work_reports wr
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if filter.UnitID != nil {
		query += fmt.Sprintf(" AND wr.unit_id = $%d", argCount)
		args = append(args, *filter.UnitID)
		argCount++
	}

	if !filter.From.IsZero() {
		query += fmt.Sprintf(" AND wr.data_from >= $%d", argCount)
		args = append(args, filter.From)
		argCount++
	}

	if !filter.To.IsZero() {
		query += fmt.Sprintf(" AND wr.data_to <= $%d", argCount)
		args = append(args, filter.To)
		argCount++
	}

	if filter.DocName != nil && *filter.DocName != "" {
		query += fmt.Sprintf(" AND wr.report_name ILIKE $%d", argCount)
		args = append(args, "%"+*filter.DocName+"%")
		argCount++
	}

	// Add ORDER BY clause
	query += " ORDER BY wr.data_from DESC"

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to query work reports: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var wr application.WorkReport
		err := rows.Scan(
			&wr.ID,
			&wr.DocName,
			&wr.From,
			&wr.To,
			&wr.Text,
			&wr.Data,
			&wr.UnitID,
		)
		if err != nil {
			return nil, metadata, fmt.Errorf("failed to scan work report: %w", err)
		}
		workReports = append(workReports, &wr)
	}

	if err = rows.Err(); err != nil {
		return nil, metadata, fmt.Errorf("error iterating work reports: %w", err)
	}

	// Get total count for metadata
	countQuery := `
		SELECT COUNT(*)
		FROM work_reports wr
		WHERE 1=1
	`
	// Reuse the WHERE conditions from the main query
	countQuery += query[strings.Index(query, "WHERE 1=1")+10 : strings.Index(query, "ORDER BY")]

	err = s.db.QueryRow(context.Background(), countQuery, args[:len(args)-2]...).Scan(&metadata.TotalCount)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to get total count: %w", err)
	}

	metadata.CurrentPage = (filter.Offset() / filter.Limit()) + 1
	metadata.PageSize = filter.Limit()
	metadata.FirstPage = 1
	metadata.LastPage = (metadata.TotalCount + filter.Limit() - 1) / filter.Limit()

	return workReports, metadata, nil

}

func (s *WorkReportService) DeleteWorkReport(id int) error {
	query := `DELETE FROM work_reports WHERE id = $1`

	result, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete work report: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return application.Errorf(application.ERRNOTFOUND, "work report not found")
	}

	return nil
}

func (s *WorkReportService) CreateWorkReportTopic(wr *application.WorkReportTopic) error {
	if err := wr.Validate(); err != nil {
		return fmt.Errorf("invalid work report topic: %w", err)
	}

	query := `
		INSERT INTO work_report_topics (
			work_report_id,
			title,
			text
		)
		VALUES ($1, $2, $3)
		RETURNING work_report_topic_id
	`

	_, err := s.db.Exec(context.Background(), query, wr.WorkReportID, wr.Title, wr.Text)
	if err != nil {
		return fmt.Errorf("failed to create work report topic: %w", err)
	}
	return nil
}

func (s *WorkReportService) DeleteDuplicatesWorkReportTopics() error {
	query := `
		WITH ranked_topics AS (
			SELECT 
				id,
				work_report_id,
				title,
				content,
				ROW_NUMBER() OVER (
					PARTITION BY work_report_id, title
					ORDER BY id
				) AS rn
			FROM work_report_topics
		)
		DELETE FROM work_report_topics
		WHERE id IN (
			SELECT id
			FROM ranked_topics
			WHERE rn > 1
		)
	`

	_, err := s.db.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to delete duplicate work report topics: %w", err)
	}

	return nil
}
func (s *WorkReportService) FindWorkReportTopicByID(id int) (*application.WorkReportTopic, error) {
	query := `
		SELECT wrt.work_report_topic_id, wrt.work_report_id, wrt.title, wrt.text,
			   wr.work_report_from, wr.work_report_to, wr.work_report_docname, wr.unit_id
		FROM work_report_topics wrt
		JOIN work_reports wr ON wrt.work_report_id = wr.work_report_id
		WHERE wrt.work_report_topic_id = $1
	`

	var topic application.WorkReportTopic
	var workReport application.WorkReport

	err := s.db.QueryRow(context.Background(), query, id).Scan(
		&topic.ID,
		&topic.WorkReportID,
		&topic.Title,
		&topic.Text,
		&workReport.From,
		&workReport.To,
		&workReport.DocName,
		&workReport.UnitID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("work report topic not found: %w", err)
		}
		return nil, fmt.Errorf("failed to find work report topic: %w", err)
	}

	workReport.ID = topic.WorkReportID
	topic.WorkReport = &workReport

	return &topic, nil
}
func (s *WorkReportService) FindWorkReportTopics(filter application.WorkReportTopicFilter) ([]*application.WorkReportTopic, application.Metadata, error) {
	query := `
		SELECT wrt.work_report_topic_id, wrt.work_report_id, wrt.title, wrt.text,
			   wr.work_report_from, wr.work_report_to, wr.work_report_docname, wr.unit_id
		FROM work_report_topics wrt
		JOIN work_reports wr ON wrt.work_report_id = wr.work_report_id
		WHERE 1=1
	`

	var topics []*application.WorkReportTopic
	var metadata application.Metadata
	args := []interface{}{}
	argCount := 1

	query += " ORDER BY wrt.work_report_topic_id DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit(), filter.Offset())

	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to query work report topics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var topic application.WorkReportTopic
		var workReport application.WorkReport

		err := rows.Scan(
			&topic.ID,
			&topic.WorkReportID,
			&topic.Title,
			&topic.Text,
			&workReport.From,
			&workReport.To,
			&workReport.DocName,
			&workReport.UnitID,
		)
		if err != nil {
			return nil, metadata, fmt.Errorf("failed to scan work report topic: %w", err)
		}
		topics = append(topics, &topic)
	}

	if err = rows.Err(); err != nil {
		return nil, metadata, fmt.Errorf("error iterating work report topics: %w", err)
	}

	// Get total count for metadata
	countQuery := `
		SELECT COUNT(*)
		FROM work_report_topics wrt
		JOIN work_reports wr ON wrt.work_report_id = wr.work_report_id
		WHERE 1=1
	`
	// Reuse the WHERE conditions from the main query
	countQuery += query[strings.Index(query, "WHERE 1=1")+10:]

	err = s.db.QueryRow(context.Background(), countQuery, args[:len(args)-2]...).Scan(&metadata.TotalCount)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to get total count: %w", err)
	}

	metadata.CurrentPage = (filter.Offset() / filter.Limit()) + 1
	metadata.PageSize = filter.Limit()
	metadata.FirstPage = 1
	metadata.LastPage = (metadata.TotalCount + filter.Limit() - 1) / filter.Limit()

	return topics, metadata, nil

}

func (s *WorkReportService) FindWorkReportTopicsAdvSearch(filter application.WRAdvSearchFilter) ([]*application.WRAdvSearchResult, application.Metadata, error) {
	var topics []*application.WRAdvSearchResult
	var metadata application.Metadata

	query := `
		SELECT wrt.work_report_topic_id, wrt.work_report_id, wrt.title, wrt.text,
			   wr.work_report_from, wr.work_report_to, wr.work_report_docname, wr.unit_id
		FROM work_report_topics wrt
		JOIN work_reports wr ON wrt.work_report_id = wr.work_report_id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	query += " ORDER BY wrt.work_report_topic_id DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to query work report topics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var topic application.WRAdvSearchResult
		var workReport application.WorkReport

		err := rows.Scan(
			&topic.WorkReportTopicID, // Change this line
			&topic.WorkReportID,
			&topic.WorkReportTopicTitle, // Change this line
			&topic.WorkReportTopicText,  // Change this line
			&workReport.From,
			&workReport.To,
			&workReport.DocName,
			&workReport.UnitID,
		)
		if err != nil {
			return nil, metadata, fmt.Errorf("failed to scan work report topic: %w", err)
		}
		topics = append(topics, &topic)
	}

	if err = rows.Err(); err != nil {
		return nil, metadata, fmt.Errorf("error iterating work report topics: %w", err)
	}

	return topics, metadata, nil
}

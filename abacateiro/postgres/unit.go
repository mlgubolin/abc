package postgres

import (
	"application"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitService struct {
	db *pgxpool.Pool
}

func NewUnitService(db *pgxpool.Pool) *UnitService {
	return &UnitService{
		db: db,
	}
}

func (s *UnitService) CreateUnit(ctx context.Context, u *application.Unit) error {

	if err := u.Validate(); err != nil {
		return fmt.Errorf("invalid work report: %w", err)
	}

	query := `INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES ($1, $2, $3, $4) RETURNING unit_id`

	err := s.db.QueryRow(context.Background(), query, u.ID, u.Name, u.StoragePath, u.PrometheusServerAddr).Scan(&u.ID)

	if err != nil {
		return fmt.Errorf("failed to create work report: %w", err)
	}

	return nil
}

func (s *UnitService) FindUnitByID(ctx context.Context, id int) (*application.Unit, error) {

	query := `SELECT unit_id, unit_name, storage_path, prometheus_server_address FROM units_v2 WHERE unit_id = $1`

	var u application.Unit
	err := s.db.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.Name, &u.StoragePath, &u.PrometheusServerAddr)

	if err != nil {
		return nil, fmt.Errorf("failed to find work report: %w", err)
	}

	return &u, nil

}

func (s *UnitService) FindUnits(ctx context.Context, filter application.UnitFilter) ([]*application.Unit, application.Metadata, error) {
	query := `SELECT unit_id, unit_name, storage_path, prometheus_server_address FROM units_v2`
	var conditions []string
	var args []interface{}
	argID := 1 // $1

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("unit_id = $%d", argID))
		args = append(args, *filter.ID)
		argID++ // $2
	}

	if len(filter.IDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("unit_id = ANY($%d)", argID))
		args = append(args, filter.IDs)
		argID++ // $3
	}

	if filter.Name != nil {
		conditions = append(conditions, fmt.Sprintf("unit_name ILIKE $%d", argID))
		args = append(args, "%"+*filter.Name+"%")
		argID++ // $4
	}

	if filter.GlobalSearch != nil {
		search := "%" + *filter.GlobalSearch + "%"
		conditions = append(conditions, fmt.Sprintf("(unit_name ILIKE $%d OR storage_path ILIKE $%d OR prometheus_server_address ILIKE $%d)", argID, argID, argID))
		args = append(args, search)
		// argID++ // $5
	}

	if len(conditions) > 0 {
		query += " WHERE 1 = 1 AND " + strings.Join(conditions, " AND ")
	}

	var us []*application.Unit
	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		return nil, application.Metadata{}, fmt.Errorf("failed to find units: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var u application.Unit
		err = rows.Scan(&u.ID, &u.Name, &u.StoragePath, &u.PrometheusServerAddr)
		if err != nil {
			return nil, application.Metadata{}, fmt.Errorf("failed to scan row: %w", err)
		}
		us = append(us, &u)
	}

	if rows.Err() != nil {
		return nil, application.Metadata{}, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return us, application.Metadata{}, nil
}

func (s *UnitService) FindAllUnits(ctx context.Context) ([]*application.UnitSimple, error) {

	query := `SELECT unit_id, unit_name FROM units_v2`

	var us []*application.UnitSimple
	rows, err := s.db.Query(context.Background(), query)

	if err != nil {
		return nil, fmt.Errorf("failed to find work reports: %w", err)
	}

	for rows.Next() {
		var u application.UnitSimple
		err = rows.Scan(&u.ID, &u.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		us = append(us, &u)
	}

	return us, nil
}

func (s *UnitService) UpdateUnit(ctx context.Context, id int, upd application.UnitUpdate) (*application.Unit, error) {
	query := `UPDATE units_v2 SET unit_name = $1, prometheus_server_address = $2 WHERE unit_id = $3 RETURNING unit_id`

	var u application.Unit
	err := s.db.QueryRow(context.Background(), query, upd.Name, upd.PrometheusServerAddr, id).Scan(&u.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to update work report: %w", err)
	}

	return &u, nil
}

func (s *UnitService) DeleteUnit(ctx context.Context, id int) error {
	query := `DELETE FROM units_v2 WHERE unit_id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	if err != nil {
		return fmt.Errorf("failed to delete work report: %w", err)
	}

	return nil
}

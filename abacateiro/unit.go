package application

import (
	"context"
	"unicode/utf8"
)

const (
	MaxUnitNameLen = 300
)

type Unit struct {
	Name                 string `json:"unit_name"`
	StoragePath          string `json:"unit_storage_path"`
	PrometheusServerAddr string `json:"prometheus_server_address"`

	ID int `json:"unit_id"`

	TotalDevices int `json:"total_devices"`
}

type UnitSimple struct {
	Name string `json:"unit_name"`
	ID   int    `json:"unit_id"`
}

// type UnitStorage struct {
// 	ID          int    `json:"unit_id"`
// 	Name        string `json:"unit_name"`
// 	StoragePath string `json:"unit_storage_path"`
// }

type UnitFilter struct {
	ID   *int    `json:"unit_id"`
	IDs  []int   `json:"unit_ids"`
	Name *string `json:"unit_name"`

	GlobalSearch     *string `json:"global_search"`
	DoNotLoadSummary *bool

	Pagination
}

type UnitUpdate struct {
	Name                 *string `json:"unit_name"`
	PrometheusServerAddr *string `json:"prometheus_server_address"`
}

func (a *Unit) Validate() error {
	if a.Name == "" {
		return Errorf(ERRINVALID, "Unit name required.")
	} else if utf8.RuneCountInString(a.Name) > MaxUnitNameLen {
		return Errorf(ERRINVALID, "Unit name too long.")
	}
	return nil
}

type UnitService interface {
	CreateUnit(ctx context.Context, u *Unit) error
	FindUnitByID(ctx context.Context, id int) (*Unit, error)
	FindUnits(ctx context.Context, filter UnitFilter) ([]*Unit, Metadata, error)
	FindAllUnits(ctx context.Context) ([]*UnitSimple, error)
	UpdateUnit(ctx context.Context, id int, upd UnitUpdate) (*Unit, error)
	DeleteUnit(ctx context.Context, id int) error
	//FindUnitStoragePathByID(ctx context.Context, id int) (*UnitStorage, error)
}

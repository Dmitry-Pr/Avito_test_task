package repositories

import (
	"database/sql"
	"log"
)

type MerchRepositoryInterface interface {
	GetAll() ([]string, error)
}

type MerchRepository struct {
	db *sql.DB
}

func NewMerchRepository(db *sql.DB) MerchRepositoryInterface {
	return &MerchRepository{db: db}
}

func (r *MerchRepository) GetAll() ([]string, error) {
	rows, err := r.db.Query("SELECT name FROM merch")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var merchList []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		merchList = append(merchList, name)
	}

	return merchList, nil
}

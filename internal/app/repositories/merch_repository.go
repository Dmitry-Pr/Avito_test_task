package repositories

import (
	"database/sql"
)

type MerchRepository struct {
	db *sql.DB
}

func NewMerchRepository(db *sql.DB) *MerchRepository {
	return &MerchRepository{db: db}
}

func (r *MerchRepository) GetAll() ([]string, error) {
	rows, err := r.db.Query("SELECT name FROM merch")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

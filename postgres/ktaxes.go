package postgres

import (
	"time"

	"github.com/NouNhaN-GitHub/assessment-tax/ktaxes"
)

type Allowance struct {
	AllowanceType string    `postgres:"allowance_type"`
	Amount        float64   `postgres:"amount"`
	CreatedAt     time.Time `postgres:"created_at"`
}

func (p *Postgres) Allowances() ([]ktaxes.Allowance, error) {
	rows, err := p.Db.Query("SELECT * FROM allowances")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allowances []ktaxes.Allowance
	for rows.Next() {
		var a Allowance
		err := rows.Scan(&a.AllowanceType,
			&a.Amount,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		allowances = append(allowances, ktaxes.Allowance{
			AllowanceType: a.AllowanceType,
			Amount:        a.Amount,
			CreatedAt:     a.CreatedAt,
		})
	}
	return allowances, nil
}

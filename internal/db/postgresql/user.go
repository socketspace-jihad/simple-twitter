package postgresql

import (
	"simple_twitter/internal/models"
)

func (p *PostgreSQL) SaveUser(u *models.User) error {
	tx, err := p.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec("INSERT INTO users (display_name, username, passwd) VALUES ($1, $2, $3)", u.DisplayName, u.Username, u.Password); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *PostgreSQL) GetUser(u *models.User) error {
	rows, err := p.DB.Query("SELECT id,display_name,username,passwd FROM users WHERE username=$1 OR id=$2", u.Username, u.ID)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.DisplayName, &u.Username, &u.Password); err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgreSQL) DeleteUser(u *models.User) error {
	return nil
}

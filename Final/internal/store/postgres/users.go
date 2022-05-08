package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"hw9/internal/models"
	"hw9/internal/store"
)

func (db *DB) Users() store.UsersRepository {
	if db.users == nil {
		db.users = NewUsersRepository(db.conn)
	}
	return db.users
}

type UsersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UsersRepository {
	return &UsersRepository{conn: conn}
}

func (c UsersRepository) Create(ctx context.Context, user *models.User) error {
	_, err := c.conn.Exec("INSERT INTO users(name, surname, email, password) VALUES ($1, $2, $3, $4)", user.Name, user.Surname, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (c UsersRepository) All(ctx context.Context, filter *models.UsersFilter) ([]*models.User, error) {
	users := make([]*models.User, 0)
	basicQuery := "SELECT * FROM users"

	if filter.Query != nil {
		basicQuery = fmt.Sprintf("%s WHERE name ILIKE $1", basicQuery)

		if err := c.conn.Select(&users, basicQuery, "%"+*filter.Query+"%"); err != nil {
			return nil, err
		}
		return users, nil
	}

	if err := c.conn.Select(&users, basicQuery); err != nil {
		return nil, err
	}

	return users, nil
}

func (c UsersRepository) ByID(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)
	if err := c.conn.Get(user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return nil, err
	}

	return user, nil
}

func (c UsersRepository) Update(ctx context.Context, user *models.User) error {
	_, err := c.conn.Exec("UPDATE users SET name = $1, surname = $2, password = $3  WHERE id = $4", user.Name, user.Surname, user.Password, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c UsersRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
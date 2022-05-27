package sqlite

import (
	"context"
	"database/sql"
	"forum/models"
	"log"

	"github.com/pkg/errors"
)

type User struct {
	ID       int
	Username string
	Password string
}

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	model := toSqlUser(user)

	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		username TEXT,
		password TEXT
	  );`)
	if err != nil {
		log.Fatalf("cannot exec file: %v", err.Error())
	}
	sqlStatement := `insert into users (username, password) values ($1, $2)`
	_, err = r.db.Exec(sqlStatement, model.Username, model.Password)
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}

	return nil
}

func (r AuthRepository) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	user := new(User)
	user.Username = username
	user.Password = password
	sqlRow := `SELECT * from users where username = $1 AND password = $2`
	rows, err := r.db.Query(sqlRow, user.Username, user.Password)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return toModel(user), nil
	}
	return nil, errors.Errorf("Not registered")
}

func toSqlUser(u *models.User) *User {
	return &User{
		Username: u.Username,
		Password: u.Password,
	}
}

func toModel(u *User) *models.User {
	return &models.User{
		Id:       u.ID,
		Username: u.Username,
		Password: u.Password,
	}
}

package models

import (
	"errors"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const passwordCost = 12

type User struct {
	ID        int       `db:"id,omitempty"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
	Activated bool      `db:"activated"`
}

func (m UsersModel) Table() string {
	return "users"
}

type UsersModel struct {
	db db.Session
}

func (m UsersModel) Get(id int) (*User, error) {

	var u User
	err := m.db.Collection(m.Table()).Find(db.Cond{"id": id}).One(&u)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrNoMoreRows
		}
		return nil, err
	}

	return &u, nil
}

func (m UsersModel) FindByEmail(email string) (*User, error) {

	var u User
	err := m.db.Collection(m.Table()).Find(db.Cond{"email": email}).One(&u)
	if err != nil {
		if errors.Is(err, db.ErrNoMoreRows) {
			return nil, ErrNoMoreRows
		}
		return nil, err
	}

	return &u, nil
}

func (m UsersModel) Insert(u *User) error {

	newHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), passwordCost)
	if err != nil {
		return err
	}

	u.Password = string(newHash)
	u.CreatedAt = time.Now()
	col := m.db.Collection(m.Table())
	res, err := col.Insert(u)
	if err != nil {
		switch {
		case errHasDuplicate(err, "users_email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	u.ID = convertUpperIDtoInt(res.ID())

	return nil
}

func (u *User) ComparePassword(plainPassword string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (m UsersModel) Authenticate(email, password string) (*User, error) {
	user, err := m.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if !user.Activated {
		return nil, ErrUserNotActive
	}

	match, err := user.ComparePassword(password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, ErrInvalidLogin
	}

	return user, nil
}

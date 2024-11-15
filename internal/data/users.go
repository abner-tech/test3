package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/abner-tech/Comments-Api.git/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var AnonymouseUser = &User{}

type User struct {
	ID         int64     `json:"id"`
	Created_At time.Time `json:"created_at"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   password  `json:"-"`
	Activated  bool      `json:"activated"`
	Version    int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

// struct setup for our model
type UserModel struct {
	DB *sql.DB
}

// insert new user to the database
func (u *UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{user.Username, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//if an email already exists, we will get a pq error message
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Created_At, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// get a user from the database based on their email provided
func (u *UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, created_at, username, email, password_hash, activated, version
	FROM users
	WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Created_At,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *UserModel) GetForToken(tokenScope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
	SELECT users.id, users.created_at, users.username, users.email, users.password_hash, users.activated, users.version
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	and tokens.expiry > $3
	`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Created_At,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//return the correct user
	return &user, nil
}

/*
update a user. If the version number is different than what was before we ran the query , it means someone

	did a previouse edit or is doing an edit, so our query will fail and we would need to try again a bit later
*/
func (u *UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET username = $1, email =$2, password_hash = $3, activated = $4,
	version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version
	`

	args := []any{user.Username, user.Email, user.Password.hash, user.Activated, user.ID, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	//check for errors during an update
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique key constraints "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConfilct
		default:
			return err
		}
	}
	return nil
}

// the set method computes the hash of the password
func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plainTextPassword
	p.hash = hash

	return nil
}

// compare if client-provided plain-text password matches saved hashed-password version
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil

		default:
			return false, nil
		}
	}
	return true, nil //when password is correct
}

// validation for the email address
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// check for password to be valid
func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 7, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "mustnot be more than 72 bytes long")
}

// validate username
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 200, "username", "must not be more than 200 bytes long")

	//validate user for email
	ValidateEmail(v, user.Email)
	//validate the plain text email
	if user.Password.plaintext != nil {
		ValidatePassword(v, *user.Password.plaintext)
	}

	//check if we messed up in our codebase
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

// check if current user is anonymous
func (u *User) IsAnonymous() bool {
	return u == AnonymouseUser
}

package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userContextKey contextKey = "auth_user"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Service struct {
	db        *sql.DB
	jwtSecret []byte
}

type Claims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewService(db *sql.DB, jwtSecret string) *Service {
	return &Service{db: db, jwtSecret: []byte(jwtSecret)}
}

func (s *Service) EnsureAdmin(username, password string) error {
	var count int
	if err := s.db.QueryRow("SELECT COUNT(1) FROM users WHERE username = ?", username).Scan(&count); err != nil {
		return fmt.Errorf("count admin users: %w", err)
	}
	if count > 0 {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	_, err = s.db.Exec(
		"INSERT INTO users (username, password_hash, role, disabled) VALUES (?, ?, 'admin', 0)",
		username,
		string(hash),
	)
	if err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}
	return nil
}

func (s *Service) Login(username, password string) (string, User, error) {
	var user User
	var passwordHash string
	var disabled bool
	err := s.db.QueryRow(
		"SELECT id, username, password_hash, role, disabled FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &passwordHash, &user.Role, &disabled)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", User{}, errors.New("invalid credentials")
		}
		return "", User{}, fmt.Errorf("find user: %w", err)
	}
	if disabled {
		return "", User{}, errors.New("user disabled")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return "", User{}, errors.New("invalid credentials")
	}

	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return "", User{}, fmt.Errorf("sign jwt: %w", err)
	}
	return token, user, nil
}

func (s *Service) ParseToken(tokenString string) (User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return User{}, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return User{}, errors.New("invalid token")
	}
	return User{ID: claims.UserID, Username: claims.Username, Role: claims.Role}, nil
}

func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey).(User)
	return user, ok
}

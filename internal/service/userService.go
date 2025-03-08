package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
	"math/rand"
	"strconv"
	"time"
	"valerii/crudbananas/internal/domain"
)

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type SessionsRepository interface {
	Create(ctx context.Context, token domain.RefreshSession) error
	Get(ctx context.Context, token string) (domain.RefreshSession, error)
}

type Users struct {
	userRepo    UsersRepository
	sessionRepo SessionsRepository
	hasher      PasswordHasher
	auditClient AuditClient

	hmacSecret []byte //Подпись для jwt токена
}

func NewUsers(userRepo UsersRepository, sessionRepo SessionsRepository, auditClient AuditClient, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		hasher:      hasher,
		auditClient: auditClient,
		hmacSecret:  secret,
	}
}

func (s *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	//Хешируем пароль
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	//Создаем объект структуры domain.User
	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	//Передаем на уровень репозитория
	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	userCred, err := s.userRepo.GetByCredentials(ctx, inp.Email, password)
	if err != nil{
		return err
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    "REGISTER",
		Entity:    "USER",
		EntityID:  userCred.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "SignUp",
		}).Error(err)
	}

	return nil
}

func (s *Users) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	//хеширование пароля
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	//Получение user ID с помощью email и passoword
	user, err := s.userRepo.GetByCredentials(ctx, inp.Email, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", domain.ErrUserNotFound
		}

		return "", "", err
	}

	//Генерируем токены
	accessToken, refreshToken, err := s.generateTokens(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	if err := s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    "LOGIN",
		Entity:    "USER",
		EntityID:  user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "SignIn",
		}).Error(err)
	}

	return accessToken, refreshToken, nil
}

func (s *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	//Разбираем токен на структуру jwt.Token
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.hmacSecret, nil
	})
	if err != nil {
		return 0, err
	}

	//Проверка на валидность
	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	//Достаем claims
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	//Достаем subject в котором хранится ID
	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	//Преобразуем ID в int
	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

/*
generateTokens() генерирует access и refresh токены
и создает их БД
*/
func (s *Users) generateTokens(ctx context.Context, userId int64) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userId)),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
	})

	//Генерация access токена
	accessToken, err := t.SignedString(s.hmacSecret)
	if err != nil {
		return "", "", err
	}

	//Генерация refresh токена
	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	//Создание токена в базе данных
	if err := s.sessionRepo.Create(ctx, domain.RefreshSession{
		UserID:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *Users) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.sessionRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return s.generateTokens(ctx, session.UserID)
}

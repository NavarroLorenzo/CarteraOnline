package auth

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"cartera-app/backend/internal/shared/normalize"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("credenciales inválidas")
var ErrEmailAlreadyInUse = errors.New("el email ya está en uso")
var ErrUsernameAlreadyInUse = errors.New("el nombre de usuario ya está en uso")
var ErrUserNotFound = errors.New("el usuario no existe")
var ErrInvalidEmail = errors.New("el email no es válido")
var ErrInvalidUsername = errors.New("el nombre de usuario no es válido")
var ErrWeakPassword = errors.New("la contraseña debe tener al menos 8 caracteres")

var usernamePattern = regexp.MustCompile(`^[a-z0-9._-]+$`)

type Service interface {
	Register(input RegisterInput) (AuthResponse, error)
	Login(input LoginInput) (AuthResponse, error)
	GetByID(id int64) (User, error)
}

type service struct {
	repo         Repository
	tokenManager *TokenManager
}

func NewService(repo Repository, tokenManager *TokenManager) Service {
	return &service{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (s *service) Register(input RegisterInput) (AuthResponse, error) {
	email, err := normalizeEmail(input.Email)
	if err != nil {
		return AuthResponse{}, err
	}

	username, err := normalizeUsername(input.Username)
	if err != nil {
		return AuthResponse{}, err
	}

	if len(strings.TrimSpace(input.Password)) < 8 {
		return AuthResponse{}, ErrWeakPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.Create(email, username, string(passwordHash))
	if err != nil {
		return AuthResponse{}, err
	}

	if err := s.repo.ClaimOrphanedData(user.ID); err != nil {
		return AuthResponse{}, err
	}

	token, err := s.tokenManager.Generate(user.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *service) Login(input LoginInput) (AuthResponse, error) {
	identifier := normalize.LowerIdentifier(input.Identifier)

	record, found, err := s.repo.GetByIdentifier(identifier)
	if err != nil {
		return AuthResponse{}, err
	}
	if !found {
		return AuthResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(record.PasswordHash), []byte(input.Password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	token, err := s.tokenManager.Generate(record.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token: token,
		User:  record.User,
	}, nil
}

func (s *service) GetByID(id int64) (User, error) {
	user, found, err := s.repo.GetByID(id)
	if err != nil {
		return User{}, err
	}
	if !found {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func normalizeEmail(value string) (string, error) {
	email := normalize.LowerIdentifier(value)
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", ErrInvalidEmail
	}

	return email, nil
}

func normalizeUsername(value string) (string, error) {
	username := normalize.LowerIdentifier(value)
	if len(username) < 3 || !usernamePattern.MatchString(username) {
		return "", ErrInvalidUsername
	}

	return username, nil
}

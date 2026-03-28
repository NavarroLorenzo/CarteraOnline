package auth

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("credenciales inválidas")
var ErrEmailAlreadyInUse = errors.New("el email ya está en uso")
var ErrUsernameAlreadyInUse = errors.New("el nombre de usuario ya está en uso")
var ErrUserNotFound = errors.New("el usuario no existe")

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
	email := normalizeIdentifier(input.Email)
	username := normalizeIdentifier(input.Username)

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
	identifier := normalizeIdentifier(input.Identifier)

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

func normalizeIdentifier(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

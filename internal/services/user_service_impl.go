package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userServiceImpl struct {
	userRepo  ports.UserRepository
	jwtSecret string
}

func NewUserService(repo ports.UserRepository, jwtSecret string) ports.UserService {
	return &userServiceImpl{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

func (s *userServiceImpl) Register(username, password string) (*domain.User, error) {
	// Verificar si el usuario ya existe
	if _, err := s.userRepo.FindByUsername(username); !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         domain.UserRole, // Por defecto, rol 'user'
	}

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userServiceImpl) Login(username, password string) (string, *domain.Role, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Crear token JWT
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}
	return tokenString, &user.Role, nil
}

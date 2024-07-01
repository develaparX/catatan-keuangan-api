package service

import (
	"database/sql"
	"fmt"
	"livecode-catatan-keuangan/models"
	"livecode-catatan-keuangan/models/dto"
	"livecode-catatan-keuangan/repository"
	"livecode-catatan-keuangan/utils"
)

type UserService interface {
	CreateNew(payload models.User) (models.User, error)
	Login(payload dto.LoginDto) (dto.LoginResponseDto, error)
}

type userService struct {
	repo       repository.UserRepository
	jwtService JwtService
}

func (c *userService) Login(payload dto.LoginDto) (dto.LoginResponseDto, error) {
	user, err := c.repo.GetByEmail(payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.LoginResponseDto{}, fmt.Errorf("email not found")
		}
		return dto.LoginResponseDto{}, fmt.Errorf("failed to fetch user: %v", err)
	}

	err = utils.ComparePasswordHash(user.Password, payload.Password)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("password incorrect")
	}

	user.Password = ""
	token, err := c.jwtService.GenerateToken(user)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed to create token")
	}

	return token, nil
}

// CreateNew implements UserService.
func (c *userService) CreateNew(payload models.User) (models.User, error) {

	passwordHash, error := utils.EncryptPassword(payload.Password)
	if error != nil {
		return models.User{}, error
	}
	payload.Password = passwordHash
	return c.repo.CreateNew(payload)
}

func NewUserService(repositori repository.UserRepository, jS JwtService) UserService {
	return &userService{repo: repositori, jwtService: jS}
}

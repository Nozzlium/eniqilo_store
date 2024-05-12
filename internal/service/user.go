package service

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   *repository.UserRepository
	secret string
	salt   int
}

func NewUserService(
	repo *repository.UserRepository,
	secret string,
	salt int,
) *UserService {
	return &UserService{
		repo:   repo,
		secret: secret,
		salt:   salt,
	}
}

func (service *UserService) Register(
	ctx context.Context,
	user model.User,
) (model.RegisterResponse, error) {
	generatedUUID, err := uuid.NewV7()
	if err != nil {
		return model.RegisterResponse{}, err
	}

	hashedPass, err := bcrypt.GenerateFromPassword(
		[]byte(user.Password),
		int(service.salt),
	)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	userResult, err := service.repo.FindByPhoneNumber(
		ctx,
		user.PhoneNumber,
	)
	if err != nil &&
		!errors.Is(
			constant.ErrNotFound,
			err,
		) {
		return model.RegisterResponse{}, err
	}

	if userResult.PhoneNumber == user.PhoneNumber {
		return model.RegisterResponse{}, constant.ErrConflict
	}

	user.ID = generatedUUID
	user.Password = string(hashedPass)
	inserted, err := service.repo.Save(
		ctx,
		user,
	)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	accessToken, err := generateJwtToken(
		service.secret,
		userResult,
	)
	if err != nil {
		return model.RegisterResponse{}, err
	}

	return model.RegisterResponse{
		PhoneNumber: inserted.PhoneNumber,
		Name:        inserted.Name,
		AccessToken: accessToken,
	}, nil
}

func (service *UserService) Login(
	ctx context.Context,
	user model.User,
) (model.LoginResponse, error) {
	userResult, err := service.repo.FindByPhoneNumber(
		ctx,
		user.PhoneNumber,
	)
	if err != nil {
		return model.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(userResult.Password),
		[]byte(user.Password),
	)
	if err != nil {
		if errors.Is(
			err,
			bcrypt.ErrMismatchedHashAndPassword,
		) {
			return model.LoginResponse{}, errors.New(
				"invalid credentials",
			)
		}
		return model.LoginResponse{}, err
	}

	accessToken, err := generateJwtToken(
		service.secret,
		userResult,
	)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		PhoneNumber: userResult.PhoneNumber,
		Name:        userResult.Name,
		AccessToken: accessToken,
	}, nil
}

func (service *UserService) ValidateUserData(
	ctx context.Context,
) (bool, error) {
	userID := ctx.Value("userID").(string)
	email := ctx.Value("email").(string)
	_, err := service.repo.FindByPhoneNumberAndID(
		ctx,
		userID,
		email,
	)
	if err != nil {
		if errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func generateJwtToken(
	secret string,
	user model.User,
) (string, error) {
	token := jwt.New(
		jwt.SigningMethodHS256,
	)

	claims := token.Claims.(jwt.MapClaims)
	userID := base64.RawStdEncoding.EncodeToString(
		[]byte(user.ID.String()),
	)
	email := base64.RawStdEncoding.EncodeToString(
		[]byte(user.PhoneNumber),
	)
	claims["ui"] = userID
	claims["ea"] = email
	claims["exp"] = time.Now().
		Add(time.Hour * 72).
		Unix()

	t, err := token.SignedString(
		[]byte(secret),
	)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return t, nil
}

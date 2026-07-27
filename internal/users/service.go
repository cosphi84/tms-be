package users

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tms-be/internal/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	HashPassword(string) (string, error)
	VerifyPassword(string, string) bool
	Create(context.Context, *CreateUserDTO) error
	AllUsers(context.Context) ([]Model, error)
	Update(context.Context, uint64, *CreateUserDTO) error
	Activate(context.Context, int64) error
	Deactivate(context.Context, int64) error
	Delete(context.Context, int64) error
}

type userService struct {
	usrRepo   Repository
	authorize *auth.RoleService
}

func NewUserService(r Repository, a *auth.RoleService) UserService {
	return &userService{
		usrRepos:  r,
		authorize: a,
	}
}

func (u *userService) HashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *userService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *userService) Create(ctx context.Context, usr *CreateUserDTO) error {
	// Cek user is registered
	exists, err := u.usrRepo.IsExists(ctx, usr.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(fmt.Sprintf("User with email %s already registered.", usr.Email))
	}

	loggedUser, err := auth.GetClaims(ctx)
	if err != nil {
		return err
	}

	// hash password
	hashed, err := u.HashPassword(usr.Password)
	if err != nil {
		return err
	}

	newUsr := Model{
		Username:    usr.Name,
		Email:       usr.Email,
		Password:    hashed,
		OfficeID:    uint64(usr.OfficeID),
		FotoProfile: &usr.Image,
		IsActive:    true,
		CreatedAt:   time.Now(),
		CreatedBy:   &loggedUser.UserID,
	}

	// set role
	_, err = u.authorize.GrantRole(usr.Email, usr.Role)
	if err != nil {
		return err
	}

	return u.usrRepo.Create(ctx, &newUsr)
}

func (u *userService) AllUsers(ctx context.Context) ([]Model, error) {
	return u.usrRepo.FindAll(ctx)
}

func (u *userService) Update(ctx context.Context, id uint64, usr *UpdateUserDTO) error {
	theUser, err := u.usrRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(fmt.Sprintf("user %s not exists", usr.Email))
		}
	}

	// check if password is change
	if usr.Password != "" {
		hashed, err := u.HashPassword(usr.Password)
		if err != nil {
			return err
		}

		pwdMatch := u.VerifyPassword(theUser.Password, hashed)
		if !pwdMatch {
			return errors.New("old password mismatch")
		}

		theUser.Password = hashed
	}
	if usr.Email != "" {
		theUser.Email = usr.Email
	}
	if usr.OfficeID > 0 {
		theUser.OfficeID = usr.OfficeID
	}
	if usr.Image != "" {
		theUser.FotoProfile = &usr.Image
	}
	usr.Name == theUser.Username ? theUser.Username = usr.Name : 
	theUser.UpdatedAt = time.Now()

	if usr.Role != "" {
		_, err = u.authorize.GrantRole(usr.Email, usr.Role)
		if err != nil {
			return err
		}
	}

	_, err = u.usrRepo.Update(ctx, theUser.ID, theUser)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) Activate(ctx context.Context, ID uint64)

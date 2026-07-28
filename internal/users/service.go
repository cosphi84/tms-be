package users

import (
	"context"
	"errors"
	"fmt"
	"time"
	"tms-be/internal/auth"
	"tms-be/internal/casbin"

	"gorm.io/gorm"
)

type UserService interface {
	Create(context.Context, *CreateUserDTO) error
	AllUsers(context.Context) ([]Model, error)
	Update(context.Context, uint64, *UpdateUserDTO) error
	Activate(context.Context, uint64) error
	Deactivate(context.Context, uint64) error
	Delete(context.Context, uint64) error
}

type userService struct {
	usrRepo   Repository
	authorize *casbin.Service
}

func NewUserService(r Repository, a *casbin.Service) UserService {
	return &userService{
		usrRepo:   r,
		authorize: a,
	}
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

	if usr.Name != theUser.Username {
		theUser.Username = usr.Name
	}

	if usr.OfficeID != theUser.OfficeID {
		theUser.OfficeID = usr.OfficeID
	}

	if usr.Image != *theUser.FotoProfile {
		theUser.FotoProfile = &usr.Image
	}

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

func (u *userService) Activate(ctx context.Context, ID uint64) error {
	usr, err := u.usrRepo.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	usr.IsActive = true
	usr.UpdatedAt = time.Now()

	_, err = u.usrRepo.Update(ctx, ID, usr)
	if err != nil {
		return err
	}

	return nil

}

func (u *userService) Deactivate(ctx context.Context, ID uint64) error {
	usr, err := u.usrRepo.FindByID(ctx, ID)
	if err != nil {
		return err
	}

	usr.IsActive = false
	usr.UpdatedAt = time.Now()

	_, err = u.usrRepo.Update(ctx, ID, usr)
	if err != nil {
		return err
	}

	return nil
}

func (u *userService) Delete(ctx context.Context, ID uint64) error {
	_, err := u.usrRepo.Delete(ctx, ID)
	return err
}

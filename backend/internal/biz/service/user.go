package service

import (
	"sun-panel/internal/biz/repository"
)

type UserService struct {
	itemGroupRepo repository.IItemIconGroupRepo
	userRepo      repository.IUserRepo
}

type IUserService interface {
	CreateUser(user *repository.User) error
}

func NewUserService(userRepo repository.IUserRepo, itemGroupRepo repository.IItemIconGroupRepo) *UserService {
	return &UserService{userRepo: userRepo, itemGroupRepo: itemGroupRepo}
}

func (s *UserService) CreateUser(user *repository.User) error {
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	defaultGroup := repository.ItemIconGroup{
		Title:  "APP",
		UserId: user.ID,
		Icon:   "material-symbols:ad-group-outline",
	}

	if err := s.itemGroupRepo.Save(&defaultGroup); err != nil {
		return err
	}

	return nil
}

package user

import (
	"github.com/col3name/images-search/pkg/common/util"
	"github.com/col3name/images-search/pkg/domain/auth"
	"github.com/col3name/images-search/pkg/domain/user"
)

type Service struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) *Service {
	s := new(Service)
	s.userRepo = userRepo

	return s
}

func (s *Service) Find(usernameOrId string) (*user.FindUserDto, error) {
	var userObject user.User
	var err error
	uuid := util.IsUUID(usernameOrId)

	if uuid {
		userObject, err = s.userRepo.FindByUserName(usernameOrId)
	} else {
		userObject, err = s.userRepo.FindByUserName(usernameOrId)
	}
	if err != nil {
		return nil, auth.ErrUserNotExist
	}

	return &user.FindUserDto{
		Id:       userObject.Id,
		Username: usernameOrId,
		Role:     userObject.Role.Values(),
	}, nil
}

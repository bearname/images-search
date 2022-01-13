package user

import (
	"photofinish/pkg/domain/auth"
	"photofinish/pkg/domain/user"
	"regexp"
)

type Service struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) *Service {
	s := new(Service)
	s.userRepo = userRepo

	return s
}

func (s *Service) Find(usernameOrId string) (user.FindUserDto, error) {
	var userObject user.User
	var err error
	uuid := s.isUUID(usernameOrId)

	if uuid {
		userObject, err = s.userRepo.FindByUserName(usernameOrId)
	} else {
		userObject, err = s.userRepo.FindByUserName(usernameOrId)
	}
	if err != nil {
		return user.FindUserDto{}, auth.ErrUserNotExist
	}

	return user.FindUserDto{Username: usernameOrId,
		Role: userObject.Role.Values(),
	}, nil
}

func (s *Service) isUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

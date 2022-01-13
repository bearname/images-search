package auth

import (
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"photofinish/pkg/domain/auth"
	"photofinish/pkg/domain/user"
	"time"
)

type ServiceImpl struct {
	userRepo user.Repository
}

func NewAuthService(userRepository user.Repository) *ServiceImpl {
	v := new(ServiceImpl)
	v.userRepo = userRepository
	return v
}

func (s *ServiceImpl) CreateUser() (auth.Token, error) {
	cred := "admin"
	//userFromDb, err := s.userRepo.FindByUserName(username)
	//if (err == nil && userFromDb.Username == newUser.Username) || (err != nil && err.Error() != "sql: no rows in result set") {
	//    log.Error(err.Error())
	//    return util.Token{}, domain.ErrDuplicateUser
	//}
	//log.Println(userFromDb.Username == newUser.Username)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(cred), bcrypt.DefaultCost)
	role := user.Admin
	accessToken, err := CreateToken(cred, role)
	if err != nil {
		log.Error(err.Error())
		return auth.Token{}, auth.ErrFailedCreateAccessToken
	}

	refreshToken, err := CreateTokenWithDuration(cred, role, time.Hour*24*365*10)
	if err != nil {
		log.Error(err.Error())
		return auth.Token{}, auth.ErrFailedUpdateAccessToken
	}

	err = s.userRepo.CreateUser(cred, passwordHash, role)
	if err != nil {
		log.Error(err.Error())
		return auth.Token{}, auth.ErrFailedSaveUser
	}

	return auth.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *ServiceImpl) Login(userDto auth.Credentials) (auth.Token, error) {
	userFromDb, err := s.userRepo.FindByUserName(userDto.Username)
	if (err == nil && userFromDb.Username != userDto.Username) || err != nil {
		if err != nil {
			log.Error(err.Error())
		}
		return auth.Token{}, auth.ErrUserNotExist
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.Password, []byte(userDto.Password))
	if err != nil {
		log.Error(err.Error())
		return auth.Token{}, auth.ErrWrongPassword
	}

	role := user.Admin
	accessToken, err := CreateTokenWithDuration(userFromDb.Username, role, time.Hour*24*365*10)
	if err != nil {
		log.Error(err.Error())
		return auth.Token{}, auth.ErrFailedCreateAccessToken
	}

	refreshToken, err := CreateTokenWithDuration(userFromDb.Username, role, time.Hour*24*365*10)
	if err != nil {
		log.Error(err.Error())
		return auth.Token{}, auth.ErrFailedUpdateAccessToken
	}

	return auth.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *ServiceImpl) ValidateToken(authorizationHeader string) (string, error) {
	tokenString, ok := ParseToken(authorizationHeader)
	if !ok {
		return "", auth.ErrInvalidAuthorizationHeader
	}

	token, ok := CheckToken(tokenString)
	log.Println("bearerToken " + tokenString)

	if !ok {
		return "", auth.ErrInvalidAccessToken
	}

	username, _, ok := s.parsePayload(token)
	if !ok {
		return "", auth.ErrInvalidAccessToken
	}

	_, err := s.userRepo.FindByUserName(username)
	if err != nil {
		return "", auth.ErrUserNotExist
	}

	return username, nil
}

func (s *ServiceImpl) RefreshToken(refreshTokenDto auth.RefreshTokenDto) (auth.Token, error) {
	username := refreshTokenDto.Username
	userFromDb, err := s.userRepo.FindByUserName(username)
	if (err == nil && userFromDb.Username != username) || err != nil {
		return auth.Token{}, auth.ErrUserNotExist
	}

	//if userFromDb.RefreshToken != refreshTokenDto.Token {
	//    return auth.Token{}, auth.ErrInvalidRefreshToken
	//}

	accessToken, err := CreateToken(username, userFromDb.Role)
	if err != nil {
		return auth.Token{}, auth.ErrFailedCreateAccessToken
	}
	refreshToken, err := CreateTokenWithDuration(userFromDb.Username, user.Admin, time.Hour*24*365*10)
	if err != nil {
		log.Error(err.Error())
		return auth.Token{}, auth.ErrFailedUpdateAccessToken
	}
	//
	//ok := s.userRepo.UpdateAccessToken(username, accessToken)
	//if !ok {
	//    return auth.Token{}, auth.ErrFailedUpdateAccessToken
	//}

	return auth.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *ServiceImpl) parsePayload(token *jwt.Token) (string, string, bool) {
	var username string
	var userId string
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username, ok = claims["username"].(string)
		if !ok {
			return "unauthorized, username not exist", "", false
		}
		//userId, ok = claims["userId"].(string)
		//if !ok {
		//    return "unauthorized, userId not exist", "", false
		//}

		_, err := s.userRepo.FindByUserName(username)
		if err != nil {
			return "unauthorized, user not exists", "", false
		}
	}

	return username, userId, true
}

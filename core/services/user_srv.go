package services

import (
	"fmt"
	"structure-golang/common/logs"
	"structure-golang/core/models"
	"structure-golang/core/repositories"
	"structure-golang/utils"
)

type userSrv struct {
	log      logs.AppLog
	userRepo repositories.UserRepository
}

func NewUserService(log logs.AppLog, userRepo repositories.UserRepository) UserService {
	return userSrv{log, userRepo}
}

func (s userSrv) Signin(username, password string) (result string, err error) {
	s.log.Info(fmt.Sprintf(`{"user_id": "%v", "message": "%v"}`, "system", "signin"))

	// Get users repository
	res, err := s.userRepo.Get(models.RepoFilterUserModel{Username: username, Password: password})
	if err != nil {
		s.log.Error(fmt.Sprintf(`{"user_id":"%v","message": "%v"}`, "system", fmt.Sprintf("signin | %v", models.ErrUnexpected)))
		return result, utils.Err_Handler{
			Code:    400,
			Message: models.ErrUnexpected,
		}
	}

	result = res.UserID

	return result, nil
}

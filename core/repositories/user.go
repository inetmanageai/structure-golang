package repositories

import "structure-golang/core/models"

type UserRepository interface {
	Get(filter models.RepoFilterUserModel) (result models.RepoUserModel, err error)
}

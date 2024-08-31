package models

import "time"

type RepoUserModel struct {
	UserID     string    `bson:"user_id"`
	Username   string    `bson:"username"`
	Password   string    `bson:"password"`
	CreateDate time.Time `bson:"create_date"`
	UpdateDate time.Time `bson:"update_date"`
}

type RepoFilterUserModel struct {
	UserID   string `bson:"user_id,omitempty"`
	Username string `bson:"username,omitempty"`
	Password string `bson:"password,omitempty"`
}

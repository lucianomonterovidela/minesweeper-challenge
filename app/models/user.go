package models

import "time"

type User struct {
	UserName   string    `bson:"_id" json:"userName"`
	Password   string    `bson:"password" json:"password"`
	CreationAt time.Time `bson:"creation_at" json:"createAt"`
	UpdateAt   time.Time `bson:"update_at" json:"createAt"`
}

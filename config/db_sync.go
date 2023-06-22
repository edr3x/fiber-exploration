package config

import "github.com/edr3x/fiber-explore/model"

func DbSync() {
	DB.AutoMigrate(&model.User{})
}

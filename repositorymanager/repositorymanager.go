package repositorymanager

import (
	"github.com/ednesic/coursemanagement/cache"
	"github.com/ednesic/coursemanagement/storage"
)

var (
	Dal storage.DataAccessLayer
	Redis cache.RedisClient
)
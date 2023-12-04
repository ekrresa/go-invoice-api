package utils

import "github.com/ekrresa/invoice-api/pkg/models"

var apiKeyCache = make(map[string]models.UserWithoutPassword)

func CacheUser(key string, user models.UserWithoutPassword) {
	apiKeyCache[key] = user
}

func GetUserFromCache(key string) (*models.UserWithoutPassword, bool) {
	var user, found = apiKeyCache[key]

	return &user, found
}

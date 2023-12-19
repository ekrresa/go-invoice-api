package helpers

import "github.com/ekrresa/invoice-api/pkg/models"

var apiKeyCache = make(map[string]models.User)

func CacheUser(key string, user models.User) {
	apiKeyCache[key] = user
}

func GetUserFromCache(key string) (models.User, bool) {
	var user, found = apiKeyCache[key]

	return user, found
}

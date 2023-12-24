package helpers

import (
	"time"

	"github.com/ekrresa/invoice-api/pkg/models"
	"github.com/patrickmn/go-cache"
)

var c = cache.New(1*time.Hour, 0)

func CacheUser(key string, user models.User) {
	c.Set(key, user, cache.NoExpiration)
}

func GetUserFromCache(key string) (models.User, bool) {
	var user models.User
	var value, found = c.Get(key)

	if found {
		user = value.(models.User)
	}

	return user, found
}

func RemoveUserFromCache(key string) {
	c.Delete(key)
}

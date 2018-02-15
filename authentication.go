package tastypie

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-pg/pg"
	"github.com/tomi77/go-pg-django/auth"
)

type Authentication interface {
	IsAuthenticated(r *http.Request) (*auth.User, error)
}

type ApiKeyAuthentication struct {
	DB *pg.DB
}

// Extract credentials from request
func (a *ApiKeyAuthentication) ExtractCredentials(r *http.Request) (string, string, error) {
	authentication := r.Header.Get("Authorization")
	if strings.Index(strings.ToLower(authentication), "apikey ") == -1 {
		return "", "", errors.New("Invalid Authorization header")
	}
	data := strings.Split(authentication, " ")
	data = strings.Split(data[1], ":")
	if len(data) != 2 {
		return "", "", errors.New("Invalid Authorization header")
	}
	return data[0], data[1], nil
}

// Check user is authenticated.
// Returns auth.User object
func (a *ApiKeyAuthentication) IsAuthenticated(r *http.Request) (*auth.User, error) {
	username, key, err := a.ExtractCredentials(r)
	if err != nil {
		return nil, err
	}

	var apikey ApiKey
    err = a.DB.Model(&apikey).
		Column("User").
        Where("api_key.key = ? and username = ?", key, username).
        Select()
	if err != nil {
		return nil, errors.New("Invalid Authentication credentials")
	}
	if apikey.User.IsActive == false {
		return nil, errors.New("User is inactive")
	}
	return apikey.User, nil
}

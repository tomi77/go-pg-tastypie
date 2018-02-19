package tastypie

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-pg/pg"
	"github.com/tomi77/go-pg-django/auth"
)

// Authentication types
const (
	TYPE_APIKEY = 1 << iota // ApiKey Authentication
)

const APIKEY_AUTHORIZATION_HEADER = "Authorization"

// Predefined errors
var (
	InvalidAuthenticationType = errors.New("Invalid authentication type")
	InvalidAuthorizationHeader = errors.New("Invalid Authorization header")
	InvalidCredentials = errors.New("Invalid credentials")
	ErrorUserIsInactive = errors.New("User is inactive")
)

// Interface to various authentication backends
type Authentication interface {
	// Extract credentials from request
	ExtractCredentials(r *http.Request) (string, string, error)

	// Checks if user is authenticated and return it
	IsAuthenticated(r *http.Request) (*auth.User, error)
}

// Get Authentication object based on type and init it
func GetAuthentication(authenticationType int, db *pg.DB) (Authentication, error) {
	switch authenticationType {
	case TYPE_APIKEY:
		return ApiKeyAuthentication{DB: db}, nil
	default:
		return nil, InvalidAuthenticationType
	}
}

type ApiKeyAuthentication struct {
	// Handler do database connection
	DB *pg.DB
}

// Extract username and tastypie apikey from request
func (a ApiKeyAuthentication) ExtractCredentials(r *http.Request) (string, string, error) {
	authentication := r.Header.Get(APIKEY_AUTHORIZATION_HEADER)
	if strings.Index(strings.ToLower(authentication), "apikey ") == -1 {
		return "", "", InvalidAuthorizationHeader
	}
	data := strings.Split(authentication, " ")
	data = strings.Split(data[1], ":")
	if len(data) != 2 {
		return "", "", InvalidAuthorizationHeader
	}
	return data[0], data[1], nil
}

// Checks if user is authenticated and return it
func (a ApiKeyAuthentication) IsAuthenticated(r *http.Request) (*auth.User, error) {
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
		return nil, InvalidCredentials
	}
	if apikey.User.IsActive == false {
		return nil, ErrorUserIsInactive
	}
	return apikey.User, nil
}

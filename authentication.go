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
	TypeAPIKey = 1 << iota // APIKey Authentication
)

// APIKeyAuthorizationHeader is a key in HTTP header where is stored authorization key
const APIKeyAuthorizationHeader = "Authorization"

// Predefined errors
var (
	ErrorInvalidAuthenticationType  = errors.New("Invalid authentication type")
	ErrorInvalidAuthorizationHeader = errors.New("Invalid Authorization header")
	ErrorInvalidCredentials         = errors.New("Invalid credentials")
	ErrorUserIsInactive             = errors.New("User is inactive")
)

// Authentication is a interface to various authentication backends
type Authentication interface {
	// Extract credentials from request
	ExtractCredentials(r *http.Request) (string, string, error)

	// Checks if user is authenticated and return it
	IsAuthenticated(r *http.Request) (*auth.User, error)
}

// GetAuthentication gets Authentication object based on type and initiates it
func GetAuthentication(authenticationType int, db *pg.DB) (Authentication, error) {
	switch authenticationType {
	case TypeAPIKey:
		return APIKeyAuthentication{DB: db}, nil
	default:
		return nil, ErrorInvalidAuthenticationType
	}
}

// APIKeyAuthentication represents authentication based on API Key
type APIKeyAuthentication struct {
	// Handler do database connection
	DB *pg.DB
}

// ExtractCredentials returns username and tastypie apikey extracted from request
func (a APIKeyAuthentication) ExtractCredentials(r *http.Request) (string, string, error) {
	authentication := r.Header.Get(APIKeyAuthorizationHeader)
	if strings.Index(strings.ToLower(authentication), "apikey ") == -1 {
		return "", "", ErrorInvalidAuthorizationHeader
	}
	data := strings.Split(authentication, " ")
	data = strings.Split(data[1], ":")
	if len(data) != 2 {
		return "", "", ErrorInvalidAuthorizationHeader
	}
	return data[0], data[1], nil
}

// IsAuthenticated checks if user is authenticated and return it
func (a APIKeyAuthentication) IsAuthenticated(r *http.Request) (*auth.User, error) {
	username, key, err := a.ExtractCredentials(r)
	if err != nil {
		return nil, err
	}

	var apikey APIKey
	err = a.DB.Model(&apikey).
		Column("User").
		Where("api_key.key = ? and username = ?", key, username).
		Select()
	if err != nil {
		return nil, ErrorInvalidCredentials
	}
	if apikey.User.IsActive == false {
		return nil, ErrorUserIsInactive
	}
	return apikey.User, nil
}

package tastypie

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetAuthentication(t *testing.T) {
	t.Run("ApiKey", func(t *testing.T) {
		authentication, err := GetAuthentication(TYPE_APIKEY, nil)
		if authentication == nil {
			t.Error("should return Authentication object")
		}
		if err != nil {
			t.Error("shouldn't return error")
		}
		_, ok := authentication.(ApiKeyAuthentication)
		if !ok {
			t.Error("should return ApiKeyAuthentication object")
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		authentication, err := GetAuthentication(-1, nil)
		if authentication != nil {
			t.Error("shouldn't return Authentication object")
		}
		if err == nil {
			t.Error("should return error")
		}
	})
}

func TestApiKeyAuthenticationExtractCredentials(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "http://www.example.com/", nil)
	authentication := ApiKeyAuthentication{}
	t.Run("NoHeader", func(t *testing.T) {
		r.Header.Del("Authorization")
		username, key, err := authentication.ExtractCredentials(r)
		if username != "" && key != "" {
			t.Error("username and key must are a empty strings")
		}
		if err == nil {
			t.Error("error must be set")
		}
	})
	t.Run("InvalidType", func(t *testing.T) {
		r.Header.Set("Authorization", "ApiKeys admin:qaz123")
		username, key, err := authentication.ExtractCredentials(r)
		if username != "" && key != "" {
			t.Error("username and key must are a empty strings")
		}
		if err == nil {
			t.Error("error must be set")
		}
	})
	t.Run("InvalidData", func(t *testing.T) {
		r.Header.Set("Authorization", "ApiKey admin-qaz123")
		username, key, err := authentication.ExtractCredentials(r)
		if username != "" && key != "" {
			t.Error("username and key must are a empty strings")
		}
		if err == nil {
			t.Error("error must be set")
		}
	})
	t.Run("Valid", func(t *testing.T) {
		r.Header.Set("Authorization", "ApiKey admin:qaz123")
		username, key, err := authentication.ExtractCredentials(r)
		if username != "admin" {
			t.Error("username must set to 'admin'")
		}
		if key != "qaz123" {
			t.Error("key must be set to 'qaz123'")
		}
		if err != nil {
			t.Error("error must not be set")
		}
	})
}

func ExampleApiKeyAuthentication_ExtractCredentials() {
	r, _ := http.NewRequest(http.MethodGet, "http://www.example.com/", nil)
	r.Header.Set("Authorization", "ApiKey admin:qaz123")
	authentication := ApiKeyAuthentication{}

	username, key, _ := authentication.ExtractCredentials(r)
	fmt.Println(username, key)
	// Output:
	// admin qaz123
}

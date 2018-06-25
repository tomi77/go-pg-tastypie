package tastypie

import (
	"fmt"
	"time"

	"github.com/tomi77/go-pg-django/auth"
)

func ExampleAPIAccess_String() {
	row := APIAccess{
		ID:            1,
		Identifier:    "test",
		URL:           "https://www.example.com/",
		RequestMethod: "GET",
		Accessed:      123456,
	}

	fmt.Println(row)
	// Output:
	// test @ 123456
}

func ExampleAPIKey_String() {
	user := auth.User{
		Username: "admin",
	}
	row := APIKey{
		ID:      1,
		User:    &user,
		Key:     "qaz123",
		Created: time.Now(),
	}

	fmt.Println(row)
	// Output:
	// qaz123 for admin
}

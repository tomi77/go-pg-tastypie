package tastypie

import (
	"fmt"
	"time"

	"github.com/tomi77/go-pg-django/auth"
)

func ExampleApiAccess_String() {
	row := ApiAccess{
		Id: 1,
		Identifier: "test",
		Url: "https://www.example.com/",
		RequestMethod: "GET",
		Accessed: 123456,
	}

	fmt.Println(row)
	// Output:
	// test @ 123456
}

func ExampleApiKey_String() {
	user := auth.User{
		Username: "admin",
	}
	row := ApiKey{
		Id: 1,
		User: &user,
		Key: "qaz123",
		Created: time.Now(),
	}

	fmt.Println(row)
	// Output:
	// qaz123 for admin
}

package tastypie

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/tomi77/go-pg-django/auth"
)

// APIAccess represents tastypie_apiaccess table
type APIAccess struct {
	TableName string `sql:"tastypie_apiaccess"`

	ID            uint16
	Identifier    string `sql:"type:varchar(255),notnull"`
	URL           string `sql:"type:varchar(255),notnull,default:''"`
	RequestMethod string `sql:"type:varchar(10),notnull,default:''"`
	Accessed      int64  `sql:",notnull"`
}

func (a APIAccess) String() string {
	return fmt.Sprintf("%s @ %d", a.Identifier, a.Accessed)
}

// BeforeInsert hook update "Accessed" field
func (a *APIAccess) BeforeInsert(db orm.DB) error {
	a.Accessed = time.Now().Unix()
	return nil
}

// BeforeUpdate hook update "Accessed" field
func (a *APIAccess) BeforeUpdate(db orm.DB) error {
	a.Accessed = time.Now().Unix()
	return nil
}

// APIKey represents tastypie_apikey table
type APIKey struct {
	TableName string `sql:"tastypie_apikey"`

	ID      uint16
	UserID  uint16 `sql:",notnull" pg:",fk:auth.User"`
	User    *auth.User
	Key     string    `sql:"type:varchar(128),notnull,default:''"`
	Created time.Time `sql:",notnull"`
}

func (a APIKey) String() string {
	return fmt.Sprintf("%s for %s", a.Key, a.User)
}

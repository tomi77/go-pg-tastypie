package tastypie

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/tomi77/go-pg-django/auth"
)

// "tastypie_apiaccess" table
type ApiAccess struct {
	TableName string `sql:"tastypie_apiaccess"`

	Id            uint16
	Identifier    string `sql:"type:varchar(255),notnull"`
	Url           string `sql:"type:varchar(255),notnull,default:''"`
	RequestMethod string `sql:"type:varchar(10),notnull,default:''"`
	Accessed      int64  `sql:",notnull"`
}

func (a ApiAccess) String() string {
	return fmt.Sprintf("%s @ %d", a.Identifier, a.Accessed)
}

// Update "Accessed" field
func (a *ApiAccess) BeforeInsert(db orm.DB) error {
	a.Accessed = time.Now().Unix()
	return nil
}

// Update "Accessed" field
func (a *ApiAccess) BeforeUpdate(db orm.DB) error {
	a.Accessed = time.Now().Unix()
	return nil
}

// "tastypie_apikey" table
type ApiKey struct {
	TableName string `sql:"tastypie_apikey"`

	Id      uint16
	UserId  uint16 `sql:",notnull" pg:",fk:auth.User"`
	User    *auth.User
	Key     string    `sql:"type:varchar(128),notnull,default:''"`
	Created time.Time `sql:",notnull"`
}

func (a ApiKey) String() string {
	return fmt.Sprintf("%s for %s", a.Key, a.User)
}

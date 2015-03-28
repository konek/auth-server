package models

import (
	"gopkg.in/mgo.v2/bson"

	"go.konek.io/auth-server/tools"
	"go.konek.io/mgo"
)

// User is the model for users
type User struct {
	ID        bson.ObjectId          `bson:"_id" json:"uid"`
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Domains   []string               `json:"domains"`
	Enable    bool                   `json:"enable"`
	Variables map[string]interface{} `json:"variables"`
	Salt      string                 `json:"-"`
}

// IDFromHex fill the ID from a string hex
func (u *User) IDFromHex(hex string) {
	u.ID = bson.ObjectIdHex(hex)
}

// Create a new user in database.
// (generates Salt)
func (u *User) Create(db *mgo.DbQueue) (bson.ObjectId, error) {
	var n int
	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error
		q := db.C("users").Find(bson.M{"username": u.Username})
		n, e = q.Count()
		ec <- e
	})
	if err != nil {
		return "", err
	}
	if n != 0 {
		return "", tools.NewError(nil, 409, "duplicate: user already exists")
	}

	u.ID = bson.NewObjectId()
	u.Salt, err = tools.GenSalt(12)
	if err != nil {
		return "", err
	}
	u.Password = tools.PasswordHash(u.Username, u.Password, u.Salt)
	err = db.Push(func (db *mgo.Database, ec chan error) {
		ec <- db.C("users").Insert(u)
	})
	return u.ID, err
}

// Delete user from database
func (u User) Delete(db *mgo.DbQueue) error {
	var n int

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error
		q := db.C("users").FindId(u.ID)
		n, e = q.Count()
		ec <- e
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: user does not exist")
	}
	err = db.Push(func (db *mgo.Database, ec chan error) {
		ec <- db.C("users").RemoveId(u.ID)
	})
	return err
}

// Update user in database. (update salt and password if needed)
func (u User) Update(db *mgo.DbQueue) error {
	var n int
	var user User

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error
		q := db.C("users").FindId(u.ID)
		n, e = q.Count()
		if e != nil {
			ec <- e
			return
		}
		ec <- q.One(&user)
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: user does not exist")
	}

	if u.Username != "" {
		if u.Password == "" {
			return tools.NewError(nil, 400, "bad request: need password to update username")
		}
		user.Username = u.Username
		user.Salt, err = tools.GenSalt(12)
		if err != nil {
			return err
		}
		user.Password = tools.PasswordHash(user.Username, u.Password, user.Salt)
	}
	if u.Password != "" {
		user.Salt, err = tools.GenSalt(12)
		if err != nil {
			return err
		}
		user.Password = tools.PasswordHash(user.Username, u.Password, user.Salt)
	}
	if u.Domains != nil {
		user.Domains = u.Domains
	}
	if u.Variables != nil {
		user.Variables = u.Variables
	}
	user.Enable = u.Enable

	err = db.Push(func (db *mgo.Database, ec chan error) {
		ec <- db.C("users").UpdateId(u.ID, user)
	})
	return err
}

// Get user from database
func (u *User) Get(db *mgo.DbQueue) error {
	var n int
	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error
		q := db.C("users").FindId(u.ID)
		n, e = q.Count()
		if e != nil {
			ec <- e
			return
		}
		ec <- q.One(u)
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: user does not exist")
	}
	return nil
}

// GetByUsername fill the User given its username
func (u *User) GetByUsername(db *mgo.DbQueue) error {
	var n int

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error
		q := db.C("users").Find(bson.M{"username": u.Username})
		n, e = q.Count()
		if e != nil {
			ec <- e
			return
		}
		ec <- q.One(u)
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: user does not exist")
	}
	return nil
}

// Check if username and password matches
func (u *User) Check(db *mgo.DbQueue) (bool, error) {
	password := u.Password

	err := u.GetByUsername(db)
	if err != nil {
		return false, err
	}
	h := tools.PasswordHash(u.Username, password, u.Salt)
	if u.Password == h {
		return true, nil
	}
	return false, nil
}

// CheckDomain Check if the user has access to domain
func (u User) CheckDomain(domain string) bool {
	return tools.CheckDomains(u.Domains, domain)
}

// ListUsers return a list of all users
func ListUsers(db *mgo.DbQueue) ([]User, error) {
	var list []User
	var n int

	err := db.Push(func (db *mgo.Database, ec chan error) {
		var e error

		q := db.C("users").Find(nil)
		n, e = q.Count()
		if e != nil {
			ec <- e
			return
		}
		list = make([]User, n)
		ec <- q.All(&list)
	})

	if err != nil {
		return nil, err
	}
	return list, nil
}

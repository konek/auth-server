package models

import (
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/konek/mgo.v1"
	"github.com/asaskevich/govalidator"
	"gopkg.in/konek/auth-server.v1/tools"
)

// User is the model for users
type User struct {
	ID        bson.ObjectId          `bson:"_id" json:"_id"`
	Username  string                 `json:"username"`
	Password  string                 `json:"password"`
	Domains   []string               `json:"domains"`
	Enable    bool                   `json:"enable"`
	Variables map[string]interface{} `json:"variables"`
	Salt      string                 `json:"salt"`
}

// IDFromHex fill the ID from a string hex
func (u *User) IDFromHex(hex string) {
	u.ID = bson.ObjectIdHex(hex)
}

// Create a new user in database.
// (generates Salt)
func (u *User) Create(db *mgo.DbQueue) (bson.ObjectId, error) {
	n, err := db.Count("users", bson.M{"username": u.Username})
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
	err = db.Insert("users", u)
	return u.ID, err
}

// Delete user from database
func (u User) Delete(db *mgo.DbQueue) error {
	n, err := db.CountID("users", u.ID)
	if err != nil {
		return err
	}
	if n == 0 {
		return tools.NewError(nil, 404, "not found: user does not exist")
	}
	err = db.RemoveID("users", u.ID)
	return err
}

// Update user in database. (update salt and password if needed)
func (u User) Update(db *mgo.DbQueue) error {
	var user User

	err := db.FindOneID("users", &user, u.ID)
	if err != nil {
		return err
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

	err = db.Push(func(db *mgo.Database, ec chan error) {
		ec <- db.C("users").UpdateId(u.ID, user)
	})
	return err
}

// Get user from database
func (u *User) Get(db *mgo.DbQueue) error {
	err := db.FindOneID("users", u, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetByUsername fill the User given its username
func (u *User) GetByUsername(db *mgo.DbQueue) error {
	err := db.FindOne("users", u, bson.M{"username": u.Username})
	if err != nil {
		return err
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

func (u User) Login(username, password, domain string, lifespan int64, db *mgo.DbQueue) (Session, error) {
	var s Session
	var err error

	u.Username, err = govalidator.NormalizeEmail(username)
	if err != nil {
		return s, tools.NewError(nil, 400, "bad request: username must be a valid email")
	}
	u.Password = password
	if password == "" {
		return s, tools.NewError(nil, 400, "bad request: password is missing")
	}

	ok, err := u.Check(db)
	if err != nil {
		return s, err
	}
	if ok == false {
		return s, tools.NewError(nil, 403, "forbidden: invalid user or password")
	}
	if u.Enable == false {
		return s, tools.NewError(nil, 403, "forbidden: user is diabled")
	}
	ok = u.CheckDomain(domain)
	if ok == false {
		return s, tools.NewError(nil, 403, "forbidden: restricted domain")
	}

	s.UserID = u.ID
	s.Domain = domain
	_, err = s.Create(db, lifespan)
	if err != nil {
		return s, err
	}

	return s, nil
}

// ListUsers return a list of all users
func ListUsers(db *mgo.DbQueue) ([]User, error) {
	var list []User
	err := db.Find("users", &list, nil)

	if err != nil {
		return nil, err
	}
	return list, nil
}

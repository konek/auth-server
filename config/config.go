package config

import (
	"os"
	"strconv"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Conf is the global configuration struct, set via environment variables
//
// TODO : add a toml reader for configuration file
type Conf struct {
	Listen            string
	Public		  bool
	DbURL             string
	DbName            string
	SessionLifespan   int64
	PasswordMinLength int
	Root              string
	Etcd              string // TODO
	Node              string // TODO
	LogToDb           bool   // TODO
}

const (
	defaultListen            = "127.0.0.1:8080"
	defaultPublic		 = false
	defaultDbURL             = "localhost:27017"
	defaultDbName            = "auth"
	defaultSessionLifespan   = 60 * 60 * 24 // 24 hours
	defaultPasswordLinlength = 8
	defaultetcd              = ""
	defaultnode              = "auth"
	defaultLogToDb           = false
	defaultRoot              = ""
)

// Config sets the Conf struct from environment.
//
// Environment variables :
//   LISTEN: either :port or ip:port
//   MONGODB_URL: mongodb server's url
//   MONGODB_DATABASE: database's name for mongodb
//   SESSION_LIFESPAN: validity time limit of sessions (in seconds)
//   PASSWD_MIN_LENGTH: minimum password length
//   ETCD: (for future use)
//   NODE: (for future use)
//   LOG2DB: (for future use)
//   ROOT: root path for the router (such as /auth/v1)
func Config() Conf {
	var ret Conf
	ret.Listen = defaultListen
	ret.Public = defaultPublic
	ret.DbURL = defaultDbURL
	ret.DbName = defaultDbName
	ret.SessionLifespan = defaultSessionLifespan
	ret.PasswordMinLength = defaultPasswordLinlength
	ret.Etcd = defaultetcd
	ret.Node = defaultnode
	ret.LogToDb = defaultLogToDb
	ret.Root = defaultRoot

	if len(os.Args) >= 2 {
		conffile := os.Args[1]
		b, err := ioutil.ReadFile(conffile)
		if err != nil {
			panic(err)
			return ret
		}
		_, err = toml.Decode(string(b), &ret)
		if err != nil {
			panic(err)
			return ret
		}
	}
	if os.Getenv("LISTEN") != "" {
		ret.Listen = os.Getenv("LISTEN")
	} else if os.Getenv("PORT") != "" {
		ret.Listen = ":" + os.Getenv("PORT")
	}
	if os.Getenv("PUBLIC") == "true" {
		ret.Public = true
	}
	if os.Getenv("MONGO_URL") != "" {
		ret.DbURL = os.Getenv("MONGO_URL")
	}
	if os.Getenv("MONGODB_DATABASE") != "" {
		ret.DbName = os.Getenv("MONGODB_DATABASE")
	}
	if os.Getenv("SESSION_LIFESPAN") != "" {
		lifespan, err := strconv.ParseInt(os.Getenv("SESSION_LIFESPAN"), 10, 64)
		if err == nil {
			ret.SessionLifespan = lifespan
		}
	}
	if os.Getenv("PASSWD_MIN_LENGTH") != "" {
		length, err := strconv.ParseInt(os.Getenv("PASSWD_MIN_LENGTH"), 10, 32)
		if err == nil {
			ret.PasswordMinLength = int(length)
		}
	}
	if os.Getenv("ETCD") != "" {
		ret.Etcd = os.Getenv("ETCD")
	}
	if os.Getenv("NODE") != "" {
		ret.Node = os.Getenv("NODE")
	}
	if os.Getenv("LOG2DB") != "" {
		if os.Getenv("LOG2DB") == "false" {
			ret.LogToDb = false
		} else {
			ret.LogToDb = true
		}
	}
	if os.Getenv("ROOT") != "" {
		ret.Root = os.Getenv("ROOT")
	}
	return ret
}

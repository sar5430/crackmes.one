package database

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	Mongo *mgo.Session
	databases Info
)

// Type is the type of database from a Type* constant
type Type string

const (
	// TypeMongoDB is MongoDB
	TypeMongoDB Type = "MongoDB"
)

// Info contains the database configurations
type Info struct {
	// Database type
	Type Type
	// MongoDB info if used
	MongoDB MongoDBInfo
}

// MongoDBInfo is the details for the database connection
type MongoDBInfo struct {
	URL      string
	Database string
}

// Connect to the database
func Connect(d Info) {
	var err error

	// Store the config
	databases = d

	// Connect to MongoDB
	if Mongo, err = mgo.DialWithTimeout(d.MongoDB.URL, 5*time.Second); err != nil {
		log.Println("MongoDB Driver Error", err)
		return
	}

	// Prevents these errors: read tcp 127.0.0.1:27017: i/o timeout
	Mongo.SetSocketTimeout(1 * time.Second)

	// Check if is alive
	if err = Mongo.Ping(); err != nil {
		log.Println("Database Error", err)
	}
}

// CheckConnection returns true if MongoDB is available
func CheckConnection() bool {
	if Mongo == nil {
		Connect(databases)
	}

	if Mongo != nil {
		return true
	}

	return false
}

// ReadConfig returns the database information
func ReadConfig() Info {
	return databases
}

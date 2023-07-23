package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Ctx       context.Context
	Mongo     *mongo.Client
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

	ctx := context.TODO()

	// Connect to MongoDB
	Mongo, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Println("MongoDB Driver Error", err)
		return
	}
	if err = Mongo.Ping(ctx, readpref.Primary()); err != nil {
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

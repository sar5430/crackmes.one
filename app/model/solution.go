package model

import (
	"time"

	"github.com/sar5430/crackmes.one/app/shared/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// *****************************************************************************
// Crackme
// *****************************************************************************

// Crackme table contains the information for each note
type Solution struct {
	ObjectId  primitive.ObjectID `bson:"_id,omitempty"`
	HexId     string             `bson:"hexid,omitempty"`
	Info      string             `bson:"info"`
	CrackmeId primitive.ObjectID `bson:"crackmeid,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	Author    string             `bson:"author,omitempty"`
	Visible   bool               `bson:"visible"`
	Deleted   bool               `bson:"deleted"`
}

type SolutionExtended struct {
	Solution      *Solution
	Crackmeshexid string
	Crackmename   string
}

func CountSolutions() (int, error) {
	var err error
	var nb int64
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")
		nb, err = collection.CountDocuments(database.Ctx, bson.D{})
	} else {
		err = ErrUnavailable
	}

	return int(nb), standardizeError(err)
}

func CountSolutionsByUser(username string) (int, error) {
	var err error
	var nb int64
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")
		nb, err = collection.CountDocuments(database.Ctx, bson.M{"author": username, "visible": true})
	} else {
		err = ErrUnavailable
	}
	return int(nb), standardizeError(err)
}

func CountSolutionsByCrackme(crackmehexid string) (int, error) {
	var err error
	var nb int64
	var objid primitive.ObjectID
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")
		objid, err = primitive.ObjectIDFromHex(crackmehexid)
		nb, err = collection.CountDocuments(database.Ctx, bson.M{"crackmeid": objid, "visible": true})
	} else {
		err = ErrUnavailable
	}
	return int(nb), standardizeError(err)
}

func SolutionByHexId(hexid string) (Solution, error) {
	var err error

	var result Solution
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")

		// Validate the object id
		err = collection.FindOne(database.Ctx, bson.M{"hexid": hexid, "visible": true}).Decode(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func SolutionsByUser(username string) ([]Solution, error) {
	var err error

	var result []Solution
	var cursor *mongo.Cursor

	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")
		opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(50)

		// Validate the object id
		cursor, err = collection.Find(database.Ctx, bson.M{"author": username, "visible": true}, opts)
		err = cursor.All(database.Ctx, &result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func SolutionsByUserAndCrackMe(username, crackmehexid string) (Solution, error) {
	var err error

	var result Solution
	crackme, _ := CrackmeByHexId(crackmehexid)
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")

		// Validate the object id
		err = collection.FindOne(database.Ctx, bson.M{"crackmeid": crackme.ObjectId, "author": username}).Decode(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func SolutionsByCrackme(crackme primitive.ObjectID) ([]Solution, error) {
	var err error

	var result []Solution
	var cursor *mongo.Cursor

	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")

		// Validate the object id
		cursor, err = collection.Find(database.Ctx, bson.M{"crackmeid": crackme, "visible": true})
		err = cursor.All(database.Ctx, &result)

	} else {
		err = ErrUnavailable
	}
	return result, err
}

// NoteCreate creates a note
func SolutionCreate(info, username, crackmehexid string) error {
	var err error
	crackme, _ := CrackmeByHexId(crackmehexid)

	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("solution")
		objId := primitive.NewObjectID()
		solution := &Solution{
			ObjectId:  objId,
			HexId:     objId.Hex(),
			Info:      info,
			CrackmeId: crackme.ObjectId,
			CreatedAt: time.Now(),
			Author:    username,
			Visible:   false,
			Deleted:   false,
		}
		_, err = collection.InsertOne(database.Ctx, solution)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

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
type Crackme struct {
	ObjectId    primitive.ObjectID `bson:"_id,omitempty"`
	HexId       string             `bson:"hexid,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Info        string             `bson:"info,omitempty"`
	Lang        string             `bson:"lang,omitempty"`
	Arch        string             `bson:"arch,omitempty"`
	Author      string             `bson:"author,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	Visible     bool               `bson:"visible"`
	Deleted     bool               `bson:"deleted"`
	Difficulty  float64            `bson:"difficulty"`
	Quality     float64            `bson:"quality"`
	NbSolutions int                // Not present in the database! Just for rendering
	NbComments  int                // Not present in the database! Just for rendering
	Platform    string             `bson:"platform,omitempty"`
}

func CountCrackmes() (int, error) {
	var err error
	var nb int64
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")
		nb, err = collection.CountDocuments(database.Ctx, bson.M{"visible": true})
	} else {
		err = ErrUnavailable
	}

	return int(nb), standardizeError(err)
}

func CountCrackmesByUser(username string) (int, error) {
	var err error
	var nb int64
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")
		nb, err = collection.CountDocuments(database.Ctx, bson.M{"author": username, "visible": true})
	} else {
		err = ErrUnavailable
	}
	return int(nb), standardizeError(err)
}

func GetAllCrackmes() ([]Crackme, error) {
	var err error
	var result []Crackme
	var cursor *mongo.Cursor

	if database.CheckConnection() {
		// Create a copy of mongo
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")

		// Validate the object id
		cursor, err = collection.Find(database.Ctx, bson.M{})
		err = cursor.All(database.Ctx, &result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func CrackmeSetFloat(hexid, champ string, nb float64) error {
	var err error
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")

		// Validate the object id
		_, err = collection.UpdateOne(database.Ctx, bson.M{"hexid": hexid}, bson.M{"$set": bson.M{champ: float64(nb)}})
	} else {
		err = ErrUnavailable
	}
	return err
}

func SearchCrackme(name, author, lang, arch, platform string, difficulty_min, difficulty_max, quality_min, quality_max int) ([]Crackme, error) {
	var err error
	var result []Crackme
	var cursor *mongo.Cursor

	if database.CheckConnection() {
		// Create a copy of mongo
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")
		opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(150)

		// Validate the object id
		cursor, err = collection.Find(database.Ctx,
			bson.D{
				{"name", primitive.Regex{Pattern: name, Options: "i"}},
				{"lang", primitive.Regex{Pattern: lang, Options: "i"}},
				{"arch", primitive.Regex{Pattern: arch, Options: "i"}},
				{"difficulty", bson.M{"$gte": difficulty_min, "$lte": difficulty_max}},
				{"quality", bson.M{"$gte": quality_min, "$lte": quality_max}},
				{"author", primitive.Regex{Pattern: author, Options: "i"}},
				{"visible", true},
				{"platform", primitive.Regex{Pattern: platform, Options: "i"}},
			}, opts)

		err = cursor.All(database.Ctx, &result)

	} else {
		err = ErrUnavailable
	}
	return result, err
}

func LastCrackMes(page int) ([]Crackme, error) {
	var err error
	var result []Crackme
	var cursor *mongo.Cursor

	if database.CheckConnection() {
		// Create a copy of mongo
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")
		opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(50).SetSkip(int64((page - 1) * 50))

		// Validate the object id
		cursor, err = collection.Find(database.Ctx, bson.M{"visible": true}, opts)
		err = cursor.All(database.Ctx, &result)

	} else {
		err = ErrUnavailable
	}
	return result, err
}

func CrackmeByHexId(hexid string) (Crackme, error) {
	var err error

	var result Crackme
	if database.CheckConnection() {
		// Create a copy of mongo
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")

		// Validate the object id
		err = collection.FindOne(database.Ctx, bson.M{"hexid": hexid, "visible": true}).Decode(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func CrackmesByUser(username string) ([]Crackme, error) {
	var err error
	var cursor *mongo.Cursor
	var result []Crackme
	if database.CheckConnection() {
		// Create a copy of mongo
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")
		opts := options.Find().SetSort(bson.D{{"created_at", -1}})

		// Validate the object id
		cursor, err = collection.Find(database.Ctx, bson.M{"author": username, "visible": true}, opts)
		err = cursor.All(database.Ctx, &result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func CrackmeByUserAndName(username, name string, visible bool) (Crackme, error) {
	var err error

	var result Crackme

	if database.CheckConnection() {
		// Create a copy of mongo
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")

		// Validate the object id
		err = collection.FindOne(database.Ctx, bson.M{"name": name, "author": username, "visible": visible}).Decode(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

// NoteCreate creates a note
func CrackmeCreate(name, info, username, lang, arch, platform string) error {
	var err error

	if database.CheckConnection() {
		objId := primitive.NewObjectID()
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("crackme")
		crackme := &Crackme{
			ObjectId:  objId,
			HexId:     objId.Hex(),
			Name:      name,
			Info:      info,
			Lang:      lang,
			Arch:      arch,
			Author:    username,
			CreatedAt: time.Now(),
			Visible:   false,
			Deleted:   false,
			Platform:  platform,
		}
		_, err = collection.InsertOne(database.Ctx, crackme)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

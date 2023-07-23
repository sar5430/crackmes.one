package model

import (
	"time"

	"github.com/sar5430/crackmes.one/app/shared/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RatingQuality struct {
	ObjectId     primitive.ObjectID `bson:"_id,omitempty"`
	Author       string             `bson:"author,omitempty"`
	CrackMeHexId string             `bson:"crackmehexid,omitempty"`
	Rating       int                `bson:"rating"`
	CreatedAt    time.Time          `bson:"created_at"`
	Visible      bool               `bson:"visible"`
	Deleted      bool               `bson:"deleted"`
}

func IsAlreadyRatedQuality(username, crackmehexid string) (bool, error) {
	var err error
	var nb int64
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("rating_quality")
		nb, err = collection.CountDocuments(database.Ctx, bson.M{"author": username, "crackmehexid": crackmehexid})
	} else {
		err = ErrUnavailable
	}

	return nb != 0, err
}

func RatingQualityByCrackme(crackmehexid string) ([]RatingQuality, error) {
	var err error
	var result []RatingQuality
	var cursor *mongo.Cursor
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("rating_quality")

		// Validate the object id
		cursor, err = collection.Find(database.Ctx, bson.M{"crackmehexid": crackmehexid})
		err = cursor.All(database.Ctx, &result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func RatingQualitySetRating(username, crackmehexid string, rating int) error {
	var err error
	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("rating_quality")

		// Validate the object id
		_, err = collection.UpdateOne(database.Ctx, bson.M{"crackmehexid": crackmehexid, "author": username}, bson.M{"$set": bson.M{"rating": rating}})
	} else {
		err = ErrUnavailable
	}
	return err
}

func RatingQualityCreate(username, crackmehexid string, rating int) error {
	var err error

	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("rating_quality")
		objId := primitive.NewObjectID()
		rating_quality := &RatingQuality{
			ObjectId:     objId,
			Rating:       rating,
			Author:       username,
			CrackMeHexId: crackmehexid,
			CreatedAt:    time.Now(),
			Visible:      true,
			Deleted:      false,
		}
		_, err = collection.InsertOne(database.Ctx, rating_quality)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

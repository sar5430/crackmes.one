package model

import (
	"app/shared/database"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Comment
// *****************************************************************************

// Comment table contains the information for each note
type RatingDifficulty struct {
	ObjectId     bson.ObjectId `bson:"_id,omitempty"`
	Author       string        `bson:"author,omitempty"`
	CrackMeHexId string        `bson:"crackmehexid,omitempty"`
    Rating       int           `bson:"rating"`
	CreatedAt    time.Time     `bson:"created_at"`
	Visible      bool          `bson:"visible"`
	Deleted      bool          `bson:"deleted"`
}

func IsAlreadyRatedDifficulty(username, crackmehexid string) (bool, error) {
	var err error
	var nb int
    var res bool
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_difficulty")
        nb, err = c.Find(bson.M{"author": username, "crackmehexid": crackmehexid}).Count()
	} else {
		err = ErrUnavailable
	}

    if nb != 0 {
        res = true
    } else {
        res = false
    }
	return res, err
}

func RatingDifficultyByCrackme(crackmehexid string) ([]RatingDifficulty, error) {
	var err error
	var result []RatingDifficulty
	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_difficulty")

		// Validate the object id
		err = c.Find(bson.M{"crackmehexid": crackmehexid}).All(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func RatingDifficultySetRating(username, crackmehexid string, rating int) error {
	var err error
	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_difficulty")

		// Validate the object id
        err = c.Update(bson.M{"crackmehexid": crackmehexid, "author": username}, bson.M{"$set": bson.M{"rating": rating}})
	} else {
		err = ErrUnavailable
	}
	return err
}

func RatingDifficultyCreate(username, crackmehexid string, rating int) error {
	var err error

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_difficulty")
		objId := bson.NewObjectId()
		rating_difficulty := &RatingDifficulty {
			ObjectId:     objId,
			Rating:       rating,
			Author:       username,
			CrackMeHexId: crackmehexid,
			CreatedAt:    time.Now(),
			Visible:      true,
			Deleted:      false,
		}
		err = c.Insert(rating_difficulty)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

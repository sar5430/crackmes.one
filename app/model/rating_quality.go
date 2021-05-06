package model

import (
    "github.com/sar5430/crackmes.one/app/shared/database"
    "time"

    "gopkg.in/mgo.v2/bson"
)

type RatingQuality struct {
    ObjectId     bson.ObjectId `bson:"_id,omitempty"`
    Author       string        `bson:"author,omitempty"`
    CrackMeHexId string        `bson:"crackmehexid,omitempty"`
    Rating       int           `bson:"rating"`
    CreatedAt    time.Time     `bson:"created_at"`
    Visible      bool          `bson:"visible"`
    Deleted      bool          `bson:"deleted"`
}

func IsAlreadyRatedQuality(username, crackmehexid string) (bool, error) {
    var err error
    var nb int
    var res bool
    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_quality")
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

func RatingQualityByCrackme(crackmehexid string) ([]RatingQuality, error) {
    var err error
    var result []RatingQuality
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_quality")

        // Validate the object id
        err = c.Find(bson.M{"crackmehexid": crackmehexid}).All(&result)
    } else {
        err = ErrUnavailable
    }
    return result, err
}

func RatingQualitySetRating(username, crackmehexid string, rating int) error {
    var err error
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_quality")

        // Validate the object id
        err = c.Update(bson.M{"crackmehexid": crackmehexid, "author": username}, bson.M{"$set": bson.M{"rating": rating}})
    } else {
        err = ErrUnavailable
    }
    return err
}

func RatingQualityCreate(username, crackmehexid string, rating int) error {
    var err error

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("rating_quality")
        objId := bson.NewObjectId()
        rating_quality:= &RatingQuality{
            ObjectId:     objId,
            Rating:       rating,
            Author:       username,
            CrackMeHexId: crackmehexid,
            CreatedAt:    time.Now(),
            Visible:      true,
            Deleted:      false,
        }
        err = c.Insert(rating_quality)
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

package model

import (
    "github.com/sar5430/crackmes.one/app/shared/database"
    "time"

    "gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Comment
// *****************************************************************************

// Comment table contains the information for each note
type Comment struct {
    ObjectId     bson.ObjectId `bson:"_id,omitempty"`
    Content      string        `bson:"info,omitempty"`
    Author       string        `bson:"author,omitempty"`
    CrackMeHexId string        `bson:"crackmehexid,omitempty"`
    CreatedAt    time.Time     `bson:"created_at"`
    Visible      bool          `bson:"visible"`
    Deleted      bool          `bson:"deleted"`
}

func CountCommentsByUser(username string) (int, error) {
    var err error
    var nb int
    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("comment")
        nb, err = c.Find(bson.M{"author": username}).Count()
    } else {
        err = ErrUnavailable
    }
    return nb, standardizeError(err)
}

func CountCommentsByCrackme(crackmehexid string) (int, error) {
    var err error
    var nb int
    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("comment")
        nb, err = c.Find(bson.M{"crackmehexid": crackmehexid, "visible": true}).Count()
    } else {
        err = ErrUnavailable
    }
    return nb, standardizeError(err)
}

func CommentsByUser(name string) ([]Comment, error) {
    var err error
    var result []Comment
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("comment")

        // Validate the object id
        err = c.Find(bson.M{"author": name, "visible": true}).Limit(50).Sort("-created_at").All(&result)
    } else {
        err = ErrUnavailable
    }
    return result, err
}

func CommentsByCrackMe(crackmehexid string) ([]Comment, error) {
    var err error
    var result []Comment
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("comment")

        // Validate the object id
        err = c.Find(bson.M{"crackmehexid": crackmehexid, "visible": true}).Sort("created_at").All(&result)
    } else {
        err = ErrUnavailable
    }
    return result, err
}

func CommentCreate(content, username, crackmehexid string) error {
    var err error

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("comment")
        objId := bson.NewObjectId()
        comment := &Comment{
            ObjectId:     objId,
            Content:      content,
            Author:       username,
            CrackMeHexId: crackmehexid,
            CreatedAt:    time.Now(),
            Visible:      true,
            Deleted:      false,
        }
        err = c.Insert(comment)
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

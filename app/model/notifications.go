package model

import (
    "time"

    "app/shared/database"
    "gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Notifications
// *****************************************************************************

// Notifications table contains the notification informations for each user
type Notification struct {
    ObjectId    bson.ObjectId   `bson:"_id,omitempty"`
    HexId       string          `bson:"hexid,omitempty"`
    User        string          `bson:"user,omitempty"`
    Text        string          `bson:"text,omitempty"`
    Time        time.Time       `bson:"time"`
    Seen        bool            `bson:"seen"`
}


// Returns all notifications of a user
func NotificationsByUser(username string) ([]Notification, error) {
    var err error

    result := []Notification{}

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("notifications")
        err = c.Find(bson.M{"user": username}).Sort("-time").All(&result)
    } else {
        err = ErrUnavailable
    }

    return result, standardizeError(err)
}

// Adds a new notification for user
func NotificationAdd(username, text string) error {
    var err error

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("notifications")
        objId := bson.NewObjectId()
        notif := &Notification{
            ObjectId:  objId,
            HexId:	   objId.Hex(),
            User:      username,
            Text:      text,
            Time:      time.Now(),
            Seen:      false,
        }
        err = c.Insert(notif)
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

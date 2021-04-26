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

// Sets these notifications to Seen in the db.
func NotificationsSetSeen(toSetSeen []Notification) error {
    var err error

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("notifications")
        b := c.Bulk();
        b.Unordered();
        for i, _ := range toSetSeen {
            if toSetSeen[i].Seen {
                continue
            }
            b.Update(
            bson.M{
                "hexid": toSetSeen[i].HexId },
            bson.M{
                "$set": bson.M{"seen": true} })
        }
        _, err = b.Run();
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

// Returns true, if there are unseen notifications for user
func NotificationsHasUnseen(username string) (bool, error) {
    var err error
    var result bool

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("notifications")
        n, err := c.Find(bson.M{"user": username, "seen": false}).Count()
        if err == nil {
            result = n != 0;
        }
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

// Removes a notification from user
func NotificationRemove(username, hexid string) error {
    var err error

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("notifications")
        err = c.Remove(bson.M{"user": username, "hexid": hexid})
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

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
// Notifications
// *****************************************************************************

// Notifications table contains the notification informations for each user
type Notification struct {
	ObjectId primitive.ObjectID `bson:"_id,omitempty"`
	HexId    string             `bson:"hexid,omitempty"`
	User     string             `bson:"user,omitempty"`
	Text     string             `bson:"text,omitempty"`
	Time     time.Time          `bson:"time"`
	Seen     bool               `bson:"seen"`
}

// Returns all notifications of a user
func NotificationsByUser(username string) ([]Notification, error) {
	var err error
	var cursor *mongo.Cursor

	result := []Notification{}
	if database.CheckConnection() {
		opts := options.Find().SetSort(bson.D{{"time", -1}})
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("notifications")
		cursor, err = collection.Find(database.Ctx, bson.M{"user": username}, opts)
		err = cursor.All(database.Ctx, &result)
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// Sets these notifications to Seen in the db.
func NotificationsSetSeen(toSetSeen []Notification) error {
	var err error

	if database.CheckConnection() {

		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("notifications")

		for i, _ := range toSetSeen {
			if toSetSeen[i].Seen {
				continue
			}

			collection.UpdateOne(database.Ctx,
				bson.M{
					"hexid": toSetSeen[i].HexId},
				bson.M{
					"$set": bson.M{"seen": true}})
		}

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
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("notifications")
		n, err := collection.CountDocuments(database.Ctx, bson.M{"user": username, "seen": false})
		if err == nil {
			result = n != 0
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
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("notifications")

		objId := primitive.NewObjectID()
		notif := &Notification{
			ObjectId: objId,
			HexId:    objId.Hex(),
			User:     username,
			Text:     text,
			Time:     time.Now(),
			Seen:     false,
		}
		_, err = collection.InsertOne(database.Ctx, notif)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// Removes a notification from user
func NotificationRemove(username, hexid string) error {
	var err error

	if database.CheckConnection() {
		collection := database.Mongo.Database(database.ReadConfig().MongoDB.Database).Collection("notifications")
		_, err = collection.DeleteOne(database.Ctx, bson.M{"user": username, "hexid": hexid})
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

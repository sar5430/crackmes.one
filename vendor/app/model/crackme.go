package model

import (
	"app/shared/database"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Crackme
// *****************************************************************************

// Crackme table contains the information for each note
type Crackme struct {
	ObjectId    bson.ObjectId `bson:"_id,omitempty"`
	HexId       string        `bson:"hexid,omitempty"`
	Name        string        `bson:"name,omitempty"`
	Info        string        `bson:"info,omitempty"`
	Lang        string        `bson:"lang,omitempty"`
	Author      string        `bson:"author,omitempty"`
	CreatedAt   time.Time     `bson:"created_at"`
	Visible     bool          `bson:"visible"`
	Deleted     bool          `bson:"deleted"`
    Difficulty  float64 `bson:"difficulty"`
    Quality     float64 `bson:"quality"`
    NbSolutions int           // Not present in the database! Just for rendering
    NbComments  int           // Not present in the database! Just for rendering
	Platform    string `bson:"platform,omitempty"`
}

func CountCrackmes() (int, error) {
	var err error
	var nb int
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")
		nb, err = c.Count()
	} else {
		err = ErrUnavailable
	}

	return nb, standardizeError(err)
}

func CountCrackmesByUser(username string) (int, error) {
	var err error
	var nb int
	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")
		nb, err = c.Find(bson.M{"author": username}).Count()
	} else {
		err = ErrUnavailable
	}
	return nb, standardizeError(err)
}

func GetAllCrackmes() ([]Crackme, error) {
	var err error
	var result []Crackme
	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")

		// Validate the object id
		err = c.Find(nil).All(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func CrackmeSetFloat(hexid, champ string, nb float64) error {
	var err error
	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")

		// Validate the object id
		err = c.Update(bson.M{"hexid": hexid}, bson.M{"$set": bson.M{champ: float64(nb)}})
	} else {
		err = ErrUnavailable
	}
	return err
}

func SearchCrackme(name, author, lang, platform string, difficulty int) ([]Crackme, error) {
	var err error
	var result []Crackme
	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")

		// Validate the object id
		err = c.Find(bson.M{"name": bson.RegEx{name, "i"}, "lang": bson.RegEx{lang, "i"}, "difficulty": bson.M{ "$gte": difficulty},"author": bson.RegEx{author, "i"}, "visible": true, "platform": bson.RegEx{platform, "i"}}).Limit(150).Sort("-created_at").All(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func LastCrackMes(page int) ([]Crackme, error) {
	var err error
	var result []Crackme
	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")

		// Validate the object id
		err = c.Find(bson.M{"visible": true}).Skip((page - 1) * 50).Limit(50).Sort("-created_at").All(&result)
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
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")

		// Validate the object id
		err = c.Find(bson.M{"hexid": hexid, "visible": true}).One(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

func CrackmesByUser(username string) ([]Crackme, error) {
	var err error

	var result []Crackme
	if database.CheckConnection() {
		// Create a copy of mongo
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")

		// Validate the object id
		err = c.Find(bson.M{"author": username, "visible": true}).Sort("-created_at").All(&result)
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
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")

		// Validate the object id
		err = c.Find(bson.M{"name": name, "author": username, "visible": visible}).One(&result)
	} else {
		err = ErrUnavailable
	}
	return result, err
}

// NoteCreate creates a note
func CrackmeCreate(name, info, username, lang, platform string) error {
	var err error

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("crackme")
		objId := bson.NewObjectId()
		crackme := &Crackme{
			ObjectId:   objId,
			HexId:      objId.Hex(),
			Name:       name,
			Info:       info,
			Lang:       lang,
			Author:     username,
			CreatedAt:  time.Now(),
			Visible:    false,
			Deleted:    false,
			Platform:   platform,
		}
		err = c.Insert(crackme)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

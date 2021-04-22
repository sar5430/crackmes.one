package model

import (
    "app/shared/database"
    "gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
    ObjectId	bson.ObjectId	`bson:"_id,omitempty"`
    HexId		string			`bson:"hexid,omitempty"`
    Name		string			`bson:"name,omitempty"`
    Email		string			`bson:"email,omitempty"`
    Password	string			`bson:"password,omitempty"`
    Visible		bool			`bson:"visible"`
    Deleted		bool			`bson:"deleted"`
    NbCrackmes	int
    NbSolutions	int
    NbComments	int
}

// Username returns the user name
func (u *User) Username() string {
    return u.Name
}

func CountUsers() (int, error) {
    var err error
    var nb int
    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
        nb, err = c.Count()
    } else {
        err = ErrUnavailable
    }

    return nb, standardizeError(err)
}

// UserByEmail gets user information from email
func UserByName(name string) (User, error) {
    var err error

    result := User{}

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
        err = c.Find(bson.M{"name": bson.RegEx{"^" + name + "$", "i"}}).One(&result)
    } else {
        err = ErrUnavailable
    }

    return result, standardizeError(err)
}

func UserByMail(email string) (User, error) {
    var err error

    result := User{}

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
        err = c.Find(bson.M{"email": bson.RegEx{"^" + email + "$", "i"}}).One(&result)
    } else {
        err = ErrUnavailable
    }

    return result, standardizeError(err)
}

func UserByHexId(hexid string) (User, error) {
    var err error

    result := User{}

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
        err = c.Find(bson.M{"hexid": hexid}).One(&result)
    } else {
        err = ErrUnavailable
    }

    return result, standardizeError(err)
}

func AllUsersVisible() ([]User, error) {
    var users []User
    var err error
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("user")

        // Validate the object id
        err = c.Find(bson.M{"visible": true}).All(&users)
    } else {
        err = ErrUnavailable
    }
    return users, err
}

// UserCreate creates user
func UserCreate(name, email, password string) error {
    var err error

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
        objId := bson.NewObjectId()
        user := &User{
            ObjectId:  objId,
            HexId:	   objId.Hex(),
            Name:      name,
            Email:     email,
            Password:  password,
            Visible:   true,
            Deleted:   false,
        }
        err = c.Insert(user)
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

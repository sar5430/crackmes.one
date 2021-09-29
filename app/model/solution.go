package model

import (
    "github.com/sar5430/crackmes.one/app/shared/database"
    "time"

    "gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Crackme
// *****************************************************************************

// Crackme table contains the information for each note
type Solution struct {
    ObjectId  bson.ObjectId `bson:"_id,omitempty"`
    HexId     string        `bson:"hexid,omitempty"`
    Info      string        `bson:"info"`
    CrackmeId bson.ObjectId `bson:"crackmeid,omitempty"`
    CreatedAt time.Time     `bson:"created_at"`
    Author    string        `bson:"author,omitempty"`
    Visible   bool          `bson:"visible"`
    Deleted   bool          `bson:"deleted"`
}

type SolutionExtended struct {
    Solution      *Solution
    Crackmeshexid string
    Crackmename   string
}

func CountSolutions() (int, error) {
    var err error
    var nb int
    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")
        nb, err = c.Count()
    } else {
        err = ErrUnavailable
    }

    return nb, standardizeError(err)
}

func CountSolutionsByUser(username string) (int, error) {
    var err error
    var nb int
    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")
        nb, err = c.Find(bson.M{"author": username, "visible": true}).Count()
    } else {
        err = ErrUnavailable
    }
    return nb, standardizeError(err)
}

func CountSolutionsByCrackme(crackmehexid string) (int, error) {
    var err error
    var nb int
    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        //obj := bson.M{"crackmeid": crackmehexid}
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")
        nb, err = c.Find(bson.M{"crackmeid": bson.ObjectIdHex(crackmehexid), "visible": true}).Count()
    } else {
        err = ErrUnavailable
    }
    return nb, standardizeError(err)
}

func SolutionByHexId(hexid string) (Solution, error) {
    var err error

    var result Solution
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")

        // Validate the object id
        err = c.Find(bson.M{"hexid": hexid, "visible": true}).One(&result)
    } else {
        err = ErrUnavailable
    }
    return result, err
}

func SolutionsByUser(username string) ([]Solution, error) {
    var err error

    var result []Solution
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")

        // Validate the object id
        err = c.Find(bson.M{"author": username, "visible": true}).Limit(50).Sort("-created_at").All(&result)
    } else {
        err = ErrUnavailable
    }
    return result, err
}

func SolutionsByUserAndCrackMe(username, crackmehexid string) (Solution, error) {
    var err error

    var result Solution
    crackme, _ := CrackmeByHexId(crackmehexid)
    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")

        // Validate the object id
        err = c.Find(bson.M{"crackmeid": crackme.ObjectId, "author": username}).One(&result)
    } else {
        err = ErrUnavailable
    }
    return result, err
}

func SolutionsByCrackme(crackme bson.ObjectId) ([]Solution, error) {
    var err error

    var result []Solution

    if database.CheckConnection() {
        // Create a copy of mongo
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")

        // Validate the object id
        err = c.Find(bson.M{"crackmeid": crackme, "visible": true}).All(&result)
    } else {
        err = ErrUnavailable
    }
    return result, err
}

// NoteCreate creates a note
func SolutionCreate(info, username, crackmehexid string) error {
    var err error
    crackme, _ := CrackmeByHexId(crackmehexid)

    if database.CheckConnection() {
        session := database.Mongo.Copy()
        defer session.Close()
        c := session.DB(database.ReadConfig().MongoDB.Database).C("solution")
        objId := bson.NewObjectId()
        solution := &Solution{
            ObjectId:  objId,
            HexId:     objId.Hex(),
            Info:      info,
            CrackmeId: crackme.ObjectId,
            CreatedAt: time.Now(),
            Author:    username,
            Visible:   false,
            Deleted:   false,
        }
        err = c.Insert(solution)
    } else {
        err = ErrUnavailable
    }

    return standardizeError(err)
}

package main

import (
	"fmt"
	"time"
	"os"
	"io/ioutil"
	"regexp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Crackme
// *****************************************************************************
// Crackme table contains the information for each note
type Crackme struct {
	ObjectId	bson.ObjectId		`bson:"_id,omitempty"`
	HexId		string			`bson:"hexid,omitempty"`
	Name		string			`bson:"name,omitempty"`
	Info		string			`bson:"info,omitempty"`
	Lang		string			`bson:"lang,omitempty"`
	Difficulty      string                  `bson:"difficulty,omitempty"`
	Author		string			`bson:"author,omitempty"`
	CreatedAt	time.Time		`bson:"created_at"`
	Visible		bool			`bson:"visible"`
	Deleted		bool			`bson:"deleted"`
}


func main(){
	info := &mgo.DialInfo{
        	Addrs:    []string{"127.0.0.1"},
        	Timeout:  60 * time.Second,
        	Database: "crackmesone",
    	}

    	session, err1 := mgo.DialWithInfo(info)
    	if err1 != nil {
		panic(err1)
    	}

	dir, _ := os.Open("../crackmesde.github.io/users")
	files, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}
	dir.Close()
	for _, file := range files{
		diruser, _ := os.Open("../crackmesde.github.io/users/" + file.Name())
		crackmes, _ := diruser.Readdir(0)
		for _, crackme := range crackmes{
			if crackme.Name() != "index.html"{
				buf, _ := ioutil.ReadFile("../crackmesde.github.io/users/" + file.Name() + "/" + crackme.Name() + "/index.html")
				r, _ := regexp.Compile("<em class=\"ae\".*<br />")
				matches := r.FindAllString(string(buf), -1)
				platform := matches[1][30:len(matches[1])-6]
				fmt.Println(platform)
				fmt.Println(crackme.Name() + " by " + file.Name())
				c := session.DB("crackmesone").C("crackme")
				err = c.Update(bson.M{"name": crackme.Name() + " by " + file.Name()}, bson.M{"$set": bson.M{"platform": platform}})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

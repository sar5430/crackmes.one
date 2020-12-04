package main

import (
	"fmt"
	"time"
	"os"
	"bytes"
	//"io/ioutil"
	//"regexp"
	"os/exec"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Solution struct {
	ObjectId	bson.ObjectId		`bson:"_id,omitempty"`
	HexId		string			`bson:"hexid,omitempty"`
	Info		string			`bson:"info"`
	CrackmeId	bson.ObjectId		`bson:"crackmeid,omitempty"`
	Author		string			`bson:"author,omitempty"`
	Visible		bool			`bson:"visible"`
	Deleted		bool			`bson:"deleted"`
}

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

func CrackmeByUserAndName(session *mgo.Session, username, name string) (Crackme, error) {
	var err error
	var result Crackme
	// Create a copy of mongo	
	c := session.DB("crackmesone").C("crackme")
	err = c.Find(bson.M{"name": name, "author": username}).One(&result)
	return result, err
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
			dirsolutions, _ := os.Open("../crackmesde.github.io/users/" + file.Name() + "/" + crackme.Name() + "/solutions/")
			userdirs, err := dirsolutions.Readdir(0)
			if err != nil {
				continue
			}
			crackmeobj, _ := CrackmeByUserAndName(session, "crackmes.de", crackme.Name() + " by " + file.Name())
			for _, userdir := range userdirs {
				fmt.Println(userdir.Name())
				objId := bson.NewObjectId()
				solutionobj := &Solution{
					ObjectId:  objId,
					HexId:     objId.Hex(),
					Info:      "This solution has been imported from crackmes.de. The original author is " + userdir.Name() + ". The password of the archive is \"crackmes.de\"",
					CrackmeId: crackmeobj.ObjectId,
					Author:    "crackmes.de",
					Visible:   true,
					Deleted:   false,
				}
				fmt.Println(solutionobj)
				c := session.DB("crackmesone").C("solution")
				err = c.Insert(solutionobj)
				if err != nil {
					fmt.Println(err)
				}

				dirdownload, _ := os.Open("../crackmesde.github.io/users/" + file.Name() + "/" + crackme.Name() + "/solutions/" + userdir.Name() + "/")
				download, _ := dirdownload.Readdir(0)
				filezip := download[0]
				var stderr bytes.Buffer
				cmd := exec.Command("cp", "../crackmesde.github.io/users/" + file.Name() + "/" + crackme.Name() + "/solutions/" + userdir.Name() + "/" + filezip.Name(), "../static/solution/" + objId.Hex() + ".zip")
			//cmd := exec.Command("ls", "../crackmesde.github.io/users/" + file.Name() + "/" + crackme.Name() + "/download/")
				cmd.Stderr = &stderr
				err := cmd.Run()
				if err != nil {
					fmt.Println("copy failed")
					fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				}
			}
		}
	}
}

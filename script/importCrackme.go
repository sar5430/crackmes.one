package main

import (
	"fmt"
	"time"
	"os"
	"bytes"
	"io/ioutil"
	"regexp"
	"os/exec"
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
				difficulty := matches[0][32:33]
				language := matches[2][30:len(matches[2])-42]
				objId := bson.NewObjectId()
				crackmeobj := &Crackme{
					ObjectId:  objId,
					HexId:     objId.Hex(),
					Name:      crackme.Name() + " by " + file.Name() ,
					Info:      "This crackme has been imported from crackmes.de. The original author is " + file.Name() + ". The password of the archive is \"crackmes.de\"",
					Lang:	   language,
					Difficulty:difficulty,
					Author:    "crackmes.de",
					CreatedAt: time.Now(),
					Visible:   true,
					Deleted:   false,
				}
				c := session.DB("crackmesone").C("crackme")
				err = c.Insert(crackmeobj)
				if err != nil {
					fmt.Println(err)
				}

				dirdownload, _ := os.Open("../crackmesde.github.io/users/" + file.Name() + "/" + crackme.Name() + "/download/")
				download, _ := dirdownload.Readdir(0)
				filezip := download[0]

				var stderr bytes.Buffer
				cmd := exec.Command("cp", "../crackmesde.github.io/users/" + file.Name() + "/" + crackme.Name() + "/download/" + filezip.Name(), "../static/crackme/" + objId.Hex() + ".zip")
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

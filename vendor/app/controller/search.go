package controller

import (
	"log"
	"strconv"
	"net/http"
	"app/model"
	"github.com/josephspurrier/csrfbanana"
	"app/shared/view"
	"app/shared/session"
)

// AboutGET displays the About page
func SearchGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	sess := session.Instance(r)
	v := view.New(r)
	v.Name = "search/search"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
	sess.Save(r, w)
}

func SearchPOST(w http.ResponseWriter, r *http.Request) {
    var difficultyint int

	sess := session.Instance(r)

	name := r.FormValue("name")
	author := r.FormValue("author")
	difficulty := r.FormValue("difficulty")
	lang := r.FormValue("lang")
	platform := r.FormValue("platform")

    if difficulty == "" {
        difficultyint = 0
    } else {
        difficultyint, _ = strconv.Atoi(difficulty)
    }

	crackmes, err := model.SearchCrackme(name, author, lang, platform, difficultyint)
	if err != nil {
                log.Println(err)
                Error500(w, r)
                return
        }

	//crackmes = CrackMeConvertDiffToImg(crackmes)

	v := view.New(r)
        v.Name = "search/search"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Vars["crackmes"] = crackmes
	sess.Save(r, w)
	v.Render(w)
}

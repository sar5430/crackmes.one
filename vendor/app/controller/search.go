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
    var difficulty_min_int, difficulty_max_int int
    var quality_min_int, quality_max_int int

	sess := session.Instance(r)

	name := r.FormValue("name")
	author := r.FormValue("author")
	difficulty_min := r.FormValue("difficulty-min")
	difficulty_max := r.FormValue("difficulty-max")
	quality_min := r.FormValue("quality-min")
	quality_max := r.FormValue("quality-max")
	lang := r.FormValue("lang")
	platform := r.FormValue("platform")

    difficulty_min_int, _ = strconv.Atoi(difficulty_min)
    difficulty_max_int, _ = strconv.Atoi(difficulty_max)
    quality_min_int, _ = strconv.Atoi(quality_min)
    quality_max_int, _ = strconv.Atoi(quality_max)

	crackmes, err := model.SearchCrackme(name, author, lang, platform, difficulty_min_int, difficulty_max_int, quality_min_int, quality_max_int)
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

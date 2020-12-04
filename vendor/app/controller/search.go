package controller

import (
	"log"
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
	sess := session.Instance(r)

	/*if validate, missingField := view.Validate(r, []string{"name"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		SearchGET(w, r)
		return
	}*/

	name := r.FormValue("name")
	author := r.FormValue("author")
	difficulty := r.FormValue("difficulty")
	lang := r.FormValue("lang")
	platform := r.FormValue("platform")
	crackmes, err := model.SearchCrackme(name, author, difficulty, lang, platform)
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

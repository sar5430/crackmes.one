package controller

import (
	"fmt"
	"log"
	"net/http"
	"app/model"
	"app/shared/recaptcha"
	"app/shared/session"
	"app/shared/view"

	"github.com/kennygrant/sanitize"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

func LeaveCommentPOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)
	var err error
	var params httprouter.Params
        params = context.Get(r, "params").(httprouter.Params)
        crackmehexid := params.ByName("hexid")

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"comment"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		CrackMeGET(w, r)
		return
	}

	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		CrackMeGET(w, r)
		return
	}

	username := fmt.Sprintf("%s", sess.Values["name"])
	comment := r.FormValue("comment")

	comment = sanitize.HTML(comment)

	err = model.CommentCreate(comment, username, crackmehexid)

	if err != nil {
		log.Println(err)
	}
	
	crackme, err := model.CrackmeByHexId(crackmehexid)	
	if err != nil {
                log.Println(err)
        }
	
	err = model.CrackmeSet(crackmehexid, "nbcomments", crackme.NbComments + 1)
	if err != nil {
                log.Println(err)
        }

	sess.AddFlash(view.Flash{"Comment uploaded!", view.FlashSuccess})
	sess.Save(r, w)
	http.Redirect(w, r, "/crackme/" + crackmehexid, http.StatusFound)
	return
}


package controller

import (
	"app/model"
	"app/shared/recaptcha"
	"app/shared/session"
	"app/shared/view"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/kennygrant/sanitize"
)

func CrackMeGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	sess := session.Instance(r)
	var params httprouter.Params

	params = context.Get(r, "params").(httprouter.Params)
	hexid := params.ByName("hexid")

	crackme, err := model.CrackmeByHexId(hexid)
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	solutions, err := model.SolutionsByCrackme(crackme.ObjectId)
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	comments, err := model.CommentsByCrackMe(crackme.HexId)
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	v := view.New(r)
	v.Name = "crackme/read"
	v.Vars["info"] = crackme.Info
	v.Vars["name"] = crackme.Name
	v.Vars["hexid"] = crackme.HexId
	v.Vars["lang"] = crackme.Lang
	v.Vars["difficulty"] = crackme.Difficulty
	v.Vars["createdat"] = crackme.CreatedAt
	v.Vars["username"] = crackme.Author
	v.Vars["platform"] = crackme.Platform
	v.Vars["solutions"] = solutions
	v.Vars["comments"] = comments
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
	sess.Save(r, w)

}

func LastCrackMesGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	crackmes, err := model.LastCrackMes()
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	v := view.New(r)
	v.Name = "crackme/lasts"
	v.Vars["crackmes"] = crackmes
	v.Render(w)
}

func UploadCrackMeGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "crackme/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
	sess.Save(r, w)
}

// NotepadCreatePOST handles the note creation form submission
func UploadCrackMePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"name", "info", "lang", "difficulty", "platform"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		UploadCrackMeGET(w, r)
		return
	}

	username := fmt.Sprintf("%s", sess.Values["name"])
	name := r.FormValue("name")
	lang := r.FormValue("lang")
	difficulty := r.FormValue("difficulty")
	info := r.FormValue("info")
	platform := r.FormValue("platform")
	file, header, err := r.FormFile("file")

	name = sanitize.HTML(name)
	lang = sanitize.HTML(lang)
	info = sanitize.HTML(info)

	diffint, _ := strconv.Atoi(difficulty)
	if diffint > 6 || diffint < 1 {
		sess.AddFlash(view.Flash{"Wrong difficulty", view.FlashError})
		sess.Save(r, w)
		UploadCrackMeGET(w, r)
		return
	}

	if !recaptcha.Verified(r) {
		sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
		sess.Save(r, w)
		UploadCrackMeGET(w, r)
		return
	}

	if err != nil {
		log.Println(err)
	}

	if header.Filename == "" {
		sess.AddFlash(view.Flash{"Field missing: file", view.FlashError})
		sess.Save(r, w)
		UploadCrackMeGET(w, r)
		return
	}

	data, err := ioutil.ReadAll(file)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	if len(data) > 5000000 {
		sess.AddFlash(view.Flash{"This file is too large !", view.FlashError})
		sess.Save(r, w)
		UploadCrackMeGET(w, r)
		return
	}

	err = model.CrackmeCreate(name, info, username, lang, difficulty, platform)

	if err != nil {
		log.Println(err)
	}

	crackme, err := model.CrackmeByUserAndName(username, name, false)

	if err != nil {
		log.Println(err)
	}

	filename := path.Join("./tmp/crackme/" + username + "+++" + crackme.HexId + "+++" + header.Filename)
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"Crackme uploaded! Should be available soon.", view.FlashSuccess})
		sess.Save(r, w)
		http.Redirect(w, r, "/user/"+username, http.StatusFound)
		return
	}

}

func UploadCount(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	crackmes, _ := model.GetAllCrackmes()
	for i, _ := range crackmes {
		nbSolutions, err := model.CountSolutionsByCrackme(crackmes[i].ObjectId)
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}

		nbComments, err := model.CountCommentsByCrackme(crackmes[i].HexId)
		if err != nil {
			log.Println(err)
			Error500(w, r)
			return
		}
		model.CrackmeSet(crackmes[i].HexId, "nbsolutions", nbSolutions)
		model.CrackmeSet(crackmes[i].HexId, "nbcomments", nbComments)

	}

	// Display the view
	v := view.New(r)
	v.Name = "crackme/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
	sess.Save(r, w)
}

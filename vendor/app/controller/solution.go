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

    "github.com/gorilla/context"
    "github.com/josephspurrier/csrfbanana"
    "github.com/julienschmidt/httprouter"
    "github.com/kennygrant/sanitize"
)

func UploadSolutionGET(w http.ResponseWriter, r *http.Request) {
    // Get session
    var params httprouter.Params
    sess := session.Instance(r)
    params = context.Get(r, "params").(httprouter.Params)
    hexidcrackme := params.ByName("hexidcrackme")

    //Get crackme and user
    crackme, _ := model.CrackmeByHexId(hexidcrackme)

    // Display the view
    v := view.New(r)
    v.Name = "solution/create"
    v.Vars["token"] = csrfbanana.Token(w, r, sess)
    v.Vars["hexidcrackme"] = hexidcrackme
    v.Vars["username"] = crackme.Author
    v.Vars["crackmename"] = crackme.Name
    v.Render(w)
    sess.Save(r, w)
}

// NotepadCreatePOST handles the note creation form submission
func UploadSolutionPOST(w http.ResponseWriter, r *http.Request) {
    var params httprouter.Params
    sess := session.Instance(r)
    params = context.Get(r, "params").(httprouter.Params)
    hexidcrackme := params.ByName("hexidcrackme")
    var solution model.Solution

    username := fmt.Sprintf("%s", sess.Values["name"])
    info := r.FormValue("info")
    file, header, err := r.FormFile("file")

    info = sanitize.HTML(info)

    solution, _ = model.SolutionsByUserAndCrackMe(username, hexidcrackme)

    emptysol := model.Solution{}
    if solution != emptysol {
        sess.AddFlash(view.Flash{"You've already submitted a solution to this crackme", view.FlashError})
        sess.Save(r, w)
        UploadSolutionGET(w, r)
        return
    }

    if !recaptcha.Verified(r) {
        sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
        sess.Save(r, w)
        UploadSolutionGET(w, r)
        return
    }

    if err != nil {
        sess.AddFlash(view.Flash{"Field missing: file", view.FlashError})
        sess.Save(r, w)
        fmt.Println("missing file")
        UploadSolutionGET(w, r)
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
        UploadSolutionGET(w, r)
        return
    }

    err = model.SolutionCreate(info, username, hexidcrackme)
    solution, _ = model.SolutionsByUserAndCrackMe(username, hexidcrackme)

    if err != nil {
        log.Println(err)
    }

    filename := path.Join("./tmp/solution/" + username + "+++" + solution.HexId + "+++" + header.Filename)
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
        sess.AddFlash(view.Flash{"Solution uploaded! Should be available soon.", view.FlashSuccess})
        sess.Save(r, w)
        http.Redirect(w, r, "/user/"+username, http.StatusFound)
        return
    }

}

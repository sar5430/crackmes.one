package controller

import (
    "app/model"
    "app/shared/session"
    "app/shared/view"
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/context"
    "github.com/julienschmidt/httprouter"
)

func RateDifficultyPOST(w http.ResponseWriter, r *http.Request) {
    // Get session
    var already_exist bool
    sess := session.Instance(r)
    var err error
    var params httprouter.Params
    params = context.Get(r, "params").(httprouter.Params)
    crackmehexid := params.ByName("hexid")

    // Validate with required fields
    if validate, missingField := view.Validate(r, []string{"difficulty"}); !validate {
        sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
        sess.Save(r, w)
        CrackMeGET(w, r)
        return
    }

    username := fmt.Sprintf("%s", sess.Values["name"])
    rating := r.FormValue("difficulty")

    ratingint, _ := strconv.Atoi(rating)

    if ratingint < 1 || ratingint > 6 {
        log.Println("Wrong rating number")
        Error500(w, r)
        return
    }

    already_exist, err = model.IsAlreadyRatedDifficulty(username, crackmehexid)

    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    if already_exist {
        err = model.RatingDifficultySetRating(username, crackmehexid, ratingint)
        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }
    } else {
        err = model.RatingDifficultyCreate(username, crackmehexid, ratingint)
        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }
    }

    if err != nil {
        log.Println(err)
        Error500(w, r)
    }

    sess.AddFlash(view.Flash{"Rated!", view.FlashSuccess})
    sess.Save(r, w)
    http.Redirect(w, r, "/crackme/" + crackmehexid, http.StatusFound)
    return
}

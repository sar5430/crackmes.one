package controller

import (
    "github.com/sar5430/crackmes.one/app/model"
    "log"
    "net/http"
    "time"

    "github.com/sar5430/crackmes.one/app/shared/session"
    "github.com/sar5430/crackmes.one/app/shared/view"
    "github.com/josephspurrier/csrfbanana"
)

func NotificationsGET(w http.ResponseWriter, r *http.Request) {
    sess := session.Instance(r)

    notifs, err := model.NotificationsByUser(sess.Values["name"].(string))
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    for i, _ := range notifs {
        if !notifs[i].Seen {
            model.NotificationsSetSeen(notifs)
            break
        }
    }

    // Display the view
    v := view.New(r)
    v.Name = "notifs/notifs"
    v.Vars["notifs"] = notifs
    v.Vars["token"] = csrfbanana.TokenWithPath(w, r, sess, "/notifications/delete")
    v.Vars["startTime"] = time.Unix(0, 0)
    v.Render(w)
}

func NotificationsDeletePOST(w http.ResponseWriter, r *http.Request) {
    sess := session.Instance(r)
    uname := sess.Values["name"].(string)
    hexid := r.FormValue("hexid");

    if hexid == "" {
        Error500(w, r)
        return
    }

    err := model.NotificationRemove(uname, hexid)
    if err != nil {
        Error500(w, r)
        return
    }

    w.WriteHeader(http.StatusOK)
}

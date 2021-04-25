package controller

import (
    "app/model"
    "log"
    "net/http"

    "app/shared/view"
    "app/shared/session"

    //"github.com/gorilla/sessions"
)

func NotificationsGET(w http.ResponseWriter, r *http.Request) {
    sess := session.Instance(r)

    notifs, err := model.NotificationsByUser(sess.Values["name"].(string))
    if err != nil {
        log.Println(err)
        Error500(w, r)
        return
    }

    // Display the view
    v := view.New(r)
    v.Name = "notifs/notifs"
    v.Vars["notifs"] = notifs
    v.Render(w)
}

package controller

import (
	"fmt"
	"log"
	"net/http"
	"app/shared/view"
	"app/model"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "index/index"
	var nbusers, nbcrackmes, nbsolutions int
	var err error

	nbusers, err = model.CountUsers()
	if err != nil {
                log.Println(err)
                Error500(w, r)
                return
        }
	
	nbcrackmes, err = model.CountCrackmes()
        if err != nil {
                log.Println(err)
                Error500(w, r)
                return
        }

	nbsolutions, err = model.CountSolutions()
        if err != nil {
                log.Println(err)
                Error500(w, r)
                return
        }

	v.Vars["nbusers"] = fmt.Sprintf("0x%x", nbusers)
	v.Vars["nbsolutions"] = fmt.Sprintf("0x%x",nbsolutions)
	v.Vars["nbcrackmes"] = fmt.Sprintf("0x%x",nbcrackmes)
	v.Render(w)
}

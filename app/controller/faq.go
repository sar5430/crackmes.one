package controller

import (
    "net/http"

    "github.com/sar5430/crackmes.one/app/shared/view"
)

// AboutGET displays the About page
func FaqGET(w http.ResponseWriter, r *http.Request) {
    // Display the view
    v := view.New(r)
    v.Name = "faq/faq"
    v.Render(w)
}

package controller

import (
	"net/http"

	"app/shared/view"
)

// AboutGET displays the About page
func FaqGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "faq/faq"
	v.Render(w)
}

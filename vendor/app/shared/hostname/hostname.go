package hostname

import(
	"net/http"
)

func Handler(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "crackmes.one" || r.Host == "www.crackmes.one"{
			next.ServeHTTP(w, r)
		} else {
			//time.Sleep(100000 * time.Millisecond)
			http.Error(w,"Not found", 404)
		}
        })
}

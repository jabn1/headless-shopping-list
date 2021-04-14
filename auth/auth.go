package auth
 
import (
	"net/http"
	"strings"
)
 
func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, rq *http.Request) {
		u, p, ok := rq.BasicAuth()
		if !ok || len(strings.TrimSpace(u)) < 1 || len(strings.TrimSpace(p)) < 1 {
			unauthorised(rw)
			return
		}
 
		// This is a dummy check for credentials.
		if u != "hello" || p != "world" {
			unauthorised(rw)
			return
		}
 
		// If required, Context could be updated to include authentication
		// related data so that it could be used in consequent steps.
		handler(rw, rq)
	}
}
 
func unauthorised(rw http.ResponseWriter) {
	rw.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	rw.WriteHeader(http.StatusUnauthorized)
}
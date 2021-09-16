package auth

import (
	"github.com/solarlabsteam/dvpn-openwrt/utilities/shadow"
	"net/http"
)

func BasicAuthForHandler(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := shadow.Lookup(username)

			if err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if err = user.VerifyPassword(password); err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if user.IsPasswordValid() {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			user, err := shadow.Lookup(username)

			if err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if err = user.VerifyPassword(password); err != nil {
				w.Header().Add("Clear-Site-Data", "*")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if user.IsPasswordValid() {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

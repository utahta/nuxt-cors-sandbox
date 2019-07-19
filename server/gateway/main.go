package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// singleJoiningSlash was copied from net/http/httputil
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// ProxyHandler returns proxy handler.
func ProxyHandler(target *url.URL) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this code was copied from net/http/httputil
		targetQuery := target.RawQuery
		director := func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
			if targetQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = targetQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
			}
		}

		proxy := &httputil.ReverseProxy{
			Director: director,
		}
		proxy.ServeHTTP(w, r)
	})
}

func CORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		const target = "http://localhost:3000"
		fmt.Printf("method:%s\n", r.Method)
		fmt.Printf("origin:%s target:%s\n", origin, target)
		fmt.Printf("Authorization:%s\n", r.Header.Get("Authorization"))
		cookie, err := r.Cookie("session_id")
		if err == nil {
			fmt.Printf("cookie:%s\n", cookie.Value)
		} else {
			fmt.Printf("no cookie\n")
		}

		switch r.Method {
		case http.MethodOptions: // preflight request
			if origin != target {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "0")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", target)
			w.WriteHeader(http.StatusOK)
		default:
			if origin != target {
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", target)
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	nuxtURL, err := url.Parse("http://localhost:3000")
	if err != nil {
		panic(err)
	}
	nuxtProxyHandler := ProxyHandler(nuxtURL)

	apiURL, err := url.Parse("https://localhost:3001")
	if err != nil {
		panic(err)
	}
	apiProxyHandler := ProxyHandler(apiURL)
	apiProxyHandler = CORSHandler(apiProxyHandler)

	mux := http.NewServeMux()
	mux.Handle("/api/", apiProxyHandler)
	mux.Handle("/auth", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "foobar",
			HttpOnly: true,
			Secure:   true,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}))
	mux.Handle("/", nuxtProxyHandler)

	err = http.ListenAndServeTLS(":8080", "./localhost.pem", "./localhost-key.pem", mux)
	//err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

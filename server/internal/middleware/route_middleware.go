package middleware

import "net/http"

func Get(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(res, req)
	}
}

func Post(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(res, req)
	}
}

func GetWithAuth(handler http.HandlerFunc) http.HandlerFunc {
	return Get(WithAuth(handler))
}

func PostWithAuth(handler http.HandlerFunc) http.HandlerFunc {
	return Post(WithAuth(handler))
}

package server

import "net/http"

func CorsMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, Accept")

		if req.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(rw, req)
	})
}

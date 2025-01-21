package user

import (
	"encoding/json"
	"net/http"
	"time"
)

func Logout(res http.ResponseWriter, req *http.Request) {
	cookie := createAuthCookie("")
	cookie.Expires = time.Now().Add(-time.Hour)
	http.SetCookie(res, cookie)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}

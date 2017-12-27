package app

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func baseauth(d string) []string {
	as := strings.SplitN(d, " ", 2)
	if len(as) != 2 {
		return nil
	}

	if as[0] != "Basic" {
		return nil
	}

	authStr, err := base64.StdEncoding.DecodeString(as[1])
	if err != nil {
		return nil
	}

	userPwd := strings.SplitN(string(authStr), ":", 2)
	if len(userPwd) != 2 {
		return nil
	}
	return userPwd
}

func checkBaseauth(d string) bool {
	userPwd := baseauth(d)
	if userPwd == nil {
		return false
	}

	if v, ok := Conf.User[userPwd[0]]; ok {
		if v.Password == userPwd[1] {
			return true
		}
	}

	return false
}

func CheckAuth(req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		return false
	}

	return checkBaseauth(auth)
}

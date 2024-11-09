package utils

import (
	"net/http"
	"strings"
)

func ContainsRole(jwtRoles []string, acceptableRoles []string) bool {
	setAcceptableRoles := make(map[string]bool)
	for _, v := range acceptableRoles {
		setAcceptableRoles[v] = true
	}

	for _, jwtRole := range jwtRoles {
		if _, exist := setAcceptableRoles[jwtRole]; exist {
			return true
		}
	}
	return false
}

type ResponseStatusCode struct {
	http.ResponseWriter
	StatusCode int
}

func (rst *ResponseStatusCode) WriteHeader(statusCode int) {
	rst.StatusCode = statusCode
	rst.ResponseWriter.WriteHeader(statusCode)
}

func GetToken(w http.ResponseWriter, r *http.Request) string {
	tokenHead := r.Header.Get("Authorization")
	if tokenHead == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return ""
	}

	tokenStr := strings.TrimPrefix(tokenHead, "Bearer ")
	return tokenStr
}

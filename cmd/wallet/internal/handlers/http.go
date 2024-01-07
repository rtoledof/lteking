package handlers

import (
	"net/http"

	"cubawheeler.io/pkg/cubawheeler"
)

func canDo(r *http.Request, roles ...cubawheeler.Role) bool {
	user := cubawheeler.UserFromContext(r.Context())
	if user == nil {
		return false
	}
	if user.Role == cubawheeler.RoleAdmin {
		return true
	}
	for _, role := range roles {
		if role == user.Role {
			return true
		}
	}
	return false
}

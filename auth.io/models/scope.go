package models

import "context"

type Scope string

const (
	ScopeIdentity     Scope = "models:*"
	ScopeIdentityMe   Scope = "models:me"
	ScopeUser         Scope = "models:user:*"
	ScopeUserCreate   Scope = "models:user:create"
	ScopeUserRead     Scope = "models:user:read"
	ScopeUserUpdate   Scope = "models:user:update"
	ScopeUserDelete   Scope = "models:user:delete"
	ScopeUserRegister Scope = "models:user:register"
	ScopeUserLogin    Scope = "models:user:login"

	ScopeClient       Scope = "models:client:*"
	ScopeClientRead   Scope = "models:client:read"
	ScopeClientCreate Scope = "models:client:create"
	ScopeClientUpdate Scope = "models:client:update"

	ScopeProfile       Scope = "models:profile:*"
	ScopeProfileUpdate Scope = "models:profile:update"
	ScopeProfileRead   Scope = "models:profile:read"
)

func HasScope(ctx context.Context, scope Scope) bool {
	scopes := ScopesFromContext(ctx)
	for _, v := range scopes {
		if v == scope {
			return true
		}
	}
	return false
}

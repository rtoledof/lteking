package identity

type Scope string

const (
	ScopeIdentity     Scope = "identity:*"
	ScopeIdentityMe   Scope = "identity:me"
	ScopeUser         Scope = "identity:user:*"
	ScopeUserCreate   Scope = "identity:user:create"
	ScopeUserRead     Scope = "identity:user:read"
	ScopeUserUpdate   Scope = "identity:user:update"
	ScopeUserDelete   Scope = "identity:user:delete"
	ScopeUserRegister Scope = "identity:user:register"
	ScopeUserLogin    Scope = "identity:user:login"

	ScopeClient       Scope = "identity:client:*"
	ScopeClientRead   Scope = "identity:client:read"
	ScopeClientCreate Scope = "identity:client:create"
	ScopeClientUpdate Scope = "identity:client:update"

	ScopeProfile       Scope = "identity:profile:*"
	ScopeProfileUpdate Scope = "identity:profile:update"
	ScopeProfileRead   Scope = "identity:profile:read"
)

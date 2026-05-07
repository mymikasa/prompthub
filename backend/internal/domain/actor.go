package domain

type Actor struct {
	UserID      int64
	WorkspaceID int64
	Role        string
}

func (a *Actor) IsOwner() bool {
	return a.Role == "owner"
}

func (a *Actor) IsAdmin() bool {
	return a.Role == "admin"
}

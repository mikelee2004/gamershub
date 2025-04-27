package types

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleGuest    Role = "guest"
	RoleEmployee Role = "employee"
)

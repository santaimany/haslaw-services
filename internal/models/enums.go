package models

type NewsStatus string

const (
	Posted  NewsStatus = "Posted"
	Drafted NewsStatus = "Drafted"
)

var ValidNewsStatuses = []NewsStatus{Posted, Drafted}

func (ns NewsStatus) String() string {
	return string(ns)
}

func (ns NewsStatus) IsValid() bool {
	for _, status := range ValidNewsStatuses {
		if ns == status {
			return true
		}
	}
	return false
}

type UserRole string

const (
	SuperAdmin UserRole = "super_admin"

	Admin UserRole = "admin"
)

var ValidUserRoles = []UserRole{SuperAdmin, Admin}

func (ur UserRole) String() string {
	return string(ur)
}

func (ur UserRole) IsValid() bool {
	for _, role := range ValidUserRoles {
		if ur == role {
			return true
		}
	}
	return false
}

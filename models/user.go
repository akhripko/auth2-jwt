package models

type User struct {
	ID         string
	Email      string
	Name       string
	NickName   string
	GivenName  string
	FamilyName string
	Picture    string
	Context    string
	Groups     map[string]string
}

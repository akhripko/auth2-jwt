package google

type User struct {
	ID         string            `json:"id"`
	Email      string            `json:"email"`
	Name       string            `json:"name"`
	NickName   string            `json:"name"` // nolint
	GivenName  string            `json:"given_name"`
	FamilyName string            `json:"family_name"`
	Picture    string            `json:"picture"`
	Context    string            `json:"-"`
	Groups     map[string]string `json:"-"`
}

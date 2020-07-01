package github

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	NickName string `json:"login"`
	Picture  string `json:"avatar_url"`
}

type Team struct {
	ID           int64        `json:"id"`
	Name         string       `json:"slug"`
	Organization Organization `json:"organization"`
}

type Organization struct {
	ID   int64  `json:"id"`
	Name string `json:"login"`
}

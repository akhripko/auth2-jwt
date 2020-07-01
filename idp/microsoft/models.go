package microsoft

//{"@odata.context":"https://graph.microsoft.com/v1.0/$metadata#users/$entity",
// "businessPhones":[],
// "displayName":"Oleksii Khrypko",
// "givenName":"Oleksii",
// "jobTitle":null,
// "mail":"Oleksii.Khrypko@90poe.io",
// "mobilePhone":null,
// "officeLocation":null,
// "preferredLanguage":"en-US",
// "surname":"Khrypko",
// "userPrincipalName":"Oleksii.Khrypko@90poe.io",
// "id":"8616cdb1-c82b-4d30-bb37-c1dc793bd667"}
type User struct {
	ID         string            `json:"id"`
	Email      string            `json:"mail"`
	Name       string            `json:"displayName"`
	NickName   string            `json:"userPrincipalName"`
	GivenName  string            `json:"givenName"`
	FamilyName string            `json:"surname"`
	Picture    string            `json:"picture"`
	Context    string            `json:"-"`
	Groups     map[string]string `json:"-"`
}

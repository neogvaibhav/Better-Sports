package helpers

type Login struct {
	Username string `json:"username" validator:"min=3,max=40,regexp=^[a-zA-Z0-9]$"`
	Passowrd string `json:"password" validator:"min=8"`
}

type SignUp struct {
	Name     string `json:"name" validator:"regexp=^[a-zA-Z0-9]$"`
	Username string `json:"username" validator:"min=3,max=40,regexp=^[a-zA-Z0-9.]$"`
	Passowrd string `json:"password" validator:"min=8"`
}

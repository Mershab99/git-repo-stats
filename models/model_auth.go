package models

type Auth struct {

	Token string `json:"token,omitempty"`

	Username string `json:"username,omitempty"`

	Password string `json:"password,omitempty"`
}

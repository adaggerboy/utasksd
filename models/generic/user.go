package generic

type User struct {
	ID         int     `json:"id,omitempty"`
	Email      *string `json:"email,omitempty"`
	Username   *string `json:"username,omitempty"`
	Firstname  *string `json:"firstname,omitempty"`
	Lastname   *string `json:"lastname,omitempty"`
	AvatarPath *string `json:"avatar_path,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	IsAdmin    *bool   `json:"is_admin,omitempty"`
	IsDirector *bool   `json:"is_director,omitempty"`
}

type Credentials struct {
	Hash string
	Salt string
	Algo string
}

func AllocateUser() *User {
	return &User{
		Username:   new(string),
		Lastname:   new(string),
		Firstname:  new(string),
		Email:      new(string),
		AvatarPath: new(string),
		IsActive:   new(bool),
		IsAdmin:    new(bool),
		IsDirector: new(bool),
	}
}

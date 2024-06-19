package generic

type Project struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	LogoPath    *string `json:"logo_path,omitempty"`
	OwnerID     *int    `json:"owner,omitempty"`

	//Recursively acquired fields
	Managers      *[]int `json:"managers,omitempty"`
	Workers       *[]int `json:"workers,omitempty"`
	SupportAgents *[]int `json:"support_agents,omitempty"`
}

func AllocateProject() *Project {
	return &Project{
		Name:        new(string),
		Description: new(string),
		LogoPath:    new(string),
		OwnerID:     new(int),

		Managers:      &[]int{},
		Workers:       &[]int{},
		SupportAgents: &[]int{},
	}
}

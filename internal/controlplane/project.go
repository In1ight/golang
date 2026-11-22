package controlplane

import "time"

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateProjectInput struct {
	Name     *string `json:"name"`
	Timezone *string `json:"timezone"`
}

type CreateProjectInput struct {
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

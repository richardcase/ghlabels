package main

type Label struct {
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Description *string `json:"description,omitempty"`
	Default     *bool   `json:"default,omitempty"`
}

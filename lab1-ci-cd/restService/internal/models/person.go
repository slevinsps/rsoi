package models

type Person struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Age     int    `json:"age,omitempty"`
	Address string `json:"address,omitempty"`
	Work    string `json:"work,omitempty"`
}

package models

type AdditionalProp struct {
	AdditionalProp1 string `json:"additionalProp1"`
	AdditionalProp2 string `json:"additionalProp2"`
	AdditionalProp3 string `json:"additionalProp3"`
}

type Message struct {
	Message string          `json:"message,omitempty"`
	Errors  *AdditionalProp `json:"errors,omitempty"`
}

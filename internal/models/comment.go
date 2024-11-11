package models

type Comment struct {
	Content     string `json:"Content,omitempty"`
	ContentType string `json:"ContentType,omitempty"`
}

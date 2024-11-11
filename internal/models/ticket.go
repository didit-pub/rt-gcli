package models

import (
	"time"
)

// CreateTicket representa los campos necesarios para crear un ticket
type TicketCreate struct {
	Subject      string            `json:"Subject,omitempty"`
	Queue        string            `json:"Queue,omitempty"`
	Status       string            `json:"Status,omitempty"`
	Priority     string            `json:"Priority,omitempty"`
	Owner        string            `json:"Owner,omitempty"`
	Requestor    string            `json:"Requestor,omitempty"`
	Content      string            `json:"Content,omitempty"`
	ContentType  string            `json:"ContentType,omitempty"`
	Parent       string            `json:"Parent,omitempty"`
	CustomFields map[string]string `json:"CustomFields,omitempty"`
}

type TicketCreateResponse struct {
	URL  string `json:"_url,omitempty"`
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

// Ticket representa un ticket en el sistema RT
type CustomField struct {
	Item
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

type Ticket struct {
	ID              int           `json:"id,omitempty"`
	Subject         string        `json:"Subject,omitempty"`
	Queue           Queue         `json:"Queue,omitempty"`
	Status          string        `json:"Status,omitempty"`
	FinalPriority   string        `json:"FinalPriority,omitempty"`
	Owner           User          `json:"Owner,omitempty"`
	Requestor       []User        `json:"Requestor,omitempty"`
	Created         *time.Time    `json:"Created,omitempty"`
	Cc              []User        `json:"Cc,omitempty"`
	Creator         User          `json:"Creator,omitempty"`
	TimeLeft        string        `json:"TimeLeft,omitempty"`
	TimeEstimated   string        `json:"TimeEstimated,omitempty"`
	AdminCc         []User        `json:"AdminCc,omitempty"`
	Starts          *time.Time    `json:"Starts,omitempty"`
	Started         *time.Time    `json:"Started,omitempty"`
	LastUpdated     *time.Time    `json:"LastUpdated,omitempty"`
	InitialPriority string        `json:"InitialPriority,omitempty"`
	Due             *time.Time    `json:"Due,omitempty"`
	LastUpdatedBy   User          `json:"LastUpdatedBy,omitempty"`
	Priority        string        `json:"Priority,omitempty"`
	Resolved        *time.Time    `json:"Resolved,omitempty"`
	EffectiveID     Item          `json:"EffectiveID,omitempty"`
	CustomFields    []CustomField `json:"CustomFields,omitempty"`
}

// TicketUpdate representa los campos actualizables de un ticket
type TicketUpdate struct {
	Status       *string           `json:"Status,omitempty"`
	CustomFields map[string]string `json:"CustomFields,omitempty"`
}

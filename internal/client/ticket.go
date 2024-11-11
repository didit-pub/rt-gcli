package client

import (
	"encoding/json"
	"fmt"

	"github.com/didit-pub/rt-gcli/internal/models"
)

// CreateTicket crea un nuevo ticket
func (c *Client) CreateTicket(ticket *models.TicketCreate) (*models.TicketCreateResponse, error) {
	resp, err := c.doRequest("POST", "ticket", ticket, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating ticket: %w", err)
	}

	var result models.TicketCreateResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &result, nil
}

// GetTicket obtiene un ticket por su ID
func (c *Client) GetTicket(id int) (*models.Ticket, error) {
	params := map[string]string{
		"fields[Queue]":   "Name",
		"fields[Owner]":   "Name,EmailAddress",
		"fields[Creator]": "Name,EmailAddress",
	}
	resp, err := c.doRequest("GET", fmt.Sprintf("ticket/%d", id), nil, params)
	if err != nil {
		return nil, fmt.Errorf("error getting ticket: %w", err)
	}

	var ticket models.Ticket
	if err := json.Unmarshal(resp, &ticket); err != nil {
		return nil, fmt.Errorf("error parsing ticket: %w", err)
	}
	// Iterate through requestors and fetch additional details
	for i := range ticket.Requestor {
		if ticket.Requestor[i].ID == "" {
			continue
		}

		user, err := c.GetUser(ticket.Requestor[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error getting requestor details: %w", err)
		}
		ticket.Requestor[i].EmailAddress = user.EmailAddress
		ticket.Requestor[i].Name = user.Name
	}
	return &ticket, nil
}

// UpdateTicket actualiza un ticket existente
func (c *Client) UpdateTicket(id int, updates *models.TicketUpdate) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("ticket/%d", id), updates, nil)
	if err != nil {
		return fmt.Errorf("error updating ticket: %w", err)
	}

	return nil
}

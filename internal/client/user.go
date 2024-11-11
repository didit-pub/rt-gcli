package client

import (
	"encoding/json"
	"fmt"

	"github.com/didit-pub/rt-gcli/internal/models"
)

func (c *Client) GetUser(userName string) (models.User, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("user/%s", userName), nil, nil)
	if err != nil {
		return models.User{}, fmt.Errorf("error getting ticket: %w", err)
	}

	var user models.User
	if err := json.Unmarshal(resp, &user); err != nil {
		return models.User{}, fmt.Errorf("error parsing user: %w", err)
	}

	return user, nil
}

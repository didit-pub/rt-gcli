package commands

import (
	"fmt"
	"strings"

	"github.com/didit-pub/rt-gcli/internal/client"
	"github.com/didit-pub/rt-gcli/internal/models"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var (
		subject      string
		queue        string
		content      string
		requestor    string
		owner        string
		parent       string
		customFields map[string]string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Crear un nuevo ticket",
		RunE: func(cmd *cobra.Command, args []string) error {
			rtClient, err := client.NewClient(&cfg)
			if err != nil {
				return fmt.Errorf("error al crear el cliente: %w", err)
			}

			if requestor == "" {
				requestor = cfg.Me
			}
			if owner == "" {
				owner = cfg.Me
			}

			ticket := &models.TicketCreate{
				Subject:      subject,
				Queue:        queue,
				Content:      strings.ReplaceAll(content, "\\n", "\n"),
				ContentType:  "text/plain",
				Requestor:    requestor,
				Owner:        owner,
				Parent:       parent,
				CustomFields: customFields,
			}

			ticketResponse, err := rtClient.CreateTicket(ticket)
			if err != nil {
				return fmt.Errorf("error al crear el ticket: %w", err)
			}

			fmt.Printf("RT_TICKET_ID=%s\n", ticketResponse.ID)
			fmt.Printf("RT_TICKET_URL=%s\n", ticketResponse.URL)
			return nil
		},
	}

	// Flags obligatorios
	cmd.Flags().StringVarP(&subject, "subject", "s", "", "Asunto del ticket")
	cmd.Flags().StringVarP(&queue, "queue", "q", "", "Cola del ticket")
	cmd.Flags().StringVarP(&requestor, "requestor", "r", "", "Solicitante del ticket")
	cmd.Flags().StringVarP(&owner, "owner", "o", "", "Propietario del ticket")
	cmd.Flags().StringVarP(&parent, "parent", "p", "", "Ticket padre")
	cmd.Flags().StringToStringVarP(&customFields, "custom-field", "f", map[string]string{}, "Campos personalizados del ticket")
	cmd.MarkFlagRequired("subject")
	cmd.MarkFlagRequired("queue")
	// Flags opcionales
	cmd.Flags().StringVarP(&content, "content", "c", "", "Contenido del ticket")

	return cmd
}

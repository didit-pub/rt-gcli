package commands

import (
	"fmt"
	"strings"

	"github.com/didit-pub/rt-gcli/internal/client"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	var ticketID int

	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Obtener un ticket por ID",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			rtClient, err := client.NewClient(&cfg)
			if err != nil {
				return fmt.Errorf("error al crear el cliente: %w", err)
			}

			ticket, err := rtClient.GetTicket(ticketID)
			if err != nil {
				return fmt.Errorf("error al obtener el ticket: %w", err)
			}

			fmt.Printf("RT_TICKET_URL=%s/Ticket/Display.html?id=%d\n", cfg.URL, ticket.ID)
			fmt.Printf("RT_TICKET_ID=%d\n", ticket.ID)
			fmt.Printf("RT_TICKET_API_URL=%s\n", ticket.EffectiveID.URL)
			fmt.Printf("RT_QUEUE_NAME=%s\n", ticket.Queue.Name)
			fmt.Printf("RT_OWNER_NAME=%s\n", ticket.Owner.Name)
			fmt.Printf("RT_OWNER_EMAIL=%s\n", ticket.Owner.EmailAddress)
			fmt.Printf("RT_REQUESTOR_NAME=%s\n", ticket.Requestor[0].Name)
			fmt.Printf("RT_REQUESTOR_EMAIL=%s\n", ticket.Requestor[0].EmailAddress)
			fmt.Printf("RT_STATUS=%s\n", ticket.Status)
			for _, customField := range ticket.CustomFields {
				if len(customField.Values) > 0 {
					fmt.Printf("RT_CUSTOM_FIELD_%s=%v\n", strings.ReplaceAll(customField.Name, " ", "-"), customField.Values)
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVarP(&ticketID, "ticket-id", "t", 0, "ID del ticket")
	cmd.MarkFlagRequired("ticket-id")

	return cmd
}

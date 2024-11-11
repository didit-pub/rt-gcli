package commands

import (
	"fmt"

	"github.com/didit-pub/rt-gcli/internal/client"
	"github.com/didit-pub/rt-gcli/internal/models"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var (
		ticketID     int
		status       string
		customFields map[string]string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Actualiza un ticket",
		RunE: func(cmd *cobra.Command, args []string) error {
			rtClient, err := client.NewClient(&cfg)
			if err != nil {
				return fmt.Errorf("error al crear el cliente: %w", err)
			}

			update := &models.TicketUpdate{
				Status:       &status,
				CustomFields: customFields,
			}

			err = rtClient.UpdateTicket(ticketID, update)
			if err != nil {
				return fmt.Errorf("error al actualizar el ticket: %w", err)
			}

			if !silent {
				fmt.Printf("RT_TICKET_ID=%d\n", ticketID)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&ticketID, "ticket-id", "t", 0, "ID del ticket")
	cmd.Flags().StringVarP(&status, "status", "s", "", "Estado del ticket")
	cmd.Flags().StringToStringVarP(&customFields, "custom-field", "f", map[string]string{}, "Campos personalizados del ticket")
	cmd.MarkFlagRequired("ticket-id")
	return cmd
}

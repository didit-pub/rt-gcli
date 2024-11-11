package commands

import (
	"fmt"

	"github.com/didit-pub/rt-gcli/internal/client"
	"github.com/didit-pub/rt-gcli/internal/models"
	"github.com/spf13/cobra"
)

func newCommentCmd() *cobra.Command {
	var (
		ticketID   int
		message    string
		correspond bool
	)

	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Agrega un comentario a un ticket",
		RunE: func(cmd *cobra.Command, args []string) error {
			rtClient, err := client.NewClient(&cfg)
			if err != nil {
				return fmt.Errorf("error al crear el cliente: %w", err)
			}

			comment := &models.Comment{
				Content:     message,
				ContentType: "text/plain",
			}

			if correspond {
				err = rtClient.CorrespondTicket(ticketID, comment)
			} else {
				err = rtClient.CommentTicket(ticketID, comment)
			}
			if err != nil {
				return fmt.Errorf("error al crear el comentario: %w", err)
			}

			if !silent {
				fmt.Printf("RT_TICKET_ID=%d\n", ticketID)
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&ticketID, "ticket-id", "t", 0, "ID del ticket")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Mensaje del comentario")
	cmd.Flags().BoolVarP(&correspond, "correspond", "c", false, "Agregar una respuesta al ticket")
	cmd.MarkFlagRequired("ticket-id")
	cmd.MarkFlagRequired("message")
	return cmd
}

package yasmtp

import (
	"context"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	ctx := context.Background()

	err := SendHTML(ctx, &Input{
		From: From{
			ServerHostPort: GmailSMTPHostPort,
			Email:          "your_gmail_acc@gmail.com",
			Auth: smtp.PlainAuth(
				"",
				"yout_gamil_password",
				"your_gmail_password",
				"smtp.gmail.com",
			),
		},
		To: []To{{
			Email: "reciver_email@example.com",
		}},
		Msg: Message{
			Subject: "Your email sub",
			Body:    "Test",
		},
	})
	require.NoError(t, err)
}

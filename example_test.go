package yasmtp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	ctx := context.Background()
	err := SendHTML(ctx, &Input{
		From: From{
			Email:          "your_gmail_acc@gmail.com",
			Password:       "your_gmail_acc_password",
			ServerHostPort: "smtp.gmail.com:587",
		},
		To: []To{{
			Email: "reciver_email@example.com",
		}},
		Msg: Message{
			Subject: "Your email subject",
			Body:    "Some email body",
		},
	})
	require.NoError(t, err)
}

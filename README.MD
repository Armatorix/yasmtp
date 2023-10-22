# Yet Another SMTP client

Created for simplare SMTP client usage, as most of the time you simply want to send email from your gmail account.

Additionally it handles context - which does not come by default from net/smtp.

## Example
```go
package main

func main() {
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
    if err != nil {
        panic(err)
    }
}

```

aaaand that's it.

## Links
* [send with gmail](https://www.123formbuilder.com/docs/setting-up-smtp-with-2-step-verification-on/)

## TODO:
[] option for local testing
[] handling non-html

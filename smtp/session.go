package smtp

import (
    "io"
    "log"
    "strings"
    "time"

    "github.com/emersion/go-smtp"
    "github.com/google/uuid"
    "github.com/mnako/letters"

    "github.com/forscht/watt/db"
)

// Session represents a SMTP session.
type Session struct {
    dataHandler DataHandler // DataHandler is an interface to handle email data.
}

// AuthPlain is a callback function for handling PLAIN authentication.
func (s *Session) AuthPlain(username, password string) error {
    return nil
}

// Mail is a callback function for handling MAIL command.
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
    return nil
}

// Rcpt is a callback function for handling RCPT command.
func (s *Session) Rcpt(to string) error {
    return nil
}

// Data is a callback function for handling DATA command. It reads the email data from the provided io.Reader,
// parses it using letters library, and calls the dataHandler with the parsed email data.
func (s *Session) Data(r io.Reader) error {
    // Generate a new UUID to use as the email ID
    id := uuid.New()

    // Parse the email data from the io.Reader using the letters library
    email, err := letters.ParseEmail(r)
    if err != nil {
        return err
    }

    // Loop through the email's "To" addresses and insert db.Message into cache
    for _, addr := range email.Headers.To {
        if addr != nil {
            message := db.Message{
                ID:      strings.ReplaceAll(id.String(), "-", ""), // Set the email ID as the UUID without dashes
                To:      addr.Address,                             // Set the "To" address of the email
                From:    email.Headers.From[0].Address,            // Set the "From" address of the email
                Date:    email.Headers.Date.Format(time.RFC3339),  // Format and set the email's date in RFC3339 format
                Subject: email.Headers.Subject,                    // Set the email's subject
                Text:    email.Text,                               // Set the plain text body of the email
                HTML:    email.HTML,                               // Set the HTML body of the email
            }
            log.Printf("new email -> to:%v from:%v subject:%v", message.To, message.From, message.Subject)
            s.dataHandler(message)
        }
    }

    return nil // Return nil indicating success
}

// Reset is a callback function for handling RSET command.
func (s *Session) Reset() {}

// Logout is a callback function for handling LOGOUT command.
func (s *Session) Logout() error {
    return nil
}

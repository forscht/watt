package smtp

import (
    "log"
    "time"

    "github.com/emersion/go-smtp"

    "github.com/forscht/watt/db"
)

type DataHandler func(message db.Message)

func Serv(
    addr string,                // Listening address of the SMTP server
    readTimeOut time.Duration,  // Read timeout duration for incoming connections
    writeTimeOut time.Duration, // Write timeout duration for outgoing data
    maxMessageBytes int,        // Maximum allowed size of an email message in bytes
    handler DataHandler,
) error {
    // Create a new backend instance
    be := &Backend{
        dataHandler: handler,
    }

    // Create a new SMTP server instance using the backend
    s := smtp.NewServer(be)

    // Set configurations for the SMTP server
    s.Addr = addr                       // Set listening address
    s.ReadTimeout = readTimeOut         // Set read timeout
    s.WriteTimeout = writeTimeOut       // Set write timeout
    s.MaxMessageBytes = maxMessageBytes // Set maximum message size
    s.MaxRecipients = 10                // Set maximum recipients
    s.AuthDisabled = true               // Set authentication disabled flag

    log.Println("SMTP server listening on", s.Addr) // Log server starting with listening address
    return s.ListenAndServe()                       // Start SMTP server and check for errors
}

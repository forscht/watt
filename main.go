package main

import (
    "fmt"
    "log"
    "os"

    "github.com/alecthomas/kingpin/v2"
    "github.com/joho/godotenv"

    "github.com/forscht/watt/db"
    "github.com/forscht/watt/http"
    "github.com/forscht/watt/smtp"
)

var (
    app             = kingpin.New("watt", "Watt: smtp wrapper for temp mail with web based interface").Version(version())
    port            = app.Flag("port", "Port number to start the HTTP server on").Default("3000").Envar("PORT").Int()
    saddr           = app.Flag("saddr", "Address to start the SMTP server on").Default(":25").Envar("SMTP_ADDR").String()
    readTimeout     = app.Flag("readtimeout", "Set the read timeout duration for the SMTP server").Default("30s").Envar("READ_TIMEOUT").Duration()
    writeTimeout    = app.Flag("writetimeout", "Set the write timeout duration for the SMTP server").Default("30s").Envar("WRITE_TIMEOUT").Duration()
    maxMessageBytes = app.Flag("maxmessagebytes", "Set the maximum email size in bytes for the SMTP server").Envar("MAX_MESSAGE_BYTES").Default("1048576").Int()
    domain          = app.Flag("domain", "Domain name for SMTP server. Example: 'spamok.org'").Required().Envar("DOMAIN").String()
    name            = app.Flag("name", "This will be replaced for 'Watt' in webpage").Default("Watt").Envar("NAME").String()
    ttl             = app.Flag("ttl", "Set the time-to-live duration for the mail cache").Default("30m").Envar("TTL").Duration()
)

func main() {
    // Load env file.
    godotenv.Load()

    // Load options
    kingpin.MustParse(app.Parse(os.Args[1:]))

    // Initialize a new TTL cache with the duration provided in the command-line flag
    cache := db.NewCache(*ttl)

    // Create a new Watt instance with the cache
    watt := NewWatt(cache)

    go func() {
        // Start the SMTP server with the provided parameters and the Watt's email handler function
        err := smtp.Serv(*saddr, *readTimeout, *writeTimeout, *maxMessageBytes, watt.HandleEmail)
        if err != nil {
            log.Fatalln(err)
        }
    }()

    // Prepare strings to be replaced in HTML
    replacer := map[string]string{
        "version":     version(),
        "domain":      *domain,
        "serviceName": *name,
    }
    // Start HTTP server
    addr := fmt.Sprintf(":%d", *port)
    err := http.Serv(addr, cache, replacer, watt.HandleNewWsCh)
    if err != nil {
        log.Fatalln(err)
    }
}

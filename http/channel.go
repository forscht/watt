package http

import (
    "encoding/json"
    "fmt"
    "net"

    "github.com/gobwas/ws/wsutil"

    "github.com/forscht/watt/db"
)

// Channel represents a websocket channel for a specific email address.
type Channel struct {
    Addr  string   // email address for this channel
    Conn  net.Conn // ws connection
    Alive bool
}

// Write writes a message to the websocket channel.
func (c *Channel) Write(message db.Message) error {
    // If there is an error during marshaling or writing to
    // the websocket connection, the connection will be closed.
    data, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("[channel] JSON marshal error", err.Error())
    }
    // Write JSON message to websocket and if failed close the connection
    err = wsutil.WriteServerText(c.Conn, data)
    if err != nil {
        return fmt.Errorf("[channel] ws write error", err.Error())
    }
    return nil
}

// Close closes the websocket channel.
func (c *Channel) Close() {
    c.Conn.Close()
    c.Alive = false
}

// Watt Package routes incoming emails to incoming websockets with ttl cache
package main

import (
    "sync"

    "github.com/forscht/watt/db"
    "github.com/forscht/watt/http"
)

// Watt represents the main application struct.
type Watt struct {
    cache *db.Cache
    store map[string][]*http.Channel
    mu    sync.Mutex
}

// NewWatt creates a new Watt instance with the given cache.
func NewWatt(cache *db.Cache) Watt {
    return Watt{
        cache: cache,
        store: make(map[string][]*http.Channel),
    }
}

// HandleEmail handles an incoming email message and stores it in the cache,
// then sends it to WebSocket channels associated with the email address.
func (c *Watt) HandleEmail(message db.Message) {
    // Lock the mutex before accessing the store map
    c.mu.Lock()
    defer c.mu.Unlock()

    // Add message to cache
    c.cache.Add(message.To, message)

    // Get the slice of WebSocket channels associated with the email address
    channels, ok := c.store[message.To]
    if ok {
        // Create a new slice to store the WebSocket channels to remove
        var channelsToRemove []*http.Channel

        // Iterate over each WebSocket channel in the slice and send the message
        for _, ch := range channels {
            err := ch.Write(message)
            if err != nil {
                // Add the WebSocket channel to the channelsToRemove slice
                channelsToRemove = append(channelsToRemove, ch)
            }
        }

        // Remove the WebSocket channels that encountered errors from the slice
        for _, ch := range channelsToRemove {
            for i, channel := range c.store[message.To] {
                if channel == ch {
                    c.store[message.To] = append(c.store[message.To][:i], c.store[message.To][i+1:]...)
                    break
                }
            }
        }
    }
}

// HandleNewWsCh handles a new WebSocket channel and adds it to the store.
func (c *Watt) HandleNewWsCh(ch *http.Channel) {
    // Lock the mutex before accessing the store map
    c.mu.Lock()
    // Check if the slice for the email address exists, if not, initialize it
    if c.store[ch.Addr] == nil {
        c.store[ch.Addr] = []*http.Channel{}
    }
    // Append the WebSocket channel to the slice
    c.store[ch.Addr] = append(c.store[ch.Addr], ch)
    // Unlock the mutex
    c.mu.Unlock()
    // Fetch all the messages from cache and push it to channel
    messages := c.cache.GetByAddr(ch.Addr)
    for _, msg := range messages {
        ch.Write(msg)
    }
}

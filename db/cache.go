package db

import (
    "time"

    "github.com/yudai/ttlslicemap"
)

type Message struct {
    ID      string
    To      string
    From    string
    Date    string
    Subject string
    Text    string
    HTML    string
}

type Cache struct {
    db *ttlslicemap.TTLSliceMap
}

// Create a new Cache struct and initialize it with the given TTL duration
func NewCache(ttl time.Duration) *Cache {
    cache := &Cache{}
    cache.Init(ttl)

    return cache
}

// Initialize the Cache for given TTL
func (c *Cache) Init(ttl time.Duration) {
    c.db = ttlslicemap.New(ttl)
}

// Add Message to cache
func (c *Cache) Add(addr string, message Message) {
    // Store message with messageId -> message
    c.db.Add(message.ID, message)
    // Store message by address messageAddr -> messageId
    c.db.Add(addr, message.ID)
}

// GetByAddr Return all []Message with matching Message.To
func (c *Cache) GetByAddr(addr string) []Message {
    // Get the list of Message IDs from the cache using the given address as the key
    msgIds, ok := c.db.Get(addr)
    if !ok {
        return make([]Message, 0) // If the key is not found, return an empty slice of Message
    }
    // Clear msgs slice
    msgs := make([]Message, 0)
    // Retrieve each Message from the cache using the Message ID as the key
    for _, msgId := range msgIds {
        msg, ok := c.db.Get(msgId.(string))
        if ok {
            // Type assert the retrieved Message from interface{} to Message and append it to the slice of Message to be returned
            msgs = append(msgs, msg[0].(Message))
        }
    }
    return msgs
}

// GetById Return Message with matching Message.ID
func (c *Cache) GetById(id string) (Message, bool) {
    msgs, ok := c.db.Get(id)
    if !ok {
        return Message{}, false
    }
    return msgs[0].(Message), true
}

package smtp

import (
    "github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct {
    dataHandler DataHandler
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
    session := &Session{dataHandler: bkd.dataHandler}
    return session, nil
}

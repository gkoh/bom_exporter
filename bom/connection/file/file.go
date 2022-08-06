package file

import (
	"io/ioutil"
)

// Connection holds a filepath based Retriever.
type Connection struct {
	filepath string
}

// New implements the Retriever interface.
func New(id string) *Connection {
	return &Connection{filepath: id}
}

// Identifier implements the Retriever interface.
func (c *Connection) Identifier() string {
	return c.filepath
}

// Retrieve implements the Retrieve interface.
func (c *Connection) Retrieve() ([]byte, error) {
	return ioutil.ReadFile(c.filepath)
}

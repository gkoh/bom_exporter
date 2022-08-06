package ftp

import (
	"bufio"
	"bytes"
	"fmt"
	ftpClient "github.com/gonutz/ftp-client/ftp"
	log "github.com/sirupsen/logrus"
	"strings"
)

const ftpBom = "ftp.bom.gov.au"
const ftpPort = 21

// Connection holds the details for an FTP based Retriever.
type Connection struct {
	id      string
	address string
	path    string
	conn    *ftpClient.Connection
}

// New implements the Retriever interface.
func New(id string) *Connection {
	return &Connection{id: id, address: ftpBom, path: "anon/gen/fwo/" + id + ".xml"}
}

// Identifier implements the Retriever interface.
func (c *Connection) Identifier() string {
	return c.id
}

// Retrieve implements the Retriever interface.
func (c *Connection) Retrieve() ([]byte, error) {
	var err error

	c.conn, err = ftpClient.Connect(ftpBom, ftpPort)
	if err != nil {
		log.Errorf("Failed to connect to '%s': %s", ftpBom, err)
		return nil, err
	}
	defer c.conn.Close()

	err = c.conn.Login("anonymous", "")
	if err != nil {
		log.Errorf("Failed to login: %s", err)
		return nil, err
	}
	defer c.conn.Quit()

	// check file path
	_, status, err := c.conn.StatusOf(c.path)
	if err != nil {
		log.Errorf("Status failed for '%s'", c.path)
		return nil, err
	}
	if !strings.Contains(status, c.id) {
		return nil, fmt.Errorf("Failed to find '%s'", c.id)
	}

	var data bytes.Buffer
	dataWriter := bufio.NewWriter(&data)
	err = c.conn.Download(c.path, dataWriter)
	if err != nil {
		log.Errorf("Failed to download '%s': %s", c.path, err)
		return nil, err
	}

	return data.Bytes(), nil
}

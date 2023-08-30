package connection

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	clientconfig "github.com/ProjectOrangeJuice/vm-manager-client/clientConfig"
	"github.com/ProjectOrangeJuice/vm-manager-client/update"
	"github.com/ProjectOrangeJuice/vm-manager-server/shared"
)

type Connection struct {
	Conn   net.Conn
	Config *clientconfig.Config
}

func NewConnection(conn net.Conn, config *clientconfig.Config) Connection {
	return Connection{
		Conn:   conn,
		Config: config,
	}
}

func (c *Connection) ProcessLines() {
	reader := bufio.NewReader(c.Conn)
	for {
		// Read a line of data
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		// Print the line
		fmt.Println(line)
		c.readLine(line)
	}
	log.Printf("Disconnected, trying again in 5 seconds")
	c.Conn.Close()
	time.Sleep(5 * time.Second)

}

func (c *Connection) readLine(line string) {
	switch strings.TrimSpace(line) {
	case "STORAGE_INFO":
		c.sendBackStorage()
	case "SYSTEM_INFO":
		c.sendBackSystem()
	case "UPDATE_NOW":
		update.UpdateIfNeeded(c.Config)
	}
}

// A generic function that creates an event.
func createEvent[R any](request string, result R) ([]byte, error) {
	resultByte, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("could not marshal result, %s", err)
	}

	evt := shared.EventData{
		Request: request,
		Result:  resultByte,
	}
	eventOut, err := json.Marshal(evt)
	if err != nil {
		return nil, fmt.Errorf("could not marshal event, %s", err)
	}
	return eventOut, nil
}

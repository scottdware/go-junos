package junos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"log"
)

// Session holds the connection information to our Junos device.
type Session struct {
	Conn *netconf.Session
}

// RollbackXML parses our configuration after requesting it via rollback.
type RollbackXML struct {
	XMLName xml.Name `xml:"rollback-information"`
	Config  string   `xml:"configuration-information>configuration-output"`
}

// NewSession establishes a new connection to a Junos device that we will use
// to run our commands against.
func NewSession(host, user, password string) *Session {
	s, err := netconf.DialSSH(host, netconf.SSHConfigPassword(user, password))
	if err != nil {
		log.Fatal(err)
	}

	return &Session{
		Conn: s,
	}
}

// Lock locks the candidate configuration.
func (s *Session) Lock() {
	lockRPC := "<rpc><lock><target><candidate/></target></lock></rpc>"
	resp, _ := s.Conn.Exec(lockRPC)
	// if err != nil {
	// log.Fatal(err)
	// }

	if resp.Ok == false {
		for _, m := range resp.Errors {
			fmt.Printf("%s\n", m.Message)
		}
	}
}

// Unlock unlocks the candidate configuration.
func (s *Session) Unlock() {
	unlockRPC := "<rpc><unlock><target><candidate/></target></unlock></rpc>"
	resp, _ := s.Conn.Exec(unlockRPC)
	// if err != nil {
	// fmt.Printf("Error: %+v\n", err)
	// }

	if resp.Ok == false {
		for _, m := range resp.Errors {
			fmt.Printf("%s\n", m.Message)
		}
	}
}

// GetRollbackConfig returns the configuration of the given rollback state.
func (s *Session) GetRollbackConfig(number int) (string, error) {
	rb := &RollbackXML{}
	rpcCommand := fmt.Sprintf("<rpc><get-rollback-information><rollback>%d</rollback><format>text</format></get-rollback-information></rpc>", number)
	reply, err := s.Conn.Exec(rpcCommand)
    
    if err != nil {
        log.Fatal(err)
    }
    
	if reply.Ok == false {
		for _, m := range reply.Errors {
			return "", errors.New(m.Message)
		}
	}

	err = xml.Unmarshal([]byte(reply.Data), rb)
	if err != nil {
		log.Fatal(err)
	}

	return rb.Config, nil
}

// Close disconnects and closes the session to our Junos device.
func (s *Session) Close() {
	s.Conn.Close()
}

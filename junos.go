package junos

import (
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"log"
)

type Session struct {
	Conn *netconf.Session
}

func NewSession(host, user, password string) *Session {
    s, err := netconf.DialSSH(host, netconf.SSHConfigPassword(user, password))
    if err != nil {
        log.Fatal(err)
    }

    return &Session{
        Conn: s,
    }
}

func (s *Session) Lock() {
	resp, err := s.Conn.Exec("<rpc><get-software-information/></rpc>")

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
    
    fmt.Printf("%+v\n", resp)
}

func (s *Session) Unlock() {
	resp, err := s.Conn.Exec("<rpc><get-software-information/></rpc>")

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
    
    fmt.Printf("%+v\n", resp)
}

func (s *Session) Close() {
    s.Conn.Close()
}
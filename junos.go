package junos

import (
	"fmt"
	"github.com/Juniper/go-netconf/netconf"
	"log"
)

type Session struct {
	Conn string
}

func NewSession(host, user, password string) *Session {
	s, err := netconf.DialSSH(host, netconf.SSHConfigPassword(user, password))
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	return &Session{
		Conn: s,
	}
}

func (s *Session) Lock() {
	resp, err := s.Conn.Exec("<rpc><lock-configuration/></rpc>")

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
}

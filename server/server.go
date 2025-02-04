package server

import (
	"github.com/ihciah/rabbit-tcp/logger"
	"github.com/ihciah/rabbit-tcp/peer"
	"github.com/ihciah/rabbit-tcp/tunnel"
	"net"
)

type Server struct {
	peerGroup peer.PeerGroup
	logger    *logger.Logger
}

func NewServer(cipher tunnel.Cipher) Server {
	return Server{
		peerGroup: peer.NewPeerGroup(cipher),
		logger:    logger.NewLogger("[Server]"),
	}
}
func ServeThread(address string, ss Server ) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			ss.logger.Errorf("Error when accept connection: %v.\n", err)
			continue
		}
		err = ss.peerGroup.AddTunnelFromConn(conn)
		if err != nil {
			ss.logger.Errorf("Error when add tunnel to tunnel pool: %v.\n", err)
		}
	}
}
func (s *Server) Serve(addresses []string) error {
	for _,address:= range addresses {
		go error=ServeThread(address,s)
	}
	return 
}

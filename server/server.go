package server

import (
	"github.com/aagun1234/rabbit-tcp/logger"
	"github.com/aagun1234/rabbit-tcp/peer"
	"github.com/aagun1234/rabbit-tcp/tunnel"
	"net"
	"sync"
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


func ServeThread(address string, ss *Server, wg *sync.WaitGroup) error {
	defer wg.Done()
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
        var wg sync.WaitGroup
	for _,address:= range addresses {
		wg.Add(1)
		s.logger.Infof("Serve on %s\n",address)
		go ServeThread(address,s, &wg)
	}
	wg.Wait()
	return nil
}

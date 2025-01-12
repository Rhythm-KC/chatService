package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/rhythm/chatservice/protocol"
)
type Server struct{
    connectionManager *connectionManager
    roomManager       *roomManager
    wg                sync.WaitGroup
    listner           net.Listener
    done              chan struct{}
}

func async(f func (), wg *sync.WaitGroup){
    wg.Add(1)
    go func(){
        defer wg.Done()
        f()

    }()
}

func NewServer(port uint32)(*Server){
    listner, _ := net.Listen("tcp", fmt.Sprintf("%d", port))
    roommangar := NewRoomManager()
    connectionManager := NewConnectionManger(roommangar)
    server := Server{connectionManager: connectionManager, 
                     roomManager: roommangar,
                     listner: listner}
    return &server

}

func (s *Server) Start(){
    async(func(){s.roomManager.Start()}, &s.wg)
    async(func(){s.connectionManager.Start()}, &s.wg)
    for {
        select{
            case <- s.done:
                return 
            default:
                conn, _  := s.listner.Accept()
                newConn := protocol.NewServerConnection(&conn)
                s.connectionManager.Submit(newConn)
        }
    }
}

func (s *Server)Close(){

    close(s.done)
    async(func(){s.connectionManager.Close()}, &s.wg)
    async(func(){s.roomManager.Stop()}, &s.wg)
    s.listner.Close()
    s.wg.Wait()
}

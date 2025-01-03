package protocol

import (
	"net"

	pe "github.com/rhythm/chatservice/protocol/error"
	protocol "github.com/rhythm/chatservice/protocol/util"
)


type ServerConnection struct{
    client net.Conn
}

func NewServerConnection(conn *net.Conn) *ServerConnection{

    return &ServerConnection{*conn}
}

func (conn *ServerConnection)Close(){
    conn.client.Close()
}

func (conn *ServerConnection) Listen() ([]byte, *pe.ProtocolError){

    var buf = make([]byte, 1024) 
    
    n, err := conn.client.Read(buf)
    if err != nil{
        return nil, pe.TimedOutError("Reading from the server timedout")
    }

    er := protocol.CheckForValidHeader(buf[0], buf[1:n])
    if er != nil{
        return nil, er
    }
    return buf[:n], nil

}





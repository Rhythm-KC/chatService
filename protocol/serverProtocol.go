package protocol

import (
	"net"
	"time"

	pe "github.com/rhythm/chatservice/protocol/error"
	protocol "github.com/rhythm/chatservice/protocol/util"
)

type ServerConnection struct{
    client net.Conn
    connectionBuffer []byte
}


func NewServerConnection(conn *net.Conn) *ServerConnection{

    return &ServerConnection{*conn, make([]byte, 1024)}

}

func (conn *ServerConnection)Close(){
    conn.client.Close()
}

func (conn *ServerConnection) Listen(timeout uint) (Message, *pe.ProtocolError){


    conn.client.SetReadDeadline(time.
                                Now().
                                Add(time.Duration(timeout) * time.Second)) 

    n, err := conn.client.Read(conn.connectionBuffer)
    if err != nil{
        return nil, pe.TimedOutError("Reading from the server timedout")
    }

    er := protocol.CheckForValidHeader(conn.connectionBuffer[0], 
                                       conn.connectionBuffer[1:])
    if er != nil{
        return nil, er
    }
    return Unmarshal(conn.connectionBuffer[:n])

}

func (conn *ServerConnection) SendResponse(msg Message){
    data, _:= msg.Marshal()
    conn.client.Write(data)
}





package protocol

import (
	"net"
	"time"

	pe "github.com/rhythm/chatservice/protocol/error"
	protocol "github.com/rhythm/chatservice/protocol/util"
)

type ServerConnection struct{
    client net.Conn
    readBuffer []byte
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

    n, err := conn.client.Read(conn.readBuffer)
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        return nil, pe.TimedOutError("Reading from the server timedout")
    }

    if err.Error() == "EOF"{
        return nil, pe.DisconectedError("Connection Lost")
    }


    er := protocol.CheckForValidHeader(conn.readBuffer[0], 
                                       conn.readBuffer[1:])
    if er != nil{
        return nil, er
    }
    return Unmarshal(conn.readBuffer[:n])

}

func (conn *ServerConnection) SendResponse(msg Message){
    data, _:= msg.Marshal()
    conn.client.Write(data)
}





package protocol

import (
	"net"
	"time"

	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	protocol "github.com/rhythm/chatservice/pkg/protocol/util"
)

type Connection struct{
    client net.Conn
    readBuffer []byte
}

var listner net.Listener

func NewConnection(conn *net.Conn) *Connection{

    return &Connection{*conn, make([]byte, 1024)}

}

func Listen(address string)(*Connection, error){
    if listner == nil{
        listner, _ = net.Listen("tcp", address)
    }
    conn, err := listner.Accept()
    if err != nil{
        return nil, err
    }
    return NewConnection(&conn), nil

}
func Close(){
    if listner != nil{
        listner.Close()
    }
}

func (conn *Connection)Close(){
    conn.client.Close()
}

func (conn *Connection) Listen(timeout uint) (Message, *pe.ProtocolError){


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

func (conn *Connection) SendResponse(msg Message){
    data, _:= msg.Marshal()
    conn.client.Write(data)
}


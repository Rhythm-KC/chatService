package client

import (
	"github.com/rhythm/chatservice/utils"

	"fmt"
	"net"
	"time"
)

type connection struct{
    server net.Conn
    readBuffer []byte
}

func (conn *connection)close(){
    conn.server.Close()
}

func (c *connection)sendMessage(userMsg utils.Message) error{
    msg, err := userMsg.Marshal()
    if err != nil{
        return err
    }
    _, err = c.server.Write(msg)
    if err != nil{
        return err
    }
    return nil

}

func (conn *connection) createRoom() (*utils.CreateRoomInst, error){
    roomInfo := getRoomNameFromUser(conn.stdinReader)
    marshaledRoomInfo, _ := roomInfo.Marshal()
    conn.server.Write(marshaledRoomInfo)
    msg, err := conn.readwithtimeout(5)
    if err != nil{
        return nil, err 
    }
    room, ok := msg.(*utils.CreateRoomInst)
    if ok{
        return room, nil
    }
    return nil, fmt.Errorf("Did not receive created room from server")

}

func (conn *connection)readwithtimeout(second uint8) (utils.Message, error){

    timeout := time.Second * time.Duration(second)
    conn.server.SetReadDeadline(time.Now().Add(timeout))
    conn.server.Read(conn.readBuffer)

    return utils.Unmarshal(conn.readBuffer)
}

func (conn *connection) registerUser() (*utils.RegisterUser, error) {
    user := getClientForRegistration(conn.stdinReader)
    seralizedUser, err := user.Marshal()
    if err != nil{
        return nil, err
    }

    conn.server.Write(seralizedUser)

    msg, err := conn.readwithtimeout(5)
    if err != nil{
        return nil, fmt.Errorf("Network timeout: Cannot confirm if the user is registered")
    }

    userWithId, ok := msg.(*utils.RegisterUser)
    if ok{
        return userWithId, nil
    }

    return nil, fmt.Errorf("Ivalid Response from the server")
}

func (conn *connection) createRoom()

package protocol

import (
	"net"
	"time"

	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/request"
	"github.com/rhythm/chatservice/pkg/protocol/response"
)



type MessageClient struct{
    tcpClient net.Conn
}

func Connect(address string) (*MessageClient, *pe.ProtocolError){

    conn, err := net.Dial("tcp", address)
    if err != nil{
        return nil, pe.NotFoundError("Could not connect to resource")
    }

    return &MessageClient{conn}, nil
}

func (client *MessageClient) Close(){
    client.tcpClient.Close()

}

func (client *MessageClient) waitResponsewithtimeout(second uint8)([]byte, 
                                                *pe.ProtocolError){
    buf := make([]byte, 1024) 
    timeout := time.Second * time.Duration(second)
    client.tcpClient.SetReadDeadline(time.Now().Add(timeout))
    n, err := client.tcpClient.Read(buf)

    if err!= nil{
        return nil, pe.TimedOutError("Response Timed Out")
    }
    return buf[:n], nil
}


func (client *MessageClient) SendMessage(senderId uint16, 
                                         roomId uint16, 
                                         msg string) (*response.SendMessageResponse,
                                                      *pe.ProtocolError) {
    
    message, err:= request.NewSendMessageRequest(senderId, roomId, msg)
    if err!= nil {
        // internal protocol error
        return nil, err
    }
    serailData, serialError := message.Marshal()
    if serialError!= nil {
        return nil, serialError
    }

    _, er := client.tcpClient.Write(serailData)
    if er!= nil {
        return nil,
               pe.TimedOutError("Could not send message to the server")
    }
    
    resp, err := client.waitResponsewithtimeout(5)
    var response response.SendMessageResponse
    err = response.Unmarshal(resp)
    if err != nil{
        return nil, err
    }
    return &response, nil
}

func (client *MessageClient) CreateRoom(senderId uint16, roomName string)(
                                        *response.CreateRoomResponse,
                                        *pe.ProtocolError){

    room, err := request.NewCreateRoomRequest(senderId,roomName)
    if err != nil{
        return nil, err
    }

    var buf []byte
    buf, err = room.Marshal()
    if err != nil{
        return nil, err
    }

    _, er := client.tcpClient.Write(buf)
    if er != nil{
        return nil, 
               pe.TimedOutError("Could not write to the server")
    }
    
    buf, err = client.waitResponsewithtimeout(5) 
    var createRoomResp response.CreateRoomResponse
    err = createRoomResp.Unmarshal(buf)
    if err != nil{
        return nil, err
    }

    return &createRoomResp, nil
}

func (client *MessageClient) JoinRoom(senderId uint16,
                                      roomId uint16)(*response.JoinRoomResponse,
                                                     *pe.ProtocolError){

    request := request.NewJoinRoomRequest(senderId, roomId)
    buf, err := request.Marshal()
    if err != nil{
        return nil, err
    }
    _, er:=client.tcpClient.Write(buf)
    if er != nil{
        return nil, pe.TimedOutError("Could not send Join Room Request")
    }

    buf, err = client.waitResponsewithtimeout(5)
    var joinRoomRsp response.JoinRoomResponse
    err = joinRoomRsp.Unmarshal(buf)

    if err != nil{
        return nil, err
    }

    return &joinRoomRsp, nil
}

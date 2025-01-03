package protocol

import (
	"fmt"

	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
	"github.com/rhythm/chatservice/protocol/request"
	"github.com/rhythm/chatservice/protocol/response"
)

type Message interface{
    Marshal()([]byte, *pe.ProtocolError)
    Unmarshal(dataReceived []byte) *pe.ProtocolError
}

func Unmarshal(data []byte) (Message, error){
    header := data[0] 

    switch header{
        //request
        case messagecode.JoinRoomRequestIdentifier:
            var obj request.JoinRoomRequest
            err := obj.Unmarshal(data)
            return &obj, err

        case messagecode.CreateRoomRequestIdentifier:
            var obj request.CreateRoomRequest
            err := obj.Unmarshal(data)
            return &obj, err

        case messagecode.SendMessageRequestIdentifier:
            var obj request.SendMessageRequest
            err := obj.Unmarshal(data)
            return &obj, err

        // Response
        case messagecode.MessageResponseIdentifier:
            var obj response.SendMessageResponse
            err := obj.Unmarshal(data)
            return &obj, err

        case messagecode.JoinRoomResponseIdentifier:
            var obj response.JoinRoomResponse
            err := obj.Unmarshal(data)
            return &obj, err

        case messagecode.CreateRoomResponseIdentifier:
            var obj response.CreateRoomResponse
            err := obj.Unmarshal(data)
            return &obj, err

        case messagecode.IsErrorMessage:
            var obj pe.ProtocolError
            err := obj.Unmarshal(data[1:])
            return &obj, err
        default:
            return nil, fmt.Errorf("Not a valid header: %d", header)
    }
}


package protocol

import (
	"fmt"

	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/messagecode"
	"github.com/rhythm/chatservice/pkg/protocol/request"
	"github.com/rhythm/chatservice/pkg/protocol/response"
)

// These are the messages that is passed between server and client
type Message interface{
    Marshal()([]byte, *pe.ProtocolError)
    Unmarshal(dataReceived []byte) *pe.ProtocolError
    GetHeader() uint8
}

// Given bytes of data we return the message if the data is valid
func Unmarshal(data []byte) (Message, *pe.ProtocolError){
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
            msg := fmt.Sprintf("%d is an invalid header", header)
            return nil, pe.InvalidMessageHeaderError(msg)
    }
}


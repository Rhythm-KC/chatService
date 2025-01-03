package request

import (
	"bytes"
	"encoding/binary"
	"fmt"
    
    p "github.com/rhythm/chatservice/protocol/util"
	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
)

type CreateRoomRequest struct{
    senderId uint16 
    roomName string
}

func NewCreateRoomRequest(senderId uint16, roomName string) (*CreateRoomRequest,
                                                          *pe.ProtocolError) {
   
    if len(roomName) > max_text_in_bytes{
        return nil, pe.InvalidMessageLengthError("Room Name is to long")
    }
    return &CreateRoomRequest{senderId: senderId, roomName: roomName}, nil

}

func (m *CreateRoomRequest) Marshal()([]byte, *pe.ProtocolError){
    roomLen := len(m.roomName)

    if roomLen > max_text_in_bytes{
        msg := fmt.Sprintf("Room Name cannot be greater than %d", 
                           max_text_in_bytes)
        return nil, 
        pe.InvalidMessageLengthError(msg)
    }

    var buf bytes.Buffer
    
    binary.Write(&buf, binary.BigEndian, 
                 uint8(messagecode.CreateRoomRequestIdentifier))
    binary.Write(&buf, binary.BigEndian, m.senderId)
    binary.Write(&buf, binary.BigEndian, uint8(roomLen))
    binary.Write(&buf, binary.BigEndian, m.roomName)

    return buf.Bytes(), nil
}

func (m *CreateRoomRequest)Unmarshal(data []byte) *pe.ProtocolError{
    header := data[0] 
    if header != messagecode.CreateRoomRequestIdentifier{
        return p.InvalidHeaderError(messagecode.CreateRoomRequestIdentifier,
                                    header,
                                    "CreateRoomRequest")    
    }
    
    l := 0
    sender := binary.BigEndian.Uint16(data[l:(l+sender_id_in_bytes)])
    l+=sender_id_in_bytes

    msgLen := int(data[l])
    l += 1
    roomName := string(data[l:msgLen])

    m.senderId = sender
    m.roomName= roomName

    return nil
}

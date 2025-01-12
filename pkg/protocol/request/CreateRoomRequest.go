package request

import (
	"bytes"
	"encoding/binary"
	"fmt"
    
    p "github.com/rhythm/chatservice/pkg/protocol/util"
	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/messagecode"
)

type CreateRoomRequest struct{
    SenderId uint16 
    header uint8
    RoomName string
}

func NewCreateRoomRequest(senderId uint16, roomName string) (*CreateRoomRequest,
                                                          *pe.ProtocolError) {
   
    if len(roomName) > max_text_in_bytes{
        return nil, pe.InvalidMessageLengthError("Room Name is to long")
    }
    return &CreateRoomRequest{SenderId: senderId,
                              header: messagecode.CreateRoomRequestIdentifier,
                              RoomName: roomName}, nil

}

func (m *CreateRoomRequest) Marshal()([]byte, *pe.ProtocolError){
    roomLen := len(m.RoomName)

    if roomLen > max_text_in_bytes{
        msg := fmt.Sprintf("Room Name cannot be greater than %d", 
                           max_text_in_bytes)
        return nil, 
        pe.InvalidMessageLengthError(msg)
    }

    var buf bytes.Buffer
    
    binary.Write(&buf, binary.BigEndian, m.header)
    binary.Write(&buf, binary.BigEndian, m.SenderId)
    binary.Write(&buf, binary.BigEndian, uint8(roomLen))
    binary.Write(&buf, binary.BigEndian, m.RoomName)

    return buf.Bytes(), nil
}

func (m *CreateRoomRequest)Unmarshal(data []byte) *pe.ProtocolError{
    header := data[0] 
    body := data[1:]

    err := p.CheckForValidHeader(header, body) 
    if err != nil{
        return err
    }
    if header != messagecode.CreateRoomRequestIdentifier{
        return p.InvalidHeaderError(messagecode.CreateRoomRequestIdentifier,
                                    header,
                                    "CreateRoomRequest")    
    }

    m.header = header
    l := 0
    sender := binary.BigEndian.Uint16(data[l:(l+sender_id_in_bytes)])
    l+=sender_id_in_bytes

    msgLen := int(data[l])
    l += 1
    roomName := string(data[l:msgLen])

    m.SenderId = sender
    m.RoomName= roomName

    return nil
}

func (m *CreateRoomRequest) GetHeader() uint8{
    return m.header
}

package response

import (
	"bytes"
	"encoding/binary"

	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/messagecode"
	"github.com/rhythm/chatservice/pkg/protocol/util"
)

type CreateRoomResponse struct{

    senderId uint16
    roomId uint16
    header uint8
    roomName string
}

func NewCreateRoomResponse(senderid uint16, 
                           roomid uint16, 
                           roomName string) *CreateRoomResponse{

    
    return &CreateRoomResponse{senderId: senderid, 
                               roomId: roomid,
                               header: messagecode.CreateRoomResponseIdentifier,
                               roomName: roomName,}
    
}

func (m *CreateRoomResponse) Marshal() ([]byte, *pe.ProtocolError){

    var buf bytes.Buffer
    binary.Write(&buf, binary.BigEndian, m.header)
    binary.Write(&buf, binary.BigEndian, m.senderId)
    binary.Write(&buf, binary.BigEndian, m.roomId)
    roomNameLen := len(m.roomName)
    if roomNameLen > max_text_in_bytes{
        return nil, pe.InvalidMessageLengthError("Room Length is longer than 255 bytes")
    }
    binary.Write(&buf, binary.BigEndian, roomNameLen)
    binary.Write(&buf, binary.BigEndian, m.roomName)
    return buf.Bytes(), nil
}

func (m *CreateRoomResponse) Unmarshal(dataReceived []byte) *pe.ProtocolError{
    
    header := dataReceived[0]
    body := dataReceived[1:]
     err := protocol.CheckForValidHeader(header, body)
     if err != nil{
         return err
     }
     m.header = header
     l := 0
     m.senderId = binary.BigEndian.Uint16(body[l:l+sender_id_in_bytes])
     l+= sender_id_in_bytes

     m.roomId = binary.BigEndian.Uint16(body[l:l+room_id_in_bytes])
     l+= room_id_in_bytes 

     msgLen := int(body[l])
     l+=1

     m.roomName = string(body[l:l+msgLen])

    return nil
}

func (m *CreateRoomResponse) GetHeader() uint8{
    return m.header
}







package response

import (
	"bytes"
	"encoding/binary"

	"github.com/rhythm/chatservice/protocol/util"
	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
)

type CreateRoomResponse struct{

    senderId uint16
    roomId uint16
    roomName string
}

func (m *CreateRoomResponse) Marshal() ([]byte, *pe.ProtocolError){

    var buf bytes.Buffer
    binary.Write(&buf, binary.BigEndian, messagecode.CreateRoomResponseIdentifier)
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







package response

import (
	"bytes"
	"encoding/binary"

	p "github.com/rhythm/chatservice/pkg/protocol/util"
	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/messagecode"
)

type JoinRoomResponse struct{
    header uint8
    roomID uint16
    roomName string
}

func NewJoinRoomResponse(roomId uint16, roomName string) *JoinRoomResponse{
    return &JoinRoomResponse{header: messagecode.JoinRoomResponseIdentifier,
                             roomID: roomId,
                             roomName: roomName}
}


func (m *JoinRoomResponse) Marshal() ([]byte, *pe.ProtocolError){

    var buf bytes.Buffer

    binary.Write(&buf, binary.BigEndian, m.header)
    binary.Write(&buf, binary.BigEndian, m.roomID)
    lenName := len(m.roomName)
    binary.Write(&buf, binary.BigEndian, uint8(lenName))
    binary.Write(&buf, binary.BigEndian, m.roomName)
    return buf.Bytes(), nil
}

func (m *JoinRoomResponse) Unmarshal(dataReceived []byte) *pe.ProtocolError{

    head := dataReceived[0]
    body := dataReceived[1:]

    err := p.CheckForValidHeader(head, body)
    if err != nil{
        return err
    }
    
    if head != messagecode.JoinRoomResponseIdentifier{
        return p.InvalidHeaderError(messagecode.JoinRoomResponseIdentifier,
                                    head,
                                    "JoinRoomResponse")
    }
    m.header = head
    m.roomID = binary.BigEndian.Uint16(body[:2])
    lenName := body[2]
    m.roomName = string(body[3:int(lenName)])
    return nil
}

func (m *JoinRoomResponse) GetHeader() uint8{
    return m.header
}


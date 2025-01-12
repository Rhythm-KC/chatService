package request

import (
	"bytes"
	"encoding/binary"

	p "github.com/rhythm/chatservice/pkg/protocol/util"
	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/messagecode"
)

type JoinRoomRequest struct{
    senderId uint16
    RoomId uint16
    header uint8
}

func NewJoinRoomRequest(senderId uint16, roomId uint16) *JoinRoomRequest{
    return &JoinRoomRequest{senderId: senderId, 
                            RoomId: roomId,
                            header: messagecode.JoinRoomRequestIdentifier}
}

func (m *JoinRoomRequest) Marshal() ([]byte, *pe.ProtocolError){

    var buf bytes.Buffer

    binary.Write(&buf, binary.BigEndian, m.header)
    binary.Write(&buf, binary.BigEndian, m.senderId)
    binary.Write(&buf, binary.BigEndian, m.RoomId)

    return buf.Bytes(), nil
}

func (m *JoinRoomRequest) Unmarshal(dataReceived []byte) *pe.ProtocolError{

    header := dataReceived[0]
    body := dataReceived[1:]

    err := p.CheckForValidHeader(header, body) 
    if err != nil{
        return err
    }
    if header != messagecode.JoinRoomRequestIdentifier{
        return p.InvalidHeaderError(messagecode.JoinRoomRequestIdentifier,
                                    header,
                                    "JoinRoomRequest")
    }
    m.header = header

    l := 0
    m.senderId = binary.BigEndian.Uint16(body[l:l+receiver_id_in_bytes])
    l += receiver_id_in_bytes

    m.RoomId = binary.BigEndian.Uint16(body[l:l+room_id_in_bytes])

    return nil
}

func (m *JoinRoomRequest) GetHeader() uint8{
    return m.header
}

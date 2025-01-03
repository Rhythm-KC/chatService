package response

import (
	"bytes"
	"encoding/binary"

	p "github.com/rhythm/chatservice/protocol/util"
	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
)

type JoinRoomResponse struct{
    status uint8
}

func (m *JoinRoomResponse) Marshal() ([]byte, *pe.ProtocolError){

    var buf bytes.Buffer

    binary.Write(&buf, binary.BigEndian, messagecode.JoinRoomResponseIdentifier)
    binary.Write(&buf, binary.BigEndian, m.status)
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

    m.status = body[1]
    return nil
}


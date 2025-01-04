package response

import (
	"bytes"
	"encoding/binary"

	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
	protocol "github.com/rhythm/chatservice/protocol/util"
)

type CreateUserResponse struct{

    userId uint16
    header uint8
}

func NewCreateUserResponse(userId uint16) *CreateUserResponse{
    return &CreateUserResponse{userId: userId, 
                              header: messagecode.CreateUserResponseIdentifier}
}

func (m *CreateUserResponse) Marshal() ([]byte, *pe.ProtocolError){
    
    var buf bytes.Buffer

    binary.Write(&buf, binary.BigEndian, m.header)
    binary.Write(&buf, binary.BigEndian, m.userId)
    return buf.Bytes(), nil
}

func (m *CreateUserResponse) Unmarshal(dataReceived []byte) *pe.ProtocolError{

    header := dataReceived[0]
    body := dataReceived[1:]

    err := protocol.CheckForValidHeader(header, body)
    if err != nil{
        return err
    }

    if header != messagecode.CreateUserResponseIdentifier{
        return protocol.InvalidHeaderError(messagecode.CreateUserResponseIdentifier,
                                           header,
                                           "CreateUserResponse")
    }

    m.header = header
    m.userId = binary.BigEndian.Uint16(body)
    return nil

}

func (m *CreateUserResponse) GetHeader() uint8{
    return m.header
}

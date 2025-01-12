package request

import (
	"bytes"
	"encoding/binary"
	"fmt"

	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/messagecode"
	protocol "github.com/rhythm/chatservice/pkg/protocol/util"
)



type CreateUserRequest struct{
    userName string
    header uint8
}

func NewCreateUserRequest(name string) (*CreateUserRequest, *pe.ProtocolError){
    if len(name) > max_text_in_bytes{
        errmsg := fmt.Sprintf("User Name is longer than %d bytes", 
                              max_text_in_bytes)
        return nil, pe.InvalidMessageLengthError(errmsg)
    }

    return &CreateUserRequest{userName: name, 
                              header: messagecode.CreateUserRequestIdentifier},
                              nil
}

func (m *CreateUserRequest) Name() string{
    return m.userName
}

func (m *CreateUserRequest) Marshal() ([]byte, *pe.ProtocolError){
    
    var buf bytes.Buffer 
    binary.Write(&buf, binary.BigEndian, m.header)
    nameLen := len(m.userName)
    binary.Write(&buf, binary.BigEndian, nameLen)
    binary.Write(&buf, binary.BigEndian, m.userName)
    return buf.Bytes(), nil
}

func (m *CreateUserRequest) Unmarshal(dataReceived []byte) *pe.ProtocolError{
    header := dataReceived[0]
    body := dataReceived[1:]
    err := protocol.CheckForValidHeader(header,body)
    if err != nil{
        return err
    }
    if header != messagecode.CreateUserRequestIdentifier{
        return protocol.InvalidHeaderError(messagecode.CreateUserRequestIdentifier,
                                           header,
                                           "CreateUserRequest")
    }

    m.header = header
    nameLen := int(body[0])
    m.userName = string(body[1:nameLen])
    return nil
}

func (m *CreateUserRequest) GetHeader() uint8{
    return m.header
}

package response

import (
	"bytes"
	"encoding/binary"

	p "github.com/rhythm/chatservice/protocol/util"
	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
)

type SendMessageResponse struct{
    status int
}

func (m *SendMessageResponse) Marshal()([]byte, *pe.ProtocolError){
    var buf bytes.Buffer

    binary.Write(&buf, binary.BigEndian, messagecode.MessageResponseIdentifier)
    binary.Write(&buf, binary.BigEndian, m.status)
    return buf.Bytes(), nil
}

func (m *SendMessageResponse) Unmarshal(dataReceived []byte) *pe.ProtocolError{

    header := dataReceived[0]
    err := p.CheckForValidHeader(header, dataReceived[1:])
    if err != nil{
        return err
    }

    if header != messagecode.MessageResponseIdentifier{
       return p.InvalidHeaderError(messagecode.MessageResponseIdentifier,
                                   header,
                                   "SendMessageResponse") 
    }

    m.status = int(dataReceived[1])
    return nil
}

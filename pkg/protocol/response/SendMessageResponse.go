package response

import (
	"bytes"
	"encoding/binary"

	p "github.com/rhythm/chatservice/pkg/protocol/util"
	pe "github.com/rhythm/chatservice/pkg/protocol/error"
	"github.com/rhythm/chatservice/pkg/protocol/messagecode"
)

type SendMessageResponse struct{
    status int
    header uint8
}

func (m *SendMessageResponse) Marshal()([]byte, *pe.ProtocolError){
    var buf bytes.Buffer

    binary.Write(&buf, binary.BigEndian, m.header)
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
    
    m.header = header
    m.status = int(dataReceived[1])
    return nil
}

func (m *SendMessageResponse) GetHeader() uint8{
    return m.header
}

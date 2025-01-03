package request

import (
	"bytes"
	"encoding/binary"

	p "github.com/rhythm/chatservice/protocol/util"
	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
)


type SendMessageRequest struct{
    sender uint16 
    receiver uint16
    message string
}

func NewSendMessageRequest(senderId uint16, 
                              receiverId uint16, 
                              message string) (*SendMessageRequest, 
                                              *pe.ProtocolError){
    
    if len(message) > max_text_in_bytes{
        return nil, pe.InvalidMessageLengthError("Invalid Message length")

    }
    return &SendMessageRequest{senderId, receiverId, message}, nil

}

func (m *SendMessageRequest) Marshal()([]byte, *pe.ProtocolError){
    
    var buf bytes.Buffer
    msgLength := len(m.message)
    if msgLength > msgLength{
        return nil, 
        pe.InvalidMessageLengthError("Message sent is longer than 255 bytes")
    }
    binary.Write(&buf, binary.BigEndian, messagecode.SendMessageRequestIdentifier)
    binary.Write(&buf, binary.BigEndian, m.sender)
    binary.Write(&buf, binary.BigEndian, m.receiver)
    binary.Write(&buf, binary.BigEndian, uint8(msgLength))
    binary.Write(&buf, binary.BigEndian, m.message)

    return buf.Bytes(), nil
}

func (m *SendMessageRequest) Unmarshal(dataReceived []byte) *pe.ProtocolError{
    header :=dataReceived[0]
    body := dataReceived[1:]
    
    err := p.CheckForValidHeader(header, body)
    if err != nil{
        return err
    }

    if header != messagecode.SendMessageRequestIdentifier{
        return p.InvalidHeaderError(messagecode.SendMessageRequestIdentifier,
                                  header,
                                  "SendMessageRequest")
    }

    l := 0
    sender := binary.BigEndian.Uint16(body[l:(l+sender_id_in_bytes)])
    l+=sender_id_in_bytes

    receiver :=binary.BigEndian.Uint16(body[l:(l+receiver_id_in_bytes)])
    l+= receiver_id_in_bytes

    msgLen := int(body[l])
    l += 1
    msg := string(body[l:msgLen])

    m.sender = sender
    m.receiver = receiver
    m.message = msg

    return nil
}


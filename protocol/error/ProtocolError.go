package protocol_error

import (
	"encoding/binary"
	"fmt"

	"github.com/rhythm/chatservice/protocol/messagecode"
)


type ProtocolError struct{

    message string
    errorCode uint8
}

func (e *ProtocolError) ErrorCode() int{
    return int(e.errorCode)
}

func (e *ProtocolError) Error() string{
    return fmt.Sprintf("Error Message: %s with error code %d\n",
                       e.message,
                       e.errorCode)
}

func (e *ProtocolError) Marshal()([]byte, *ProtocolError){
    
    buf := make([]byte, 258)
    binary.Encode(buf, binary.BigEndian, messagecode.IsErrorMessage)
    binary.Encode(buf, binary.BigEndian, e.errorCode)
    messageLen := len(e.message)
    if messageLen > 255{
        return nil, InvalidMessageLengthError("ERROR Message too long")
    }
    binary.Encode(buf, binary.BigEndian, uint8(messageLen))

    binary.Encode(buf, binary.BigEndian, e.message)

    return buf[:2+messageLen], nil
}


func (e *ProtocolError) Unmarshal(dataReceived []byte) *ProtocolError{ 
    
    code := dataReceived[0]
    msgLen := int(dataReceived[1])
    message := string(dataReceived[2:msgLen])

    e.errorCode = code
    e.message = message

    return nil

}

func (e *ProtocolError) GetHeader() uint8{
    return messagecode.IsErrorMessage
} 



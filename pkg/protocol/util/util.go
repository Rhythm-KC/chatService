package protocol

import (
	"fmt"

	pe "github.com/rhythm/chatservice/pkg/protocol/error"
)

func InvalidHeaderError(expectedHeader uint8, 
                        receivedHeader uint8, 
                        messageName string) *pe.ProtocolError{

    msg :=fmt.Sprintf("Invalid Header Received for %s.\n Expected: %d, Found %d",
                      messageName,
                      int(expectedHeader),
                      int(receivedHeader))
    return pe.InvalidMessageHeaderError(msg)
}

func CheckForValidHeader(header uint8, body []byte) *pe.ProtocolError{
    if pe.IsErrorCode(header){
        var err pe.ProtocolError
        err.Unmarshal(body)
        return &err
    }
    return nil
}

package protocol_error

import "fmt"

const(

    INVALID_MESSAGE_LENGTH uint8 = 100
    NULL_ERROR uint8 = 110
    RESOURCE_NOT_FOUND uint8 = 120
    INVALID_HEADER uint8 = 130
    TIMED_OUT uint8 = 140
    DISCONNECTED uint8 = 150


    //
    errorMsgLength = 255

)
func New(errorCode int, errMsg string) *ProtocolError{

    if len(errMsg) > 255{
        panic("Error Message can only be 255 bytes in lenght")
    }

    switch uint8(errorCode){
    case NULL_ERROR:
        return NullError(errMsg)
    case INVALID_MESSAGE_LENGTH:
        return InvalidMessageLengthError(errMsg)
    case RESOURCE_NOT_FOUND:
        return NotFoundError(errMsg)
    case INVALID_HEADER:
        return InvalidMessageHeaderError(errMsg)
    case TIMED_OUT:
        return TimedOutError(errMsg)
    case DISCONNECTED:
        return TimedOutError(errMsg)
    default:
        panic(fmt.Sprintf("%d is not a defined error code", errorCode))
    }
}

func IsErrorCode(code uint8) bool {
    if INVALID_MESSAGE_LENGTH   == code ||
       NULL_ERROR               == code || 
       RESOURCE_NOT_FOUND       == code ||
       INVALID_HEADER           == code ||
       TIMED_OUT                == code ||
       DISCONNECTED             == code {
           return true
       }
    return false
}
func validateMessageValid(msg string) {
    if len(msg) > errorMsgLength{
        panic("Message can only be upto 255 bytes")
    }
}

func NullError(errMsg string)*ProtocolError{
    
    validateMessageValid(errMsg)
    return &ProtocolError{errorCode: NULL_ERROR, message: errMsg}
}

func InvalidMessageLengthError(errMsg string) *ProtocolError{

    validateMessageValid(errMsg)
    return &ProtocolError{errorCode: INVALID_MESSAGE_LENGTH, message: errMsg}
}

func NotFoundError(errMsg string) *ProtocolError{

    validateMessageValid(errMsg)
    return &ProtocolError{errorCode: RESOURCE_NOT_FOUND, message: errMsg}
}

func InvalidMessageHeaderError(errMsg string) *ProtocolError{

    validateMessageValid(errMsg)
    return &ProtocolError{errorCode: INVALID_HEADER, message: errMsg}
}

func TimedOutError(errMsg string) *ProtocolError{

    validateMessageValid(errMsg)
    return &ProtocolError{errorCode: TIMED_OUT, message: errMsg}
}
func DisconectedError(errMsg string) *ProtocolError{

    validateMessageValid(errMsg)
    return &ProtocolError{errorCode: DISCONNECTED, message: errMsg}
}

package messagecode

const (
    //  request are going to be 1 to 10 
    SendMessageRequestIdentifier uint8 = 1
    CreateRoomRequestIdentifier uint8 = 2
    JoinRoomRequestIdentifier uint8 = 3


    // Response Status
    MessageResponseIdentifier uint8 = 51
    CreateRoomResponseIdentifier uint8 = 52
    JoinRoomResponseIdentifier uint8 = 53

    IsErrorMessage uint8 = 0
)



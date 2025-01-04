package messagecode

const (
    //  request are going to be 1 to 10 
    SendMessageRequestIdentifier uint8 = 1
    CreateRoomRequestIdentifier uint8 = 2
    JoinRoomRequestIdentifier uint8 = 3
    CreateUserRequestIdentifier uint8 = 4


    // Response Status
    MessageResponseIdentifier uint8 = 51
    CreateRoomResponseIdentifier uint8 = 52
    JoinRoomResponseIdentifier uint8 = 53
    CreateUserResponseIdentifier uint8 = 54

    IsErrorMessage uint8 = 0
)



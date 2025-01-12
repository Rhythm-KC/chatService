package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/rhythm/chatservice/utils"
)

const(
    register_room_option = "C"
    join_room_option = "J"
)



func Run(){

    var serverAdd = "localhost:8080"
    conn, err := net.Dial("tcp", serverAdd)
    if err != nil{
        fmt.Printf("Could not connet to host @ %s", serverAdd)
    }
    consoleReader := bufio.NewReader(os.Stdin)
    userConnection := connection{server: conn, 
                                stdinReader: bufio.NewReader(os.Stdin),
                                readBuffer: make([]byte, 1024),}  

    defer userConnection.close()
    var user *utils.RegisterUser
    user, err = userConnection.registerUser()
    

    for{
        fmt.Print("Enter your message here: ")
        line, err := consoleReader.ReadString('\n')
        if err != nil{
            fmt.Println("Cannot read from stdin")
            break
        }
        conn.Write([]byte(line))
        fmt.Printf("Here is your string %s\n", line)
    }
}

func getClientForRegistration(consoleReader *bufio.Reader) *utils.RegisterUser{
    var username string
    var err error
    for {
        fmt.Print("Enter your Name for Registraion: ")
        username, err = consoleReader.ReadString('\n')
        if err != nil{
            fmt.Println("Cannot read from stdin")
            break
        }
        if len(username) > 255{
            fmt.Println("Name is too long")
            continue
        }
        break
    }

    return &utils.RegisterUser{UserName: username, 
        UserID: uint16(0)}
}


func getRoomNameFromUser(consoleReader *bufio.Reader) (*utils.CreateRoomInst){
    var roomName string
    var err error
    for {
        fmt.Print("Enter your Name for Registraion: ")
        roomName, err = consoleReader.ReadString('\n')
        if err != nil{
            fmt.Println("Error Reading from your text")
            continue
        }
        if len(roomName) > 255{
            fmt.Println("Room Name is Too long Try again")
            continue
        }
        break
    }
    return &utils.CreateRoomInst{RoomName: roomName}
}

func askUserToCreateOrJoinRoom(consoleReader *bufio.Reader, conn connection){
    for{
        fmt.Printf("Enter %s to register a new Room or %s to join existing: ",
                    register_room_option,
                    join_room_option)
        roomOption, err := consoleReader.ReadString('\n')
        if err != nil && (roomOption != register_room_option || 
                          roomOption != join_room_option){
            fmt.Printf("Enter the correct option only")
            continue
        }
        if roomOption == register_room_option{
            conn.createRoom()
            break
        } else if roomOption == join_room_option{
            break
        }else {
            fmt.Printf("%s is not a valid option. Try again\n", roomOption)
        }
    }
}

package server

import (
	"fmt"
	"sync"

	"github.com/rhythm/chatservice/protocol"
	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
	"github.com/rhythm/chatservice/protocol/request"
)

type Room struct{
    ID uint16
    users map[uint16]*User
    MessageQueue chan protocol.Message
    userChannel  chan *User
    done       chan struct{}
    RoomName string
    wg       sync.WaitGroup
    mu       sync.Mutex
}

func NewRoom(roomID uint16, roomName string)*Room{
   return &Room{ID: roomID, RoomName: roomName} 
}

func generateRoomId(existingRoom map[uint16]*Room) uint16{
    var i uint16 = 0
    for {
        if _, exists := existingRoom[i]; !exists{
            return i
        }
        i++
    }
}

func (r *Room) Subscribe(user *User){
    r.mu.Lock()
    r.userChannel <- user
    r.mu.Lock()
}

func (r *Room)removeUser(userId uint16){

    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.users, userId)
}

func (r *Room) asyncListen(function func()){
    r.wg.Add(1)
    go func() {
        defer r.wg.Done()
        function()
    }()
    
}

func (r *Room) listen(usr *User){

    for {
        select{
            case <- r.done:
                //close connection
                usr.conn.Close()
                r.removeUser(usr.Userid)
                return

            default:
                message, err := usr.conn.Listen(10)
                if err.ErrorCode() == int(pe.DISCONNECTED){
                    usr.conn.Close()
                    r.removeUser(usr.Userid)
                    return
                }
                if err.ErrorCode() == int(pe.TIMED_OUT){
                    continue
                }
                r.MessageQueue <- message

        }
    }

}

func (r *Room) pushMessageToAll(message *request.SendMessageRequest){
    for userId, usr := range r.users{
        if userId == message.SenderId(){
            continue
        }
        usr.conn.SendResponse(message)

    }
}

func (r *Room) handleMessage(message protocol.Message){
    header := message.GetHeader()
    if header == messagecode.SendMessageRequestIdentifier{
        reqmessage, _ := message.(*request.SendMessageRequest)
        r.pushMessageToAll(reqmessage)

    }
}

func (r *Room) brodcast(){
    for {
        select{
            case <- r.done:
                return
            case message := <- r.MessageQueue:
                r.handleMessage(message)


        }



    }
    
}

func(r *Room) Run(){
   for {
       select{
            case <-r.done:
                return
            case user := <- r.userChannel:
                r.users[user.Userid] = user
                r.asyncListen(func() {r.listen(user)})

            default:
                fmt.Printf("DOING THINGS")
            
       }
   }

}

func (r *Room) Close(){
    close(r.done)
    r.wg.Wait()
}




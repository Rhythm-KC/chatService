package server

import (
	"fmt"
	"sync"
)

type Room struct{
    ID uint16
    users map[uint16]*User
    MessageQueue chan Message
    WorkerPool chan struct{}
    done       chan struct{}
    RoomName string
    wg       sync.WaitGroup
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

func (r *Room) AddUser(user *User){
    r.users[user.Userid] = user
}

func(r *Room) Run(){
   r.wg.Add(1) 
   defer r.wg.Done()

   for {
       select{
            case <-r.done:
                return
            default:
                fmt.Printf("DOING THINGS")
            
       }
   }

}

func (r *Room) Close(){
    close(r.done)
    r.wg.Wait()
}




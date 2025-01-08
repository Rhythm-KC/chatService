package server

import (
	"fmt"
	"sync"

	"github.com/rhythm/chatservice/protocol"
    pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
	"github.com/rhythm/chatservice/protocol/request"
	"github.com/rhythm/chatservice/protocol/response"
)

type joinRoom struct{
    roomID uint16
    user *User
}

type roomManager struct{

    Rooms      map[uint16]*Room
    Register   chan *User
    done       chan struct{}
    mu         sync.Mutex
    wg         sync.WaitGroup

}

func newRoomManager() *roomManager{

    return &roomManager{
        Rooms:      make(map[uint16]*Room),
        Register:   make(chan *User),
    }
}

func (rm *roomManager) Start(){
    rm.run()
}

func (rm *roomManager) createRoom(roomName string) *Room{
    id := generateRoomId(rm.Rooms)  
    return NewRoom(id, roomName)
}
func (rm *roomManager) Stop(){
    close(rm.done)
}

func (rm *roomManager) waitForUserToJoin(user *User){
    defer rm.wg.Done()
    for {
        req, _ :=  user.conn.Listen(10)
        header := req.GetHeader()
        switch header{
            case messagecode.CreateRoomRequestIdentifier:

                roomReq := req.(*request.CreateRoomRequest)
                newRoom := rm.createRoom(roomReq.RoomName)
                newRoom.Subscribe(user)
                rm.mu.Lock()
                rm.Rooms[newRoom.ID] = newRoom
                rm.mu.Unlock()
                resp := response.NewCreateRoomResponse(user.Userid,
                                                       newRoom.ID,
                                                       roomReq.RoomName)
                user.conn.SendResponse(resp)
                return 
            case messagecode.JoinRoomRequestIdentifier:
                roomReq := req.(*request.JoinRoomRequest)
                rm.mu.Lock()
                room, found := rm.Rooms[roomReq.RoomId]
                var roomResp protocol.Message
                if found{
                    room.Subscribe(user)
                    roomResp = response.NewJoinRoomResponse(room.ID,
                                                             room.RoomName)
                }else{
                    msg := fmt.Sprintf("Cannot find room with id %d",
                                        roomReq.RoomId)
                    roomResp = pe.NotFoundError(msg)
                } 
                rm.mu.Unlock()
                user.conn.SendResponse(roomResp)

                return 
        }
    }
}

func (rm *roomManager)run(){
    for {
        select{
            case user := <-rm.Register:
                // wait for user to to join a room
                rm.wg.Add(1)
                go rm.waitForUserToJoin(user)
            case <- rm.done:
                rm.wg.Wait()
                return
        }
    }
}

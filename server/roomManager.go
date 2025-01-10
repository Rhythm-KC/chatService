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
    rooms      map[uint16]*Room
    Register   chan *User
    done       chan struct{}
    mu         sync.Mutex
    wg         sync.WaitGroup

}

func NewRoomManager() *roomManager{

    return &roomManager{
        rooms:      make(map[uint16]*Room),
        Register:   make(chan *User),
    }
}

func (rm *roomManager) Start(){
    rm.run()
}

func (rm *roomManager)run(){
    for {
        select{
            case user := <-rm.Register:
                rm.async(func() {rm.waitForUserToJoin(user)})
            case <- rm.done:
                return
        }
    }
}

func (rm *roomManager) async(f func()){
    rm.wg.Add(1)
    go func() {
        defer rm.wg.Done()
        f()
    }()
}

func (rm *roomManager) waitForUserToJoin(user *User){
    for {
        select{
            case <- rm.done:
                return
            default:
                resp, err := rm.handleUser(user)
                if err != nil && err.ErrorCode() == int(pe.TIMED_OUT){
                    continue
                }
                if err != nil && err.ErrorCode() == int(pe.DISCONNECTED){
                    user.Delete()
                    return
                }
                if err != nil{
                    user.conn.SendResponse(err)
                    return
                }
                user.conn.SendResponse(resp)
            }
        }
}

func (rm *roomManager) handleUser(user *User)(protocol.Message, *pe.ProtocolError){
    req, err :=  user.conn.Listen(10)
    if err != nil{
        return nil, err
    }
    header := req.GetHeader()
    switch header{

    case messagecode.CreateRoomRequestIdentifier:
        roomReq, _:= req.(*request.CreateRoomRequest) 
        newRoom := rm.createUserRequestedRoom(*roomReq)
        newRoom.Subscribe(user)
        rm.startRoom(newRoom)
        resp := response.NewCreateRoomResponse(user.Userid, newRoom.ID,
                                               roomReq.RoomName)
        return resp, nil 

    case messagecode.JoinRoomRequestIdentifier:
         roomReq := req.(*request.JoinRoomRequest)
         rm.mu.Lock()
         room, found := rm.rooms[roomReq.RoomId]
         var roomResp protocol.Message
         if found{
            room.Subscribe(user)
            roomResp = response.NewJoinRoomResponse(room.ID, room.RoomName)
         }else{
            msg := fmt.Sprintf("Cannot find room with id %d", roomReq.RoomId)
            roomResp = pe.NotFoundError(msg)
        } 
        rm.mu.Unlock()
        return roomResp, nil
    }
    return nil, pe.InvalidMessageHeaderError("Not a valid message for joining room")
}

func (rm *roomManager) createUserRequestedRoom(req request.CreateRoomRequest)*Room{

    newRoom := rm.createRoom(req.RoomName)
    rm.mu.Lock()
    rm.rooms[newRoom.ID] = newRoom
    rm.mu.Unlock()
    return newRoom
}

func (rm *roomManager) startRoom(room *Room){
    rm.mu.Lock()
    rm.async(func(){room.Start()})
    rm.mu.Unlock()
}

func (rm *roomManager) createRoom(roomName string) *Room{
    id := rm.generateRoomId()  
    return NewRoom(id, roomName)
}

func (rm *roomManager) generateRoomId() uint16{
    var i uint16 = 0
    for {
        if _, exists := rm.rooms[i]; !exists{
            return i
        }
        i++
    }
}

func (rm *roomManager) Stop(){
    close(rm.done)
    rm.wg.Wait()
}


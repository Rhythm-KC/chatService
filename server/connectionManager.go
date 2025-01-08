package server

import (
	"fmt"
	"sync"
	"sync/atomic"

	p "github.com/rhythm/chatservice/protocol"
	"github.com/rhythm/chatservice/protocol/messagecode"
	"github.com/rhythm/chatservice/protocol/request"
	"github.com/rhythm/chatservice/protocol/response"
)
const (
    closingFalse = 1
    closingTrue = 0
)

type connectionManager struct{

    newConnection chan *p.ServerConnection
    done chan struct{}
    roomManager *roomManager
    wg sync.WaitGroup
    mu sync.Mutex
    closing int32

}

var singleInstance *connectionManager

func NewConnectionManger(roomManager *roomManager) *connectionManager{
    return &connectionManager{newConnection: make(chan *p.ServerConnection,10),
                              roomManager: roomManager,
                              mu: sync.Mutex{},
                              closing: 0,} 
}

func (cm *connectionManager) Start() error{
    if singleInstance != nil{
        return fmt.Errorf("connection manager already exists")
    }
    singleInstance = cm
    go cm.run()
    return nil
}

func (cm *connectionManager) Close(){
    close(cm.done)
}

func (cm *connectionManager) setClosing(val int32){
    atomic.StoreInt32(&cm.closing, val)
}

func  (cm *connectionManager) isClosing()bool{
    return atomic.LoadInt32(&cm.closing) == 1
}

func (cm *connectionManager) waitForUserRegistration(conn *p.ServerConnection){
    defer cm.wg.Done()
    for {
        msg, err := conn.Listen(10) 
        if err != nil{
            // deal with inactive user
           continue 
        }
        if cm.isClosing(){
            // need to send user a response saying that the server shut down
            return
        }
        
        if msg.GetHeader() == messagecode.CreateUserRequestIdentifier{
            req, _:= msg.(*request.CreateUserRequest)
            user := NewUser(req.Name(), conn)
            Resp := response.NewCreateUserResponse(user.userid)
            conn.SendResponse(Resp)

            cm.mu.Lock()
            cm.roomManager.Register <- user
            cm.mu.Unlock()

            return 
        }
    }
}


func (cm *connectionManager) run(){

    defer cm.wg.Wait()
    for{
        select{
        case newConn := <- cm.newConnection:
            cm.wg.Add(1)
            go cm.waitForUserRegistration(newConn)
        case <- cm.done:
            cm.setClosing(closingFalse)
            return
        }
    }
}


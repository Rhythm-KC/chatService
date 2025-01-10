package server

import (
	"fmt"
	"sync"

	p "github.com/rhythm/chatservice/protocol"
	pe "github.com/rhythm/chatservice/protocol/error"
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
    cm.run()
    return nil
}
func (cm *connectionManager) async(f func()){
    cm.wg.Add(1)
    go func() {
        defer cm.wg.Done()
        f()
    }()
}

func (cm *connectionManager) run(){

    for{
        select{
        case newConn := <- cm.newConnection:
            cm.async(func(){cm.waitForUserRegistration(newConn)})
        case <- cm.done:
            return
        }
    }
}

func (cm *connectionManager) waitForUserRegistration(conn *p.ServerConnection){
    for {
        select{
        case <- cm.done:
            return
        default:
            msg, err := conn.Listen(10) 
            if err != nil && err.ErrorCode() == int(pe.TIMED_OUT){
               continue 
            }
            if err != nil && err.ErrorCode() == int(pe.DISCONNECTED){
                conn.Close()
                return
            }
            if err != nil{
                conn.SendResponse(err)
            }
            
            if msg.GetHeader() == messagecode.CreateUserRequestIdentifier{
                req, _:= msg.(*request.CreateUserRequest)
                user := NewUser(req.Name(), conn)
                Resp := response.NewCreateUserResponse(user.Userid)
                conn.SendResponse(Resp)

                cm.mu.Lock()
                cm.roomManager.Register <- user
                cm.mu.Unlock()
                return 
            }else{
                conn.SendResponse(pe.InvalidMessageHeaderError("User has not been registered"))

            }

        }
    }
}

func (cm *connectionManager) Submit(newConn *p.ServerConnection){
    cm.newConnection <- newConn
}

func (cm *connectionManager) Close(){
    close(cm.done)
    cm.wg.Wait()
}


package server

import (
	"sync"

	"github.com/rhythm/chatservice/protocol"
	pe "github.com/rhythm/chatservice/protocol/error"
	"github.com/rhythm/chatservice/protocol/messagecode"
	"github.com/rhythm/chatservice/protocol/request"
	"github.com/rhythm/chatservice/thread"
)

type Room struct{
    ID               uint16
    users            sync.Map
    broadcastPool    thread.ThreadPool
    MessageQueue     chan protocol.Message
    userChannel      chan *User
    done             chan struct{}
    RoomName         string
    wg               sync.WaitGroup
}

func NewRoom(roomID uint16, roomName string)*Room{
   return &Room{ID: roomID, 
                RoomName: roomName, 
                broadcastPool: *thread.NewThreadPool(10)} 
}


// Starts lisiting to the users and processing user mesages
func (r *Room) Start(){
    r.async(func() {r.brodcast()})
    r.run()
}

func (r *Room) Subscribe(user *User){
    r.userChannel <- user
}

func (r *Room)removeUser(userId uint16){
    user, _:= r.users.Load(userId)
    usr, _ := user.(*User) 
    usr.Delete()
    r.users.Delete(userId)
}

func (r *Room) async(function func()){
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
                r.removeUser(usr.Userid)
                return

            default:
                message, err := usr.conn.Listen(10)
                if err.ErrorCode() == int(pe.DISCONNECTED){
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
    r.users.Range(func(key, value interface{}) bool {
        userID := key.(uint16)  
        user := value.(*User)  

        if userID != message.SenderId() {
            user.conn.SendResponse(message)
        }
        return true 
    })


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
                r.broadcastPool.Submit(func (){r.handleMessage(message)})
        }
    }
    
}

func(r *Room) run(){
   for {
       select{
            case <-r.done:
                return
            case user := <- r.userChannel:
                r.users.Store(user.Userid, user)
                r.async(func() {r.listen(user)})
            
       }
   }

}

func (r *Room) Close(){
    close(r.done)
    r.broadcastPool.Shutdown()
    r.wg.Wait()
}

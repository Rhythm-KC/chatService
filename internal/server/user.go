package server

import (

	p "github.com/rhythm/chatservice/protocol"
)

type User struct{
    Userid uint16
    userName string
    conn *p.ServerConnection
}
var usedId map[uint16]struct{} = make(map[uint16]struct{})

func NewUser(name string, conn *p.ServerConnection) *User{
    newId := generateValidUserId()
    usedId[newId] =struct{}{} 
    return &User{Userid: newId, userName: name, conn: conn}
}

func (u *User) Delete(){
    delete(usedId, u.Userid)
    (*u.conn).Close()
}

func generateValidUserId() uint16{
    var i uint16 = 1
    for{
        if _, exists := usedId[i]; !exists{
            return i
        }
        i+=1
    }
}



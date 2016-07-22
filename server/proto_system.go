// Copyright © 2016.6 Claus Chen
//
//

package server

import (
	"fmt"

	"github.com/clauschen123/Therion/protobuf/system"
	//"github.com/clauschen123/Therion/server"
)

//TODO
var ProtoSystem IProtocol = MakeProtocol(GetProtID())

func init() {
	fmt.Println("call proto_system.go init")
	GetServer().RegisterProtocol(ProtoSystem)
}

func GetProtID() ProtoID {
	return ProtoID(system.Proto_id)
}

//!+ New a protocol
func MakeProtocol(protid ProtoID) IProtocol {
	proto := &Protocol{

		Pid: protid,

		Handler: map[MsgID]func(*Message){

			1: func(msg *Message) {
				fmt.Println("Get connected: ", msg.protoid, msg.msgid, string(msg.data))
			},

			2: func(msg *Message) {
				fmt.Println("Get disconnect: ", msg.protoid, msg.msgid, string(msg.data))
			},

			3: func(msg *Message) {
				fmt.Println("Get echo: ", msg.protoid, msg.msgid, string(msg.data))
			},
		},
	}
	return proto
}

// Copyright © 2016.6 Claus Chen
//
//

package server

import (
	"fmt"

	"github.com/clauschen123/Therion/protobuf/system"
	"github.com/golang/protobuf/proto"
)

//TODO
//type ProtocolSystem struct {
//	Protocol
//}

//TODO
var ProtoSystem IProtocol = makeProtocol(GetProtID())

func init() {
	fmt.Println("call proto_system.go init")
	//TODO
	//GetServer().RegisterProtocol(ProtoSystem)
}

func GetProtID() ProtoID {
	return ProtoID(system.Proto_id)
}

func (protocol *Protocol) makeMsgAuthClient() *Message {

	msg := &system.MsgAuthClient{MsgId: proto.Uint32(1)}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("Masharl MsgAuthClient error:", err)
		return nil
	}
	return MakeMessage(protocol.GetProtID(), 1, data)
}

func (protocol *Protocol) makeMsgAuthServer(sid uint32) *Message {

	msg := &system.MsgAuthServer{MsgId: proto.Uint32(2), SvrId: proto.Uint32(sid)}

	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("Masharl MsgAuthServerv error:", err)
		return nil
	}
	return MakeMessage(protocol.GetProtID(), 2, data)
}

func (protocol *Protocol) makeMsgEcho(text string) *Message {

	protmsg := &system.MsgEcho{MsgId: proto.Uint32(3), Text: proto.String(text)}

	data, err := proto.Marshal(protmsg)
	if err != nil {
		fmt.Println("Masharl MsgEcho error:", err)
		return nil
	}
	return MakeMessage(protocol.GetProtID(), 3, data)
}

func makeProtocol(protid ProtoID) IProtocol {
	prot := &Protocol{
		Pid:        protid,
		MsgHandler: map[MsgID]func(*Message){},
	}

	prot.RegisterMessage(2, func(msg *Message) {
		fmt.Println("Get MsgAuthServer: ", msg.protoid, msg.msgid)
		authServer := &system.MsgAuthServer{}
		if err := proto.Unmarshal(msg.data, authServer); err == nil {
			fmt.Println("   SvrID=", authServer.GetSvrId())
			prot.Send2Server_MsgEcho(SID(0), "Hi GATE!")
		}
	})

	prot.RegisterMessage(3, func(msg *Message) {
		fmt.Println("Get MsgEcho: ", msg.protoid, msg.msgid)
		authEcho := &system.MsgEcho{}
		if err := proto.Unmarshal(msg.data, authEcho); err == nil {
			fmt.Println("   ", authEcho.GetText())
			prot.Send2Server_MsgEcho(SID(0), authEcho.GetText())
		}
	})

	return prot
}

//Send Message
func (proto *Protocol) Send2Server_MsgAuthServer(sid SID) {
	msg := proto.makeMsgAuthServer(uint32(sid))
	proto.GetConn().Send2Server(sid, msg)
}

//func (proto *Protocol) Send2Servers_MsgAuthServer(sids SIDS) {
//	msg := proto.makeMsgAuthServer()
//	proto.GetConn().Send2Servers(sids, msg)
//}

func (proto *Protocol) Send2Server_MsgEcho(sid SID, echo string) {
	msg := proto.makeMsgEcho(echo)
	proto.GetConn().Send2Server(sid, msg)
}

//func (proto *Protocol) Send2Servers_MsgEcho(sids SIDS, echo string) {
//	msg := proto.makeMsgEcho(echo)
//	proto.GetConn().Send2Servers(sids, msg)
//}

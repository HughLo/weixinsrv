package wxsrv

import (
	"testing"
	"database/sql"
	"log"
)

func TestHandleMsg(t *testing.T) {
	ConnString := "root:hugh1984lou@/weixin_hugh"

	bm := BaseMsg {
		ToUserName: "toUser", 
		FromUserName: "fromUser",
		CreateTime: "1234556",
		MsgType: "text",
	}

 	m := &UserMsg{BaseMsg: bm, 
 		Content:"TYDL -t 30 -e 250", 
 		MsgId:"1234567890123456",
 	}

 	mh := &UserMsgHandler{m, nil, nil}

 	err = mh.Handle(nil)
 	if err != nil {
 		t.Log(err)
 	}
}
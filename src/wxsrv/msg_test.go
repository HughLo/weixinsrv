package wxsrv

import (
	"log"
	"net/http"
	"testing"
)

type ResponseWriterMock struct{}

func (rwm *ResponseWriterMock) Header() http.Header {
	log.Println("calling the Header() function")
	return nil
}

func (rwm ResponseWriterMock) Write(content []byte) (int, error) {
	log.Println(string(content))
	return len(content), nil
}

func (rwm ResponseWriterMock) WriteHeader(statusCode int) {
	log.Printf("calling WriteHeader(%d) function", statusCode)
}

func TestHandleMsg(t *testing.T) {
	ConnString = "root:hugh1984lou@/weixin_hugh"

	bm := BaseMsg{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   "1234556",
		MsgType:      "text",
	}

	m := &UserMsg{BaseMsg: bm,
		Content: "TYDL -t 30 -e 250",
		MsgId:   "1234567890123456",
	}

	mh := &UserMsgHandler{m, &ResponseWriterMock{}}
	err := mh.Handle()
	if err != nil {
		t.Log(err)
	}
}

func TestHandleHelpMsg(t *testing.T) {
	bm := BaseMsg{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   "1234556",
		MsgType:      "text",
	}

	m := &UserMsg{BaseMsg: bm,
		Content: "help",
		MsgId:   "1234567890123456",
	}

	mh := &UserMsgHandler{m, &ResponseWriterMock{}}
	err := mh.Handle()
	if err != nil {
		t.Log(err)
	}
}

func TestReportAll(t *testing.T) {
	ConnString = "root:hugh1984lou@/weixin_hugh"
	bm := BaseMsg{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   "1234556",
		MsgType:      "text",
	}

	m := &UserMsg{BaseMsg: bm,
		Content: "Report -all",
		MsgId:   "1234567890123456",
	}

	mh := &UserMsgHandler{m, &ResponseWriterMock{}}
	err := mh.Handle()
	if err != nil {
		t.Log(err)
	}
}

func TestReportThisWeek(t *testing.T) {
	ConnString = "root:hugh1984lou@/weixin_hugh"
	bm := BaseMsg{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   "1234556",
		MsgType:      "text",
	}

	m := &UserMsg{BaseMsg: bm,
		Content: "Report -since thisweek",
		MsgId:   "1234567890123456",
	}

	mh := &UserMsgHandler{m, &ResponseWriterMock{}}
	err := mh.Handle()
	if err != nil {
		t.Log(err)
	}
}

func TestRawQuery(t *testing.T) {
	db := CreateDBMgr("root:hugh1984lou@/weixin_hugh")
	if db == nil {
		t.FailNow()
	}

	defer db.Close()

	errHandler := func(err error, t *testing.T) {
		if err != nil {
			log.Println(err)
			t.FailNow()
		}
	}

	err := db.UseDB("weixin_hugh")
	errHandler(err, t)

	r, err := db.RawQuery(`select sum(er.execisetime) as all_execise_time, sum(er.execiseenergy) as all_execise_energy from execise_records as er;`)
	errHandler(err, t)

	defer r.Rows.Close()

	if !r.Rows.Next() {
		log.Println("empty result set")
		return
	}

	var rd ReportData
	err = r.Rows.Scan(&(rd.TotalTime), &(rd.TotalEnergy))
	errHandler(err, t)

	log.Println(rd)
}

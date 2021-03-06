package wxsrv

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func WithTestDB(t *testing.T, fn func(*ExeciseDB)) {
	ConnString := "root:hugh1984lou@/weixin_hugh"
	db := CreateExeciseDB()
	if db == nil {
		t.Logf("failed to create execise db with conn string: %s", ConnString)
		t.FailNow()
	}
	defer db.Close()
	fn(db)
}

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

func WithHandleMsg(t *testing.T, msgc string) {
	ConnString = "root:hugh1984lou@/weixin_hugh"

	bm := BaseMsg{
		ToUserName:   "toUser",
		FromUserName: "fromUser",
		CreateTime:   "1234556",
		MsgType:      "text",
	}

	m := &UserMsg{BaseMsg: bm,
		Content: msgc,
		MsgId:   "1234567890123456",
	}

	mh := &UserMsgHandler{m, &ResponseWriterMock{}}
	err := mh.Handle()
	if err != nil {
		t.Log(err)
	}
}

func TestHandleMsg(t *testing.T) {
	WithHandleMsg(t, "TYDL -t 30 -e 250")
}

func TestHandleHelpMsg(t *testing.T) {
	WithHandleMsg(t, "help")
}

func TestHandleReportAllMsg(t *testing.T) {
	WithHandleMsg(t, "Report -all")
}

func TestHandleReportThisWeekMsg(t *testing.T) {
	WithHandleMsg(t, "Report -since thisweek")
}

func TestReportSinceWeek(t *testing.T) {
	WithTestDB(t, func(db *ExeciseDB) {
		year, wk := time.Now().ISOWeek()
		t.Logf("current week: %d\n", wk)
		//test since last week
		rd, err := db.ReportSinceWeek(year, wk-1)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		t.Logf("report since last week: t(%d), e(%d)", rd.TotalTime, rd.TotalEnergy)
	})
}

func TestReportSinceLastWeek(t *testing.T) {
	WithTestDB(t, func(db *ExeciseDB) {
		rd, err := db.ReportSinceLastWeek()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		t.Logf("report since last week: t(%d), e(%d)", rd.TotalTime, rd.TotalEnergy)
	})
}

func TestReportSinceMonth(t *testing.T) {
	WithTestDB(t, func(db *ExeciseDB) {
		now := time.Now()
		rd, err := db.ReportSinceMonth(now.Year(), int(now.Month()))
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		t.Logf("report since this month: t(%d), e(%d)", rd.TotalTime, rd.TotalEnergy)
	})
}

func TestReportSinceThisYear(t *testing.T) {
	WithTestDB(t, func(db *ExeciseDB) {
		rd, err := db.ReportSinceThisYear()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		t.Logf("report since this year: t(%d), e(%d)", rd.TotalTime, rd.TotalEnergy)
	})
}

func TestReportSinceLastYear(t *testing.T) {
	WithTestDB(t, func(db *ExeciseDB) {
		rd, err := db.ReportSinceLastYear()
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		t.Logf("report since last year: t(%d), e(%d)", rd.TotalTime, rd.TotalEnergy)
	})
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

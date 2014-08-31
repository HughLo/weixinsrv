package wxsrv

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	NOT_HANDLED_MSG_TYPE = errors.New("the msg type is not handled")
)

func SendResponseMsg(m *UserMsg, d []byte, w http.ResponseWriter) error {
	bm := BaseMsg{
		ToUserName:   fmt.Sprintf("%s", m.FromUserName),
		FromUserName: fmt.Sprintf("%s", m.ToUserName),
		CreateTime:   fmt.Sprintf("%d", time.Now().Unix()),
		MsgType:      "text",
	}

	rm := ResponseMsg{
		BaseMsg: bm,
		Content: string(d),
	}

	response, err := xml.Marshal(&rm)
	if err != nil {
		return err
	}

	_, err = w.Write(response)

	return err
}

type MsgHandler interface {
	Handle(http.ResponseWriter) error
}

type BaseMsg struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
}

type UserMsg struct {
	BaseMsg
	Content string `xml:"Content"`
	MsgId   string `xml:"MsgId"`
}

type ResponseMsg struct {
	XMLName xml.Name `xml:"xml"`
	BaseMsg
	Content string `xml:"Content"`
}

type UserMsgHandler struct {
	Msg *UserMsg
	W   http.ResponseWriter
}

func (m *UserMsgHandler) Handle() error {
	log.Println("inside user_msg::Handle")
	switch m.Msg.BaseMsg.MsgType {
	case "text": // handle text message
		log.Println("handle text msg")
		return m.HandleTextMsg()
	default:
		return NOT_HANDLED_MSG_TYPE
	}
}

func (m *UserMsgHandler) HandleTextMsg() error {
	args := strings.Split(m.Msg.Content, " ")
	switch strings.ToLower(args[0]) {
	case "help":
		h := &HelpMsgHandler{m.W, m.Msg}
		return h.Handle()

	case "tydl":
		fs := flag.NewFlagSet("TYDL", flag.ContinueOnError)
		t := fs.Int("t", 0, "execise time. the merit is minute.")
		e := fs.Int("e", 0, "execise enegry. the merit is kcal.")
		err := fs.Parse(args[1:])
		if err != nil {
			log.Println(err)
			return err
		}

		h := &TYDLMsgHandler{m.Msg, *t, *e}
		return h.Handle()

	case "report":
		fs := flag.NewFlagSet("report", flag.ContinueOnError)
		since := fs.String("since", "thisweek", "specify the time scope for report")
		all := fs.Bool("all", false, "report from the very begining")
		err := fs.Parse(args[1:])
		if err != nil {
			log.Println(err)
			return err
		}

		h := &ReportMsgHandler{m.W, *since, *all, m.Msg}
		return h.Handle()
	default:
		return errors.New("not recognized commands")
	}

	return nil
}

type HelpMsgHandler struct {
	w http.ResponseWriter
	m *UserMsg
}

func (h *HelpMsgHandler) Handle() error {
	//@TODO: shall develop a file cache mechanism so that the file can be read
	//from the disk at the first using. The subsequent call will get the data
	//from the memory.

	gp := os.Getenv("GOPATH")
	f, err := os.Open(filepath.Join(gp, "bin/help.txt"))
	if err != nil {
		return err
	}

	defer f.Close()

	buf := make([]byte, 2048)
	var readCnt int = 0
	for {
		readCnt, err = f.Read(buf)
		if err != io.EOF && readCnt > 0 {
			sndErr := SendResponseMsg(h.m, buf[:readCnt], h.w)
			if sndErr != nil {
				return sndErr
			}
		}

		if err == io.EOF && readCnt <= 0 {
			break
		}
	}

	return nil
}

type ReportMsgHandler struct {
	w     http.ResponseWriter
	since string
	all   bool
	m     *UserMsg
}

func (rmh *ReportMsgHandler) Handle() error {
	var rd *ReportData
	var err error
	db := CreateExeciseDB()
	if db == nil {
		return errors.New("create execise db failure")
	}
	defer db.Close()

	if rmh.all {
		rd, err = db.ReportAll()
	} else {
		switch rmh.since {
		case "thisweek":
			rd, err = db.ReportSinceThisWeek()
		case "lastweek":
			rd, err = db.ReportSinceLastWeek()
		case "thisyear":
			rd, err = db.ReportSinceThisYear()
		case "lastyear":
			rd, err = db.ReportSinceLastYear()
		default:
			return errors.New("unrecognized since string")
		}
	}

	if err != nil {
		return err
	}

	var rs string
	if rd == nil {
		rs = "no report generated"
	} else {
		rs = fmt.Sprintf("total time: %d, total energy: %d", rd.TotalTime, rd.TotalEnergy)
	}

	return SendResponseMsg(rmh.m, []byte(rs), rmh.w)
}

type TYDLMsgHandler struct {
	m *UserMsg
	t int
	e int
}

func (h *TYDLMsgHandler) Handle() error {
	db := CreateExeciseDB()
	if db == nil {
		return errors.New("failed to create execise db")
	}

	r := &ExeciseRecord{
		UserName:      h.m.BaseMsg.FromUserName,
		ExeciseTime:   h.t,
		ExeciseEnergy: h.e,
	}

	_, err := db.Insert(r)
	return err
}

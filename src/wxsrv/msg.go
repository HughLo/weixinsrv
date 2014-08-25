package wxsrv

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	NOT_HANDLED_MSG_TYPE = errors.New("the msg type is not handled")
)

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
	switch args[0] {
	case "TYDL":
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

	case "Report":
		fs := flag.NewFlagSet("Report", flag.ContinueOnError)
		since := fs.String("since", "thisweek", "specify the time scope for report")
		err := fs.Parse(args[1:])
		if err != nil {
			log.Println(err)
			return err
		}

		h := &ReportMsgHandler{m.W, *since, m.Msg}
		return h.Handle()
	default:
	}

	return nil
}

type ReportMsgHandler struct {
	w     http.ResponseWriter
	since string
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

	switch rmh.since {
	case "thisweek":
		rd, err = db.ReportThisWeek()
	case "all":
		rd, err = db.ReportAll()
	default:
		return errors.New("unrecognized since string")
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

	bm := BaseMsg{
		ToUserName:   fmt.Sprintf("%s", rmh.m.FromUserName),
		FromUserName: fmt.Sprintf("%s", rmh.m.ToUserName),
		CreateTime:   fmt.Sprintf("%d", time.Now().Unix()),
		MsgType:      "text",
	}

	rm := ResponseMsg{
		BaseMsg: bm,
		Content: rs,
	}

	response, err := xml.Marshal(&rm)
	if err != nil {
		return err
	}

	_, err = rmh.w.Write(response)

	return err
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

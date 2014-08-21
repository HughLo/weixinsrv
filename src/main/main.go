package main
import(
	"net/http"
	"log"
	"sort"
	"strings"
	"crypto/sha1"
	"fmt"
	"io"
	"encoding/xml"
	"wxsrv"
)

//@todo: msg duplication elimination

type normal_msg struct {
	XMLName xml.Name `xml:"Xml"`
	ToUserName string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime string `xml:"CreateTime"`
	MsgType string `xml:"MsgType"`
	Content string `xml:"Content"`
	MsgId string `xml:"MsgId"`
}

func hugh_test_handler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("test success"))
	if err != nil {
		log.Println(err)
	}
}

func checkSignature(req *http.Request) bool {
	const token = "hugh_weixin_account"
	sign := req.FormValue("signature")
	nonc := req.FormValue("nonce")
	tsta := req.FormValue("timestamp")

	sorted_string := []string{token, nonc, tsta}
	sort.Strings(sorted_string)
	joined_string := strings.Join(sorted_string, "")

	//log.Println("joined string: ", strings.Join(sorted_string, ","))
	
	h := sha1.New()
	
	_, err := io.WriteString(h, joined_string)
	if err != nil {
		log.Fatal(err)
	}

	r := fmt.Sprintf("%x", h.Sum(nil))
	return string(r) == sign
}

func firstBind(w http.ResponseWriter, req *http.Request) {
	estr := req.FormValue("echostr")

	log.Println("recv a message")

	if checkSignature(req) {
		log.Println("verification success")
		w.Write([]byte(estr))
	}
}

func extract_body_content(req *http.Request) string {
	buf := make([]byte, 1024)
	ret_str := ""

	for {
		rd_cnt, err := req.Body.Read(buf)
		if rd_cnt > 0 {
			ret_str = strings.Join([]string{ret_str, string(buf[:rd_cnt])}, "")
		}

		if err == io.EOF {
			break;
		}
	}

	return ret_str
}

func extract_user_message(req *http.Request) *wxsrv.UserMsg {
	bd_content := extract_body_content(req)

	log.Println(bd_content)
	
	msg := wxsrv.UserMsg{}
	err := xml.Unmarshal([]byte(bd_content), &msg)

	if err != nil {
		log.Println(err)
		return nil
	}

	return &msg
}

func hugh_weixin_handler(w http.ResponseWriter, req *http.Request) {
	const verified = true
	if !verified {
		firstBind(w, req)
	} else {
		if req.Method == "POST" {
			msg := extract_user_message(req)
			if msg == nil {
				log.Println("cannot extract msg")
				return
			}

			log.Println(*msg)

			wxsrv.ConnString = "root:hughroot@/weixin_hugh"
			msgHandler := &wxsrv.UserMsgHandler{msg, w}

			err := msgHandler.Handle()
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func main() {
	http.HandleFunc("/hugh_test", hugh_test_handler)
	http.HandleFunc("/hugh_weixin", hugh_weixin_handler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}
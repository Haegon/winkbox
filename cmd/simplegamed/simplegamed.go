package main

import (
	"fmt"
	"gos/simplegamed"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

const (
	// http 리슨 주소
	addr = "localhost:9999"
)

var (
	uMgr = simplegamed.NewUserManager() // 유저를 관리하는 매니저를 초기화 한다.
	log  = logrus.New()                 // 로그러스 초기화
)

func main() {
	// 시작했다는 로그하나 찍어주자.
	log.WithField("Listen Address", addr).Info("Simple Game Server Start")

	// 라우터를 하나 생성하고
	router := httprouter.New()
	// 핸들러를 등록한다.
	router.POST("/signup/:id", signup)
	router.POST("/login/:id", login)
	router.GET("/action/:uid/:action", action)

	// 서비스 한다.
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.WithField("err", err).Error("ListenAndServe Error")
	}
}

// 에러났을때 유저에게 실패 메세지 전송 및 서버 로깅.
func fail(w http.ResponseWriter, msg string, err error) {
	io.WriteString(w, fmt.Sprintf("%v: %v\n", msg, err))
	log.WithField("err", err).Error(msg)
}

func signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 패킷에서 아이디와 비밀번호를 가져온다.
	id := ps.ByName("id")
	pw, err := getPW(r.Body)
	if err != nil {
		fail(w, "Read Body Error", err)
		return
	}

	// 회원 가입 신청을 한다.
	err = uMgr.Signup(id, pw)
	if err != nil {
		fail(w, "SignUp Error", err)
		return
	}

	// 유저에게 환영 인사를 한다.
	fmt.Fprintf(w, "Hello, %s!\n", id)
	log.WithFields(logrus.Fields{"id": id, "pw": pw}).
		Info("Joined a New User")
}

func login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 패킷에서 아이디와 비밀번호를 가져온다.
	id := ps.ByName("id")
	pw, err := getPW(r.Body)
	if err != nil {
		fail(w, "Read Body Error", err)
		return
	}

	// 유저 매니저를 통해 로그인한다.
	info, err := uMgr.Login(id, pw)
	if err != nil {
		fail(w, "Login Error", err)
		return
	}

	// 쿠키를 생성한다.
	http.SetCookie(w, &http.Cookie{
		Name:  "sessionid",
		Value: info.Secure.SessionID,
		Path:  "/",
	})

	// 유저에게 uid를 알려준다.
	fmt.Fprintf(w, "Login Success\nUID : %v\n", info.UID)
	log.WithFields(logrus.Fields{"id": id, "uid": info.UID, "sessionid": info.Secure.SessionID}).Info("User Logined")
}

func action(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// 쿠키에서 세션 아이디를 가져온다.
	sid, err := getSessionID(r)
	if err != nil {
		fail(w, "Cookie Error", err)
		return
	}

	// uid를 int로 변환
	u, err := strconv.ParseInt(ps.ByName("uid"), 0, 64)
	if err != nil {
		fail(w, "Int Parse Error", err)
		return
	}

	// 유저 매니저를 통해 Action처리를 한다.
	action := ps.ByName("action")
	cnt, err := uMgr.Action(action, u, sid)
	if err != nil {
		fail(w, "Ping Error", err)
		return
	}

	// 정상 처리 됐다.
	fmt.Fprintf(w, "[Ping:%v] Done Action %v\n", cnt, action)
	log.WithFields(logrus.Fields{"uid": u, "count": cnt}).
		Info("Done Action: " + action)
}

// getPW는 바디에서 비밀번호를 가져온다.
func getPW(r io.Reader) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// getSessionID는 쿠키에서 세션 아이디를 가져온다.
func getSessionID(r *http.Request) (string, error) {
	sid, err := r.Cookie("sessionid")
	if err != nil {
		return "", err
	}
	return sid.Value, nil
}

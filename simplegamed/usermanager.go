package simplegamed

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

var (
	ErrDuplicateID    = errors.New("Duplicate ID")    // 아이디 중복 에러
	ErrUnknownID      = errors.New("Unknown ID")      // 존재 하지 않는 아이디
	ErrUnknownUID     = errors.New("Unknown UID")     // 존재 하지 않는 UID
	ErrWrongPassword  = errors.New("Wrong Password")  // 비밀번호 틀림
	ErrWrongSessionID = errors.New("Wrong SessionID") // 세션아이디 이상함
	ErrUnknownAction  = errors.New("Unknown Action")  // 알 수 없는 액션
)

const (
	// 액션 타입 정의
	ACTION_RESET Action = "reset"
	ACTION_PING  Action = "ping"
)

// Action은 action 패킷을 적절히 처리하기 위한 타입
type Action string

// UserManager는 유저를 관리 하는 매니저이다.
type UserManager struct {
	role *SafeMap
	uid  *SafeMap
	seq  int64
}

// UserInfo는 비밀번호, 고유 ID, 세션을 가진다.
type UserInfo struct {
	password string
	UID      int64
	Secure   Session
}

// NewUserManager는 유저 매니저를 생성해서 리턴한다.
func NewUserManager() *UserManager {
	return &UserManager{
		role: NewSafeMap(),
		uid:  NewSafeMap(),
	}
}

// Signup은 회원 가입 할 때 사용하는 함수.
func (u *UserManager) Signup(id, pw string) error {
	// 중복아이디 인지 체크 한다.
	// 에러가 없으면 이 아이디로 계정이 이미 있다.
	_, err := u.getUserInfo(id)
	if err == nil {
		return ErrDuplicateID
	}

	// 유저를 추가 한다.
	u.seq++
	u.role.set(id, UserInfo{password: pw, UID: u.seq})
	u.uid.set(u.seq, id)

	return nil
}

// Login은 로그인 할 때 사용하는 함수.
func (u *UserManager) Login(id, pw string) (UserInfo, error) {
	// id로 회원 명부를 찾아서 없으면 에러
	info, err := u.getUserInfo(id)
	if err != nil {
		return UserInfo{}, err
	}

	// 비밀번호를 검증한다.
	if info.password != pw {
		return UserInfo{}, ErrWrongPassword
	}

	// 세션키를 랜덤하게 가져온다.(16자)
	s, err := newSecureID(16)
	if err != nil {
		return UserInfo{}, err
	}

	// 유저의 세션을 생성한다.
	info.Secure = Session{SessionID: s, PingCount: 0}
	u.role.set(id, info)

	return info, nil
}

// Action은 유저가 action패킷을 보냈을때 사용하는 함수.
func (u *UserManager) Action(action string, uid int64, sid string) (int64, error) {
	// uid로 id를 가져옴 없으면 에러
	id, err := u.getID(uid)
	if err != nil {
		return 0, err
	}

	// id로 회원 명부를 찾아서 없으면 에러
	info, err := u.getUserInfo(id)
	if err != nil {
		return 0, err
	}

	// 세션 체크
	if !info.Secure.checkSessionID(sid) {
		return 0, ErrWrongSessionID
	}

	// 액션 별로 처리한다.
	switch Action(action) {
	case ACTION_RESET:
		info.Secure.reset()
		break
	case ACTION_PING:
		info.Secure.increase(1)
		break
	default:
		return 0, ErrUnknownAction
	}

	// 유저 정보를 갱신 한다.
	u.role.set(id, info)

	return info.Secure.PingCount, nil
}

// getID는 uid의 ID를 리턴해주는 함수.
func (u *UserManager) getID(uid int64) (string, error) {
	val, ok := u.uid.get(uid)
	if !ok {
		return "", ErrUnknownUID
	}
	return val.(string), nil
}

// getUserInfo는 id에 해당하는 유저 정보를 리턴해주는 함수.
func (u *UserManager) getUserInfo(id string) (UserInfo, error) {
	val, ok := u.role.get(id)
	if !ok {
		return UserInfo{}, ErrUnknownID
	}
	return val.(UserInfo), nil
}

// newSecureID는 랜덤 문자열을 생성해 준다.
func newSecureID(length int) (string, error) {

	r := make([]byte, length)
	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(r), nil
}

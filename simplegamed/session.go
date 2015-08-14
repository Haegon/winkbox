package simplegamed

// Session은 세션아이디와 핑 수를 가진다.
type Session struct {
	SessionID string
	PingCount int64
}

// increase는 핑수를 올려주는 함수.
func (s *Session) increase(i int) {
	s.PingCount += int64(i)
}

// reset은 핑수를 초기화 시켜주는 함수.
func (s *Session) reset() {
	s.PingCount = 0
}

// checkSessionID는 세션 아이디를 확인해주는 함수.
func (s *Session) checkSessionID(sid string) bool {
	if s.SessionID != sid {
		return false
	}
	return true
}

package rest

type sErr struct {
	Error string `json:"error,omitempty"`
}

type sEmail struct {
	Email string `json:"email,omitempty"`
}

type sToken struct {
	Token string `json:"token,omitempty"`
}

type sUser struct {
	User string `json:"user,omitempty"`
}

// status - используется для корректного вывода возвращаемых значений в swagger
type status struct {
	sErr
	sEmail
	sToken
	sUser
}

func s() *status {
	return &status{}
}

func (s *status) SetError(value error) sErr {
	s.Error = value.Error()
	return s.sErr
}

func (s *status) SetEmail(email string) sEmail {
	s.Email = email
	return s.sEmail
}

func (s *status) SetToken(token string) sToken {
	s.Token = token
	return s.sToken
}

func (s *status) SetUser(User string) sUser {
	s.User = User
	return s.sUser
}

func (s *status) Status() *status {
	return s
}

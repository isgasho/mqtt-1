package mqtt

////////////////////Interface//////////////////////////////

type SessionState int

type Session interface {
	GetRetransmitTimer() int
	SetRetransmitTimer(retransmitTimer int)

	GetState() SessionState

	Terminate()

	GetApplicationData() interface{}
	SetApplicationData(interface{})
}

////////////////////Implementation////////////////////////
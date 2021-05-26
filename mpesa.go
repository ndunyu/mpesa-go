package mpesa_go

type Mpesa struct {
	Live           bool
	ConsumerKey    string
	ConsumerSecret string
}

func NewMpesa(ConsumerKey, ConsumerSecret string, live bool) Mpesa {
	return Mpesa{
		Live:           live,
		ConsumerKey:    ConsumerKey,
		ConsumerSecret: ConsumerSecret,
	}

}

func (m *Mpesa) SetMode(mode bool) {
	m.Live = mode

}
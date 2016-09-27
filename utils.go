package npmp

import "errors"

func NewMessage(mt MessageType) Message {
	m := Message(make([]byte, 4))
	m.SetVersion(0)
	m.SetCookie(MagicCookie)
	m.SetMessageType(mt)
	return m
}

func NewRegisterMessage() *RegisterMessage {
	p := Message(make([]byte, 21))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(Register)
	return &RegisterMessage{Message: p}
}

func NewSettingsMessage() *SettingsMessage {
	return &SettingsMessage{Message: NewMessage(Settings)}
}

func NewDisconnectMessage() Message {
	return NewMessage(Disconnect)
}

func NewStartMessage() StartMessage {
	p := Message(make([]byte, 8))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(Start)
	return StartMessage{p}
}

func NewEndMessage() EndMessage {
	p := Message(make([]byte, 8))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(End)
	return EndMessage{p}
}

func NewDataMessage() DataMessage {
	p := Message(make([]byte, 9))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(Data)
	return DataMessage{p}
}

func NewInformMessage() InformMessage {
	return InformMessage{NewMessage(Inform)}
}

func NewVersionMessage() Message {
	return NewMessage(Version)
}

func NewACKMessage() Message {
	return NewMessage(ACK)
}

func NewNAKMessage() NAKMessage {
	p := Message(make([]byte, 5))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(NAK)
	return NAKMessage{p}
}

func ConvertToRegister(p Message) (*RegisterMessage, error) {
	if p.MessageType() != Register {
		return nil, errors.New("Incorrect message type")
	}

	r := &RegisterMessage{Message: p}
	if err := r.Process(); err != nil {
		return nil, err
	}
	return r, nil
}

func ConvertToSettings(p Message) (*SettingsMessage, error) {
	if p.MessageType() != Settings {
		return nil, errors.New("Incorrect message type")
	}

	r := &SettingsMessage{Message: p}
	if err := r.Process(); err != nil {
		return nil, err
	}
	return r, nil
}

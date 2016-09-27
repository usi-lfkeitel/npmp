package npmp

import "errors"

// newMessage will return a base Message of type mt. This function
// is for internal use only. Many message types require a slightly
// larger base.
func newMessage(mt MessageType) Message {
	m := Message(make([]byte, 4))
	m.SetVersion(0)
	m.SetCookie(MagicCookie)
	m.SetMessageType(mt)
	return m
}

// NewRegisterMessage returns a RegisterMessage with no interface information.
func NewRegisterMessage() *RegisterMessage {
	p := Message(make([]byte, 21))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(Register)
	return &RegisterMessage{Message: p}
}

// NewSettingsMessage returns a Settings Message with no Options.
func NewSettingsMessage() *SettingsMessage {
	return &SettingsMessage{Message: newMessage(Settings)}
}

// NewDisconnectMessage returns a Message of type Disconnect.
func NewDisconnectMessage() Message {
	return newMessage(Disconnect)
}

// NewStartMessage returns a StartMessage with a zeroed job ID.
func NewStartMessage() StartMessage {
	p := Message(make([]byte, 8))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(Start)
	return StartMessage{p}
}

// NewEndMessage returns a EndMessage with a zeroed job ID.
func NewEndMessage() EndMessage {
	p := Message(make([]byte, 8))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(End)
	return EndMessage{p}
}

// NewDataMessage returns a DataMessage with no data.
func NewDataMessage() DataMessage {
	p := Message(make([]byte, 9))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(Data)
	return DataMessage{p}
}

// NewInformMessage returns an InformMessage with no option codes.
func NewInformMessage() InformMessage {
	return InformMessage{newMessage(Inform)}
}

// NewVersionMessage returns a Message with type Version.
func NewVersionMessage() Message {
	return newMessage(Version)
}

// NewACKMessage returns a Message with type ACK.
func NewACKMessage() Message {
	return newMessage(ACK)
}

// NewNAKMessage returns an NAKMessage with a response code of 0.
func NewNAKMessage() NAKMessage {
	p := Message(make([]byte, 5))
	p.SetVersion(0)
	p.SetCookie(MagicCookie)
	p.SetMessageType(NAK)
	return NAKMessage{p}
}

// ConvertToRegister will take a Message and convert it into a RegisterMessage.
// It calles the Process() method on the RegisterMessage which parses the
// network interface information within the message. This is more of a
// convenience function.
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

// ConvertToSettings will take a Message and convert it into a SettingsMessage.
// It calles the Process() method on the SettingsMessage which parses the
// Option data within the message. This is more of a convenience function.
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

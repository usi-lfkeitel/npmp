package npmp

import (
	"encoding/binary"
	"errors"
	"net"
)

//go:generate stringer -type=OptionCode,MessageType,DataType,NACKResponseCode,NetType

var MagicCookie = []byte(`PM`)

type Option struct {
	Code  OptionCode
	Value []byte
}

type OptionCode byte
type MessageType byte
type DataType byte
type NACKResponseCode byte

type Messanger interface {
	Bytes() []byte
}

type Message []byte

func (p Message) Version() byte            { return p[0] }
func (p Message) Cookie() []byte           { return p[1:3] }
func (p Message) MessageType() MessageType { return MessageType(p[3]) }

func (p Message) SetVersion(version byte)       { p[0] = version }
func (p Message) SetCookie(cookie []byte)       { copy(p.Cookie(), cookie) }
func (p Message) SetMessageType(mt MessageType) { p[3] = byte(mt) }

func (p Message) Bytes() []byte { return p }

type RegisterMessage struct {
	Message
	Interfaces []*NetInterface
}

type NetType uint8

const (
	WiredEthernet    NetType = 0
	WirelessEthernet NetType = 1
)

type NetInterface struct {
	Type   NetType
	Haddr  net.HardwareAddr
	IPAddr net.IP
}

func (p *RegisterMessage) ClientID() []byte      { return p.Message[4:20] }
func (p *RegisterMessage) SetClientID(id []byte) { copy(p.ClientID(), id) }
func (p *RegisterMessage) IfCount() byte         { return p.Message[20] }
func (p *RegisterMessage) Process() error {
	c := int(p.IfCount())

	if len(p.Message) < 4+17+(11*c) {
		return errors.New("REGISTER message too small")
	}

	p.Interfaces = make([]*NetInterface, c)

	for i := 0; i < c; i++ {
		netif := &NetInterface{
			Type:   NetType(p.Message[21+(i*11)]),
			Haddr:  net.HardwareAddr(p.Message[22+(i*11) : 28+(i*11)]),
			IPAddr: net.IP(p.Message[28+(i*11) : 32+(i*11)]),
		}
		p.Interfaces[i] = netif
	}
	return nil
}

func (p *RegisterMessage) AddInterface(i *NetInterface) {
	if p.Interfaces == nil {
		p.Interfaces = make([]*NetInterface, 0)
	}
	p.Interfaces = append(p.Interfaces, i)
	p.Message[20] = byte(len(p.Interfaces))
}

func (p *RegisterMessage) Bytes() []byte {
	ret := p.Message[:4]
	ret = append(ret, p.ClientID()...)
	ret = append(ret, byte(len(p.Interfaces)))
	for _, i := range p.Interfaces {
		ret = append(ret, byte(i.Type))
		ret = append(ret, []byte(i.Haddr)...)
		ret = append(ret, []byte(i.IPAddr.To4())...)
	}
	return ret
}

type StartMessage struct {
	Message
}

func (p StartMessage) JobID() []byte      { return p.Message[4:8] }
func (p StartMessage) SetJobID(id []byte) { copy(p.JobID(), id) }

type EndMessage struct {
	Message
}

func (p EndMessage) JobID() []byte      { return p.Message[4:8] }
func (p EndMessage) SetJobID(id []byte) { copy(p.JobID(), id) }

type DataMessage struct {
	Message
}

func (p DataMessage) JobID() []byte      { return p.Message[4:8] }
func (p DataMessage) SetJobID(id []byte) { copy(p.JobID(), id) }

func (p DataMessage) Type() DataType         { return DataType(p.Message[8]) }
func (p DataMessage) SetDataType(t DataType) { p.Message[8] = byte(t) }

func (p DataMessage) Data() []byte {
	return p.Message[9:]
}
func (p *DataMessage) SetData(d []byte) {
	p.Message = append(p.Message[:9], d...)
}

type InformMessage struct {
	Message
}

func (p InformMessage) Options() []OptionCode {
	opts := p.Message[4:]
	o := make([]OptionCode, len(opts))
	for i, opt := range opts {
		o[i] = OptionCode(opt)
	}
	return o
}
func (p *InformMessage) SetOption(o OptionCode) { p.Message = append(p.Message, byte(o)) }
func (p *InformMessage) SetOptions(o []OptionCode) {
	p.Message = p.Message[:4] // Remove everything after the header
	for _, oc := range o {
		p.SetOption(oc)
	}
}

type NAKMessage struct {
	Message
}

func (p NAKMessage) ResponseCode() NACKResponseCode        { return NACKResponseCode(p.Message[4]) }
func (p NAKMessage) SetResponseCode(code NACKResponseCode) { p.Message[4] = byte(code) }

type SettingsMessage struct {
	Message
	Options []Option
}

func (p *SettingsMessage) Process() error {
	start := 4
	p.Options = make([]Option, 0)
	for start < len(p.Message) {
		len := binary.LittleEndian.Uint32(p.Message[5:9])
		p.Options = append(p.Options, Option{
			Code:  OptionCode(p.Message[start]),
			Value: p.Message[9:len],
		})
		start = start + 5 + int(len)
	}
	return nil
}

func (p *SettingsMessage) AddOption(o Option) {
	if p.Options == nil {
		p.Options = make([]Option, 0)
	}
	p.Options = append(p.Options, o)
}

func (p *SettingsMessage) StripOptions() {
	p.Message = p.Message[:4]
	p.Options = nil
}

func (p *SettingsMessage) Bytes() []byte {
	ret := p.Message[:4]
	l := make([]byte, 4)
	for _, o := range p.Options {
		// Lengths are 4 bytes long, must encode a slice of bytes
		binary.LittleEndian.PutUint32(l, uint32(len(o.Value)))
		ret = append(ret, byte(o.Code))
		ret = append(ret, l...)
		ret = append(ret, o.Value...)
	}
	return ret
}

// NPMP Message Types
const (
	Null       MessageType = 0
	Register   MessageType = 1
	Disconnect MessageType = 2
	Start      MessageType = 3
	End        MessageType = 4
	Data       MessageType = 5
	Inform     MessageType = 6
	Version    MessageType = 7
	ACK        MessageType = 8
	NAK        MessageType = 9
	Settings   MessageType = 10
)

// NPMP Option Codes
const (
	OpEnd                 OptionCode = 255
	Pad                   OptionCode = 0
	ServerIP              OptionCode = 1
	IperfServerAddress    OptionCode = 2
	IperfServerPort       OptionCode = 3
	IperfServerVersion    OptionCode = 4
	JobResourceDeadline   OptionCode = 5
	ProtocolVersion       OptionCode = 6
	ClientSoftwareVersion OptionCode = 7
	ClientSoftwareRepo    OptionCode = 8
	JobSpec               OptionCode = 9
	VendorOptions         OptionCode = 10
	HeartbeatDuration     OptionCode = 11
)

// NPMP Data Message Types
const (
	Ping   DataType = 0
	Iperf2 DataType = 1
	Iperf3 DataType = 2
)

// NPMP NACK Response Codes
const (
	GeneralError       NACKResponseCode = 0
	NotAuthorized      NACKResponseCode = 1
	UnsupportedVersion NACKResponseCode = 2
	NoPortsAvailable   NACKResponseCode = 3
	InvalidData        NACKResponseCode = 4
)

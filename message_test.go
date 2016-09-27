package npmp

import (
	"bytes"
	"net"
	"testing"
)

func TestBaseMessage(t *testing.T) {
	m := newMessage(Null)
	if m.MessageType() != Null {
		t.Fatalf("Incorrect message type. Expected Null, got %s", m.MessageType().String())
	}
	if int(m.Version()) != 0 {
		t.Fatalf("Incorrect protocol version. Expected 0, got %d", int(m.Version()))
	}
	if !bytes.Equal(m.Cookie(), MagicCookie) {
		t.Fatalf("Incorrect magic cookie. Expected %v, got %v", MagicCookie, m.Cookie())
	}

	newCookie := []byte{'A', 'M'}
	m.SetVersion(2)
	m.SetMessageType(Register)
	m.SetCookie(newCookie)
	if m.MessageType() != Register {
		t.Fatalf("Incorrect message type. Expected Register, got %s", m.MessageType().String())
	}
	if int(m.Version()) != 2 {
		t.Fatalf("Incorrect protocol version. Expected 2, got %d", int(m.Version()))
	}
	if !bytes.Equal(m.Cookie(), newCookie) {
		t.Fatalf("Incorrect magic cookie. Expected %v, got %v", newCookie, m.Cookie())
	}
}

func TestRegisterMessage(t *testing.T) {
	m := NewRegisterMessage()
	// Test empty message
	if m.MessageType() != Register {
		t.Fatalf("Incorrect message type. Expected Register, got %s", m.MessageType().String())
	}
	if int(m.IfCount()) != 0 {
		t.Fatalf("Incorrect interface length. Expected 0, got %d", int(m.IfCount()))
	}
	if !bytes.Equal(m.ClientID(), []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) {
		t.Fatalf("Incorrect client ID. Expected 0, got %v", m.ClientID())
	}
	if len(m.Interfaces) != 0 {
		t.Fatalf("Incorrect number of interfaces. Expected 0, got %d", len(m.Interfaces))
	}

	newClientID := []byte{99, 226, 170, 251, 37, 41, 43, 236, 249, 80, 159, 109, 149, 85, 244, 19}
	haddr := net.HardwareAddr([]byte{0xab, 0xcd, 0xef, 0x12, 0x34, 0x56})
	ipaddr := net.IP([]byte{192, 168, 0, 1})
	m.SetClientID(newClientID)
	m.AddInterface(&NetInterface{
		Type:   WirelessEthernet,
		Haddr:  haddr,
		IPAddr: ipaddr,
	})

	if !bytes.Equal(m.ClientID(), newClientID) {
		t.Fatalf("Incorrect client ID. Expected %v, got %v", newClientID, m.ClientID())
	}
	if int(m.IfCount()) != 1 {
		t.Fatalf("Incorrect interface length. Expected 1, got %d", int(m.IfCount()))
	}

	iface := m.Interfaces[0]
	if !bytes.Equal(iface.Haddr, haddr) {
		t.Fatalf("Incorrect interface MAC address. Expected %s, got %s", haddr.String(), iface.Haddr.String())
	}
	if !bytes.Equal(iface.IPAddr, ipaddr) {
		t.Fatalf("Incorrect interface IP address. Expected %s, got %s", ipaddr.String(), iface.IPAddr.String())
	}
	if iface.Type != WirelessEthernet {
		t.Fatalf("Incorrect interface type. Expected %s, got %s", WirelessEthernet.String(), iface.Type.String())
	}
}

func TestDisconnectMessage(t *testing.T) {
	m := NewDisconnectMessage()
	if m.MessageType() != Disconnect {
		t.Fatalf("Incorrect message type. Expected Disconnect, got %s", m.MessageType().String())
	}
}

func TestStartMessage(t *testing.T) {
	m := NewStartMessage()
	if m.MessageType() != Start {
		t.Fatalf("Incorrect message type. Expected Start, got %s", m.MessageType().String())
	}

	if !bytes.Equal(m.JobID(), []byte{0, 0, 0, 0}) {
		t.Fatalf("Incorrect job ID. Expected %v, got %v", []byte{0, 0, 0, 0}, m.JobID())
	}

	newID := []byte{250, 67, 39, 62}
	m.SetJobID(newID)
	if !bytes.Equal(m.JobID(), newID) {
		t.Fatalf("Incorrect job ID. Expected %v, got %v", newID, m.JobID())
	}
}

func TestEndMessage(t *testing.T) {
	m := NewEndMessage()
	if m.MessageType() != End {
		t.Fatalf("Incorrect message type. Expected End, got %s", m.MessageType().String())
	}

	if !bytes.Equal(m.JobID(), []byte{0, 0, 0, 0}) {
		t.Fatalf("Incorrect job ID. Expected %v, got %v", []byte{0, 0, 0, 0}, m.JobID())
	}

	newID := []byte{250, 67, 39, 62}
	m.SetJobID(newID)
	if !bytes.Equal(m.JobID(), newID) {
		t.Fatalf("Incorrect job ID. Expected %v, got %v", newID, m.JobID())
	}
}

func TestDataMessage(t *testing.T) {
	m := NewDataMessage()
	if m.MessageType() != Data {
		t.Fatalf("Incorrect message type. Expected Data, got %s", m.MessageType().String())
	}

	if !bytes.Equal(m.JobID(), []byte{0, 0, 0, 0}) {
		t.Fatalf("Incorrect job ID. Expected %v, got %v", []byte{0, 0, 0, 0}, m.JobID())
	}

	newID := []byte{250, 67, 39, 62}
	m.SetJobID(newID)
	if !bytes.Equal(m.JobID(), newID) {
		t.Fatalf("Incorrect job ID. Expected %v, got %v", newID, m.JobID())
	}

	if m.Type() != Ping {
		t.Fatalf("Incorrect data type. Expected %s, got %s", Ping.String(), m.Type().String())
	}
	m.SetDataType(Iperf3)
	if m.Type() != Iperf3 {
		t.Fatalf("Incorrect data type. Expected %s, got %s", Iperf3.String(), m.Type().String())
	}

	if len(m.Data()) != 0 {
		t.Fatalf("Data has length. Expected 0, got %d", len(m.Data()))
	}

	message := `The cow jumped over the moon`
	m.SetData([]byte(message))
	if string(m.Data()) != message {
		t.Fatalf("Wrong data. Expected %s, got %s", message, m.Data())
	}
}

func TestInformMessage(t *testing.T) {
	m := NewInformMessage()
	if m.MessageType() != Inform {
		t.Fatalf("Incorrect message type. Expected Inform, got %s", m.MessageType().String())
	}

	if len(m.Options()) != 0 {
		t.Fatalf("Options has length. Expected 0, got %d", len(m.Options()))
	}

	m.SetOption(ProtocolVersion)
	if len(m.Options()) != 1 {
		t.Fatalf("Options has length. Expected 1, got %d", len(m.Options()))
	}
	if m.Options()[0] != ProtocolVersion {
		t.Fatalf("Incorrect option code. Expected %s, got %s", ProtocolVersion.String(), (m.Options()[0]).String())
	}

	m.SetOptions([]OptionCode{ClientSoftwareVersion, ClientSoftwareRepo})
	if len(m.Options()) != 2 {
		t.Fatalf("Options has length. Expected 2, got %d", len(m.Options()))
	}
	if m.Options()[0] != ClientSoftwareVersion {
		t.Fatalf("Incorrect option code. Expected %s, got %s", ClientSoftwareVersion.String(), (m.Options()[0]).String())
	}
	if m.Options()[1] != ClientSoftwareRepo {
		t.Fatalf("Incorrect option code. Expected %s, got %s", ClientSoftwareRepo.String(), (m.Options()[1]).String())
	}
}

func TestVersionMessage(t *testing.T) {
	m := NewVersionMessage()
	if m.MessageType() != Version {
		t.Fatalf("Incorrect message type. Expected Version, got %s", m.MessageType().String())
	}
}

func TestAckMessage(t *testing.T) {
	m := NewACKMessage()
	if m.MessageType() != ACK {
		t.Fatalf("Incorrect message type. Expected Ack, got %s", m.MessageType().String())
	}
}

func TestNackMessage(t *testing.T) {
	m := NewNAKMessage()
	if m.MessageType() != NAK {
		t.Fatalf("Incorrect message type. Expected Nack, got %s", m.MessageType().String())
	}

	if m.ResponseCode() != GeneralError {
		t.Fatalf("Incorrect response code. Expected %s, got %s", GeneralError.String(), m.ResponseCode().String())
	}

	m.SetResponseCode(UnsupportedVersion)
	if m.ResponseCode() != UnsupportedVersion {
		t.Fatalf("Incorrect response code. Expected %s, got %s", UnsupportedVersion.String(), m.ResponseCode().String())
	}
}

func TestSettingsMessage(t *testing.T) {
	m := NewSettingsMessage()
	if m.MessageType() != Settings {
		t.Fatalf("Incorrect message type. Expected Settings, got %s", m.MessageType().String())
	}

	if len(m.Options) != 0 {
		t.Fatalf("Incorrect options length. Expected 0, got %d", len(m.Options))
	}

	repo := []byte(`http://repo.example.com/client/latest`)
	m.AddOption(Option{
		Code:  ClientSoftwareRepo,
		Value: repo,
	})
	if len(m.Options) != 1 {
		t.Fatalf("Incorrect options length. Expected 1, got %d", len(m.Options))
	}

	o := m.Options[0]
	if o.Code != ClientSoftwareRepo {
		t.Fatalf("Incorrect Option Code. Expected %s, for %s", ClientSoftwareRepo.String(), o.Code.String())
	}
	if !bytes.Equal(o.Value, repo) {
		t.Fatalf("Incorrect Option Value. Expected %s, got %s", repo, o.Value)
	}
}

package npmp

import (
	"bytes"
	"net"
	"testing"
)

func TestBaseMessage(t *testing.T) {
	m := NewMessage(Null)
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

}

func TestEndMessage(t *testing.T) {

}

func TestDataMessage(t *testing.T) {

}

func TestInformMessage(t *testing.T) {

}

func TestVersionMessage(t *testing.T) {

}

func TestAckMessage(t *testing.T) {

}

func TestNackMessage(t *testing.T) {

}

func TestSettingsMessage(t *testing.T) {

}

package network

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/sayan9168/sayanox-chain/internal/wallet"
)

func TestFramedMessageRoundTrip(t *testing.T) {
	left, right := net.Pipe()
	defer left.Close()
	defer right.Close()

	done := make(chan error, 1)
	go func() {
		done <- writeMessage(left, message{Type: "ping", Payload: []byte("hello")})
	}()

	msg, err := readMessage(bufio.NewReader(right))
	if err != nil {
		t.Fatal(err)
	}
	if msg.Type != "ping" || string(msg.Payload) != "hello" {
		t.Fatalf("unexpected message: %+v", msg)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestAuthenticatedHandshake(t *testing.T) {
	a, err := wallet.NewWallet()
	if err != nil { t.Fatal(err) }
	b, err := wallet.NewWallet()
	if err != nil { t.Fatal(err) }

	left, right := net.Pipe()
	defer left.Close()
	defer right.Close()

	type result struct {
		info handshakeInfo
		err  error
	}
	serverResult := make(chan result, 1)
	go func() {
		info, err := performHandshake(right, "sayanox-test", b.PrivateKey, false)
		serverResult <- result{info, err}
	}()

	clientResult := make(chan result, 1)
	go func() {
		info, err := performHandshake(left, "sayanox-test", a.PrivateKey, true)
		clientResult <- result{info, err}
	}()

	select {
	case r := <-clientResult:
		if r.err != nil { t.Fatal(r.err) }
		if r.info.PeerID == "" { t.Fatal("missing peer id") }
	case <-time.After(2 * time.Second):
		t.Fatal("client handshake timeout")
	}

	select {
	case r := <-serverResult:
		if r.err != nil { t.Fatal(r.err) }
		if r.info.PeerID == "" { t.Fatal("missing peer id") }
	case <-time.After(2 * time.Second):
		t.Fatal("server handshake timeout")
	}
}

func TestHandshakeRejectsWrongChain(t *testing.T) {
	a, err := wallet.NewWallet()
	if err != nil { t.Fatal(err) }
	b, err := wallet.NewWallet()
	if err != nil { t.Fatal(err) }

	left, right := net.Pipe()
	defer left.Close()
	defer right.Close()

	serverResult := make(chan error, 1)
	go func() {
		_, err := performHandshake(right, "chain-a", b.PrivateKey, false)
		serverResult <- err
	}()

	_, clientErr := performHandshake(left, "chain-b", a.PrivateKey, true)
	if clientErr == nil {
		t.Fatal("expected chain mismatch")
	}
	select {
	case <-serverResult:
	case <-time.After(2 * time.Second):
		t.Fatal("server handshake did not terminate")
	}
}

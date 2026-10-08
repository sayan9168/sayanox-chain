package network

import (
	"bufio"
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	protocolVersion  = 1
	maxFrameSize     = 2 << 20
	handshakeTimeout = 5 * time.Second
)

type message struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload,omitempty"`
}

type hello struct {
	Version   int    `json:"version"`
	ChainID   string `json:"chain_id"`
	NodeID    string `json:"node_id"`
	PublicKey []byte `json:"public_key"`
	Nonce     []byte `json:"nonce"`
	Signature []byte `json:"signature"`
}

type challenge struct {
	Nonce     []byte `json:"nonce"`
	Signature []byte `json:"signature"`
}

type handshakeInfo struct {
	PeerID  string
	Address string
	Version int
}

func writeMessage(w io.Writer, msg message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if len(data) == 0 || len(data) > maxFrameSize {
		return errors.New("message exceeds maximum frame size")
	}
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(data)))
	if _, err := w.Write(header[:]); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func readMessage(r *bufio.Reader) (message, error) {
	var header [4]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return message{}, err
	}
	size := binary.BigEndian.Uint32(header[:])
	if size == 0 || size > maxFrameSize {
		return message{}, errors.New("invalid message size")
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(r, data); err != nil {
		return message{}, err
	}
	var msg message
	if err := json.Unmarshal(data, &msg); err != nil {
		return message{}, err
	}
	return msg, nil
}

func makeNonce() ([]byte, error) {
	n := make([]byte, 32)
	_, err := rand.Read(n)
	return n, err
}

func nodeID(publicKey ed25519.PublicKey) string {
	sum := sha256.Sum256(publicKey)
	return hex.EncodeToString(sum[:])
}

func helloSigningBytes(h *hello) []byte {
	return []byte(fmt.Sprintf("sayanox/handshake/%d/%s/%s/%x", h.Version, h.ChainID, h.NodeID, h.Nonce))
}

func makeHello(chainID string, privateKey ed25519.PrivateKey) (hello, error) {
	nonce, err := makeNonce()
	if err != nil {
		return hello{}, err
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)
	h := hello{
		Version: protocolVersion,
		ChainID: chainID,
		NodeID: nodeID(publicKey),
		PublicKey: append([]byte(nil), publicKey...),
		Nonce: nonce,
	}
	h.Signature = ed25519.Sign(privateKey, helloSigningBytes(&h))
	return h, nil
}

func verifyHello(h hello, expectedChainID string) error {
	if h.Version != protocolVersion {
		return errors.New("unsupported protocol version")
	}
	if h.ChainID != expectedChainID {
		return errors.New("chain id mismatch")
	}
	if len(h.PublicKey) != ed25519.PublicKeySize {
		return errors.New("invalid public key")
	}
	if nodeID(ed25519.PublicKey(h.PublicKey)) != h.NodeID {
		return errors.New("invalid node id")
	}
	if len(h.Nonce) != 32 || len(h.Signature) != ed25519.SignatureSize {
		return errors.New("invalid handshake credentials")
	}
	if !ed25519.Verify(ed25519.PublicKey(h.PublicKey), helloSigningBytes(&h), h.Signature) {
		return errors.New("invalid handshake signature")
	}
	return nil
}

func performHandshake(conn net.Conn, chainID string, privateKey ed25519.PrivateKey, outbound bool) (handshakeInfo, error) {
	_ = conn.SetDeadline(time.Now().Add(handshakeTimeout))
	reader := bufio.NewReader(conn)

	localHello, err := makeHello(chainID, privateKey)
	if err != nil {
		return handshakeInfo{}, err
	}

	var remoteHello hello
	if outbound {
		if err := writeMessage(conn, message{Type: "hello", Payload: mustJSON(localHello)}); err != nil {
			return handshakeInfo{}, err
		}
		msg, err := readMessage(reader)
		if err != nil || msg.Type != "hello" {
			return handshakeInfo{}, errors.New("expected peer hello")
		}
		if err := json.Unmarshal(msg.Payload, &remoteHello); err != nil {
			return handshakeInfo{}, err
		}
		if err := verifyHello(remoteHello, chainID); err != nil {
			return handshakeInfo{}, err
		}

		challengeNonce, err := makeNonce()
		if err != nil {
			return handshakeInfo{}, err
		}
		ch := challenge{Nonce: challengeNonce, Signature: ed25519.Sign(privateKey, challengeNonce)}
		if err := writeMessage(conn, message{Type: "challenge", Payload: mustJSON(ch)}); err != nil {
			return handshakeInfo{}, err
		}
		msg, err = readMessage(reader)
		if err != nil || msg.Type != "challenge_response" {
			return handshakeInfo{}, errors.New("expected challenge response")
		}
		var response challenge
		if err := json.Unmarshal(msg.Payload, &response); err != nil {
			return handshakeInfo{}, err
		}
		if !bytes.Equal(response.Nonce, challengeNonce) ||
			!ed25519.Verify(ed25519.PublicKey(remoteHello.PublicKey), challengeNonce, response.Signature) {
			return handshakeInfo{}, errors.New("peer challenge verification failed")
		}
	} else {
		msg, err := readMessage(reader)
		if err != nil || msg.Type != "hello" {
			return handshakeInfo{}, errors.New("expected peer hello")
		}
		if err := json.Unmarshal(msg.Payload, &remoteHello); err != nil {
			return handshakeInfo{}, err
		}
		if err := verifyHello(remoteHello, chainID); err != nil {
			return handshakeInfo{}, err
		}
		if err := writeMessage(conn, message{Type: "hello", Payload: mustJSON(localHello)}); err != nil {
			return handshakeInfo{}, err
		}

		msg, err = readMessage(reader)
		if err != nil || msg.Type != "challenge" {
			return handshakeInfo{}, errors.New("expected challenge")
		}
		var ch challenge
		if err := json.Unmarshal(msg.Payload, &ch); err != nil {
			return handshakeInfo{}, err
		}
		if len(ch.Nonce) != 32 || len(ch.Signature) != ed25519.SignatureSize {
			return handshakeInfo{}, errors.New("invalid challenge")
		}
		if !ed25519.Verify(ed25519.PublicKey(remoteHello.PublicKey), ch.Nonce, ch.Signature) {
			return handshakeInfo{}, errors.New("invalid challenge signature")
		}
		response := challenge{Nonce: ch.Nonce, Signature: ed25519.Sign(privateKey, ch.Nonce)}
		if err := writeMessage(conn, message{Type: "challenge_response", Payload: mustJSON(response)}); err != nil {
			return handshakeInfo{}, err
		}
	}

	_ = conn.SetDeadline(time.Time{})
	return handshakeInfo{
		PeerID:  remoteHello.NodeID,
		Version: remoteHello.Version,
	}, nil
}

func mustJSON(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}

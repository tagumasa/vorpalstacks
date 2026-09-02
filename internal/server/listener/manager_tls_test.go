package listener

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// testTLSCertificate builds a minimal self-signed certificate for the SNI
// resolution tests.
func testTLSCertificate(t *testing.T) *tls.Certificate {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "listener-test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{"sni.example.com"},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
	if err != nil {
		t.Fatal(err)
	}
	return &tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}
}

// freePort reserves an ephemeral port and releases it for the listener.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

// waitForListen polls until the listener goroutine has bound the port; the
// manager starts servers asynchronously.
func waitForListen(t *testing.T, port int) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+strconv.Itoa(port), 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("listener on port %d never became reachable", port)
}

func TestTLSListenerServesSNI(t *testing.T) {
	cert := testTLSCertificate(t)
	port := freePort(t)

	var seen []string
	m := NewManager(50080)
	m.Register(ListenerConfig{
		Name:        "tls-test",
		DefaultPort: port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("tls-ok"))
		}),
		TLS: &ListenerTLSConfig{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				seen = append(seen, info.ServerName)
				return cert, nil
			},
		},
	})
	m.Start()
	defer m.Shutdown(context.Background())
	waitForListen(t, port)

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         "sni.example.com",
			},
		},
	}
	resp, err := client.Get("https://127.0.0.1:" + strconv.Itoa(port))
	if err != nil {
		t.Fatalf("TLS request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if len(seen) != 1 || seen[0] != "sni.example.com" {
		t.Fatalf("SNI names seen = %v, want [sni.example.com]", seen)
	}
}

func TestTLSListenerFailsWithoutCertificate(t *testing.T) {
	port := freePort(t)
	m := NewManager(50080)
	m.Register(ListenerConfig{
		Name:        "tls-nil-test",
		DefaultPort: port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("unreachable"))
		}),
		TLS: &ListenerTLSConfig{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return nil, nil
			},
		},
	})
	m.Start()
	defer m.Shutdown(context.Background())

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         "sni.example.com",
			},
		},
	}
	if _, err := client.Get("https://127.0.0.1:" + strconv.Itoa(port)); err == nil {
		t.Fatal("a handshake without a certificate must fail")
	}
}

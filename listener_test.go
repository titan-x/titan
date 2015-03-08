package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestLen(t *testing.T) {
	a, _ := strconv.Atoi("12344324")
	t.Log(a)
}

func TestListener(t *testing.T) {
	var wg sync.WaitGroup
	cert, privKey := genCert(t)
	listener, err := Listen(cert, privKey, "localhost:8091", true)
	if err != nil {
		t.Fatal(err)
	}

	go listener.Accept(func(msg []byte) {
		wg.Add(1)
		defer wg.Done()
		t.Logf("Incoming message to listener from a client: %v", string(msg))
	})

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(cert)
	if !ok {
		panic("failed to parse root certificate")
	}

	conn, err := tls.Dial("tcp", "localhost:8091", &tls.Config{RootCAs: roots})
	if err != nil {
		t.Fatal(err)
	}

	send(t, conn, "4   ping")
	send(t, conn, "56  Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
	send(t, conn, "49  In sit amet lectus felis, at pellentesque turpis.")
	send(t, conn, "64  Nunc urna enim, cursus varius aliquet ac, imperdiet eget tellus.")
	send(t, conn, "8671Lorem ipsum dolor sit amet, consectetur adipiscing elit. In sit amet lectus felis, at pellentesque turpis. Nunc urna enim, cursus varius aliquet ac, imperdiet eget tellus. Morbi ultricies ullamcorper sollicitudin. Nam varius ante condimentum massa elementum ac elementum risus vehicula. Nunc euismod convallis vestibulum. Nulla facilisi. Nulla facilisi. Sed quis arcu ipsum, eu luctus orci. Nullam et mauris sagittis ligula accumsan elementum non ullamcorper magna. Donec imperdiet quam et nibh vulputate quis eleifend ligula commodo. Quisque ac augue dui. Curabitur sodales quam in nulla consectetur tincidunt. Nam non diam nisl, ut placerat risus. Curabitur ac leo tellus, sed mollis urna. Nunc vulputate, nunc ut semper sodales, diam risus fringilla felis, sed tristique ipsum nunc eget lorem. Pellentesque molestie bibendum ante ac ultrices. Suspendisse rutrum ante et nulla volutpat nec ultrices quam eleifend. Sed id nisi mauris. Ut elementum ultricies velit et posuere. Morbi eu commodo eros. Aliquam non urna nibh, nec laoreet lectus. Maecenas id tellus dui. Curabitur a lacus ac augue viverra elementum. Cras molestie tempus sodales. Aenean nibh est, blandit at luctus sed, porttitor non odio. Curabitur mi ligula, porta quis imperdiet sed, laoreet at felis. Donec ac mi auctor sapien consequat luctus id nec neque. Nullam varius nibh eu risus condimentum facilisis et non elit. Proin at mi justo, sit amet fermentum dolor. Aliquam pharetra tincidunt dapibus. Donec eget rutrum velit. Vivamus mauris dui, posuere nec luctus vitae, mollis pharetra mauris. Vivamus tincidunt tellus dictum orci auctor quis congue massa euismod. Donec pharetra enim vitae sapien auctor dictum. Etiam a leo et quam malesuada dignissim. Donec vestibulum, neque non gravida adipiscing, nunc lectus iaculis sem, a vehicula leo dolor vel nibh. Etiam ac augue arcu. In libero massa, aliquam id egestas sed, ullamcorper ut mauris. Donec sem nisl, tincidunt eu viverra vel, vulputate in nisi. Nunc facilisis congue iaculis. Quisque interdum, nibh ut aliquam pulvinar, elit nisi consectetur odio, in adipiscing felis massa nec neque. Sed risus nulla, pretium sed ullamcorper ac, faucibus eget mi. Aenean fermentum orci sit amet tortor tincidunt mattis. Nam venenatis, odio sit amet volutpat pharetra, massa ipsum laoreet neque, in volutpat orci tellus ut nisl. Morbi interdum mattis ipsum sed consectetur. Nam dapibus ipsum sed ligula feugiat scelerisque. Maecenas ac commodo velit. Sed justo mauris, ultricies vel varius nec, egestas non neque. Maecenas purus justo, molestie a imperdiet vel, venenatis ut purus. Mauris fermentum, felis placerat ullamcorper fringilla, turpis velit consequat mauris, eu cursus erat sapien sed arcu. Donec vitae leo felis, et lacinia neque. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Vivamus tristique, justo eu iaculis ornare, sem eros viverra felis, fermentum blandit enim lacus at elit. Nam quis risus in leo commodo rutrum. Aliquam tempor nisl quis dui vulputate et sagittis lectus lobortis. Aenean malesuada consequat ante sit amet laoreet. Aenean porta tristique imperdiet. Ut diam massa, lobortis a eleifend fermentum, fringilla non orci. Curabitur eleifend elementum tortor, ut porttitor arcu blandit sed. In sed est vel orci sagittis sollicitudin. Vivamus molestie eros sit amet ante volutpat eget egestas eros volutpat. In vulputate augue in elit sagittis venenatis. Nulla facilisi. Etiam vulputate quam id neque luctus facilisis. Praesent ac justo ante. Nam fermentum diam vel libero pharetra bibendum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis mauris massa, dapibus et eleifend ut, aliquam eget eros. Integer quam orci, vehicula at dapibus vel, eleifend ac tortor. Mauris sagittis gravida risus nec facilisis. Sed lorem sapien, congue et consectetur nec, tincidunt malesuada neque. Integer elementum, sem ut dignissim pulvinar, odio elit pharetra urna, et placerat massa tellus a nunc. Mauris porta aliquam consequat. Aliquam aliquet placerat molestie. Cras hendrerit ullamcorper velit ut ultrices. Sed consequat, neque in tristique sodales, quam ante egestas orci, id placerat diam risus sed ante. Phasellus diam nisi, pellentesque nec dictum in, tristique eu ligula. Aliquam mollis nulla sit amet libero luctus nec porttitor augue euismod. Lorem ipsum dolor sit amet, consectetur adipiscing elit. In sit amet lectus felis, at pellentesque turpis. Nunc urna enim, cursus varius aliquet ac, imperdiet eget tellus. Morbi ultricies ullamcorper sollicitudin. Nam varius ante condimentum massa elementum ac elementum risus vehicula. Nunc euismod convallis vestibulum. Nulla facilisi. Nulla facilisi. Sed quis arcu ipsum, eu luctus orci. Nullam et mauris sagittis ligula accumsan elementum non ullamcorper magna. Donec imperdiet quam et nibh vulputate quis eleifend ligula commodo. Quisque ac augue dui. Curabitur sodales quam in nulla consectetur tincidunt. Nam non diam nisl, ut placerat risus. Curabitur ac leo tellus, sed mollis urna. Nunc vulputate, nunc ut semper sodales, diam risus fringilla felis, sed tristique ipsum nunc eget lorem. Pellentesque molestie bibendum ante ac ultrices. Suspendisse rutrum ante et nulla volutpat nec ultrices quam eleifend. Sed id nisi mauris. Ut elementum ultricies velit et posuere. Morbi eu commodo eros. Aliquam non urna nibh, nec laoreet lectus. Maecenas id tellus dui. Curabitur a lacus ac augue viverra elementum. Cras molestie tempus sodales. Aenean nibh est, blandit at luctus sed, porttitor non odio. Curabitur mi ligula, porta quis imperdiet sed, laoreet at felis. Donec ac mi auctor sapien consequat luctus id nec neque. Nullam varius nibh eu risus condimentum facilisis et non elit. Proin at mi justo, sit amet fermentum dolor. Aliquam pharetra tincidunt dapibus. Donec eget rutrum velit. Vivamus mauris dui, posuere nec luctus vitae, mollis pharetra mauris. Vivamus tincidunt tellus dictum orci auctor quis congue massa euismod. Donec pharetra enim vitae sapien auctor dictum. Etiam a leo et quam malesuada dignissim. Donec vestibulum, neque non gravida adipiscing, nunc lectus iaculis sem, a vehicula leo dolor vel nibh. Etiam ac augue arcu. In libero massa, aliquam id egestas sed, ullamcorper ut mauris. Donec sem nisl, tincidunt eu viverra vel, vulputate in nisi. Nunc facilisis congue iaculis. Quisque interdum, nibh ut aliquam pulvinar, elit nisi consectetur odio, in adipiscing felis massa nec neque. Sed risus nulla, pretium sed ullamcorper ac, faucibus eget mi. Aenean fermentum orci sit amet tortor tincidunt mattis. Nam venenatis, odio sit amet volutpat pharetra, massa ipsum laoreet neque, in volutpat orci tellus ut nisl. Morbi interdum mattis ipsum sed consectetur. Nam dapibus ipsum sed ligula feugiat scelerisque. Maecenas ac commodo velit. Sed justo mauris, ultricies vel varius nec, egestas non neque. Maecenas purus justo, molestie a imperdiet vel, venenatis ut purus. Mauris fermentum, felis placerat ullamcorper fringilla, turpis velit consequat mauris, eu cursus erat sapien sed arcu. Donec vitae leo felis, et lacinia neque. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Vivamus tristique, justo eu iaculis ornare, sem eros viverra felis, fermentum blandit enim lacus at elit. Nam quis risus in leo commodo rutrum. Aliquam tempor nisl quis dui vulputate et sagittis lectus lobortis. Aenean malesuada consequat ante sit amet laoreet. Aenean porta tristique imperdiet. Ut diam massa, lobortis a eleifend fermentum, fringilla non orci. Curabitur eleifend elementum tortor, ut porttitor arcu blandit sed. In sed est vel orci sagittis sollicitudin. Vivamus molestie eros sit amet ante volutpat eget egestas eros volutpat. In vulputate augue in elit sagittis venenatis. Nulla facilisi. Etiam vulputate quam id neque luctus facilisis. Praesent ac justo ante. Nam fermentum diam vel libero pharetra bibendum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis mauris massa, dapibus et eleifend ut, aliquam eget eros. Integer quam orci, vehicula at dapibus vel, eleifend ac tortor. Mauris sagittis gravida risus nec facilisis. Sed lorem sapien, congue et consectetur nec, tincidunt malesuada neque. Integer elementum, sem ut dignissim pulvinar, odio elit pharetra urna, et placerat massa tellus a nunc. Mauris porta aliquam consequat. Aliquam aliquet placerat molestie. Cras hendrerit ullamcorper velit ut ultrices. Sed consequat, neque in tristique sodales, quam ante egestas orci, id placerat diam risus sed ante. Phasellus diam nisi, pellentesque nec dictum in, tristique eu ligula. Aliquam mollis nulla sit amet libero luctus nec porttitor augue euismod.")
	send(t, conn, "5   close")

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	conn.Close()
	listener.Close()
}

func send(t *testing.T, conn *tls.Conn, msg string) {
	n, err := io.WriteString(conn, msg)
	if err != nil {
		t.Fatalf("Error while writing message to connection %v", err)
	}
	t.Logf("Sending message to listener from client: %v (%v bytes)", msg, n)
}

// Generate a self-signed PEM encoded X.509 certificate and private key pair (i.e. 'cert.pem', 'key.pem').
// Based on the sample from http://golang.org/src/crypto/tls/generate_cert.go (taken at Jan 30, 2015).
func genCert(t *testing.T) (pemBytes, privBytes []byte) {
	hosts := []string{"localhost"}
	privKey, err := rsa.GenerateKey(rand.Reader, 512)
	notBefore := time.Now()
	notAfter := notBefore.Add(290 * 365 * 24 * time.Hour) //290 years
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		t.Fatalf("failed to generate serial number: %s", err)
	}

	cert := x509.Certificate{
		IsCA:         true,
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"localhost"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			cert.IPAddresses = append(cert.IPAddresses, ip)
		} else {
			cert.DNSNames = append(cert.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, &privKey.PublicKey, privKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %s", err)
	}

	pemBytes = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	privBytes = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	return
}

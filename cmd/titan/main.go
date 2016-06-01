package main

import (
	"flag"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/titan-x/titan"
)

const (
	addr     = "127.0.0.1:3000"
	testAddr = "127.0.0.1:3001"
)

var (
	run   = flag.Bool("run", false, "Start the Titan server.")
	caddr = flag.String("addr", addr, "Specifies a network address to start the server on. If not specific, default address will be used: "+addr)
	ext   = flag.Bool("ext", false, "Run external client test case. Titan server will run at address: "+testAddr)
)

func main() {
	flag.Parse()

	switch {
	case *run:
		startServer(*caddr)
	case *ext:
		startExtTest(testAddr)
	default:
		flag.PrintDefaults()
	}
}

func startServer(addr string) {
	s, err := titan.NewServer(addr)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	defer func() {
		if s.Close(); err != nil {
			log.Printf("error closing server: %v", err)
		}
	}()

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("error listening for connections: %v", err)
	}
}

func startExtTest(addr string) {
	log.Printf("-ext flag is provided, starting external client test case.")
	titan.InitConf("test")

	now := time.Now().Unix()
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims["userid"] = "1"
	t.Claims["created"] = now
	ts, err := t.SignedString([]byte(titan.Conf.App.JWTPass()))
	if err != nil {
		log.Fatalf("failed to sign JWT token: %v", err)
	}
	log.Printf("Sample valid user JWT token for testing: %v", ts)

	startServer(addr)
}

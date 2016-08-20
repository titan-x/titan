package main

import (
	"flag"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/titan-x/titan"
	"github.com/titan-x/titan/data/aws"
)

const (
	addr     = "127.0.0.1:3000"
	testAddr = "127.0.0.1:3001"
)

var (
	defaultFlag = flag.Bool("default", false, "Start Titan server at default address: "+addr)
	addrFlag    = flag.String("addr", "", "Start Titan server with specified address parameter.")
	testFlag    = flag.Bool("test", false, "Start Titan server for external client integration test at address: "+testAddr)
	awsFlag     = flag.Bool("aws", false, "Enable Amazon Web Services support.")
)

func main() {
	flag.Parse()

	switch {
	case *testFlag:
		startExtTest(testAddr)
	case *defaultFlag:
		startServer(addr)
	case *addrFlag != "":
		startServer(*addrFlag)
	default:
		flag.PrintDefaults()
	}
}

func startServer(addr string) {
	s, err := titan.NewServer(addr)
	if err != nil {
		log.Fatalf("error creating server: %v", err)
	}

	if *awsFlag {
		s.SetDB(aws.NewDynamoDB("", ""))
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

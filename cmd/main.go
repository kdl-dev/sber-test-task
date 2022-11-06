//go:build debug

package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/kdl-dev/sber-test-task/pkg/global"
	"github.com/kdl-dev/sber-test-task/pkg/test"
)

var addr string
var secureProtocol string = "https://"
var unSecureProtocol string = "http://"
var isSecureConn bool

func init() {
	flag.StringVar(&addr, "u", "", "URL of the test start.")
	flag.BoolVar(&isSecureConn, "s", false, "Secure remote host.")
	flag.BoolVar(&global.IsVerbose, "v", false, "Verbose information about the test result")
	flag.Parse()

	if addr == "" {
		log.Fatal("test start url not provided")
	}

	if isSecureConn {
		addr = secureProtocol + addr
	}

	if !isSecureConn {
		addr = unSecureProtocol + addr
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	newTest, err := test.NewTest(addr)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	log.Printf("Test start\n\n")
	global.PrintVerboseInfo(global.CLI_Border)
	successMsg, err := newTest.SolveTest()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	log.Printf("%s\n\n", successMsg)
	global.PrintVerboseInfo(global.CLI_Border)
	log.Println("Test finish")
}

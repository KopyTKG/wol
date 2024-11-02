package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"thekrew.app/WoL/help"
)

func main() {
	args := os.Args[1:]
	sMac := ""

	if len(args) < 1 {
		help.PrintHelp()
		os.Exit(0)
	}

	processArguments(args, &sMac)
	wol(&sMac)

}

func processArguments(args []string, sMac *string) {
	prev := ""
	for _, arg := range args {
		switch {
		case arg == "-m" || arg == "--mac":
			prev = arg

		case prev == "-m" || prev == "--mac":
			prev = ""
			*sMac = arg

		case arg == "-h" || arg == "--help":
			help.PrintHelp()
			os.Exit(1)
		case arg == "--version":
			help.Version()
			os.Exit(1)

		default:
			log.Fatalf("Unknowed switch found (%s) run \"wol -h\"\n", arg)
			os.Exit(1)
		}
	}
}
func assebleMac(slMac []string) (mac []byte) {
	mac = make([]byte, 6)

	if len(slMac) != 6 {
		log.Fatalln("invalid MAC address format")
		os.Exit(1)
	}

	for i, hexStr := range slMac {
		b, err := hex.DecodeString(hexStr)
		if err != nil {
			log.Fatalf("invalid MAC address octet: %s\n", hexStr)
			os.Exit(1)
		}
		mac[i] = b[0]
	}
	return mac
}

func assemble(mac []byte) (msg []byte) {
	l := 102
	msg = make([]byte, l)

	for i := range int(6) {
		msg[i] = 0xFF
	}

	for i := range int(l - 6) {
		msg[i+6] = mac[i%6]
	}

	return msg
}


func wol(sMac *string) {
	mac := assebleMac(strings.Split(*sMac, ":"))	
	msg := assemble(mac)


	Baddr := net.UDPAddr{
		Port: 40000,
		IP:   net.ParseIP("255.255.255.255"),
	}

	conn, err := net.DialUDP("udp", nil, &Baddr)
	if err != nil {
		log.Fatalln("Error while connection to UDP node")
		os.Exit(1)
	}


	_, err = conn.Write(msg)
	if err != nil {
		log.Fatalf("Error sending message: %s\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Sending WoL to %s\n", *sMac)
	conn.Close()
}

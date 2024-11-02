package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"thekrew.app/WoL/help"
)

func assemble(mac []byte) []byte {
	l := 102
	msg := make([]byte, l)

	for i := range int(6) {
		msg[i] = 0xFF
	}

	for i := range int(l - 6) {
		msg[i+6] = mac[i%6]
	}

	return msg
}

func join(arr []int, el string) (item string) {
	for i := range len(arr) - 1 {
		item += strconv.Itoa(arr[i]) + el
	}
	return item + strconv.Itoa(arr[len(arr)-1])
}

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
			log.Fatalf("Unknowed switch found (%s) run \"wol -h\"", arg)
			os.Exit(1)
		}
	}
}

func wol(sMac *string) {
	slMac := strings.Split((*sMac), ":")
	mac := make([]byte, 6)

	if len(slMac) != 6 {
		log.Fatalln("invalid MAC address format")
		return
	}

	for i, hexStr := range slMac {
		b, err := hex.DecodeString(hexStr)
		if err != nil {
			log.Fatalf("invalid MAC address octet: %s", hexStr)
			return
		}
		mac[i] = b[0]
	}


	Baddr := net.UDPAddr{
		Port: 40000,
		IP:   net.ParseIP("255.255.255.255"),
	}
	conn, err := net.DialUDP("udp", nil, &Baddr)
	if err != nil {
		return
	}

	msg := assemble(mac)

	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	conn.Close()
}

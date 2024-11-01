package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"

	"thekrew.app/WoL/help"
)

func getAddrMask(intf string) (ip []int, mask []int) {

	ip = make([]int, 4)
	mask = make([]int, 4)

	ief, err := net.InterfaceByName(intf)
	if err != nil { // get interface
		log.Fatalln("Error while getting interface")
		os.Exit(1)
	}
	addrs, err := ief.Addrs()
	if err != nil { // get addresses
		log.Fatalln("Error while getting addresses")
		os.Exit(1)
	}

	first := ""

	for _, addr := range addrs {
		if strings.Contains(addr.String(), ".") {
			first = addr.String()
			break
		}
	}

	slices := strings.Split(first, "/")
	sIp := strings.Split(slices[0], ".")
	iMask, err := strconv.Atoi(slices[1])
	if err != nil {
		os.Exit(1)
	}
	for i := range len(sIp) {
		ip[i], err = strconv.Atoi(sIp[i])
		if err != nil {
		log.Fatalln("Error while converting string to number")
		os.Exit(1)
		}

	}
	if iMask == 0 {
		log.Fatalln("invalid mask")
		os.Exit(1)
	}
	left := iMask

	for i := 0; i < iMask; i += 8 {
		left -= 8
		mask[i/8] = 255
	}

	if left > 0 {
		mask[3] = int(math.Pow(2, float64(left)))
	}

	return ip, mask
}

func getBroadcast(ip []int, mask []int) (brd []int) {
	brd = make([]int, 4)

	full := []int{255, 255, 255, 255}
	inv := make([]int, 4)

	for i := range len(full) {
		inv[i] = full[i] - mask[i]
	}

	for i := range len(ip) {
		brd[i] = ip[i] | inv[i]
	}

	return brd
}

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
	intf, sMac := "", ""

	if len(args) < 1 {
		help.PrintHelp()
		os.Exit(0)
	}

	processArguments(args, &intf, &sMac)
	wol(&intf, &sMac)

}

func processArguments(args []string, intf, sMac *string) {
	prev := ""
	for _, arg := range args {
		switch {
		case arg == "-i" || arg == "--interface":
			prev = arg
		case arg == "-m" || arg == "--mac":
			prev = arg

		case prev == "-i" || prev == "--interface":
			prev = ""
			*intf = arg
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

func wol(intf, sMac *string) {
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

	ip, mask := getAddrMask(*intf)
	brd := getBroadcast(ip, mask)

	host := join(brd, ".")

	Baddr := net.UDPAddr{
		Port: 40000,
		IP:   net.ParseIP(host),
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

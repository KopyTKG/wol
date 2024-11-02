package help

import "fmt"

func PrintHelp() {
	helpText := `Usage: wol -m [xx:xx:xx:xx:xx:xx]

  * -m, --mac              		Mac address for magic packet
    -h, --help                  	Display this help message
    --version               		Display the version of wol

  * marks required switch

Any errors please report to: <https://github.com/kopytkg/wol/issues>

usage
$ wol -m [xx:xx:xx:xx:xx:xx]
`
	fmt.Print(helpText)
}

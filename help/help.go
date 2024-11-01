package help

import "fmt"

func PrintHelp() {
	helpText := `Usage: wol -i [name] -m [xx:xx:xx:xx:xx:xx]

  * -i, --interface               	Network interface
  * -m, --mac              		Mac address for magic packet
    -h, --help                  	Display this help message
    --version               		Display the version of wol

  * marks required switch

Any errors please report to: <https://github.com/kopytkg/wol/issues>

usage
$ wol -i [name] -m [xx:xx:xx:xx:xx:xx]
`
	fmt.Print(helpText)
}

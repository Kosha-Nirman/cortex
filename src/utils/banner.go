package utils

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	version = "1.0.0"
)

func PrintBanner() {
	banner := `
	▄████▄   ▒█████   ██▀███  ▄███████▓ ▓█████ ▒██   ██▒
	▒██▀ ▀█  ▒██▒  ██▒▓██   ██▒▓  ██▒ ▓▒▓█   ▀ ▒▒ █ █ ▒░
	▒▓█    ▄ ▒██░  ██▒▓██  ░█ ▒▒  ██░ ▒░▒███   ░░  █   ░
	▒▓▓▄ ▄██▒▒██   ██░▒██▀▀█▄  ░  ██▓ ░ ▒▓█  ▄  ░ █ █ ▒ 
	▒ ▓███▀ ░░ ████▓▒░░██▓ ▒██▒   ██▒ ░ ░▒████▒▒██▒ ▒██▒
	░ ░▒ ▒  ░░ ▒░▒░▒░ ░██▓ ░▒▓░  ▒ ░░   ░░ ▒░ ░▒▒ ░ ░▓ ░
		░  ▒     ░ ▒ ▒░   ░▒ ░ ▒░    ░     ░ ░  ░░░   ░▒ ░
	░        ░ ░ ░ ▒    ░░   ░   ░         ░    ░    ░  
	░ ░          ░ ░     ░                 ░  ░ ░    ░  
		`
	fmt.Printf("%s\n", color.RedString(banner))
	fmt.Printf("%s\n", color.WhiteString("         Subdomain Resolver & Reconnaissance Tool v"+version))
	fmt.Printf("%s\n\n", color.YellowString("              https://github.com/Kosha-Nirman/cortex"))
}

package config

import "github.com/fatih/color"

var Banner = `
██████╗ ██╗███╗   ██╗ ██████╗ ██████╗ ███╗   ██╗
██╔══██╗██║████╗  ██║██╔════╝██╔═══██╗████╗  ██║
██████╔╝██║██╔██╗ ██║██║     ██║   ██║██╔██╗ ██║
██╔══██╗██║██║╚██╗██║██║     ██║   ██║██║╚██╗██║
██║  ██║██║██║ ╚████║╚██████╗╚██████╔╝██║ ╚████║
╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝
`

func PrintStartupBanner() {
	banner := color.New(color.Bold, color.FgHiBlue).PrintlnFunc()
	banner(Banner)
	version := color.New(color.Bold, color.FgBlue).PrintlnFunc()
	if Env == "DEV" {
		version("Running v" + Version + " in Development mode.")
	} else {
		version("Running v" + Version + " in Production mode.")
	}
	println()
}

// Package banner renders the URLShine ASCII art header.
package banner

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

const Version = "2.0.0"
const Author  = "URLShine Team"

func Print() {
	cyan   := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	green  := color.New(color.FgGreen).SprintFunc()
	faint  := color.New(color.Faint).SprintFunc()
	white  := color.New(color.FgWhite, color.Bold).SprintFunc()

	fmt.Println()
	fmt.Println(cyan(`  ██╗   ██╗██████╗ ██╗      ███████╗██╗  ██╗██╗███╗   ██╗███████╗`))
	fmt.Println(cyan(`  ██║   ██║██╔══██╗██║       ██╔════╝██║  ██║██║████╗  ██║██╔════╝`))
	fmt.Println(cyan(`  ██║   ██║██████╔╝██║       ███████╗███████║██║██╔██╗ ██║█████╗  `))
	fmt.Println(cyan(`  ██║   ██║██╔══██╗██║       ╚════██║██╔══██║██║██║╚██╗██║██╔══╝  `))
	fmt.Println(cyan(`  ╚██████╔╝██║  ██║███████╗  ███████║██║  ██║██║██║ ╚████║███████╗`))
	fmt.Println(cyan(`   ╚═════╝ ╚═╝  ╚═╝╚══════╝  ╚══════╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝`))
	fmt.Println()
	fmt.Printf("  %s   %s\n",
		yellow("URLShine v"+Version),
		faint("— Advanced URL Enumeration & Attack Surface Mapper"),
	)
	fmt.Printf("  %s\n",
		green("  GAU · Gospider · Katana · Waymore · Waybackurls · Hakrawler · XnLinkFinder"),
	)
	fmt.Printf("  %s  %s  %s\n",
		faint("go/"+runtime.Version()),
		faint("·"),
		faint(runtime.GOOS+"/"+runtime.GOARCH),
	)
	fmt.Println()
	fmt.Printf("  %s\n", faint(strings.Repeat("─", 67)))
	fmt.Println()
	_ = white // reserved
}

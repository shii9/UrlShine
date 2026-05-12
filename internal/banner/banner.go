// Package banner renders the URLShine ASCII art header and session information.
package banner

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

const Version = "2.0.1"
const Author = "URLShine Team"

// Print displays the professional URLShine banner with system information.
func Print() {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	faint := color.New(color.Faint).SprintFunc()

	fmt.Println()
	fmt.Println(cyan(`  ██╗   ██╗██████╗ ██╗      ███████╗██╗  ██╗██╗███╗   ██╗███████╗`))
	fmt.Println(cyan(`  ██║   ██║██╔══██╗██║       ██╔════╝██║  ██║██║████╗  ██║██╔════╝`))
	fmt.Println(cyan(`  ██║   ██║██████╔╝██║       ███████╗███████║██║██╔██╗ ██║█████╗  `))
	fmt.Println(cyan(`  ██║   ██║██╔══██╗██║       ╚════██║██╔══██║██║██║╚██╗██║██╔══╝  `))
	fmt.Println(cyan(`  ╚██████╔╝██║  ██║███████╗  ███████║██║  ██║██║██║ ╚████║███████╗`))
	fmt.Println(cyan(`   ╚═════╝ ╚═╝  ╚═╝╚══════╝  ╚══════╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝`))
	fmt.Println()
	fmt.Printf("  %s   %s\n",
		yellow("v"+Version),
		faint("URL Enumeration & Attack Surface Mapper"),
	)
	fmt.Printf("  %s\n\n",
		faint("go/"+runtime.Version()+" · "+runtime.GOOS+"/"+runtime.GOARCH),
	)
}

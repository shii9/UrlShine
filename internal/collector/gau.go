package collector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shii9/UrlShine/internal/utils"
)

// runGAU collects URLs via gau with aggressive parameters using multiple providers and options.
func runGAU(target, outDir string, cfg Config) ([]string, error) {
	threads := cfg.Threads
	if threads < 100 {
		threads = 100
	}

	outFile := filepath.Join(outDir, fmt.Sprintf("gau_%s.txt", utils.SanitizeFilename(target)))
	_ = os.Remove(outFile)

	// Run GAU with comprehensive parameters
	args := []string{
		"gau", target,
		"--threads", fmt.Sprintf("%d", threads),
		"--providers", "wayback,commoncrawl,urlscan,otx",
		"--blacklist", "png,jpg,jpeg,gif,bmp,svg,ico,webp,css,eot,ttf,woff,woff2,pdf,zip,rar,tar,gz,mp4,mp3,avi,webm,mkv,mov,flv,swf,wma,wav,aac,flac,m4a,ogg,webm,3gp,mkv,mov,mp4,mpeg,mpg,wmv,avi,m3u8,m3u,pls",
	}
	if cfg.Subs {
		args = append(args, "--subs")
	}

	lines, err := runCmd(args...)
	if err == nil && len(lines) > 0 {
		if err := utils.WriteLines(outFile, lines); err != nil {
			return lines, nil
		}
	}

	return lines, nil
}

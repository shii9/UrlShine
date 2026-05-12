package collector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shii9/UrlShine/internal/utils"
)

// runWaymore collects URLs via Waymore with aggressive mode 3 and targeted filters.
func runWaymore(target, outDir string, cfg Config) ([]string, error) {
	conc := 50
	if cfg.Threads > 100 {
		conc = 100
	}

	outFile := filepath.Join(outDir, fmt.Sprintf("waymore_%s.txt", utils.SanitizeFilename(target)))
	_ = os.Remove(outFile)

	args := []string{
		"waymore",
		"-i", target,
		"-mode", "3",
		"-p",
		"-formatted",
		"-exclude", "png,jpg,jpeg,gif,bmp,svg,ico,webp,css,woff,woff2,eot,ttf,pdf,zip,rar,tar,gz,mp4,mp3,avi,webm,mkv,mov,flv,swf,wma,wav",
		"-include", "api,auth,admin,user,v1,v2,v3,endpoint,callback,redirect",
		"-concurrency", fmt.Sprintf("%d", conc),
		"-oU", outFile,
	}

	lines, err := runCmd(args...)
	if err == nil && len(lines) > 0 {
		return lines, nil
	}

	// Try to read output file if command succeeded
	if utils.FileExists(outFile) {
		lines, _ := utils.ReadLines(outFile)
		return lines, nil
	}

	return nil, err
}

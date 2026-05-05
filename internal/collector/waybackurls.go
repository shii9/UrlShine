package collector

// runWaybackurls collects URLs via waybackurls.
func runWaybackurls(target, _ string, _ Config) ([]string, error) {
	return runCmdStdin(target+"\n", "waybackurls")
}

package cmd

import (
	"github.com/chomosuke/cf-tool/client"
)

// Watch command
func Watch(args ParsedArgs) (err error) {
	cln := client.Instance
	info := args.Info
	n := 10
	if args.All {
		n = -1
	}
	if _, err = cln.WatchSubmission(info, n, false); err != nil {
		if err = loginAgain(cln, err); err == nil {
			_, err = cln.WatchSubmission(info, n, false)
		}
	}
	return
}

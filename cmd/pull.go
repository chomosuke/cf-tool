package cmd

import (
	"os"

	"github.com/chomosuke/cf-tool/client"
)

// Pull command
func Pull(args ParsedArgs) (err error) {
	cln := client.Instance
	info := args.Info
	ac := args.Accepted
	rootPath, err := os.Getwd()
	if err != nil {
		return
	}
	if err = cln.Pull(info, rootPath, ac); err != nil {
		if err = loginAgain(cln, err); err == nil {
			err = cln.Pull(info, rootPath, ac)
		}
	}
	return
}

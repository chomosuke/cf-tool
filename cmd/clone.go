package cmd

import (
	"os"

	"github.com/chomosuke/cf-tool/client"
)

// Clone command
func Clone(args ParsedArgs) (err error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return
	}
	cln := client.Instance
	ac := args.Accepted
	handle := args.Handle

	if err = cln.Clone(handle, currentPath, ac); err != nil {
		if err = loginAgain(cln, err); err == nil {
			err = cln.Clone(handle, currentPath, ac)
		}
	}
	return
}

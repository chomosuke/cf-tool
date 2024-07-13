package cmd

import (
	"github.com/chomosuke/cf-tool/client"
)

// Setup command
func Setup(parsedArgs ParsedArgs) error {
	return client.Instance.Setup(parsedArgs.Handle, parsedArgs.Password)
}

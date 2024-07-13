package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/chomosuke/cf-tool/client"
	"github.com/fatih/color"
)

// Submit command
func Submit(args ParsedArgs) (err error) {
	cln := client.Instance
	info := args.Info
	filename := args.File
	langId := args.LangId

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.New("Failed to read source code from " + filename)
	}
	source := string(bytes)

	_, ok := client.Langs[langId]
	if !ok {
		return errors.New("No languages with ID " + langId)
	}

	err = cln.Submit(info, langId, source)
	if err == nil {
		// Submitted successfully in the first try
		return nil
	}

	// Only retry if the client is not logged in
	if err.Error() != client.ErrorNotLogged {
		return
	}

	color.Red("Not logged in. Trying to login\n")
	if err = cln.Login(); err != nil {
		return errors.New("Failed to login")
	}

	if err = cln.Submit(info, langId, source); err != nil {
		return errors.New("Failed to resubmit even after reauthentication")
	}

	return nil;
}

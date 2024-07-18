package cmd

import (
	"github.com/fatih/color"
	"github.com/skratchdot/open-golang/open"
	"github.com/chomosuke/cf-tool/client"
	"github.com/chomosuke/cf-tool/config"
)

func openURL(url string) error {
	color.Green("Open %v", url)
	return open.Run(url)
}

// Open command
func Open(args ParsedArgs) (err error) {
	URL, err := args.Info.OpenURL(config.Instance.Host)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Stand command
func Stand(args ParsedArgs) (err error) {
	URL, err := args.Info.StandingsURL(config.Instance.Host)
	if err != nil {
		return
	}
	return openURL(URL)
}

// Sid command
func Sid(args ParsedArgs) (err error) {
	info := args.Info
	if info.SubmissionID == "" && client.Instance.LastSubmission != nil {
		info = *client.Instance.LastSubmission
	}
	URL, err := info.SubmissionURL(config.Instance.Host)
	if err != nil {
		return
	}
	return openURL(URL)
}

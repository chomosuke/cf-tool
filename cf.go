package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chomosuke/cf-tool/client"
	"github.com/chomosuke/cf-tool/cmd"
	"github.com/chomosuke/cf-tool/config"
	"github.com/fatih/color"
	ansi "github.com/k0kubun/go-ansi"
	"github.com/mitchellh/go-homedir"

	docopt "github.com/docopt/docopt-go"
)

const version = "v1.0.0"

const usage = `
Codeforces Tool $%version%$ (cf). https://github.com/chomosuke/cf-tool

You should run "cf config" to configure your handle, password and code
templates at first.

If you want to compete, the best command is "cf race"

Usage:
  cf setup (-u <handle>) (-p <pass>)
  cf submit (-u <handle>) (-l <lang>) [-f <file>] [<specifier>...]
  cf upgrade

Options:
  -h --help            Show this screen.
  --version            Show version.
  -f <file>, --file <file>, <file>
                       Path to file. E.g. "a.cpp", "./temp/a.cpp".
  -u <handle>, --handle <handle>, <handle>
  -p <pass>, --pass <pass>, <pass>
  -l <lang>, --lang <lang>, <lang>
  <specifier>          Any useful text. E.g.
                       "https://codeforces.com/contest/100",
                       "https://codeforces.com/contest/180/problem/A",
                       "https://codeforces.com/group/Cw4JRyRGXR/contest/269760"
                       "1111A", "1111", "a", "Cw4JRyRGXR"
                       You can combine multiple specifiers to specify what you
                       want.
  <alias>              Template's alias. E.g. "cpp"
  ac                   The status of the submission is Accepted.

Examples:
  cf setup -u handle -p password      Set up the codeforces account.
  cf submit -u handle -l 54 -f a.cpp
  cf submit -u handle -l 54 https://codeforces.com/contest/100/A
  cf submit -u handle -l 54 -f a.cpp 100A
  cf submit -u handle -l 54 -f a.cpp 100 a
  cf submit -u handle -l 54 contest 100 a
  cf submit -u handle -l 54 gym 100001 a
  cf upgrade           Upgrade the "cf" to the latest version from GitHub.

File:
  cf will save some data in some files:

  "~/.cf/config/${handle}"        Configuration file, including templates, etc.
  "~/.cf/session/${handle}"       Session file, including cookies, handle, password, etc.

  "~" is the home directory of current user in your system.
`
const configPathPrefix = "~/.cf/config"
const sessionPathPrefix = "~/.cf/session"

func main() {
	color.Output = ansi.NewAnsiStdout()

	opts, _ := docopt.ParseArgs(
    strings.Replace(usage, `$%version%$`, version, 1),
    os.Args[1:], fmt.Sprintf("Codeforces Tool (cf) %v",
    version))
	opts[`{version}`] = version

  handle, _ := opts["--handle"].(string)

	cfgPath, _ := homedir.Expand(filepath.Join(configPathPrefix, handle))
	clnPath, _ := homedir.Expand(filepath.Join(sessionPathPrefix, handle))

	config.Init(cfgPath)
	client.Init(clnPath, config.Instance.Host, config.Instance.Proxy)

	err := cmd.Eval(opts)
	if err != nil {
    color.Red(err.Error())
    os.Exit(1)
	}
}

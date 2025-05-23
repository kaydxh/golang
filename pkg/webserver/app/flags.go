/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package app

import (
	"fmt"
	"os"
	"path/filepath"

	flag_ "github.com/kaydxh/golang/go/flag"
	term_ "github.com/kaydxh/golang/go/os/term"
	"github.com/spf13/cobra"
)

// defaultConfigPath returns config file's default path
func defaultConfigPath() string {
	return fmt.Sprintf("./conf/%s.yaml", filepath.Base(os.Args[0]))
}

func defaultUseConfigPath() string {
	return fmt.Sprintf(
		filepath.Join(filepath.Dir(defaultConfigPath()), ".use.%s.yaml"), filepath.Base(os.Args[0]),
	)
}

type AppFlags struct {
	ConfigFile    string
	UseConfigFile string
	cmd           *cobra.Command
	flags         flag_.NamedFlagSets
}

//  NewAppFlags if cmd is nil, then use default cobra.Command
func NewAppFlags(cmd *cobra.Command) *AppFlags {
	appFlags := &AppFlags{
		ConfigFile:    defaultConfigPath(),
		UseConfigFile: defaultUseConfigPath(),
		cmd:           cmd,
	}

	if appFlags.cmd == nil {
		appFlags.cmd = &cobra.Command{}
	}

	appFlags.initFlags()
	return appFlags
}

func (f *AppFlags) Apply() {
	fs := f.cmd.Flags()
	for _, flag := range f.flags.FlagSets {
		fs.AddFlagSet(flag)
	}
}

/*
func (f *AppFlags) AddFlags(mainfs *pflag.FlagSet) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	fs.StringVar(&f.SeaConfigFile, "config", f.SeaConfigFile, "sea config file")
}
*/

func (f *AppFlags) initFlags() {

	fs := f.flags.FlagSet("misc")
	fs.StringVarP(&f.ConfigFile, "config", "c", f.ConfigFile, "The path to the configuration file.")
	fs.StringVar(
		&f.UseConfigFile,
		"use-config",
		f.UseConfigFile,
		"If set, write the configuration values to this file and exit.",
	)

}

func (f *AppFlags) SetUsageAndHelpFunc() {
	cols, _, _ := term_.TerminalSize(f.cmd.OutOrStdout())
	flag_.SetUsageAndHelpFunc(f.cmd, f.flags, cols)
}

func (f *AppFlags) Install() {
	f.SetUsageAndHelpFunc()
	f.Apply()
}

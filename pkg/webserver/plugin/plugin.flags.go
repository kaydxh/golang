package plugin

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

type SealetFlags struct {
	ConfigFile    string
	UseConfigFile string
	cmd           *cobra.Command
	flags         flag_.NamedFlagSets
}

//  NewSealetFlags if cmd is nil, then use default cobra.Command
func NewSealetFlags(cmd *cobra.Command) *SealetFlags {
	sealetFlags := &SealetFlags{
		ConfigFile:    defaultConfigPath(),
		UseConfigFile: defaultUseConfigPath(),
		cmd:           cmd,
	}

	if sealetFlags.cmd == nil {
		sealetFlags.cmd = &cobra.Command{}
	}

	sealetFlags.initFlags()
	return sealetFlags
}

func (f *SealetFlags) Apply() {
	fs := f.cmd.Flags()
	for _, flag := range f.flags.FlagSets {
		fs.AddFlagSet(flag)
	}
}

/*
func (f *SealetFlags) AddFlags(mainfs *pflag.FlagSet) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	fs.StringVar(&f.SeaConfigFile, "config", f.SeaConfigFile, "sea config file")
}
*/

func (f *SealetFlags) initFlags() {

	fs := f.flags.FlagSet("misc")
	fs.StringVarP(&f.ConfigFile, "config", "c", f.ConfigFile, "The path to the configuration file.")
	fs.StringVar(
		&f.UseConfigFile,
		"use-config",
		f.UseConfigFile,
		"If set, write the configuration values to this file and exit.",
	)

}

func (f *SealetFlags) SetUsageAndHelpFunc() {
	cols, _, _ := term_.TerminalSize(f.cmd.OutOrStdout())
	flag_.SetUsageAndHelpFunc(f.cmd, f.flags, cols)
}

func (f *SealetFlags) Install() {
	f.SetUsageAndHelpFunc()
	f.Apply()
}

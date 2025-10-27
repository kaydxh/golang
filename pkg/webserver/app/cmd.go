/*
 *Copyright (c) 2023, kaydxh
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
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// NewSeaDateCommand creates a *cobra.Command object with default parameters
func NewCommand(ctx context.Context, runCommand func(ctx context.Context, cmd *cobra.Command) error) *cobra.Command {
	name := GetVersion().AppName
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s Public HTTP/2 and GRPC APIs", name),
		// stop printing usage when the command errors
		Long: fmt.Sprintf(`%s is a gateway serve which you can use curl over HTTP 1.1 or grpc protocal on the same host:port.
Example: curl -X POST -k https://localhost:port/healthz
See [Sea](https://github.com/kaydxh/sea/blob/master/README.md) for more information.`, name),
		//SilenceUsage: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(ctx, cmd)
		},

		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("server exit")
			return nil
		},

		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					//%q a single-quoted character literal safely escaped with Go syntax
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	appFlag := NewAppFlags(cmd)
	appFlag.Install()
	return cmd
}

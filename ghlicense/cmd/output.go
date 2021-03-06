/*
Sniperkit-Bot
- Status: analyzed
*/

// Copyright 2016 Nick Miyake. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package cmd

import (
	"io/ioutil"
	"time"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/sniperkit/snk.fork.ghcli/common"
	"github.com/sniperkit/snk.fork.ghcli/license"
)

const (
	authorFlagName   = "author"
	outputFlagName   = "output"
	licenseParamName = "license"
)

var (
	authorFlag = flag.StringFlag{
		Name:  authorFlagName,
		Usage: "author to use for copyright",
	}
	outputFlag = flag.StringFlag{
		Name:  outputFlagName,
		Usage: "file to which to write license",
		Value: "LICENSE",
	}
	licenseParam = flag.StringParam{
		Name:  licenseParamName,
		Usage: "license type",
	}
)

func Write() cli.Command {
	cmd := outputCommand("write", "write the content of a license to a file", func(ctx cli.Context, license string) error {
		return ioutil.WriteFile(ctx.String(outputFlagName), []byte(license), 0644)
	})
	cmd.Flags = append(cmd.Flags, outputFlag)
	return cmd
}

func Print() cli.Command {
	return outputCommand("print", "print the content of a license", func(ctx cli.Context, license string) error {
		ctx.Printf(license)
		return nil
	})
}

func outputCommand(name, usage string, f func(ctx cli.Context, license string) error) cli.Command {
	tokenFlag := common.GitHubTokenFlag
	tokenFlag.Required = false
	return cli.Command{
		Name:  name,
		Usage: usage,
		Flags: []flag.Flag{
			tokenFlag,
			authorFlag,
			licenseParam,
		},
		Action: func(ctx cli.Context) error {
			year := time.Now().Year()
			cache := license.NewCache(common.CachingOAuthGitHubClient(ctx.String(common.GitHubTokenFlagName), ""))
			license, err := license.Create(ctx.String(licenseParamName), cache, license.NewAuthorInfo(ctx.String(authorFlagName), year, year))
			if err != nil {
				return err
			}
			if err := f(ctx, license); err != nil {
				return err
			}
			return nil
		},
	}
}

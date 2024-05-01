//go:build !solution

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal/app/common"
	"gitlab.com/slon/shad-go/gitfame/pkg/flag"
)

func unquote(s string) string {
	if strings.HasPrefix(s, "'") {
		return s[1 : len(s)-1]
	}

	return s
}

func validateParams(p *common.Params) error {
	possibleOrderBy := map[string]struct{}{
		"lines":   {},
		"commits": {},
		"files":   {},
	}

	if _, ok := possibleOrderBy[p.OrderBy]; !ok {
		return fmt.Errorf("bad order by param")
	}

	return nil
}

func getParams() (*common.Params, error) {
	p := &common.Params{}
	if err := flag.Commandline.ParseStruct(p); err != nil {
		return nil, fmt.Errorf("invalid flags format")
	}

	if p.ExcludeArg != "" {
		p.Exclude = strings.Split(unquote(p.ExcludeArg), ",")
		for i := range p.Exclude {
			p.Exclude[i] = p.RepositoryPath + "/" + p.Exclude[i]
		}
	}

	if p.RestrictToArg != "" {
		p.RestrictTo = strings.Split(unquote(p.RestrictToArg), ",")
		for i := range p.RestrictTo {
			if strings.Contains(p.RestrictTo[i], "/") {
				p.RestrictTo[i] = p.RepositoryPath + "/" + p.RestrictTo[i]
			}
		}
	}

	if p.LanguagesArg != "" {
		p.Languages = strings.Split(unquote(p.LanguagesArg), ",")
	}

	if p.ExtensionsArg != "" {
		p.Extensions = strings.Split(unquote(p.ExtensionsArg), ",")
	}

	return p, nil
}

func main() {
	ctx := context.Background()
	p, err := getParams()
	if err != nil {
		os.Exit(1)
	}

	if err := validateParams(p); err != nil {
		os.Exit(1)
	}

	app := &App{
		params: p,
	}
	if err := app.Init(ctx); err != nil {
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		os.Exit(1)
	}
}

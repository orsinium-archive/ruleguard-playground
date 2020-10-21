package main

import (
	"context"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Runner struct {
	loader     *lint.ContextLoader
	reportData *report.Data
	cfg        *config.Config
	log        logutils.Log
}

func RunLinters() ([]result.Issue, error) {
	cfg := config.NewDefault()
	reportData := report.Data{}
	log := report.NewLogWrapper(logutils.NewMockLog(), &reportData)
	DBManager := lintersdb.NewManager(cfg, log).WithCustomLinters()
	linters := DBManager.GetAllSupportedLinterConfigs()

	ctx := context.TODO()

	issues := make([]result.Issue, 0)
	for _, lc := range linters {
		linterIssues, err := lc.Linter.Run(ctx, nil)
		if err != nil {
			return nil, err
		}
		issues = append(issues, linterIssues...)
	}

	return issues, nil
}

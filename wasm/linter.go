package main

import (
	"fmt"
	"time"

	"context"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis/load"
	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
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
	loadGuard := load.NewGuard()
	goenv := goutil.NewEnv(log)
	fileCache := fsutils.NewFileCache()
	lineCache := fsutils.NewLineCache(fileCache)
	contextLoader := lint.NewContextLoader(cfg, log, goenv, lineCache, fileCache, nil, loadGuard)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	DBManager := lintersdb.NewManager(cfg, log).WithCustomLinters()
	validator := lintersdb.NewValidator(DBManager)
	EnabledLintersSet := lintersdb.NewEnabledSet(DBManager, validator, log, cfg)

	lintersToRun, err := EnabledLintersSet.GetOptimizedLinters()
	if err != nil {
		return nil, fmt.Errorf("cannot get optimized linters: %v", err)
	}

	enabledLintersMap, err := EnabledLintersSet.GetEnabledLintersMap()
	if err != nil {
		return nil, fmt.Errorf("cannot get enabled linters map: %v", err)
	}

	for _, lc := range DBManager.GetAllSupportedLinterConfigs() {
		isEnabled := enabledLintersMap[lc.Name()] != nil
		reportData.AddLinter(lc.Name(), isEnabled, lc.EnabledByDefault)
	}

	lintCtx, err := contextLoader.Load(ctx, lintersToRun)
	if err != nil {
		return nil, fmt.Errorf("context loading failed: %v", err)
	}
	lintCtx.Log = log

	runner, err := lint.NewRunner(cfg, log,
		goenv, EnabledLintersSet, lineCache, DBManager, lintCtx.Packages)
	if err != nil {
		return nil, fmt.Errorf("cannot get linters runner: %v", err)
	}
	issues, err := runner.Run(ctx, lintersToRun, lintCtx)
	if err != nil {
		return nil, fmt.Errorf("cannot run linters: %v", err)
	}

	fixer := processors.NewFixer(cfg, log, fileCache)
	return fixer.Process(issues), nil
}

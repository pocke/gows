package main

import (
	"encoding/json"
	"os"
	"sync"
)

type AnalyzeResult struct {
	Executions []Command
	Program    []Command
}

type Analyzer struct {
	ch     chan Command
	result *AnalyzeResult
}

func StartAnalysis(cs []Command) (chan<- Command, *sync.WaitGroup) {
	wg := &sync.WaitGroup{}
	if os.Getenv("GOWS_ANALYSIS") == "" {
		wg.Add(0)
		return nil, wg
	}
	wg.Add(1)

	ch := make(chan Command, 10)
	res := &AnalyzeResult{Program: cs}
	a := &Analyzer{ch: ch, result: res}
	go a.watch(wg)
	return ch, wg
}

func (a *Analyzer) watch(wg *sync.WaitGroup) {
	for v := range a.ch {
		a.result.Executions = append(a.result.Executions, v)
	}

	f, err := os.OpenFile("./gows-analysis.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(a.result)
	if err != nil {
		panic(err)
	}
	wg.Done()
}

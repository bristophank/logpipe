package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logpipe/internal/config"
	"github.com/yourorg/logpipe/internal/filter"
	"github.com/yourorg/logpipe/internal/formatter"
	"github.com/yourorg/logpipe/internal/metrics"
	"github.com/yourorg/logpipe/internal/pipeline"
	"github.com/yourorg/logpipe/internal/ratelimit"
	"github.com/yourorg/logpipe/internal/router"
)

func main() {
	cfgPath := flag.String("config", "logpipe.json", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logpipe: %v\n", err)
		os.Exit(1)
	}

	fmt_ := formatter.New(cfg.Format)
	limiter := ratelimit.New(cfg.RateLimit)
	col := metrics.New()

	rt := router.New()
	for _, s := range cfg.Sinks {
		var w *os.File
		if s.Type == "file" && s.Target != "" {
			w, err = os.OpenFile(s.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "logpipe: open sink %s: %v\n", s.Name, err)
				os.Exit(1)
			}
			defer w.Close()
		} else {
			w = os.Stdout
		}
		rt.Add(s.Name, w)
	}

	var rules []filter.Rule
	if len(cfg.Routes) > 0 {
		for k, v := range cfg.Routes[0].Filter {
			rules = append(rules, filter.Rule{Field: k, Value: v})
		}
	}

	f := filter.New(rules)
	pl := pipeline.New(f, rt, fmt_, col, limiter)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		pl.Process(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "logpipe: read stdin: %v\n", err)
	}

	snap := col.Snapshot()
	fmt.Fprintf(os.Stderr, "logpipe: processed=%d routed=%d dropped=%d\n",
		snap.Processed, snap.Routed, snap.Dropped)
}

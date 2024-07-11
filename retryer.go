package retryer

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-retry"
)

type OptionRetries struct {
	Type           string         `yaml:"type" json:"type"`
	Base           OptionDuration `yaml:"base" json:"base"`
	CappedDuration OptionDuration `yaml:"capped_duration" json:"capped_duration"`
	JitterPercent  uint64         `yaml:"jitter_percent" json:"jitter_percent"`
	MaxDuration    OptionDuration `yaml:"max_duration" json:"max_duration"`
}

var (
	backoffFuncMap = map[string](func(t time.Duration) retry.Backoff){
		"constant":    retry.NewConstant,
		"exponential": retry.NewExponential,
		"fibonacci":   retry.NewFibonacci,
	}
)

func getBackoff(options *OptionRetries) (b retry.Backoff, err error) {
	if f, found := backoffFuncMap[options.Type]; found {
		b = f(time.Duration(options.Base))
	} else {
		validTypes := []string{}
		for validType := range backoffFuncMap {
			validTypes = append(validTypes, validType)
		}
		return nil, fmt.Errorf("unknown retry type '%s', valid types are [%#v]", options.Type, validTypes)
	}
	if options.CappedDuration > 0 {
		b = retry.WithCappedDuration(time.Duration(options.CappedDuration), b)
	}
	if options.JitterPercent > 0 {
		b = retry.WithJitterPercent(options.JitterPercent, b)
	}
	if options.MaxDuration > 0 {
		b = retry.WithMaxDuration(time.Duration(options.MaxDuration), b)
	}
	return
}

func RunWithRetries(ctx context.Context, options *OptionRetries, f retry.RetryFunc) (err error) {
	if options == nil {
		return f(ctx)
	}
	b, err := getBackoff(options)
	if err != nil {
		return fmt.Errorf("unable to create the backoff object from retries config '%#v: %w", *options, err)
	}
	if err := retry.Do(ctx, b, f); err != nil {
		return fmt.Errorf("no more retries: %w", err)
	}
	return
}

package internal

import "github.com/mytokenio/go/metrics/internal/lv"

type observeFunc func(name string, lvs lv.LabelValues, value float64)
type valFunc func(name string) float64

package internal

import (
	"github.com/mytokenio/go/metrics"
	"github.com/mytokenio/go/metrics/internal/lv"
)

func NewGauge(name string, obs observeFunc, add observeFunc, val valFunc, lvs lv.LabelValues) metrics.Gauge {
	return &Gauge{
		name: name,
		lvs:  lvs,
		obs:  obs,
		add:  add,
		val:  val,
	}
}

type Gauge struct {
	name string
	lvs  lv.LabelValues
	obs  observeFunc
	add  observeFunc
	val  valFunc
}

func (g *Gauge) With(labelValues ...string) metrics.Gauge {
	return &Gauge{
		name: g.name,
		lvs:  g.lvs.With(labelValues...),
		obs:  g.obs,
		add:  g.add,
	}
}

func (g *Gauge) Set(value int64) {
	g.obs(g.name, g.lvs, float64(value))
}

func (g *Gauge) Value() int64 {
	return int64(g.val(g.name))
}

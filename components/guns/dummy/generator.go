package dummy

import (
	"time"

	"github.com/yandex/pandora/core"
	"github.com/yandex/pandora/core/aggregator/netsample"
	"github.com/yandex/pandora/core/warmup"
)

type GunConfig struct {
	Sleep time.Duration `config:"sleep"`
}

type Gun struct {
	DebugLog bool
	Conf     GunConfig
	Aggr     core.Aggregator
	core.GunDeps
}

func DefaultGunConfig() GunConfig {
	return GunConfig{}
}

func (g *Gun) WarmUp(_ *warmup.Options) (any, error) {
	return nil, nil
}

func NewGun(conf GunConfig) *Gun {
	return &Gun{Conf: conf}
}

func (g *Gun) Bind(aggr core.Aggregator, deps core.GunDeps) error {
	g.Aggr = aggr
	g.GunDeps = deps
	return nil
}

func (g *Gun) Shoot(_ core.Ammo) {
	g.shoot()
}

func (g *Gun) shoot() {
	code := 0
	sample := netsample.Acquire("")
	defer func() {
		sample.SetProtoCode(code)
		g.Aggr.Report(sample)
	}()

	time.Sleep(g.Conf.Sleep)
	code = 200
}

var _ warmup.WarmedUp = (*Gun)(nil)

package acceptance

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
	"github.com/yandex/pandora/core/engine"
	"github.com/yandex/pandora/lib/testutil"
	"go.uber.org/zap"
)

func TestDummyGunSuite(t *testing.T) {
	suite.Run(t, new(DummyGunSuite))
}

type DummyGunSuite struct {
	suite.Suite
	fs      afero.Fs
	log     *zap.Logger
	metrics engine.Metrics
}

func (s *DummyGunSuite) SetupSuite() {
	s.fs = afero.NewOsFs()
	testOnce.Do(importDependencies(s.fs))

	s.log = testutil.NewNullLogger()
	s.metrics = engine.NewMetrics("dummy_suite")
}

func (s *DummyGunSuite) Test_Shoot() {
	tests := []struct {
		name           string
		filecfg        string
		isTLS          bool
		preStartSrv    func(srv *httptest.Server)
		wantErrContain string
		wantCnt        int
	}{
		{
			name:    "dummy",
			filecfg: "testdata/dummy/dummy.yaml",
			wantCnt: 6,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {

			conf := parseConfigFile(s.T(), tt.filecfg, "")
			s.Require().Equal(1, len(conf.Engine.Pools))
			aggr := &aggregator{}
			conf.Engine.Pools[0].Aggregator = aggr
			pandora := engine.New(s.log, s.metrics, conf.Engine)

			err := pandora.Run(context.Background())
			if tt.wantErrContain != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tt.wantErrContain)
				return
			}
			s.Require().NoError(err)
			s.Require().Equal(int64(tt.wantCnt), int64(len(aggr.samples)))
		})
	}
}

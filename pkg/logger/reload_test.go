package logger

import (
	"context"
	"testing"
	"time"

	"github.com/ProtocolONE/go-core/v2/pkg/config"
	"github.com/ProtocolONE/go-core/v2/pkg/invoker"
)

// DATA RACE in realoding
func TestReload(t *testing.T) {
	inv := invoker.NewInvoker()
	var initial = config.Initial{
		Viper: config.NewViper(),
	}
	initial.Viper.Set("logger.level", "info")

	configurator, _, err := config.Provider(initial, inv)

	cfg, _, err := ProviderCfg(configurator)
	if err != nil {
		t.FailNow()
	}

	zap := NewZap(context.Background(), cfg)

	zap.Log(LevelInfo, "some log")
	go func() {
		inv.Reload(context.Background())
	}()

	zap.Log(LevelInfo, "some")
	time.Sleep(time.Second * 5)
}

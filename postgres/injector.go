// +build wireinject

package postgres

import (
	"context"
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/ProtocolONE/go-core/provider"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

// Build returns GORM instance with resolved dependencies
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (*gorm.DB, func(), error) {
	panic(wire.Build(WireSet, provider.Set))
}

// BuildTest returns stub/mock instance GORM with resolved dependencies
func BuildTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*gorm.DB, func(), error) {
	panic(wire.Build(WireTestSet))
}

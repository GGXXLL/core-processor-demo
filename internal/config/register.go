package config

import (
	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/config"
	"github.com/DoNewsCode/core/otgorm"
	"github.com/DoNewsCode/core/otkafka"
	"github.com/DoNewsCode/core/srvhttp"
	"github.com/GGXXLL/core-process/handler"
	"github.com/GGXXLL/core-process/internal/process"
)

// Register the global options includes modules, module constructors and global dependencies
func Register() []Option {
	return []Option{
		/* Dependencies */
		Dependencies(
			otkafka.Providers(),
			handler.Provides(),
			otgorm.Providers(),
		),

		/* Module Constructors */
		Constructors(
			config.New,          // config module
			core.NewServeModule, // server module
			process.NewProcess,
			otgorm.New,
		),

		/* Modules */
		Modules(
			srvhttp.HealthCheckModule{}, // health check module (http demo)
		),
	}
}

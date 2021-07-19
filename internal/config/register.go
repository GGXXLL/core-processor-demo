package config

import (
	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/config"
	"github.com/DoNewsCode/core/otgorm"
	"github.com/DoNewsCode/core/otkafka"
	"github.com/DoNewsCode/core/otkafka/processor"
	"github.com/DoNewsCode/core/srvhttp"
	"github.com/GGXXLL/core-processor-demo/handler"
)

// Register the global options includes modules, module constructors and global dependencies
func Register() []Option {
	return []Option{
		/* Dependencies */
		Dependencies(
			otkafka.Providers(),
			otgorm.Providers(),
			handler.Provides(),
		),

		/* Module Constructors */
		Constructors(
			config.New,          // config module
			core.NewServeModule, // server module
			processor.New,
		),

		/* Modules */
		Modules(
			srvhttp.HealthCheckModule{}, // health check module (http demo)
		),
	}
}

module worker

go 1.23.0

toolchain go1.24.6

require (
	github.com/Darckfast/workers-go v0.1.1-0.20250902115741-aca215e9a4be
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mailru/easyjson v0.9.0
	go.opentelemetry.io/contrib/bridges/otelslog v0.13.0
)

require (
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/log v0.14.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
)

replace github.com/Darckfast/workers-go => ../workers-go

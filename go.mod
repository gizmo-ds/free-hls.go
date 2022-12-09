module free-hls.go

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/grafov/m3u8 v0.11.1
	github.com/labstack/echo/v4 v4.9.0
	github.com/rs/xid v1.2.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	service-provider v1.0.0
)

replace service-provider => ../service-provider

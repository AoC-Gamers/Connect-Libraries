module github.com/AoC-Gamers/connect-libraries/apikey

go 1.23.0

toolchain go1.24.3

require (
	github.com/AoC-Gamers/connect-libraries/errors v1.0.0
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/ajg/form v1.5.1 // indirect
	github.com/go-chi/render v1.0.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.35.0 // indirect
)

replace github.com/AoC-Gamers/connect-libraries/errors => ../errors

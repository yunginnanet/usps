module github.com/yunginnanet/usps

go 1.22.3

require (
	git.tcp.direct/kayos/common v0.9.7
	github.com/davecgh/go-spew v1.1.1
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
)

retract (
	v0.0.0-20240603082808-20232c5ffb07
	v0.0.1
)

require nullprogram.com/x/rng v1.1.0 // indirect

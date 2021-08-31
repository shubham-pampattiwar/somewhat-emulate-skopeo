module somewhat-emulate-skopeo/main

go 1.17

replace github.com/coreos/bbolt v1.3.6 => go.etcd.io/bbolt v1.3.6

require (
	github.com/containers/image v3.0.2+incompatible
	github.com/containers/storage v1.35.0
)


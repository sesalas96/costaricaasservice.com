module github.com/devsebas/saascr/libs/cri-lib-interop-client

go 1.24

require github.com/devsebas/saascr/libs/cri-lib-shared v0.0.0

require github.com/oklog/ulid/v2 v2.1.0 // indirect

replace (
	github.com/devsebas/saascr/libs/cri-lib-crypto => ../cri-lib-crypto
	github.com/devsebas/saascr/libs/cri-lib-shared => ../cri-lib-shared
)

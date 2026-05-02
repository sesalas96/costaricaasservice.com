module github.com/devsebas/costaricaasservice/libs/cri-lib-http

go 1.24

require (
	github.com/devsebas/costaricaasservice/libs/cri-lib-shared v0.0.0
	github.com/go-chi/chi/v5 v5.1.0
)

require github.com/oklog/ulid/v2 v2.1.0 // indirect

replace github.com/devsebas/costaricaasservice/libs/cri-lib-shared => ../cri-lib-shared

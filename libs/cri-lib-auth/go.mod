module github.com/devsebas/costaricaasservice/libs/cri-lib-auth

go 1.24

require (
	github.com/devsebas/costaricaasservice/libs/cri-lib-shared v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.1
)

replace github.com/devsebas/costaricaasservice/libs/cri-lib-shared => ../cri-lib-shared

// Package auth define el modelo de roles y el Principal que se propaga
// por toda la cadena de servicios. Convención: roles en MAYÚSCULAS_CON_GUION_BAJO.
package auth

// Role identifica un rol asignable. Estable; no renombrar sin ADR.
type Role string

const (
	// Citizen — ciudadano natural del realm.
	RoleCitizen Role = "CITIZEN"

	// Member operator — operador interno de una institución.
	RoleMemberOperator Role = "MEMBER_OPERATOR"
	// Member admin — administrador de una institución.
	RoleMemberAdmin Role = "MEMBER_ADMIN"

	// Realm admin — administra un realm completo (jurisdicción).
	RoleRealmAdmin Role = "REALM_ADMIN"

	// Costaricaasservice admin — control plane del SaaS, opera todos los realms.
	RoleCostaricaasserviceAdmin Role = "COSTARICAASSERVICE_ADMIN"
)

// Has verifica si la lista contiene el rol pedido.
func Has(roles []Role, want Role) bool {
	for _, r := range roles {
		if r == want {
			return true
		}
	}
	return false
}

// HasAny verifica si la lista contiene al menos uno de los roles pedidos.
func HasAny(roles []Role, wants ...Role) bool {
	for _, w := range wants {
		if Has(roles, w) {
			return true
		}
	}
	return false
}

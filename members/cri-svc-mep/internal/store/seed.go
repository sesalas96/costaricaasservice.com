package store

func SeedDemo(s *Store) {
	enrollments := []Enrollment{
		{ID: "enr-001", StudentCedula: "1-1234-5678", Year: 2026, Grade: "9no_secundaria", SchoolName: "Liceo de Costa Rica", SchoolKind: "publico", Status: "active"},
		{ID: "enr-002", StudentCedula: "1-9876-5432", Year: 2026, Grade: "5to_primaria", SchoolName: "Escuela República de Argentina", SchoolKind: "publico", Status: "active"},
	}
	certificates := []Certificate{
		{ID: "cert-001", StudentCedula: "1-1234-5678", Kind: "primaria_aprobada", IssuedAt: "2022-12-15", IssuedBy: "Escuela Mata Redonda"},
	}
	s.Seed("demo", enrollments, certificates)
}

package store

func SeedDemo(s *Store) {
	statuses := []ImmigrationStatus{
		{Cedula: "1-1234-5678", Nationality: "CRI", Category: "ciudadano", DocumentType: "cedula_nacional", IssuedAt: "2018-05-04", ExpiresAt: "2028-05-04", Status: "active"},
		{Cedula: "1-9876-5432", Nationality: "CRI", Category: "ciudadano", DocumentType: "cedula_nacional", IssuedAt: "2020-11-12", ExpiresAt: "2030-11-12", Status: "active"},
		{Cedula: "1869900012345", Nationality: "NIC", Category: "residente_permanente", DocumentType: "dimex", IssuedAt: "2023-02-20", ExpiresAt: "2028-02-19", Status: "active"},
	}
	movements := []Movement{
		{ID: "mov-001", Cedula: "1-1234-5678", BorderPost: "SJO", Direction: "exit", When: "2026-01-12T08:30:00-06:00"},
		{ID: "mov-002", Cedula: "1-1234-5678", BorderPost: "SJO", Direction: "entry", When: "2026-01-26T19:15:00-06:00"},
		{ID: "mov-003", Cedula: "1869900012345", BorderPost: "PENAS_BLANCAS", Direction: "entry", When: "2023-02-15T11:00:00-06:00"},
	}
	s.Seed("demo", statuses, movements)
}

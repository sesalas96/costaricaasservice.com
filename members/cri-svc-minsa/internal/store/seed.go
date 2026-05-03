package store

func SeedDemo(s *Store) {
	permits := []HealthPermit{
		{Number: "PSF-2026-1042", HolderCedula: "1-1234-5678", Kind: "farmacia", TradeName: "Farmacia La Bendición", Address: "San José, Mata Redonda", IssuedAt: "2026-01-15", ExpiresAt: "2031-01-14", Status: "active"},
		{Number: "PROF-2024-0917", HolderCedula: "1-9876-5432", Kind: "profesional_salud", TradeName: "Odontóloga colegiada CCDCR", IssuedAt: "2024-08-01", ExpiresAt: "2029-07-31", Status: "active"},
	}
	s.Seed("demo", permits)
}

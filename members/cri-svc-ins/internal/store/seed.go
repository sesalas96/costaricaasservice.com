package store

func SeedDemo(s *Store) {
	policies := []Policy{
		{Number: "SOA-2026-0001", HolderCedula: "1-1234-5678", Kind: "soa", VehiclePlate: "BLZ-145", PremiumColons: 65_000, ValidFrom: "2026-01-01", ValidUntil: "2026-12-31", Status: "active"},
		{Number: "RT-2026-0042", HolderCedula: "1-1234-5678", Kind: "riesgos_trabajo", PremiumColons: 480_000, ValidFrom: "2026-03-01", ValidUntil: "2027-02-28", Status: "active"},
		{Number: "VS-2026-0099", HolderCedula: "1-9876-5432", Kind: "voluntario_salud", PremiumColons: 1_200_000, ValidFrom: "2026-02-15", ValidUntil: "2027-02-14", Status: "active"},
	}
	s.Seed("demo", policies)
}

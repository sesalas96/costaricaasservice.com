package store

func SeedDemo(s *Store) {
	properties := []RealEstate{
		{Folio: "1-456789-000", OwnerCedula: "1-1234-5678", Province: "San José", Canton: "San José", District: "Mata Redonda", AreaSquareMts: 215.5, FiscalValue: 95_000_000},
		{Folio: "7-098765-000", OwnerCedula: "1-9876-5432", Province: "Heredia", Canton: "Belén", District: "San Antonio", AreaSquareMts: 380.0, FiscalValue: 180_000_000},
	}
	vehicles := []Vehicle{
		{Plate: "BLZ-145", OwnerCedula: "1-1234-5678", Make: "Toyota", Model: "Hilux", Year: 2022, FiscalValue: 22_000_000, Status: "active"},
		{Plate: "BMW-009", OwnerCedula: "1-9876-5432", Make: "BMW", Model: "X1", Year: 2024, FiscalValue: 38_500_000, Status: "active"},
	}
	s.Seed("demo", properties, vehicles)
}

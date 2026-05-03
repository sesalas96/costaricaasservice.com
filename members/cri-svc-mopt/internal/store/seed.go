package store

func SeedDemo(s *Store) {
	licenses := []DriverLicense{
		{Number: "1-1234-5678-B1", Cedula: "1-1234-5678", Categories: []string{"B1"}, IssuedAt: "2018-06-10", ExpiresAt: "2028-06-10", Status: "active", Points: 50},
		{Number: "1-9876-5432-A3", Cedula: "1-9876-5432", Categories: []string{"A3", "B1"}, IssuedAt: "2020-09-22", ExpiresAt: "2030-09-22", Status: "active", Points: 35},
	}
	fines := []TrafficFine{
		{ID: "fine-001", Cedula: "1-1234-5678", VehiclePlate: "BLZ-145", Reason: "Exceso de velocidad (120 km/h en zona 80)", AmountColons: 280_000, IssuedAt: "2026-03-12", Status: "pending"},
		{ID: "fine-002", Cedula: "1-9876-5432", VehiclePlate: "BMW-009", Reason: "Estacionamiento en zona prohibida", AmountColons: 45_000, IssuedAt: "2026-04-02", Status: "paid"},
	}
	s.Seed("demo", licenses, fines)
}

package store

func SeedDemo(s *Store) {
	prescriptions := []Prescription{
		{ID: "rx-001", Cedula: "1-1234-5678", Drug: "Losartán 50mg", Dosage: "1 cada 12 horas", IssuedAt: "2026-04-12", IssuedBy: "Dr. Solano (HSJD)", ValidUntil: "2026-10-12", Status: "active"},
		{ID: "rx-002", Cedula: "1-1234-5678", Drug: "Atorvastatina 20mg", Dosage: "1 al día", IssuedAt: "2026-04-12", IssuedBy: "Dr. Solano (HSJD)", ValidUntil: "2026-10-12", Status: "active"},
		{ID: "rx-003", Cedula: "1-9876-5432", Drug: "Metformina 850mg", Dosage: "1 con cada comida", IssuedAt: "2026-03-20", IssuedBy: "Dra. Vargas (HCG)", ValidUntil: "2026-09-20", Status: "active"},
	}
	appointments := []Appointment{
		{ID: "apt-001", Cedula: "1-1234-5678", Specialty: "Cardiología", Hospital: "HSJD San José", When: "2026-06-15 09:30", Status: "scheduled"},
		{ID: "apt-002", Cedula: "1-1234-5678", Specialty: "EBAIS Mata Redonda", Hospital: "EBAIS", When: "2026-05-20 14:00", Status: "scheduled"},
		{ID: "apt-003", Cedula: "1-9876-5432", Specialty: "Endocrinología", Hospital: "HCG San José", When: "2026-07-02 10:15", Status: "scheduled"},
	}
	s.Seed("demo", prescriptions, appointments)
}

package store

func SeedDemo(s *Store) {
	employments := []Employment{
		{ID: "emp-001", WorkerCedula: "1-1234-5678", EmployerName: "Cooperativa Dos Pinos R.L.", EmployerCID: "3-004-045002", Position: "Ingeniero de procesos", Contract: "indefinido", MonthlyColons: 1_850_000, StartedAt: "2022-04-01", Status: "active"},
		{ID: "emp-002", WorkerCedula: "1-9876-5432", EmployerName: "Banco Nacional de Costa Rica", EmployerCID: "4-000-001021", Position: "Ejecutiva de cuenta", Contract: "indefinido", MonthlyColons: 1_400_000, StartedAt: "2019-08-15", Status: "active"},
		{ID: "emp-003", WorkerCedula: "1-9876-5432", EmployerName: "Universidad de Costa Rica", EmployerCID: "4-000-042149", Position: "Docente media jornada", Contract: "plazo_fijo", MonthlyColons: 720_000, StartedAt: "2024-02-01", EndedAt: "2025-12-15", Status: "ended"},
	}
	s.Seed("demo", employments)
}

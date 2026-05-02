package store

// SeedDemo inserta data mock para el realm `demo`. Llamado al arrancar el
// binario en local.
func SeedDemo(s *Store) {
	persons := []Person{
		{Cedula: "1-1234-5678", FullName: "Sebastián Rojas Mora", BirthDate: "1992-04-12", Address: "San José, Mata Redonda", Email: "srojas@example.cr", Status: "alive"},
		{Cedula: "1-9876-5432", FullName: "María Castro Fernández", BirthDate: "1988-09-03", Address: "Heredia, San Rafael", Email: "mcastro@example.cr", Status: "alive"},
		{Cedula: "2-0001-0001", FullName: "Luis Fernández Solís", BirthDate: "1975-12-22", Address: "Alajuela, Centro", Email: "lfernandez@example.cr", Status: "alive"},
	}
	events := []VitalEvent{
		{ID: "ev-001", Cedula: "1-1234-5678", Type: "birth", Date: "1992-04-12", Notes: "Nacido en Hospital México"},
		{ID: "ev-002", Cedula: "1-1234-5678", Type: "marriage", Date: "2018-06-30", Notes: "Matrimonio civil"},
		{ID: "ev-003", Cedula: "1-9876-5432", Type: "birth", Date: "1988-09-03", Notes: "Nacida en Hospital Calderón Guardia"},
	}
	s.Seed("demo", persons, events)
}

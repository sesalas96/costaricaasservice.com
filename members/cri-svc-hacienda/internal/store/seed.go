package store

func SeedDemo(s *Store) {
	s.Seed("demo", []TaxRecord{
		{Cedula: "1-1234-5678", Year: 2025, GrossIncome: 18_000_000, WithheldTax: 2_700_000, Deductions: 350_000, HasDependents: false},
		{Cedula: "1-9876-5432", Year: 2025, GrossIncome: 24_500_000, WithheldTax: 3_900_000, Deductions: 1_200_000, HasDependents: true},
		{Cedula: "2-0001-0001", Year: 2025, GrossIncome: 12_000_000, WithheldTax: 1_440_000, Deductions: 0, HasDependents: false},
	})
}

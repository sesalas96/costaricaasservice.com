import type { Lang } from "../i18n/translations";

export type InstitutionStatus = "live" | "planned";

export type InstitutionId =
  | "registro-civil"
  | "hacienda"
  | "ccss"
  | "ins"
  | "mep"
  | "mopt"
  | "mtss"
  | "ice"
  | "imas"
  | "migracion";

export type InstitutionLocale = {
  name: string;
  short: string;
  custodia: string[];
  expone: string[];
  consume: string[];
};

export type InstitutionDef = {
  id: InstitutionId;
  tag: string;
  status: InstitutionStatus;
  accent: string;
  es: InstitutionLocale;
  en: InstitutionLocale;
};

export const INSTITUTIONS: InstitutionDef[] = [
  {
    id: "registro-civil",
    tag: "tse",
    status: "live",
    accent: "var(--color-cr-blue-bright)",
    es: {
      name: "Registro Civil",
      short: "TSE · fuente de oro de identidad",
      custodia: [
        "Personas: cédula, nombre, fecha de nacimiento",
        "Eventos vitales: nacimientos, defunciones, matrimonios",
        "Estado civil y parentescos",
        "Domicilio electoral",
      ],
      expone: [
        "GET /persons/{cedula}",
        "GET /persons/{cedula}/vital-events",
        "GET /persons/{cedula}/household",
      ],
      consume: [],
    },
    en: {
      name: "Civil Registry",
      short: "TSE · golden source for identity",
      custodia: [
        "Persons: ID, full name, date of birth",
        "Vital events: births, deaths, marriages",
        "Civil status and family relationships",
        "Electoral domicile",
      ],
      expone: [
        "GET /persons/{id}",
        "GET /persons/{id}/vital-events",
        "GET /persons/{id}/household",
      ],
      consume: [],
    },
  },
  {
    id: "hacienda",
    tag: "hac",
    status: "live",
    accent: "var(--color-area-members)",
    es: {
      name: "Ministerio de Hacienda",
      short: "Tributación, declaraciones y compras públicas",
      custodia: [
        "Estado tributario por contribuyente",
        "Declaraciones presentadas e historial fiscal",
        "Padrón de empresas y representantes",
      ],
      expone: [
        "GET /tax-status/{cedula}",
        "GET /prefilled-return/{cedula}/{year}",
        "POST /returns/{cedula}/submit",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula} (composición familiar)",
        "ccss · GET /contributions/{cedula} (cargas sociales) — planificado",
      ],
    },
    en: {
      name: "Tax Authority",
      short: "Taxation, returns and public procurement",
      custodia: [
        "Tax status per taxpayer",
        "Filed returns and fiscal history",
        "Roster of companies and legal representatives",
      ],
      expone: [
        "GET /tax-status/{id}",
        "GET /prefilled-return/{id}/{year}",
        "POST /returns/{id}/submit",
      ],
      consume: [
        "registro-civil · GET /persons/{id} (household composition)",
        "ccss · GET /contributions/{id} (social charges) — planned",
      ],
    },
  },
  {
    id: "ccss",
    tag: "ccss",
    status: "planned",
    accent: "var(--color-cr-red-bright)",
    es: {
      name: "Caja Costarricense del Seguro Social",
      short: "Salud pública, recetas y aseguramiento",
      custodia: [
        "Aseguramiento y régimen contributivo",
        "Historial clínico y citas médicas",
        "Recetas electrónicas firmadas",
        "Cargas sociales por patrono",
      ],
      expone: [
        "GET /insurance-status/{cedula}",
        "GET /clinical-history/{cedula}",
        "POST /prescriptions (firmada por médico)",
        "GET /contributions/{cedula}",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula} (validar identidad)",
        "iduc-signing · firma de recetas (Ed25519)",
      ],
    },
    en: {
      name: "Costa Rican Social Security Fund",
      short: "Public healthcare, prescriptions and insurance",
      custodia: [
        "Insurance and contribution regime",
        "Clinical history and appointments",
        "Signed electronic prescriptions",
        "Employer social charges",
      ],
      expone: [
        "GET /insurance-status/{id}",
        "GET /clinical-history/{id}",
        "POST /prescriptions (physician-signed)",
        "GET /contributions/{id}",
      ],
      consume: [
        "registro-civil · GET /persons/{id} (identity check)",
        "iduc-signing · prescription signature (Ed25519)",
      ],
    },
  },
  {
    id: "ins",
    tag: "ins",
    status: "planned",
    accent: "#ff7a59",
    es: {
      name: "Instituto Nacional de Seguros",
      short: "Pólizas, riesgos del trabajo y SOA",
      custodia: [
        "Pólizas vigentes (vehículos, vida, salud)",
        "Cobertura de riesgos del trabajo",
        "Histórico de siniestros",
      ],
      expone: [
        "GET /policies/{cedula}",
        "GET /soa-status/{placa}",
        "POST /claims",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "mopt · GET /vehicles/{placa} — planificado",
      ],
    },
    en: {
      name: "National Insurance Institute",
      short: "Policies, work-related risks and mandatory auto insurance",
      custodia: [
        "Active policies (vehicle, life, health)",
        "Work-related risk coverage",
        "Claims history",
      ],
      expone: [
        "GET /policies/{id}",
        "GET /soa-status/{plate}",
        "POST /claims",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "mopt · GET /vehicles/{plate} — planned",
      ],
    },
  },
  {
    id: "mep",
    tag: "mep",
    status: "planned",
    accent: "#34d399",
    es: {
      name: "Ministerio de Educación Pública",
      short: "Matrícula, certificaciones y becas",
      custodia: [
        "Matrícula estudiantil por nivel",
        "Certificaciones académicas (notas, conducta)",
        "Becas y subsidios educativos",
      ],
      expone: [
        "GET /students/{cedula}",
        "GET /certifications/{cedula}",
        "POST /enrollments",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "imas · GET /eligibility/{cedula} (becas) — planificado",
      ],
    },
    en: {
      name: "Ministry of Public Education",
      short: "Enrollment, academic records and scholarships",
      custodia: [
        "Student enrollment by level",
        "Academic certifications (grades, conduct)",
        "Scholarships and educational aid",
      ],
      expone: [
        "GET /students/{id}",
        "GET /certifications/{id}",
        "POST /enrollments",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "imas · GET /eligibility/{id} (scholarships) — planned",
      ],
    },
  },
  {
    id: "mopt",
    tag: "mopt",
    status: "planned",
    accent: "#fbbf24",
    es: {
      name: "Ministerio de Obras Públicas y Transportes",
      short: "Vehículos, licencias y transporte público",
      custodia: [
        "Padrón de vehículos por placa",
        "Licencias de conducir",
        "Concesiones de transporte público",
      ],
      expone: [
        "GET /vehicles/{placa}",
        "GET /drivers-license/{cedula}",
        "POST /vehicle-transfers",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "ins · GET /soa-status/{placa}",
      ],
    },
    en: {
      name: "Ministry of Public Works and Transport",
      short: "Vehicles, driver licenses and public transport",
      custodia: [
        "Vehicle registry by plate",
        "Driver licenses",
        "Public transport concessions",
      ],
      expone: [
        "GET /vehicles/{plate}",
        "GET /drivers-license/{id}",
        "POST /vehicle-transfers",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "ins · GET /soa-status/{plate}",
      ],
    },
  },
  {
    id: "mtss",
    tag: "mtss",
    status: "planned",
    accent: "#22d3ee",
    es: {
      name: "Ministerio de Trabajo y Seguridad Social",
      short: "Empleo, planillas y resoluciones laborales",
      custodia: [
        "Planillas patronales validadas",
        "Resoluciones de denuncias laborales",
        "Programas de empleo activos",
      ],
      expone: [
        "GET /employment-status/{cedula}",
        "GET /payroll/{cedula-patrono}",
        "POST /labor-complaints",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "ccss · GET /contributions/{cedula}",
      ],
    },
    en: {
      name: "Ministry of Labor and Social Security",
      short: "Employment, payroll and labor rulings",
      custodia: [
        "Validated employer payrolls",
        "Rulings on labor complaints",
        "Active employment programs",
      ],
      expone: [
        "GET /employment-status/{id}",
        "GET /payroll/{employer-id}",
        "POST /labor-complaints",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "ccss · GET /contributions/{id}",
      ],
    },
  },
  {
    id: "ice",
    tag: "ice",
    status: "planned",
    accent: "#a78bfa",
    es: {
      name: "Instituto Costarricense de Electricidad",
      short: "Servicios eléctricos y telecomunicaciones",
      custodia: [
        "Conexiones eléctricas por NIS",
        "Consumos históricos",
        "Servicios de telecomunicaciones",
      ],
      expone: [
        "GET /power-account/{nis}",
        "GET /consumption/{nis}/{period}",
      ],
      consume: ["registro-civil · GET /persons/{cedula}"],
    },
    en: {
      name: "Costa Rican Institute of Electricity",
      short: "Electric services and telecommunications",
      custodia: [
        "Power connections per service ID",
        "Historical consumption",
        "Telecommunications services",
      ],
      expone: [
        "GET /power-account/{nis}",
        "GET /consumption/{nis}/{period}",
      ],
      consume: ["registro-civil · GET /persons/{id}"],
    },
  },
  {
    id: "imas",
    tag: "imas",
    status: "planned",
    accent: "#f472b6",
    es: {
      name: "Instituto Mixto de Ayuda Social",
      short: "Subsidios, bonos y programas de bienestar",
      custodia: [
        "Padrón de beneficiarios",
        "Subsidios entregados (Avancemos, bono familiar)",
        "Calificación socioeconómica",
      ],
      expone: [
        "GET /eligibility/{cedula}",
        "GET /benefits/{cedula}",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula} (familia)",
        "hacienda · GET /tax-status/{cedula} (ingresos)",
        "ccss · GET /insurance-status/{cedula}",
      ],
    },
    en: {
      name: "Joint Institute of Social Aid",
      short: "Subsidies, allowances and welfare programs",
      custodia: [
        "Beneficiary roster",
        "Disbursed subsidies (school aid, family allowance)",
        "Socio-economic scoring",
      ],
      expone: [
        "GET /eligibility/{id}",
        "GET /benefits/{id}",
      ],
      consume: [
        "registro-civil · GET /persons/{id} (household)",
        "hacienda · GET /tax-status/{id} (income)",
        "ccss · GET /insurance-status/{id}",
      ],
    },
  },
  {
    id: "migracion",
    tag: "mig",
    status: "planned",
    accent: "#94a3b8",
    es: {
      name: "Dirección General de Migración y Extranjería",
      short: "Pasaportes, permisos y movimientos migratorios",
      custodia: [
        "Pasaportes emitidos",
        "Permisos de residencia",
        "Movimientos migratorios (entradas/salidas)",
      ],
      expone: [
        "GET /passport/{cedula}",
        "GET /immigration-status/{cedula}",
        "POST /passport-applications",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "hacienda · GET /tax-status/{cedula} (cobro arancel)",
      ],
    },
    en: {
      name: "Directorate of Migration and Foreign Affairs",
      short: "Passports, permits and migratory movements",
      custodia: [
        "Issued passports",
        "Residence permits",
        "Migratory movements (entries/exits)",
      ],
      expone: [
        "GET /passport/{id}",
        "GET /immigration-status/{id}",
        "POST /passport-applications",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "hacienda · GET /tax-status/{id} (fee collection)",
      ],
    },
  },
];

export const getInstitution = (id: InstitutionId): InstitutionDef | undefined =>
  INSTITUTIONS.find((i) => i.id === id);

export const localized = (def: InstitutionDef, lang: Lang): InstitutionLocale =>
  lang === "es" ? def.es : def.en;

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
  | "migracion"
  | "bccr"
  | "registro-nacional"
  | "poder-judicial"
  | "muni-sjo"
  | "aya"
  | "micitt"
  | "minae"
  | "minsa"
  | "pani"
  | "sugef"
  | "fuerza-publica"
  | "bomberos"
  | "procomer";

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
        "ccss · GET /contributions/{cedula} (cargas sociales)",
        "bccr · POST /sinpe/charge (cobro de aranceles)",
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
        "ccss · GET /contributions/{id} (social charges)",
        "bccr · POST /sinpe/charge (fee collection)",
      ],
    },
  },
  {
    id: "bccr",
    tag: "bcr",
    status: "live",
    accent: "#10b981",
    es: {
      name: "Banco Central · SINPE",
      short: "Rails financieros del Estado · pagos en tiempo real",
      custodia: [
        "Cuentas SINPE Móvil por cédula",
        "Identificadores de cuenta cliente (IBAN)",
        "Tipos de cambio oficiales",
        "Padrón financiero (KYC base)",
      ],
      expone: [
        "POST /sinpe/charge (cobro instantáneo)",
        "POST /sinpe/transfer (acreditación)",
        "GET /accounts/{cedula}",
        "GET /fx/{date}",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula} (KYC)",
        "sugef · GET /aml-flags/{cedula} (pre-cobro)",
      ],
    },
    en: {
      name: "Central Bank · SINPE",
      short: "State financial rails · real-time payments",
      custodia: [
        "SINPE Mobile accounts per ID",
        "Customer account identifiers (IBAN)",
        "Official FX rates",
        "Base KYC roster",
      ],
      expone: [
        "POST /sinpe/charge (instant collection)",
        "POST /sinpe/transfer (credit)",
        "GET /accounts/{id}",
        "GET /fx/{date}",
      ],
      consume: [
        "registro-civil · GET /persons/{id} (KYC)",
        "sugef · GET /aml-flags/{id} (pre-charge)",
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
        "minsa · GET /vaccination-record/{cedula}",
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
        "minsa · GET /vaccination-record/{id}",
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
        "mopt · GET /vehicles/{placa}",
        "bomberos · GET /inspections/{placa-comercio}",
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
        "mopt · GET /vehicles/{plate}",
        "bomberos · GET /inspections/{business-id}",
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
        "imas · GET /eligibility/{cedula} (becas)",
        "minsa · GET /vaccination-record/{cedula} (matrícula)",
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
        "imas · GET /eligibility/{id} (scholarships)",
        "minsa · GET /vaccination-record/{id} (enrollment)",
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
        "registro-nacional · GET /property/vehicle/{placa}",
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
        "registro-nacional · GET /property/vehicle/{plate}",
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
        "registro-nacional · GET /entity/{cedula-juridica}",
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
        "registro-nacional · GET /entity/{legal-entity-id}",
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
        "POST /service-requests",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "muni-sjo · GET /property/{cedula}",
        "minae · GET /environmental-clearance/{folio}",
      ],
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
        "POST /service-requests",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "muni-sjo · GET /property/{id}",
        "minae · GET /environmental-clearance/{folio}",
      ],
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
        "pani · GET /minors/{cedula} (menores en hogar)",
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
        "pani · GET /minors/{id} (minors in household)",
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
        "poder-judicial · GET /criminal-record/{cedula}",
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
        "poder-judicial · GET /criminal-record/{id}",
      ],
    },
  },
  {
    id: "registro-nacional",
    tag: "rn",
    status: "planned",
    accent: "#6366f1",
    es: {
      name: "Registro Nacional",
      short: "Propiedades, mercantil, prendas y vehículos titulares",
      custodia: [
        "Propiedades inmuebles (folio real)",
        "Padrón mercantil de personas jurídicas",
        "Prendas y gravámenes",
        "Titularidad vehicular (parte registral, distinto a MOPT)",
        "Marcas y patentes",
      ],
      expone: [
        "GET /property/{folio-real}",
        "GET /entity/{cedula-juridica}",
        "GET /property/vehicle/{placa}",
        "POST /property/transfer (escritura firmada)",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "iduc-signing · firma de escrituras (Ed25519)",
        "hacienda · POST /fees/transfer-tax (impuesto traspaso)",
      ],
    },
    en: {
      name: "National Registry",
      short: "Property, commercial entities, liens and titled vehicles",
      custodia: [
        "Real estate (folio real)",
        "Commercial roster of legal entities",
        "Liens and encumbrances",
        "Vehicle titling (registral, distinct from MOPT)",
        "Trademarks and patents",
      ],
      expone: [
        "GET /property/{folio-real}",
        "GET /entity/{legal-entity-id}",
        "GET /property/vehicle/{plate}",
        "POST /property/transfer (signed deed)",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "iduc-signing · deed signature (Ed25519)",
        "hacienda · POST /fees/transfer-tax",
      ],
    },
  },
  {
    id: "poder-judicial",
    tag: "pju",
    status: "planned",
    accent: "#facc15",
    es: {
      name: "Poder Judicial",
      short: "Antecedentes, expedientes y resoluciones",
      custodia: [
        "Hoja de delincuencia (antecedentes penales)",
        "Expedientes judiciales activos",
        "Sentencias firmes",
        "Pensiones alimentarias",
      ],
      expone: [
        "GET /criminal-record/{cedula}",
        "GET /cases/{cedula}",
        "GET /alimony/{cedula}",
        "POST /complaints",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "fuerza-publica · GET /reports/{cedula}",
        "iduc-signing · firma de demandas",
      ],
    },
    en: {
      name: "Judicial Branch",
      short: "Criminal records, case files and rulings",
      custodia: [
        "Criminal background record",
        "Active judicial case files",
        "Final rulings",
        "Alimony orders",
      ],
      expone: [
        "GET /criminal-record/{id}",
        "GET /cases/{id}",
        "GET /alimony/{id}",
        "POST /complaints",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "fuerza-publica · GET /reports/{id}",
        "iduc-signing · complaint signature",
      ],
    },
  },
  {
    id: "muni-sjo",
    tag: "sjo",
    status: "planned",
    accent: "#f97316",
    es: {
      name: "Municipalidad de San José",
      short: "Bienes inmuebles, patentes y permisos territoriales",
      custodia: [
        "Catastro municipal (bienes inmuebles)",
        "Patentes comerciales",
        "Permisos de construcción",
        "Servicios urbanos (basura, alcantarillado pluvial)",
      ],
      expone: [
        "GET /property/{cedula}",
        "GET /patents/{cedula}",
        "POST /construction-permits",
        "GET /municipal-bills/{cedula}",
      ],
      consume: [
        "registro-nacional · GET /property/{folio-real}",
        "bomberos · GET /inspections/{folio-real}",
        "minae · GET /environmental-clearance/{folio}",
        "bccr · POST /sinpe/charge (servicios)",
      ],
    },
    en: {
      name: "San José Municipality",
      short: "Real estate tax, business permits and territorial filings",
      custodia: [
        "Municipal cadaster (real estate)",
        "Business patents",
        "Construction permits",
        "Urban services (waste, storm drainage)",
      ],
      expone: [
        "GET /property/{id}",
        "GET /patents/{id}",
        "POST /construction-permits",
        "GET /municipal-bills/{id}",
      ],
      consume: [
        "registro-nacional · GET /property/{folio-real}",
        "bomberos · GET /inspections/{folio-real}",
        "minae · GET /environmental-clearance/{folio}",
        "bccr · POST /sinpe/charge (services)",
      ],
    },
  },
  {
    id: "aya",
    tag: "aya",
    status: "planned",
    accent: "#06b6d4",
    es: {
      name: "Acueductos y Alcantarillados",
      short: "Agua potable y saneamiento",
      custodia: [
        "Conexiones de agua por NIS",
        "Consumos por periodo",
        "Calidad y disponibilidad por zona",
        "Cortes programados",
      ],
      expone: [
        "GET /water-account/{nis}",
        "GET /consumption/{nis}/{period}",
        "POST /service-requests",
        "GET /quality/{cantón-distrito}",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "muni-sjo · GET /property/{cedula}",
        "minae · GET /water-concessions/{folio}",
      ],
    },
    en: {
      name: "Water & Sewerage Authority",
      short: "Drinking water and sanitation",
      custodia: [
        "Water connections per service ID",
        "Consumption per period",
        "Quality and availability by zone",
        "Scheduled outages",
      ],
      expone: [
        "GET /water-account/{nis}",
        "GET /consumption/{nis}/{period}",
        "POST /service-requests",
        "GET /quality/{canton-district}",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "muni-sjo · GET /property/{id}",
        "minae · GET /water-concessions/{folio}",
      ],
    },
  },
  {
    id: "micitt",
    tag: "mic",
    status: "planned",
    accent: "#8b5cf6",
    es: {
      name: "Ministerio de Ciencia, Tecnología y Telecomunicaciones",
      short: "Rector digital · gobernanza tecnológica",
      custodia: [
        "Estrategia digital nacional",
        "Registro de proveedores tecnológicos del Estado",
        "Política de datos abiertos",
        "Espectro radioeléctrico (vía SUTEL)",
      ],
      expone: [
        "GET /digital-strategy/{realm}",
        "GET /open-data-catalog",
        "POST /tech-procurement",
      ],
      consume: [
        "procomer · GET /tech-exporters",
        "sugef · GET /fintech-registry",
      ],
    },
    en: {
      name: "Ministry of Science, Technology and Telecom",
      short: "Digital steward · technology governance",
      custodia: [
        "National digital strategy",
        "Registry of state tech vendors",
        "Open data policy",
        "Radio spectrum (via SUTEL)",
      ],
      expone: [
        "GET /digital-strategy/{realm}",
        "GET /open-data-catalog",
        "POST /tech-procurement",
      ],
      consume: [
        "procomer · GET /tech-exporters",
        "sugef · GET /fintech-registry",
      ],
    },
  },
  {
    id: "minae",
    tag: "mie",
    status: "planned",
    accent: "#84cc16",
    es: {
      name: "Ministerio de Ambiente y Energía",
      short: "Carbono neutralidad · concesiones · áreas protegidas",
      custodia: [
        "Áreas silvestres protegidas",
        "Concesiones de agua y energía",
        "Inventario nacional de carbono",
        "Permisos de impacto ambiental",
      ],
      expone: [
        "GET /environmental-clearance/{folio}",
        "GET /water-concessions/{folio}",
        "GET /carbon-credits/{cedula-juridica}",
        "POST /environmental-complaints",
      ],
      consume: [
        "registro-nacional · GET /property/{folio-real}",
        "muni-sjo · GET /zoning/{distrito}",
      ],
    },
    en: {
      name: "Ministry of Environment and Energy",
      short: "Carbon neutrality · concessions · protected areas",
      custodia: [
        "Protected wildlands",
        "Water and energy concessions",
        "National carbon inventory",
        "Environmental impact permits",
      ],
      expone: [
        "GET /environmental-clearance/{folio}",
        "GET /water-concessions/{folio}",
        "GET /carbon-credits/{legal-entity-id}",
        "POST /environmental-complaints",
      ],
      consume: [
        "registro-nacional · GET /property/{folio-real}",
        "muni-sjo · GET /zoning/{district}",
      ],
    },
  },
  {
    id: "minsa",
    tag: "msa",
    status: "planned",
    accent: "#ec4899",
    es: {
      name: "Ministerio de Salud",
      short: "Salud pública · vigilancia epidemiológica · permisos sanitarios",
      custodia: [
        "Carné nacional de vacunación",
        "Vigilancia epidemiológica",
        "Permisos sanitarios de funcionamiento",
        "Registro de medicamentos y dispositivos",
      ],
      expone: [
        "GET /vaccination-record/{cedula}",
        "GET /sanitary-permit/{cedula-juridica}",
        "GET /epidemiological-alerts/{distrito}",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "ccss · GET /clinical-history/{cedula}",
        "muni-sjo · GET /property/{cedula} (inspección)",
      ],
    },
    en: {
      name: "Ministry of Health",
      short: "Public health · epidemiology · sanitary permits",
      custodia: [
        "National vaccination card",
        "Epidemiological surveillance",
        "Sanitary operating permits",
        "Drug and device registry",
      ],
      expone: [
        "GET /vaccination-record/{id}",
        "GET /sanitary-permit/{legal-entity-id}",
        "GET /epidemiological-alerts/{district}",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "ccss · GET /clinical-history/{id}",
        "muni-sjo · GET /property/{id} (inspection)",
      ],
    },
  },
  {
    id: "pani",
    tag: "pni",
    status: "planned",
    accent: "#fb923c",
    es: {
      name: "Patronato Nacional de la Infancia",
      short: "Niñez y adolescencia · protección integral",
      custodia: [
        "Registro de menores con expediente",
        "Medidas de protección activas",
        "Adopciones y procesos especiales",
      ],
      expone: [
        "GET /minors/{cedula}",
        "GET /protection-measures/{cedula}",
        "POST /child-welfare-reports",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}/household",
        "poder-judicial · GET /alimony/{cedula}",
        "ccss · GET /insurance-status/{cedula}",
      ],
    },
    en: {
      name: "Children & Adolescents Welfare",
      short: "Children & youth · integral protection",
      custodia: [
        "Registry of minors with files",
        "Active protection measures",
        "Adoptions and special procedures",
      ],
      expone: [
        "GET /minors/{id}",
        "GET /protection-measures/{id}",
        "POST /child-welfare-reports",
      ],
      consume: [
        "registro-civil · GET /persons/{id}/household",
        "poder-judicial · GET /alimony/{id}",
        "ccss · GET /insurance-status/{id}",
      ],
    },
  },
  {
    id: "sugef",
    tag: "sug",
    status: "planned",
    accent: "#14b8a6",
    es: {
      name: "Superintendencia de Entidades Financieras",
      short: "Supervisión financiera · KYC · AML",
      custodia: [
        "Padrón de entidades supervisadas",
        "Centro de información crediticia",
        "Listas de cumplimiento (PEP, sancionados)",
        "Registro de fintechs",
      ],
      expone: [
        "GET /credit-info/{cedula}",
        "GET /aml-flags/{cedula}",
        "GET /supervised-entities",
        "POST /suspicious-activity-reports",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "registro-nacional · GET /entity/{cedula-juridica}",
        "poder-judicial · GET /criminal-record/{cedula}",
      ],
    },
    en: {
      name: "Financial Entities Superintendence",
      short: "Financial oversight · KYC · AML",
      custodia: [
        "Roster of supervised entities",
        "Credit information bureau",
        "Compliance lists (PEP, sanctioned)",
        "Fintech registry",
      ],
      expone: [
        "GET /credit-info/{id}",
        "GET /aml-flags/{id}",
        "GET /supervised-entities",
        "POST /suspicious-activity-reports",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "registro-nacional · GET /entity/{legal-entity-id}",
        "poder-judicial · GET /criminal-record/{id}",
      ],
    },
  },
  {
    id: "fuerza-publica",
    tag: "msp",
    status: "planned",
    accent: "#475569",
    es: {
      name: "Fuerza Pública · Ministerio de Seguridad",
      short: "Denuncias ciudadanas · operativos · seguridad pública",
      custodia: [
        "Denuncias ciudadanas",
        "Reportes de incidentes",
        "Bitácora operativa",
      ],
      expone: [
        "POST /citizen-reports",
        "GET /reports/{cedula}",
        "GET /alerts/{distrito}",
      ],
      consume: [
        "registro-civil · GET /persons/{cedula}",
        "poder-judicial · POST /complaints",
        "iduc-signing · firma de denuncias",
      ],
    },
    en: {
      name: "Public Force · Ministry of Security",
      short: "Citizen reports · operations · public safety",
      custodia: [
        "Citizen reports",
        "Incident reports",
        "Operational logbook",
      ],
      expone: [
        "POST /citizen-reports",
        "GET /reports/{id}",
        "GET /alerts/{district}",
      ],
      consume: [
        "registro-civil · GET /persons/{id}",
        "poder-judicial · POST /complaints",
        "iduc-signing · report signature",
      ],
    },
  },
  {
    id: "bomberos",
    tag: "bom",
    status: "planned",
    accent: "#dc2626",
    es: {
      name: "Cuerpo de Bomberos",
      short: "Inspecciones · prevención · respuesta a emergencias",
      custodia: [
        "Inspecciones de seguridad humana",
        "Permisos contra incendios",
        "Reportes de emergencia",
      ],
      expone: [
        "GET /inspections/{folio-real}",
        "GET /fire-permit/{cedula-juridica}",
        "POST /emergency-reports",
      ],
      consume: [
        "muni-sjo · GET /construction-permits/{folio}",
        "registro-nacional · GET /property/{folio-real}",
      ],
    },
    en: {
      name: "Fire Department",
      short: "Inspections · prevention · emergency response",
      custodia: [
        "Human-safety inspections",
        "Fire permits",
        "Emergency reports",
      ],
      expone: [
        "GET /inspections/{folio-real}",
        "GET /fire-permit/{legal-entity-id}",
        "POST /emergency-reports",
      ],
      consume: [
        "muni-sjo · GET /construction-permits/{folio}",
        "registro-nacional · GET /property/{folio-real}",
      ],
    },
  },
  {
    id: "procomer",
    tag: "pro",
    status: "planned",
    accent: "#0ea5e9",
    es: {
      name: "Promotora de Comercio Exterior",
      short: "Exportadores · zonas francas · comercio internacional",
      custodia: [
        "Padrón de exportadores",
        "Empresas en régimen de zona franca",
        "Trámites VUCE (ventanilla única)",
      ],
      expone: [
        "GET /exporters/{cedula-juridica}",
        "GET /free-zone-status/{cedula-juridica}",
        "POST /export-permits",
      ],
      consume: [
        "registro-nacional · GET /entity/{cedula-juridica}",
        "hacienda · GET /tax-status/{cedula-juridica}",
        "minae · GET /environmental-clearance/{folio}",
      ],
    },
    en: {
      name: "Foreign Trade Promoter",
      short: "Exporters · free zones · international trade",
      custodia: [
        "Exporter roster",
        "Free-zone regime companies",
        "Single-window filings (VUCE)",
      ],
      expone: [
        "GET /exporters/{legal-entity-id}",
        "GET /free-zone-status/{legal-entity-id}",
        "POST /export-permits",
      ],
      consume: [
        "registro-nacional · GET /entity/{legal-entity-id}",
        "hacienda · GET /tax-status/{legal-entity-id}",
        "minae · GET /environmental-clearance/{folio}",
      ],
    },
  },
];

export const getInstitution = (id: InstitutionId): InstitutionDef | undefined =>
  INSTITUTIONS.find((i) => i.id === id);

export const localized = (def: InstitutionDef, lang: Lang): InstitutionLocale =>
  lang === "es" ? def.es : def.en;

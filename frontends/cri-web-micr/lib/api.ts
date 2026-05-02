// Cliente del cri-bff-citizen para MiCR.
// Toda llamada propaga el header X-CRI-Realm al BFF (que a su vez lo
// propaga a los upstreams).

const BFF_BASE = process.env.NEXT_PUBLIC_BFF_BASE_URL || "http://localhost:18001";
const REALM = process.env.NEXT_PUBLIC_REALM || "demo";

type Envelope<T> = { data: T; meta: { requestId: string; realm?: string } };

async function fetchEnvelope<T>(path: string): Promise<T> {
  const res = await fetch(`${BFF_BASE}${path}`, {
    headers: { "X-CRI-Realm": REALM },
    cache: "no-store",
  });
  if (!res.ok) {
    const body = await res.text();
    throw new Error(`BFF ${res.status}: ${body}`);
  }
  const env: Envelope<T> = await res.json();
  return env.data;
}

export type Person = {
  cedula: string;
  fullName: string;
  address: string;
  email: string;
};

export type TaxPrefilled = {
  person: Person;
  year: number;
  grossIncome: number;
  withheldTax: number;
  deductions: number;
  estimatedDue: number;
  hasDependents: boolean;
  _onceOnlyTrace: string;
};

export type AccessLogEntry = {
  id: string;
  realm: string;
  ts: string;
  requesterMember: string;
  targetMember: string;
  service: string;
  version: string;
  citizenId?: string;
  purpose: string;
  requestId: string;
  status: number;
  prevHash: string;
  entryHash: string;
};

export type AccessLog = {
  citizenId: string;
  realm: string;
  count: number;
  entries: AccessLogEntry[];
};

export type Prescription = {
  id: string;
  drug: string;
  dosage: string;
  issuedAt: string;
  issuedBy: string;
  validUntil: string;
  status: string;
};

export type Appointment = {
  id: string;
  specialty: string;
  hospital: string;
  when: string;
  status: string;
};

export type HealthProfile = {
  patient: { cedula: string; fullName: string; address: string };
  activePrescriptions: Prescription[];
  nextAppointment: Appointment | null;
  _onceOnlyTrace: string;
};

export type Dashboard = {
  cedula: string;
  realm: string;
  year: number;
  taxPrefilled: Envelope<TaxPrefilled>;
  healthProfile: Envelope<HealthProfile>;
  accessLog: Envelope<AccessLog>;
  upstreamErrors?: string[] | null;
};

export const api = {
  dashboard: (cedula: string, year = 2025) =>
    fetchEnvelope<Dashboard>(`/v1/citizens/${encodeURIComponent(cedula)}/dashboard?year=${year}`),
  accessLog: (cedula: string) =>
    fetchEnvelope<AccessLog>(`/v1/citizens/${encodeURIComponent(cedula)}/access-log`),
  taxPrefilled: (cedula: string, year = 2025) =>
    fetchEnvelope<TaxPrefilled>(`/v1/citizens/${encodeURIComponent(cedula)}/tax/prefilled?year=${year}`),
  healthProfile: (cedula: string) =>
    fetchEnvelope<HealthProfile>(`/v1/citizens/${encodeURIComponent(cedula)}/health/profile`),
};

export const formatCRC = (n: number) =>
  new Intl.NumberFormat("es-CR", { style: "currency", currency: "CRC", maximumFractionDigits: 0 }).format(n);

export const formatDateTime = (iso: string) => {
  try {
    return new Date(iso).toLocaleString("es-CR", {
      dateStyle: "medium",
      timeStyle: "short",
    });
  } catch {
    return iso;
  }
};

// Mock data para correr la app sin BFF (demos en iPhone, presentaciones, etc).
// Activado por defecto; se desactiva con --dart-define=USE_MOCK=false.

library;

import 'api.dart';

const bool kUseMock = bool.fromEnvironment('USE_MOCK', defaultValue: true);

/// Pequeña latencia simulada para que los spinners se sientan reales.
Future<T> _delayed<T>(T value) =>
    Future.delayed(const Duration(milliseconds: 380), () => value);

class _Profile {
  final String cedula;
  final String fullName;
  final String address;
  final String email;
  final num grossIncome;
  final num withheldTax;
  final num deductions;
  final num estimatedDue;
  final bool hasDependents;
  final List<Prescription> prescriptions;
  final Appointment? nextAppointment;
  final List<AccessLogEntry> log;
  const _Profile({
    required this.cedula,
    required this.fullName,
    required this.address,
    required this.email,
    required this.grossIncome,
    required this.withheldTax,
    required this.deductions,
    required this.estimatedDue,
    required this.hasDependents,
    required this.prescriptions,
    required this.nextAppointment,
    required this.log,
  });

  Person get person => Person(
        cedula: cedula,
        fullName: fullName,
        address: address,
        email: email,
      );
}

// Hash determinístico (no criptográfico) — solo para que los hex se vean
// estables entre runs y formen una cadena plausible visualmente.
String _h(String seed) {
  var hash = 0xcbf29ce484222325;
  for (final c in seed.codeUnits) {
    hash ^= c;
    hash = (hash * 0x100000001b3) & 0xFFFFFFFFFFFFFFFF;
  }
  final base = hash.toRadixString(16).padLeft(16, '0');
  return (base + base + base + base).substring(0, 64);
}

List<AccessLogEntry> _buildLog(
  String cedula,
  List<({DateTime ts, String requester, String target, String svc, String ver, String purpose})> raw,
) {
  String prev = '0' * 64;
  final entries = <AccessLogEntry>[];
  for (var i = 0; i < raw.length; i++) {
    final r = raw[i];
    final id = '01HXMOCK${cedula.replaceAll('-', '')}${i.toString().padLeft(2, '0')}';
    final entry = _h('$id|${r.ts.toIso8601String()}|$prev|$cedula');
    entries.add(AccessLogEntry(
      id: id,
      ts: r.ts,
      requesterMember: r.requester,
      targetMember: r.target,
      service: r.svc,
      version: r.ver,
      purpose: r.purpose,
      entryHash: entry,
      prevHash: prev,
    ));
    prev = entry;
  }
  // Más reciente primero para que la card del dashboard tome el "último".
  return entries.reversed.toList();
}

final Map<String, _Profile> _profiles = {
  '1-1234-5678': _Profile(
    cedula: '1-1234-5678',
    fullName: 'Sebastián Rojas Mora',
    address: 'San Pedro, Montes de Oca, San José',
    email: 'sebastian.rojas@correo.cr',
    grossIncome: 18000000,
    withheldTax: 1450000,
    deductions: 2100000,
    estimatedDue: 312500,
    hasDependents: false,
    prescriptions: [
      Prescription(
        id: 'rx-001',
        drug: 'Loratadina 10mg',
        dosage: '1 tableta diaria',
        issuedAt: '2026-04-12',
        issuedBy: 'Clínica Carlos Durán',
        validUntil: '2026-07-12',
        status: 'active',
      ),
    ],
    nextAppointment: Appointment(
      id: 'apt-001',
      specialty: 'Medicina General',
      hospital: 'Clínica Carlos Durán',
      when: '2026-05-18 09:30',
      status: 'scheduled',
    ),
    log: _buildLog('1-1234-5678', [
      (
        ts: DateTime(2026, 4, 30, 11, 7),
        requester: 'cri-svc-hacienda',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'prefill_tax_return',
      ),
      (
        ts: DateTime(2026, 4, 28, 14, 22),
        requester: 'cri-svc-ccss',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'verify_insurance_holder',
      ),
      (
        ts: DateTime(2026, 4, 21, 9, 4),
        requester: 'cri-svc-banco-popular',
        target: 'cri-svc-registro-civil',
        svc: 'persons.verify',
        ver: 'v1',
        purpose: 'kyc_account_opening',
      ),
      (
        ts: DateTime(2026, 4, 14, 16, 51),
        requester: 'cri-svc-hacienda',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'prefill_tax_return',
      ),
    ]),
  ),
  '1-9876-5432': _Profile(
    cedula: '1-9876-5432',
    fullName: 'María Castro Fernández',
    address: 'Curridabat, Granadilla Norte, San José',
    email: 'maria.castro@correo.cr',
    grossIncome: 24500000,
    withheldTax: 2380000,
    deductions: 3650000,
    estimatedDue: 587000,
    hasDependents: true,
    prescriptions: [
      Prescription(
        id: 'rx-101',
        drug: 'Levotiroxina 50mcg',
        dosage: '1 tableta en ayunas',
        issuedAt: '2026-03-02',
        issuedBy: 'Hospital México',
        validUntil: '2026-09-02',
        status: 'active',
      ),
      Prescription(
        id: 'rx-102',
        drug: 'Atorvastatina 20mg',
        dosage: '1 tableta en la noche',
        issuedAt: '2026-04-08',
        issuedBy: 'Hospital México',
        validUntil: '2026-10-08',
        status: 'active',
      ),
    ],
    nextAppointment: Appointment(
      id: 'apt-101',
      specialty: 'Endocrinología',
      hospital: 'Hospital México',
      when: '2026-06-04 08:00',
      status: 'scheduled',
    ),
    log: _buildLog('1-9876-5432', [
      (
        ts: DateTime(2026, 5, 1, 8, 12),
        requester: 'cri-svc-hacienda',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'prefill_tax_return',
      ),
      (
        ts: DateTime(2026, 5, 1, 8, 12),
        requester: 'cri-svc-hacienda',
        target: 'cri-svc-registro-civil',
        svc: 'dependents.read',
        ver: 'v1',
        purpose: 'verify_dependents',
      ),
      (
        ts: DateTime(2026, 4, 27, 10, 35),
        requester: 'cri-svc-ccss',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'medical_record_lookup',
      ),
      (
        ts: DateTime(2026, 4, 19, 15, 41),
        requester: 'cri-svc-mep',
        target: 'cri-svc-registro-civil',
        svc: 'dependents.read',
        ver: 'v1',
        purpose: 'school_enrollment',
      ),
      (
        ts: DateTime(2026, 4, 12, 12, 0),
        requester: 'cri-svc-banco-nacional',
        target: 'cri-svc-registro-civil',
        svc: 'persons.verify',
        ver: 'v1',
        purpose: 'kyc_credit_application',
      ),
      (
        ts: DateTime(2026, 4, 5, 9, 18),
        requester: 'cri-svc-hacienda',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'prefill_tax_return',
      ),
    ]),
  ),
  '2-0001-0001': _Profile(
    cedula: '2-0001-0001',
    fullName: 'Luis Fernández Solís',
    address: 'Alajuela centro, Alajuela',
    email: 'luis.fernandez@correo.cr',
    grossIncome: 12000000,
    withheldTax: 720000,
    deductions: 1450000,
    estimatedDue: 0,
    hasDependents: false,
    prescriptions: const [],
    nextAppointment: null,
    log: _buildLog('2-0001-0001', [
      (
        ts: DateTime(2026, 4, 29, 17, 33),
        requester: 'cri-svc-hacienda',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'prefill_tax_return',
      ),
      (
        ts: DateTime(2026, 4, 22, 13, 5),
        requester: 'cri-svc-mtss',
        target: 'cri-svc-registro-civil',
        svc: 'persons.read',
        ver: 'v1',
        purpose: 'unemployment_benefit',
      ),
    ]),
  ),
};

_Profile _profileFor(String cedula) =>
    _profiles[cedula] ?? _profiles['1-1234-5678']!;

class MockApi {
  static Future<Dashboard> dashboard(String cedula, {int year = 2025}) {
    final p = _profileFor(cedula);
    return _delayed(Dashboard(
      cedula: cedula,
      year: year,
      taxPrefilled: _taxFor(p, year),
      healthProfile: HealthProfile(
        activePrescriptions: p.prescriptions,
        nextAppointment: p.nextAppointment,
      ),
      accessLog: AccessLog(
        citizenId: cedula,
        count: p.log.length,
        entries: p.log,
      ),
    ));
  }

  static Future<AccessLog> accessLog(String cedula) {
    final p = _profileFor(cedula);
    return _delayed(AccessLog(
      citizenId: cedula,
      count: p.log.length,
      entries: p.log,
    ));
  }

  static Future<TaxPrefilled> taxPrefilled(String cedula, {int year = 2025}) {
    final p = _profileFor(cedula);
    return _delayed(_taxFor(p, year));
  }

  static TaxPrefilled _taxFor(_Profile p, int year) {
    return TaxPrefilled(
      person: p.person,
      year: year,
      grossIncome: p.grossIncome,
      withheldTax: p.withheldTax,
      deductions: p.deductions,
      estimatedDue: p.estimatedDue,
      hasDependents: p.hasDependents,
      onceOnlyTrace:
          'cri-svc-hacienda → cri-svc-registro-civil [GET /persons/${p.cedula}] · '
          'request_id=01HXMOCK${p.cedula.replaceAll('-', '')}TX$year · '
          'sig=ed25519 · 142ms · status=200',
    );
  }
}

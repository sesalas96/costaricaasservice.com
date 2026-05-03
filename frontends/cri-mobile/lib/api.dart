// Cliente del cri-bff-citizen para MiCR mobile.
//
// En dev apunta a localhost:18001 (BFF). En Android Emulator se cambia a
// http://10.0.2.2:18001 (alias del host); en device real, IP de la LAN.
library;

import 'dart:convert';
import 'dart:io' show Platform;

import 'package:http/http.dart' as http;

import 'mock.dart';

/// Base URL del BFF. Se resuelve en tiempo de runtime según plataforma:
///   - iOS Simulator / macOS: http://localhost:18001
///   - Android Emulator: http://10.0.2.2:18001
String get bffBaseUrl {
  // Override por compile-time --dart-define=BFF_BASE_URL=...
  const fromEnv = String.fromEnvironment('BFF_BASE_URL', defaultValue: '');
  if (fromEnv.isNotEmpty) return fromEnv;
  try {
    if (Platform.isAndroid) return 'http://10.0.2.2:18001';
  } catch (_) {
    // Web / unknown platform
  }
  return 'http://localhost:18001';
}

const String realm = 'demo';

class ApiException implements Exception {
  final int? status;
  final String message;
  ApiException(this.message, {this.status});
  @override
  String toString() => 'ApiException(${status ?? '-'}): $message';
}

Future<Map<String, dynamic>> _fetchEnvelope(String path) async {
  final uri = Uri.parse('$bffBaseUrl$path');
  final resp = await http.get(uri, headers: {'X-CRI-Realm': realm});
  if (resp.statusCode >= 400) {
    throw ApiException('BFF ${resp.statusCode}: ${resp.body}', status: resp.statusCode);
  }
  final decoded = jsonDecode(resp.body) as Map<String, dynamic>;
  return decoded['data'] as Map<String, dynamic>;
}

class Person {
  final String cedula;
  final String fullName;
  final String address;
  final String email;
  Person({required this.cedula, required this.fullName, required this.address, required this.email});
  factory Person.fromJson(Map<String, dynamic> j) => Person(
        cedula: (j['cedula'] ?? '') as String,
        fullName: (j['fullName'] ?? '') as String,
        address: (j['address'] ?? '') as String,
        email: (j['email'] ?? '') as String,
      );
}

class TaxPrefilled {
  final Person person;
  final int year;
  final num grossIncome;
  final num withheldTax;
  final num deductions;
  final num estimatedDue;
  final bool hasDependents;
  final String onceOnlyTrace;
  TaxPrefilled({
    required this.person,
    required this.year,
    required this.grossIncome,
    required this.withheldTax,
    required this.deductions,
    required this.estimatedDue,
    required this.hasDependents,
    required this.onceOnlyTrace,
  });
  factory TaxPrefilled.fromJson(Map<String, dynamic> j) => TaxPrefilled(
        person: Person.fromJson(j['person'] as Map<String, dynamic>),
        year: (j['year'] as num).toInt(),
        grossIncome: j['grossIncome'] as num,
        withheldTax: j['withheldTax'] as num,
        deductions: j['deductions'] as num,
        estimatedDue: j['estimatedDue'] as num,
        hasDependents: (j['hasDependents'] ?? false) as bool,
        onceOnlyTrace: (j['_onceOnlyTrace'] ?? '') as String,
      );
}

class AccessLogEntry {
  final String id;
  final DateTime ts;
  final String requesterMember;
  final String targetMember;
  final String service;
  final String version;
  final String purpose;
  final String entryHash;
  final String prevHash;
  AccessLogEntry({
    required this.id,
    required this.ts,
    required this.requesterMember,
    required this.targetMember,
    required this.service,
    required this.version,
    required this.purpose,
    required this.entryHash,
    required this.prevHash,
  });
  factory AccessLogEntry.fromJson(Map<String, dynamic> j) => AccessLogEntry(
        id: j['id'] as String,
        ts: DateTime.parse(j['ts'] as String),
        requesterMember: j['requesterMember'] as String,
        targetMember: j['targetMember'] as String,
        service: j['service'] as String,
        version: j['version'] as String,
        purpose: (j['purpose'] ?? '') as String,
        entryHash: (j['entryHash'] ?? '') as String,
        prevHash: (j['prevHash'] ?? '') as String,
      );
}

class AccessLog {
  final String citizenId;
  final int count;
  final List<AccessLogEntry> entries;
  AccessLog({required this.citizenId, required this.count, required this.entries});
  factory AccessLog.fromJson(Map<String, dynamic> j) => AccessLog(
        citizenId: j['citizenId'] as String,
        count: (j['count'] as num).toInt(),
        entries: ((j['entries'] ?? []) as List)
            .map((e) => AccessLogEntry.fromJson(e as Map<String, dynamic>))
            .toList(),
      );
}

class Prescription {
  final String id;
  final String drug;
  final String dosage;
  final String issuedAt;
  final String issuedBy;
  final String validUntil;
  final String status;
  Prescription({
    required this.id,
    required this.drug,
    required this.dosage,
    required this.issuedAt,
    required this.issuedBy,
    required this.validUntil,
    required this.status,
  });
  factory Prescription.fromJson(Map<String, dynamic> j) => Prescription(
        id: j['id'] as String,
        drug: j['drug'] as String,
        dosage: j['dosage'] as String,
        issuedAt: (j['issuedAt'] ?? '') as String,
        issuedBy: (j['issuedBy'] ?? '') as String,
        validUntil: (j['validUntil'] ?? '') as String,
        status: j['status'] as String,
      );
}

class Appointment {
  final String id;
  final String specialty;
  final String hospital;
  final String when;
  final String status;
  Appointment({
    required this.id,
    required this.specialty,
    required this.hospital,
    required this.when,
    required this.status,
  });
  factory Appointment.fromJson(Map<String, dynamic> j) => Appointment(
        id: j['id'] as String,
        specialty: j['specialty'] as String,
        hospital: j['hospital'] as String,
        when: j['when'] as String,
        status: j['status'] as String,
      );
}

class HealthProfile {
  final List<Prescription> activePrescriptions;
  final Appointment? nextAppointment;
  HealthProfile({required this.activePrescriptions, this.nextAppointment});
  factory HealthProfile.fromJson(Map<String, dynamic> j) => HealthProfile(
        activePrescriptions: ((j['activePrescriptions'] ?? []) as List)
            .map((e) => Prescription.fromJson(e as Map<String, dynamic>))
            .toList(),
        nextAppointment: j['nextAppointment'] == null
            ? null
            : Appointment.fromJson(j['nextAppointment'] as Map<String, dynamic>),
      );
}

class Dashboard {
  final String cedula;
  final int year;
  final TaxPrefilled? taxPrefilled;
  final HealthProfile? healthProfile;
  final AccessLog? accessLog;
  final List<String> upstreamErrors;
  Dashboard({
    required this.cedula,
    required this.year,
    this.taxPrefilled,
    this.healthProfile,
    this.accessLog,
    this.upstreamErrors = const [],
  });
  factory Dashboard.fromJson(Map<String, dynamic> j) {
    final taxRaw = j['taxPrefilled'] as Map<String, dynamic>?;
    final healthRaw = j['healthProfile'] as Map<String, dynamic>?;
    final logRaw = j['accessLog'] as Map<String, dynamic>?;
    return Dashboard(
      cedula: j['cedula'] as String,
      year: (j['year'] as num).toInt(),
      taxPrefilled: taxRaw != null && taxRaw['data'] != null
          ? TaxPrefilled.fromJson(taxRaw['data'] as Map<String, dynamic>)
          : null,
      healthProfile: healthRaw != null && healthRaw['data'] != null
          ? HealthProfile.fromJson(healthRaw['data'] as Map<String, dynamic>)
          : null,
      accessLog: logRaw != null && logRaw['data'] != null
          ? AccessLog.fromJson(logRaw['data'] as Map<String, dynamic>)
          : null,
      upstreamErrors: ((j['upstreamErrors'] ?? []) as List).cast<String>(),
    );
  }
}

class Api {
  static Future<Dashboard> dashboard(String cedula, {int year = 2025}) async {
    if (kUseMock) return MockApi.dashboard(cedula, year: year);
    final data = await _fetchEnvelope('/v1/citizens/${Uri.encodeComponent(cedula)}/dashboard?year=$year');
    return Dashboard.fromJson(data);
  }

  static Future<AccessLog> accessLog(String cedula) async {
    if (kUseMock) return MockApi.accessLog(cedula);
    final data = await _fetchEnvelope('/v1/citizens/${Uri.encodeComponent(cedula)}/access-log');
    return AccessLog.fromJson(data);
  }

  static Future<TaxPrefilled> taxPrefilled(String cedula, {int year = 2025}) async {
    if (kUseMock) return MockApi.taxPrefilled(cedula, year: year);
    final data = await _fetchEnvelope('/v1/citizens/${Uri.encodeComponent(cedula)}/tax/prefilled?year=$year');
    return TaxPrefilled.fromJson(data);
  }
}

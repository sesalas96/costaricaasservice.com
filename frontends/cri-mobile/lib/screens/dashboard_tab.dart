import 'package:flutter/material.dart';

import '../api.dart';
import '../format.dart';
import '../theme.dart';

class DashboardTab extends StatefulWidget {
  final String cedula;
  const DashboardTab({super.key, required this.cedula});
  @override
  State<DashboardTab> createState() => _DashboardTabState();
}

class _DashboardTabState extends State<DashboardTab> {
  Future<Dashboard>? _future;

  @override
  void initState() {
    super.initState();
    _refresh();
  }

  void _refresh() {
    setState(() => _future = Api.dashboard(widget.cedula, year: 2025));
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () async => _refresh(),
      child: FutureBuilder<Dashboard>(
        future: _future,
        builder: (context, snap) {
          if (snap.connectionState != ConnectionState.done) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snap.hasError) {
            return ListView(
              padding: const EdgeInsets.all(16),
              children: [
                ApiErrorView(error: snap.error!, onRetry: _refresh),
              ],
            );
          }
          final d = snap.data!;
          final tax = d.taxPrefilled;
          final health = d.healthProfile;
          final log = d.accessLog;
          final last = (log != null && log.entries.isNotEmpty) ? log.entries.first : null;

          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              const Text('Mi escritorio',
                  style: TextStyle(fontSize: 22, fontWeight: FontWeight.w700, color: CrColors.text)),
              Text('Año fiscal ${d.year} · cédula ${d.cedula}',
                  style: const TextStyle(fontSize: 12, color: CrColors.muted)),
              const SizedBox(height: 18),

              _Card(
                title: 'Mis datos',
                badge: ('Registro Civil', CrColors.areaIducDim, CrColors.areaIduc),
                footer:
                    'Estos datos viajan firmados desde Registro Civil. Hacienda los consulta cada vez sin guardarlos (once-only).',
                child: tax?.person != null
                    ? Column(
                        children: [
                          _Row(label: 'Nombre', value: tax!.person.fullName),
                          _Row(label: 'Cédula', value: tax.person.cedula, mono: true),
                          _Row(label: 'Domicilio', value: tax.person.address),
                          _Row(label: 'Email', value: tax.person.email),
                        ],
                      )
                    : const Text('Sin datos disponibles', style: TextStyle(color: CrColors.muted)),
              ),

              _Card(
                title: 'Mi declaración 2025',
                badge: ('Hacienda', CrColors.crBlueBrightDim, CrColors.crBlueBright),
                child: tax != null
                    ? Column(
                        children: [
                          _Row(label: 'Ingreso bruto', value: formatCRC(tax.grossIncome)),
                          _Row(label: 'Retenido', value: formatCRC(tax.withheldTax)),
                          _Row(label: 'Deducciones', value: formatCRC(tax.deductions)),
                          _Row(label: 'A pagar', value: formatCRC(tax.estimatedDue), bold: true),
                          _Row(label: 'Dependientes', value: tax.hasDependents ? 'Sí' : 'No'),
                        ],
                      )
                    : const Text('Sin declaración', style: TextStyle(color: CrColors.muted)),
              ),

              _Card(
                title: 'Mi salud',
                badge: ('CCSS', CrColors.crRedBrightDim, CrColors.crRedBright),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('${health?.activePrescriptions.length ?? 0}',
                        style: const TextStyle(
                            fontSize: 32, fontWeight: FontWeight.w700, color: CrColors.text)),
                    const Text('recetas activas', style: TextStyle(color: CrColors.muted)),
                    if (health != null) ...[
                      const SizedBox(height: 8),
                      ...health.activePrescriptions.take(2).map((rx) => Container(
                            margin: const EdgeInsets.only(top: 6),
                            padding: const EdgeInsets.all(8),
                            decoration: BoxDecoration(
                              color: CrColors.surface2,
                              borderRadius: BorderRadius.circular(6),
                              border: Border.all(color: CrColors.border),
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(rx.drug,
                                    style: const TextStyle(
                                        fontSize: 12,
                                        fontWeight: FontWeight.w600,
                                        color: CrColors.text)),
                                Text(rx.dosage,
                                    style: const TextStyle(
                                        fontSize: 11, color: CrColors.muted)),
                              ],
                            ),
                          )),
                      if (health.nextAppointment != null) ...[
                        const SizedBox(height: 12),
                        const Divider(height: 1, color: CrColors.border),
                        const SizedBox(height: 8),
                        const Text('Próxima cita',
                            style: TextStyle(fontSize: 11, color: CrColors.muted)),
                        Text(health.nextAppointment!.specialty,
                            style: const TextStyle(
                                fontSize: 13,
                                fontWeight: FontWeight.w600,
                                color: CrColors.text)),
                        Text(
                          '${health.nextAppointment!.when} · ${health.nextAppointment!.hospital}',
                          style: const TextStyle(fontSize: 11, color: CrColors.muted),
                        ),
                      ],
                    ],
                  ],
                ),
              ),

              _Card(
                title: 'Mi bitácora',
                badge: ('Audit chain', CrColors.areaPlatformDim, CrColors.areaPlatform),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('${log?.count ?? 0}',
                        style: const TextStyle(
                            fontSize: 32, fontWeight: FontWeight.w700, color: CrColors.text)),
                    const Text('accesos a sus datos',
                        style: TextStyle(color: CrColors.muted)),
                    if (last != null) ...[
                      const SizedBox(height: 12),
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: CrColors.surface2,
                          borderRadius: BorderRadius.circular(8),
                          border: Border.all(color: CrColors.border),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text('Último acceso',
                                style: TextStyle(
                                    fontWeight: FontWeight.w600, color: CrColors.text)),
                            const SizedBox(height: 4),
                            Text(
                                '${last.requesterMember} consultó ${last.service} en ${last.targetMember}',
                                style: const TextStyle(fontSize: 12, color: CrColors.text)),
                            Text(formatDateTime(last.ts),
                                style: const TextStyle(fontSize: 12, color: CrColors.muted)),
                            Text('propósito: ${last.purpose}',
                                style: const TextStyle(fontSize: 12, color: CrColors.muted)),
                          ],
                        ),
                      ),
                    ],
                  ],
                ),
              ),

              if (d.upstreamErrors.isNotEmpty)
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: const Color(0x33F59E0B),
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: CrColors.areaMembers),
                  ),
                  child: Text('Avisos: ${d.upstreamErrors.join('; ')}',
                      style: const TextStyle(color: CrColors.areaMembers, fontSize: 12)),
                ),
            ],
          );
        },
      ),
    );
  }
}

class _Card extends StatelessWidget {
  final String title;
  final (String, Color, Color)? badge; // (label, bg, fg)
  final Widget child;
  final String? footer;
  const _Card({required this.title, required this.child, this.badge, this.footer});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 14),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: CrColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: CrColors.border),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(title,
                  style: const TextStyle(
                      fontSize: 16, fontWeight: FontWeight.w600, color: CrColors.text)),
              if (badge != null)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: badge!.$2,
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(badge!.$1,
                      style: TextStyle(
                          color: badge!.$3, fontSize: 11, fontWeight: FontWeight.w600)),
                ),
            ],
          ),
          const SizedBox(height: 12),
          child,
          if (footer != null) ...[
            const SizedBox(height: 12),
            Text(footer!, style: const TextStyle(fontSize: 11, color: CrColors.muted)),
          ],
        ],
      ),
    );
  }
}

class _Row extends StatelessWidget {
  final String label;
  final String value;
  final bool mono;
  final bool bold;
  const _Row({required this.label, required this.value, this.mono = false, this.bold = false});
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: CrColors.muted)),
          Flexible(
            child: Text(
              value,
              textAlign: TextAlign.right,
              style: TextStyle(
                fontFamily: mono ? 'monospace' : null,
                fontWeight: bold ? FontWeight.w700 : null,
                color: CrColors.text,
              ),
            ),
          ),
        ],
      ),
    );
  }
}


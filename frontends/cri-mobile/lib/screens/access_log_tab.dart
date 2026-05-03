import 'package:flutter/material.dart';

import '../api.dart';
import '../format.dart';
import '../theme.dart';

class AccessLogTab extends StatefulWidget {
  final String cedula;
  const AccessLogTab({super.key, required this.cedula});
  @override
  State<AccessLogTab> createState() => _AccessLogTabState();
}

class _AccessLogTabState extends State<AccessLogTab> {
  Future<AccessLog>? _future;

  @override
  void initState() {
    super.initState();
    _refresh();
  }

  void _refresh() {
    setState(() => _future = Api.accessLog(widget.cedula));
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () async => _refresh(),
      child: FutureBuilder<AccessLog>(
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
          final log = snap.data!;
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              const Text('Mi bitácora',
                  style: TextStyle(fontSize: 22, fontWeight: FontWeight.w700, color: CrColors.text)),
              const SizedBox(height: 4),
              const Text(
                'Accesos del Estado a sus datos personales — quién, cuándo y por qué. Cada entry está '
                'encadenada criptográficamente; modificar una línea pasada rompería todas las posteriores.',
                style: TextStyle(fontSize: 12, color: CrColors.muted),
              ),
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: CrColors.surface,
                  borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
                  border: Border.all(color: CrColors.border),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text.rich(TextSpan(children: [
                      TextSpan(
                          text: '${log.count}',
                          style: const TextStyle(
                              fontSize: 22, fontWeight: FontWeight.w700, color: CrColors.text)),
                      TextSpan(
                          text: ' acceso${log.count == 1 ? '' : 's'} registrado${log.count == 1 ? '' : 's'}',
                          style: const TextStyle(fontSize: 13, color: CrColors.muted)),
                    ])),
                    Text('cédula ${widget.cedula}',
                        style: const TextStyle(
                          fontSize: 11,
                          fontFamily: 'monospace',
                          color: CrColors.muted,
                        )),
                  ],
                ),
              ),
              if (log.entries.isEmpty)
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: CrColors.surface,
                    borderRadius: const BorderRadius.vertical(bottom: Radius.circular(12)),
                    border: Border.all(color: CrColors.border),
                  ),
                  alignment: Alignment.center,
                  child: const Text('No hay accesos registrados aún.',
                      style: TextStyle(color: CrColors.muted)),
                )
              else
                Container(
                  decoration: BoxDecoration(
                    color: CrColors.surface,
                    borderRadius: const BorderRadius.vertical(bottom: Radius.circular(12)),
                    border: Border.all(color: CrColors.border),
                  ),
                  child: Column(
                    children: log.entries
                        .map((e) => _Entry(e: e, last: e == log.entries.last))
                        .toList(),
                  ),
                ),
              const SizedBox(height: 16),
              const Text(
                'El audit log es append-only y verificado por hash chain (SHA-256). Cualquier intento de '
                'borrado o modificación retroactiva sería detectable en menos de un segundo.',
                style: TextStyle(fontSize: 11, color: CrColors.muted),
              ),
            ],
          );
        },
      ),
    );
  }
}

class _Entry extends StatelessWidget {
  final AccessLogEntry e;
  final bool last;
  const _Entry({required this.e, required this.last});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        border: Border(
          bottom: last ? BorderSide.none : const BorderSide(color: CrColors.border),
        ),
      ),
      padding: const EdgeInsets.all(14),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text.rich(TextSpan(
            style: const TextStyle(color: CrColors.text),
            children: [
              TextSpan(
                  text: e.requesterMember,
                  style: const TextStyle(fontWeight: FontWeight.w600, fontFamily: 'monospace')),
              const TextSpan(text: ' consultó '),
              TextSpan(
                  text: '${e.service}/${e.version}',
                  style: const TextStyle(fontFamily: 'monospace')),
              const TextSpan(text: ' en '),
              TextSpan(
                  text: e.targetMember,
                  style: const TextStyle(fontWeight: FontWeight.w600, fontFamily: 'monospace')),
            ],
          )),
          const SizedBox(height: 4),
          Text(formatDateTime(e.ts),
              style: const TextStyle(fontSize: 11, color: CrColors.muted)),
          const SizedBox(height: 6),
          Text('Propósito: ${e.purpose}',
              style: const TextStyle(fontSize: 12, color: CrColors.text)),
          const SizedBox(height: 6),
          Wrap(
            spacing: 16,
            runSpacing: 2,
            children: [
              Text('id: ${e.id}',
                  style: const TextStyle(
                      fontSize: 10, fontFamily: 'monospace', color: CrColors.muted)),
              Text('hash: ${shortHash(e.entryHash)}',
                  style: const TextStyle(
                      fontSize: 10, fontFamily: 'monospace', color: CrColors.muted)),
              Text('prev: ${shortHash(e.prevHash)}',
                  style: const TextStyle(
                      fontSize: 10, fontFamily: 'monospace', color: CrColors.muted)),
            ],
          ),
        ],
      ),
    );
  }
}

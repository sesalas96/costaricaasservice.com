import 'package:flutter/material.dart';

import '../api.dart';
import '../format.dart';

class TaxTab extends StatefulWidget {
  final String cedula;
  const TaxTab({super.key, required this.cedula});
  @override
  State<TaxTab> createState() => _TaxTabState();
}

class _TaxTabState extends State<TaxTab> {
  Future<TaxPrefilled>? _future;
  static const int _year = 2025;

  @override
  void initState() {
    super.initState();
    _refresh();
  }

  void _refresh() {
    setState(() => _future = Api.taxPrefilled(widget.cedula, year: _year));
  }

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: () async => _refresh(),
      child: FutureBuilder<TaxPrefilled>(
        future: _future,
        builder: (context, snap) {
          if (snap.connectionState != ConnectionState.done) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snap.hasError) {
            return ListView(padding: const EdgeInsets.all(20), children: [Text(snap.error.toString())]);
          }
          final t = snap.data!;
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              const Text('Mi declaración 2025',
                  style: TextStyle(fontSize: 22, fontWeight: FontWeight.w700)),
              const SizedBox(height: 4),
              const Text(
                'Pre-llenada por Hacienda. Sus datos personales fueron consultados al Registro Civil '
                'vía Conecta CR con propósito declarado (prefill_tax_return).',
                style: TextStyle(fontSize: 12, color: Colors.black54),
              ),
              const SizedBox(height: 14),
              _section('Identificación', [
                _row('Nombre', t.person.fullName),
                _rowMono('Cédula', t.person.cedula),
                _row('Domicilio', t.person.address),
                _row('Email', t.person.email),
              ]),
              _section('Detalle del año $_year', [
                _row('Ingreso bruto', formatCRC(t.grossIncome)),
                _row('Impuesto retenido', formatCRC(t.withheldTax)),
                _row('Deducciones', formatCRC(t.deductions)),
                _row('Dependientes', t.hasDependents ? 'Sí' : 'No'),
              ]),
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: const Color(0xFFE2E8F0)),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Estimado a pagar', style: TextStyle(color: Colors.black54)),
                    Text(formatCRC(t.estimatedDue),
                        style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w700)),
                  ],
                ),
              ),
              const SizedBox(height: 12),
              SizedBox(
                width: double.infinity,
                child: FilledButton(
                  onPressed: null,
                  style: FilledButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    backgroundColor: const Color(0xFF1D4ED8),
                  ),
                  child: const Text('Firmar y presentar (próximo sprint)'),
                ),
              ),
              const SizedBox(height: 16),
              Container(
                padding: const EdgeInsets.all(14),
                decoration: BoxDecoration(
                  color: const Color(0xFFF8FAFC),
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: const Color(0xFFE2E8F0)),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('Trazabilidad del pre-llenado',
                        style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13)),
                    const SizedBox(height: 6),
                    const Text(
                      'Para construir esta declaración, Hacienda consultó al Registro Civil sus datos '
                      'básicos. Esta consulta fue firmada digitalmente por el security-server de Hacienda y '
                      'verificada por el security-server de Registro Civil antes de devolver los datos. '
                      'Quedó registrada en su bitácora.',
                      style: TextStyle(fontSize: 11, color: Colors.black54),
                    ),
                    const SizedBox(height: 8),
                    Text(t.onceOnlyTrace,
                        style: const TextStyle(
                          fontSize: 10,
                          fontFamily: 'monospace',
                          color: Color(0xFF64748B),
                        )),
                  ],
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}

Widget _section(String title, List<Widget> rows) {
  return Container(
    margin: const EdgeInsets.only(bottom: 14),
    padding: const EdgeInsets.all(16),
    decoration: BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(12),
      border: Border.all(color: const Color(0xFFE2E8F0)),
    ),
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
        const SizedBox(height: 10),
        ...rows,
      ],
    ),
  );
}

Widget _row(String label, String value) => Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.black54)),
          Flexible(child: Text(value, textAlign: TextAlign.right)),
        ],
      ),
    );

Widget _rowMono(String label, String value) => Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.black54)),
          Flexible(
            child: Text(value,
                textAlign: TextAlign.right,
                style: const TextStyle(fontFamily: 'monospace')),
          ),
        ],
      ),
    );

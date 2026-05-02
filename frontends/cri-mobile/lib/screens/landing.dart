import 'package:flutter/material.dart';

import '../session.dart';
import 'home.dart';

const _demoCitizens = [
  ('1-1234-5678', 'Sebastián Rojas Mora', 'ingreso 18M'),
  ('1-9876-5432', 'María Castro Fernández', 'ingreso 24.5M, con dependientes'),
  ('2-0001-0001', 'Luis Fernández Solís', 'ingreso 12M'),
];

class LandingScreen extends StatelessWidget {
  const LandingScreen({super.key});

  Future<void> _select(BuildContext context, String cedula) async {
    await Session.setActiveCedula(cedula);
    if (!context.mounted) return;
    Navigator.of(context).pushReplacement(
      MaterialPageRoute(builder: (_) => const HomeScreen()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    width: 44,
                    height: 44,
                    decoration: BoxDecoration(
                      color: const Color(0xFF1D4ED8),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    alignment: Alignment.center,
                    child: const Text('CR', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                  ),
                  const SizedBox(width: 12),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: const [
                      Text('MiCR', style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold)),
                      Text('Carpeta Digital Ciudadana', style: TextStyle(fontSize: 12, color: Colors.black54)),
                    ],
                  ),
                ],
              ),
              const SizedBox(height: 28),
              const Text(
                'Bienvenido a su Carpeta Digital',
                style: TextStyle(fontSize: 22, fontWeight: FontWeight.w600),
              ),
              const SizedBox(height: 8),
              const Text(
                'Acá puede ver sus datos, su declaración pre-llenada de Hacienda, y la bitácora de '
                'accesos del Estado a su información personal — quién consultó sus datos, cuándo y '
                'con qué propósito.',
                style: TextStyle(color: Colors.black54),
              ),
              const SizedBox(height: 28),
              const Text(
                'Demo: seleccione un ciudadano',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
              ),
              const SizedBox(height: 4),
              const Text(
                'Esta carpeta opera sobre el realm "demo". En producción se accede con identidad digital.',
                style: TextStyle(fontSize: 12, color: Colors.black54),
              ),
              const SizedBox(height: 16),
              ..._demoCitizens.map((c) => _CitizenCard(
                    cedula: c.$1,
                    name: c.$2,
                    hint: c.$3,
                    onTap: () => _select(context, c.$1),
                  )),
            ],
          ),
        ),
      ),
    );
  }
}

class _CitizenCard extends StatelessWidget {
  final String cedula;
  final String name;
  final String hint;
  final VoidCallback onTap;
  const _CitizenCard({required this.cedula, required this.name, required this.hint, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 0,
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: const BorderSide(color: Color(0xFFE2E8F0)),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(name, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
              const SizedBox(height: 4),
              Text(cedula,
                  style: const TextStyle(
                    fontSize: 12,
                    fontFamily: 'monospace',
                    color: Colors.black54,
                  )),
              const SizedBox(height: 8),
              Text(hint, style: const TextStyle(fontSize: 12, color: Color(0xFF94A3B8))),
            ],
          ),
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';

import '../session.dart';
import '../theme.dart';
import 'home.dart';

const _demoCitizens = [
  ('1-1234-5678', 'Sebastián Rojas Mora', 'ingreso 18M'),
  ('1-9876-5432', 'María Castro Fernández', 'ingreso 24.5M, con dependientes'),
  ('2-0001-0001', 'Luis Fernández Solís', 'ingreso 12M'),
];

const _guarantees = [
  (
    '01',
    'Once-only',
    'El ciudadano da un dato al sistema una sola vez. Hacienda consulta a Registro Civil vía interop, no le pide al ciudadano.',
    CrColors.areaIduc,
  ),
  (
    '02',
    'Descentralización',
    'Cada institución es dueña de sus datos. No hay base de datos central; los servicios se exponen y consumen vía X-Road.',
    CrColors.crBlueBright,
  ),
  (
    '03',
    'Trazabilidad ciudadana',
    'El ciudadano ve quién consultó sus datos, cuándo y con qué propósito. Bitácora hash-chained, verificable.',
    CrColors.areaPlatform,
  ),
  (
    '04',
    'Firma digital',
    'Equivalente a firma manuscrita. Claves Ed25519 custodiadas en KMS, firma server-side, sellado de tiempo RFC 3161.',
    CrColors.crRedBright,
  ),
];

class LandingScreen extends StatefulWidget {
  const LandingScreen({super.key});
  @override
  State<LandingScreen> createState() => _LandingScreenState();
}

class _LandingScreenState extends State<LandingScreen> {
  final _demoKey = GlobalKey();

  Future<void> _select(BuildContext context, String cedula) async {
    await Session.setActiveCedula(cedula);
    if (!context.mounted) return;
    Navigator.of(context).pushReplacement(
      MaterialPageRoute(builder: (_) => const HomeScreen()),
    );
  }

  void _scrollToDemo() {
    final ctx = _demoKey.currentContext;
    if (ctx == null) return;
    Scrollable.ensureVisible(
      ctx,
      duration: const Duration(milliseconds: 500),
      curve: Curves.easeOutCubic,
      alignment: 0.0,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: GlowBackground(
        child: SafeArea(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const SizedBox(height: 8),
                _hero(),
                const SizedBox(height: 28),
                _guaranteesSection(),
                const SizedBox(height: 28),
                _demoSection(),
                const SizedBox(height: 24),
                _footer(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _hero() {
    const titleStyle = TextStyle(
      fontSize: 38,
      fontWeight: FontWeight.w600,
      height: 1.08,
      letterSpacing: -0.5,
      color: CrColors.text,
    );
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Eyebrow('costaricaasservice · núcleo del ecosistema'),
        const SizedBox(height: 28),
        const Text('Gobierno digital,', style: titleStyle),
        const GradientText('como producto.', style: titleStyle),
        const SizedBox(height: 20),
        Text.rich(
          TextSpan(
            style: const TextStyle(
              fontSize: 15,
              color: CrColors.muted,
              height: 1.5,
            ),
            children: const [
              TextSpan(text: 'Estonia tardó 25 años. Nosotros, 3 meses. '),
              TextSpan(
                text: 'Las instituciones se hablan entre ellas',
                style: TextStyle(color: CrColors.text, fontWeight: FontWeight.w600),
              ),
              TextSpan(text: ' y no te piden lo mismo dos veces. '),
              TextSpan(
                text: 'Firmás desde el celular',
                style: TextStyle(color: CrColors.text, fontWeight: FontWeight.w600),
              ),
              TextSpan(text: ' con valor legal. '),
              TextSpan(
                text: 'Ves quién consulta tus datos',
                style: TextStyle(color: CrColors.text, fontWeight: FontWeight.w600),
              ),
              TextSpan(text: ' y por qué.'),
            ],
          ),
        ),
        const SizedBox(height: 24),
        CrPrimaryButton(
          label: 'Probar como ciudadano',
          onPressed: _scrollToDemo,
        ),
      ],
    );
  }

  Widget _guaranteesSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Eyebrow('Bajo el capó', dotColor: CrColors.crCream),
        const SizedBox(height: 14),
        const Text(
          'Cuatro garantías\nno negociables.',
          style: TextStyle(
            fontSize: 26,
            fontWeight: FontWeight.w600,
            height: 1.15,
            letterSpacing: -0.3,
            color: CrColors.text,
          ),
        ),
        const SizedBox(height: 18),
        ..._guarantees.map(
          (g) => Padding(
            padding: const EdgeInsets.only(bottom: 10),
            child: _GuaranteeCard(
              num: g.$1,
              title: g.$2,
              body: g.$3,
              accent: g.$4,
            ),
          ),
        ),
      ],
    );
  }

  Widget _demoSection() {
    return Column(
      key: _demoKey,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Eyebrow('Demo · realm "demo"', dotColor: CrColors.crRedBright),
        const SizedBox(height: 14),
        const Text(
          'Entrá como un ciudadano\nde prueba.',
          style: TextStyle(
            fontSize: 22,
            fontWeight: FontWeight.w600,
            height: 1.2,
            letterSpacing: -0.3,
            color: CrColors.text,
          ),
        ),
        const SizedBox(height: 8),
        const Text(
          'En producción se accede con identidad digital. Acá entramos directo con un perfil mock para mostrar la carpeta.',
          style: TextStyle(fontSize: 13, color: CrColors.muted, height: 1.5),
        ),
        const SizedBox(height: 16),
        ..._demoCitizens.map(
          (c) => _CitizenCard(
            cedula: c.$1,
            name: c.$2,
            hint: c.$3,
            onTap: () => _select(context, c.$1),
          ),
        ),
      ],
    );
  }

  Widget _footer() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: const [
        Divider(color: CrColors.border, height: 1),
        SizedBox(height: 16),
        Text(
          'Iniciativa privada inspirada en e-Estonia. No es proyecto oficial del Estado costarricense.',
          style: TextStyle(fontSize: 11, color: CrColors.muted, height: 1.5),
        ),
        SizedBox(height: 6),
        Text(
          'By @devsebas — sebashian961@gmail.com',
          style: TextStyle(
            fontFamily: 'monospace',
            fontSize: 11,
            color: CrColors.muted,
          ),
        ),
        SizedBox(height: 16),
      ],
    );
  }
}

class _GuaranteeCard extends StatelessWidget {
  final String num;
  final String title;
  final String body;
  final Color accent;
  const _GuaranteeCard({
    required this.num,
    required this.title,
    required this.body,
    required this.accent,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
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
              Text(
                num,
                style: TextStyle(
                  fontFamily: 'monospace',
                  fontSize: 12,
                  letterSpacing: 1.6,
                  color: accent,
                  fontWeight: FontWeight.w600,
                ),
              ),
              Container(
                width: 8,
                height: 8,
                decoration: BoxDecoration(color: accent, shape: BoxShape.circle),
              ),
            ],
          ),
          const SizedBox(height: 14),
          Text(
            title,
            style: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w600,
              color: CrColors.text,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            body,
            style: const TextStyle(
              fontSize: 13,
              color: CrColors.muted,
              height: 1.5,
            ),
          ),
        ],
      ),
    );
  }
}

class _CitizenCard extends StatelessWidget {
  final String cedula;
  final String name;
  final String hint;
  final VoidCallback onTap;
  const _CitizenCard({
    required this.cedula,
    required this.name,
    required this.hint,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      decoration: BoxDecoration(
        color: CrColors.surface,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: CrColors.border),
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(12),
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        name,
                        style: const TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w600,
                          color: CrColors.text,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        cedula,
                        style: const TextStyle(
                          fontSize: 12,
                          fontFamily: 'monospace',
                          color: CrColors.muted,
                        ),
                      ),
                      const SizedBox(height: 6),
                      Text(
                        hint,
                        style: const TextStyle(
                          fontSize: 12,
                          color: CrColors.muted,
                        ),
                      ),
                    ],
                  ),
                ),
                const Icon(
                  Icons.arrow_forward,
                  size: 16,
                  color: CrColors.crBlueBright,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

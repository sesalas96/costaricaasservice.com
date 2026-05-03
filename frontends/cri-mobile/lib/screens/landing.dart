import 'package:flutter/material.dart';

import '../session.dart';
import '../theme.dart';
import 'home.dart';

const _demoCitizens = [
  ('1-1234-5678', 'Sebastián Rojas Mora', 'ingreso 18M'),
  ('1-9876-5432', 'María Castro Fernández', 'ingreso 24.5M, con dependientes'),
  ('2-0001-0001', 'Luis Fernández Solís', 'ingreso 12M'),
];

class LandingScreen extends StatefulWidget {
  const LandingScreen({super.key});
  @override
  State<LandingScreen> createState() => _LandingScreenState();
}

class _LandingScreenState extends State<LandingScreen> {
  final _emailCtrl = TextEditingController();
  final _passCtrl = TextEditingController();
  final _formKey = GlobalKey<FormState>();
  bool _obscure = true;
  bool _submitting = false;

  @override
  void dispose() {
    _emailCtrl.dispose();
    _passCtrl.dispose();
    super.dispose();
  }

  Future<void> _enterAs(String cedula) async {
    await Session.setActiveCedula(cedula);
    if (!mounted) return;
    Navigator.of(context).pushReplacement(
      MaterialPageRoute(builder: (_) => const HomeScreen()),
    );
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _submitting = true);
    await Future.delayed(const Duration(milliseconds: 400));
    if (!mounted) return;
    setState(() => _submitting = false);
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        backgroundColor: CrColors.surface2,
        behavior: SnackBarBehavior.floating,
        margin: const EdgeInsets.all(16),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(10),
          side: const BorderSide(color: CrColors.border),
        ),
        content: const Text(
          'Login en línea aún no está habilitado.\nUsá el menú ⋮ para entrar como ciudadano de prueba.',
          style: TextStyle(color: CrColors.text, fontSize: 13, height: 1.4),
        ),
      ),
    );
  }

  Future<void> _showQuickAccess() async {
    await showModalBottomSheet<void>(
      context: context,
      backgroundColor: CrColors.surface,
      showDragHandle: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(18)),
      ),
      builder: (ctx) {
        return SafeArea(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(20, 4, 20, 20),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Eyebrow('Accesos rápidos · realm "demo"',
                    dotColor: CrColors.crRedBright),
                const SizedBox(height: 12),
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
                const SizedBox(height: 6),
                const Text(
                  'En producción se entra con identidad digital. Acá saltamos directo con un perfil mock.',
                  style: TextStyle(
                    fontSize: 12,
                    color: CrColors.muted,
                    height: 1.5,
                  ),
                ),
                const SizedBox(height: 16),
                ..._demoCitizens.map(
                  (c) => _CitizenCard(
                    cedula: c.$1,
                    name: c.$2,
                    hint: c.$3,
                    onTap: () {
                      Navigator.of(ctx).pop();
                      _enterAs(c.$1);
                    },
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: Row(
          children: [
            Container(
              width: 30,
              height: 30,
              decoration: BoxDecoration(
                color: CrColors.crBlueBright,
                borderRadius: BorderRadius.circular(6),
              ),
              alignment: Alignment.center,
              child: const Text(
                'CR',
                style: TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: 11,
                ),
              ),
            ),
            const SizedBox(width: 10),
            const Text(
              'MiCR',
              style: TextStyle(fontSize: 16, fontWeight: FontWeight.w700),
            ),
          ],
        ),
        actions: [
          PopupMenuButton<String>(
            tooltip: 'Accesos rápidos de prueba',
            icon: const Icon(Icons.more_vert, color: CrColors.text),
            color: CrColors.surface2,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(10),
              side: const BorderSide(color: CrColors.border),
            ),
            onSelected: (value) {
              if (value == '__sheet__') {
                _showQuickAccess();
              } else {
                _enterAs(value);
              }
            },
            itemBuilder: (ctx) => [
              const PopupMenuItem<String>(
                enabled: false,
                height: 32,
                child: Text(
                  'ENTRAR COMO (DEMO)',
                  style: TextStyle(
                    fontFamily: 'monospace',
                    fontSize: 10,
                    letterSpacing: 1.4,
                    color: CrColors.muted,
                  ),
                ),
              ),
              ..._demoCitizens.map(
                (c) => PopupMenuItem<String>(
                  value: c.$1,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        c.$2,
                        style: const TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.w600,
                          color: CrColors.text,
                        ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        c.$1,
                        style: const TextStyle(
                          fontFamily: 'monospace',
                          fontSize: 11,
                          color: CrColors.muted,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              const PopupMenuDivider(),
              const PopupMenuItem<String>(
                value: '__sheet__',
                child: Row(
                  children: [
                    Icon(Icons.tune_rounded,
                        size: 16, color: CrColors.crBlueBright),
                    SizedBox(width: 10),
                    Text(
                      'Ver detalle de perfiles',
                      style: TextStyle(fontSize: 13, color: CrColors.text),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(width: 4),
        ],
      ),
      extendBodyBehindAppBar: true,
      body: GlowBackground(
        child: SafeArea(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 20),
            keyboardDismissBehavior:
                ScrollViewKeyboardDismissBehavior.onDrag,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const SizedBox(height: 32),
                _intro(),
                const SizedBox(height: 28),
                _loginCard(),
                const SizedBox(height: 28),
                _footer(),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _intro() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Eyebrow('MiCR · tu carpeta ciudadana'),
        const SizedBox(height: 22),
        const Text(
          'Identificate.',
          style: TextStyle(
            fontSize: 34,
            fontWeight: FontWeight.w600,
            height: 1.05,
            letterSpacing: -0.5,
            color: CrColors.text,
          ),
        ),
        const SizedBox(height: 10),
        const Text(
          'Una identidad. Todos los servicios. Vos en control de tus datos.',
          style: TextStyle(
            fontSize: 14,
            color: CrColors.muted,
            height: 1.5,
          ),
        ),
      ],
    );
  }

  Widget _loginCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: CrColors.surface,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: CrColors.border),
      ),
      child: Form(
        key: _formKey,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _fieldLabel('Correo'),
            const SizedBox(height: 8),
            TextFormField(
              controller: _emailCtrl,
              keyboardType: TextInputType.emailAddress,
              autocorrect: false,
              enableSuggestions: false,
              textInputAction: TextInputAction.next,
              style: const TextStyle(color: CrColors.text, fontSize: 14),
              decoration: _inputDecoration(
                hint: 'tu@correo.cr',
                icon: Icons.alternate_email_rounded,
              ),
              validator: (v) {
                final s = (v ?? '').trim();
                if (s.isEmpty) return 'Ingresá tu correo';
                if (!s.contains('@') || !s.contains('.')) {
                  return 'Correo no válido';
                }
                return null;
              },
            ),
            const SizedBox(height: 14),
            _fieldLabel('Contraseña'),
            const SizedBox(height: 8),
            TextFormField(
              controller: _passCtrl,
              obscureText: _obscure,
              autocorrect: false,
              enableSuggestions: false,
              textInputAction: TextInputAction.done,
              onFieldSubmitted: (_) => _submit(),
              style: const TextStyle(color: CrColors.text, fontSize: 14),
              decoration: _inputDecoration(
                hint: '••••••••',
                icon: Icons.lock_outline_rounded,
                suffix: IconButton(
                  splashRadius: 18,
                  icon: Icon(
                    _obscure
                        ? Icons.visibility_off_outlined
                        : Icons.visibility_outlined,
                    size: 18,
                    color: CrColors.muted,
                  ),
                  onPressed: () => setState(() => _obscure = !_obscure),
                ),
              ),
              validator: (v) {
                if ((v ?? '').isEmpty) return 'Ingresá tu contraseña';
                if ((v ?? '').length < 6) return 'Mínimo 6 caracteres';
                return null;
              },
            ),
            const SizedBox(height: 18),
            Row(
              children: [
                Expanded(
                  child: CrPrimaryButton(
                    label: _submitting ? 'Verificando…' : 'Ingresar',
                    onPressed: _submitting ? null : _submit,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Center(
              child: TextButton(
                onPressed: () {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Recuperación de contraseña: próximamente'),
                    ),
                  );
                },
                style: TextButton.styleFrom(
                  foregroundColor: CrColors.muted,
                  textStyle: const TextStyle(fontSize: 12),
                ),
                child: const Text('¿Olvidaste tu contraseña?'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _footer() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: const [
        Divider(color: CrColors.border, height: 1),
        SizedBox(height: 14),
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

  Widget _fieldLabel(String label) {
    return Text(
      label.toUpperCase(),
      style: const TextStyle(
        fontFamily: 'monospace',
        fontSize: 10,
        letterSpacing: 1.6,
        color: CrColors.muted,
        fontWeight: FontWeight.w600,
      ),
    );
  }

  InputDecoration _inputDecoration({
    required String hint,
    required IconData icon,
    Widget? suffix,
  }) {
    return InputDecoration(
      hintText: hint,
      hintStyle: const TextStyle(color: CrColors.muted, fontSize: 14),
      prefixIcon: Icon(icon, size: 18, color: CrColors.muted),
      suffixIcon: suffix,
      filled: true,
      fillColor: CrColors.surface2,
      contentPadding:
          const EdgeInsets.symmetric(horizontal: 12, vertical: 14),
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: const BorderSide(color: CrColors.border),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: const BorderSide(color: CrColors.border),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: const BorderSide(color: CrColors.crBlueBright, width: 1.4),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: const BorderSide(color: CrColors.crRedBright),
      ),
      focusedErrorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
        borderSide: const BorderSide(color: CrColors.crRedBright, width: 1.4),
      ),
      errorStyle: const TextStyle(color: CrColors.crRedBright, fontSize: 11),
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
        color: CrColors.surface2,
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

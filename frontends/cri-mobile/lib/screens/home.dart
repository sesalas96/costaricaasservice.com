import 'package:flutter/material.dart';

import '../session.dart';
import '../theme.dart';
import 'access_log_tab.dart';
import 'chat_tab.dart';
import 'landing.dart';
import 'tax_tab.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});
  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _index = 0;
  String? _cedula;

  @override
  void initState() {
    super.initState();
    _loadCedula();
  }

  Future<void> _loadCedula() async {
    final c = await Session.activeCedula();
    if (!mounted) return;
    if (c == null) {
      Navigator.of(context).pushReplacement(MaterialPageRoute(builder: (_) => const LandingScreen()));
      return;
    }
    setState(() => _cedula = c);
  }

  Future<void> _logout() async {
    await Session.clear();
    if (!mounted) return;
    Navigator.of(context).pushReplacement(MaterialPageRoute(builder: (_) => const LandingScreen()));
  }

  Future<void> _switchCitizen() async {
    await Session.clear();
    if (!mounted) return;
    Navigator.of(context).pushReplacement(MaterialPageRoute(builder: (_) => const LandingScreen()));
  }

  String _initials(String cedula) {
    final digits = cedula.replaceAll(RegExp(r'\D'), '');
    if (digits.length >= 2) {
      return digits.substring(0, 2);
    }
    return cedula.substring(0, cedula.length >= 2 ? 2 : 1).toUpperCase();
  }

  @override
  Widget build(BuildContext context) {
    final cedula = _cedula;
    if (cedula == null) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    final tabs = [
      ChatTab(cedula: cedula),
      AccessLogTab(cedula: cedula),
      TaxTab(cedula: cedula),
    ];
    return Scaffold(
      appBar: AppBar(
        backgroundColor: CrColors.bg,
        foregroundColor: CrColors.text,
        titleSpacing: 16,
        title: Row(
          children: [
            ClipRRect(
              borderRadius: BorderRadius.circular(7),
              child: Image.asset(
                'assets/app_icon.png',
                width: 32,
                height: 32,
                fit: BoxFit.cover,
                filterQuality: FilterQuality.medium,
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
          _ProfileButton(
            cedula: cedula,
            initials: _initials(cedula),
            onSwitch: _switchCitizen,
            onLogout: _logout,
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: tabs[_index],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _index,
        onDestinationSelected: (i) => setState(() => _index = i),
        destinations: const [
          NavigationDestination(
            icon: Icon(Icons.forum_outlined),
            selectedIcon: Icon(Icons.forum),
            label: 'Asistente',
          ),
          NavigationDestination(
            icon: Icon(Icons.history),
            selectedIcon: Icon(Icons.history_toggle_off),
            label: 'Bitácora',
          ),
          NavigationDestination(
            icon: Icon(Icons.receipt_long_outlined),
            selectedIcon: Icon(Icons.receipt_long),
            label: 'Declaración',
          ),
        ],
      ),
    );
  }
}

class _ProfileButton extends StatelessWidget {
  final String cedula;
  final String initials;
  final VoidCallback onSwitch;
  final VoidCallback onLogout;

  const _ProfileButton({
    required this.cedula,
    required this.initials,
    required this.onSwitch,
    required this.onLogout,
  });

  @override
  Widget build(BuildContext context) {
    return PopupMenuButton<String>(
      tooltip: 'Mi perfil',
      offset: const Offset(0, 44),
      color: CrColors.surface,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: const BorderSide(color: CrColors.border),
      ),
      onSelected: (v) {
        if (v == 'switch') onSwitch();
        if (v == 'logout') onLogout();
      },
      itemBuilder: (ctx) => [
        PopupMenuItem<String>(
          enabled: false,
          padding: EdgeInsets.zero,
          child: Container(
            width: 240,
            padding: const EdgeInsets.fromLTRB(14, 12, 14, 12),
            decoration: const BoxDecoration(
              border: Border(bottom: BorderSide(color: CrColors.border)),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'CIUDADANO ACTIVO',
                  style: TextStyle(
                    color: CrColors.muted,
                    fontFamily: 'monospace',
                    fontSize: 9.5,
                    letterSpacing: 1.4,
                    fontWeight: FontWeight.w700,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  cedula,
                  style: const TextStyle(
                    color: CrColors.text,
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    fontFamily: 'monospace',
                  ),
                ),
                const SizedBox(height: 2),
                const Text(
                  'realm demo',
                  style: TextStyle(color: CrColors.muted, fontSize: 11),
                ),
              ],
            ),
          ),
        ),
        const PopupMenuItem<String>(
          value: 'switch',
          child: Row(
            children: [
              Icon(Icons.swap_horiz_rounded, size: 16, color: CrColors.text),
              SizedBox(width: 10),
              Text('Cambiar ciudadano', style: TextStyle(color: CrColors.text, fontSize: 13)),
            ],
          ),
        ),
        const PopupMenuItem<String>(
          value: 'logout',
          child: Row(
            children: [
              Icon(Icons.logout_rounded, size: 16, color: CrColors.crRedBright),
              SizedBox(width: 10),
              Text(
                'Cerrar sesión',
                style: TextStyle(color: CrColors.crRedBright, fontSize: 13),
              ),
            ],
          ),
        ),
      ],
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 4),
        child: Container(
          width: 36,
          height: 36,
          decoration: BoxDecoration(
            gradient: const LinearGradient(
              colors: [CrColors.crBlue, CrColors.crBlueBright],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
            shape: BoxShape.circle,
            border: Border.all(color: CrColors.border, width: 1),
          ),
          alignment: Alignment.center,
          child: Text(
            initials,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 12,
              fontWeight: FontWeight.w700,
              fontFamily: 'monospace',
            ),
          ),
        ),
      ),
    );
  }
}

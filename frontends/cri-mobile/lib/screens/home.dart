import 'package:flutter/material.dart';

import '../session.dart';
import 'access_log_tab.dart';
import 'dashboard_tab.dart';
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

  @override
  Widget build(BuildContext context) {
    final cedula = _cedula;
    if (cedula == null) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    final tabs = [
      DashboardTab(cedula: cedula),
      AccessLogTab(cedula: cedula),
      TaxTab(cedula: cedula),
    ];
    return Scaffold(
      appBar: AppBar(
        backgroundColor: const Color(0xFF0F172A),
        foregroundColor: Colors.white,
        title: Row(
          children: [
            Container(
              width: 32,
              height: 32,
              decoration: BoxDecoration(
                color: const Color(0xFF1D4ED8),
                borderRadius: BorderRadius.circular(6),
              ),
              alignment: Alignment.center,
              child: const Text('CR',
                  style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 12)),
            ),
            const SizedBox(width: 10),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text('MiCR', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w700)),
                Text(cedula, style: const TextStyle(fontSize: 11, color: Color(0xFFCBD5E1))),
              ],
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.logout),
            tooltip: 'Cerrar sesión',
            onPressed: _logout,
          ),
        ],
      ),
      body: tabs[_index],
      bottomNavigationBar: NavigationBar(
        selectedIndex: _index,
        onDestinationSelected: (i) => setState(() => _index = i),
        destinations: const [
          NavigationDestination(icon: Icon(Icons.home_outlined), selectedIcon: Icon(Icons.home), label: 'Escritorio'),
          NavigationDestination(icon: Icon(Icons.history), selectedIcon: Icon(Icons.history_toggle_off), label: 'Bitácora'),
          NavigationDestination(icon: Icon(Icons.receipt_long_outlined), selectedIcon: Icon(Icons.receipt_long), label: 'Declaración'),
        ],
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:intl/date_symbol_data_local.dart';

import 'screens/home.dart';
import 'screens/landing.dart';
import 'session.dart';
import 'theme.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await initializeDateFormatting('es_CR', null);
  final cedula = await Session.activeCedula();
  runApp(MyApp(initialCedula: cedula));
}

class MyApp extends StatelessWidget {
  final String? initialCedula;
  const MyApp({super.key, this.initialCedula});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MiCR',
      debugShowCheckedModeBanner: false,
      theme: crTheme(),
      home: initialCedula != null ? const HomeScreen() : const LandingScreen(),
    );
  }
}

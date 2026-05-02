// Sesión-MVP: la cédula del ciudadano activo se guarda en shared_preferences.
// En producción esto vendría de un JWT validado por el gateway.

import 'package:shared_preferences/shared_preferences.dart';

class Session {
  static const _key = 'micr_cedula';

  static Future<String?> activeCedula() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_key);
  }

  static Future<void> setActiveCedula(String cedula) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_key, cedula);
  }

  static Future<void> clear() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_key);
  }
}

import 'package:intl/intl.dart';

final _crc = NumberFormat.currency(locale: 'es_CR', symbol: '₡', decimalDigits: 0);
final _dt = DateFormat("d MMM yyyy 'a las' HH:mm", 'es_CR');

String formatCRC(num n) => _crc.format(n);
String formatDateTime(DateTime t) => _dt.format(t.toLocal());
String shortHash(String h) => h.length > 16 ? '${h.substring(0, 16)}…' : h;

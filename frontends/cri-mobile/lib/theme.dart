import 'package:flutter/material.dart';

class CrColors {
  static const bg = Color(0xFF07090D);
  static const surface = Color(0xFF0F1218);
  static const surface2 = Color(0xFF161A22);
  static const border = Color(0xFF232936);
  static const text = Color(0xFFE7EAF0);
  static const muted = Color(0xFF8B93A3);

  static const crBlue = Color(0xFF002B7F);
  static const crBlueBright = Color(0xFF1D4FC4);
  static const crRed = Color(0xFFCE1126);
  static const crRedBright = Color(0xFFEF3145);
  static const crCream = Color(0xFFF7F5EF);

  static const areaIduc = Color(0xFF10B981);
  static const areaPlatform = Color(0xFFA78BFA);
  static const areaInterop = Color(0xFF38BDF8);
  static const areaMembers = Color(0xFFF59E0B);
  static const areaGateway = Color(0xFFF43F5E);

  static const areaIducDim = Color(0x2610B981);
  static const crBlueBrightDim = Color(0x261D4FC4);
  static const crRedBrightDim = Color(0x26EF3145);
  static const areaPlatformDim = Color(0x26A78BFA);
}

ThemeData crTheme() {
  const scheme = ColorScheme(
    brightness: Brightness.dark,
    primary: CrColors.crBlueBright,
    onPrimary: Colors.white,
    secondary: CrColors.crRedBright,
    onSecondary: Colors.white,
    error: CrColors.crRedBright,
    onError: Colors.white,
    surface: CrColors.surface,
    onSurface: CrColors.text,
    onSurfaceVariant: CrColors.muted,
    outline: CrColors.muted,
    outlineVariant: CrColors.border,
    surfaceContainerHighest: CrColors.surface2,
  );
  return ThemeData(
    useMaterial3: true,
    brightness: Brightness.dark,
    colorScheme: scheme,
    scaffoldBackgroundColor: CrColors.bg,
    appBarTheme: const AppBarTheme(
      backgroundColor: CrColors.bg,
      foregroundColor: CrColors.text,
      surfaceTintColor: Colors.transparent,
      elevation: 0,
    ),
    navigationBarTheme: NavigationBarThemeData(
      backgroundColor: CrColors.surface,
      indicatorColor: CrColors.crBlueBrightDim,
      surfaceTintColor: Colors.transparent,
      labelTextStyle: WidgetStatePropertyAll(
        const TextStyle(fontSize: 11, color: CrColors.text),
      ),
      iconTheme: WidgetStatePropertyAll(
        const IconThemeData(color: CrColors.text),
      ),
    ),
    dividerTheme: const DividerThemeData(color: CrColors.border, space: 1),
    progressIndicatorTheme: const ProgressIndicatorThemeData(
      color: CrColors.crBlueBright,
    ),
  );
}

class Eyebrow extends StatelessWidget {
  final String text;
  final Color dotColor;
  const Eyebrow(this.text, {super.key, this.dotColor = CrColors.crBlueBright});

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(color: dotColor, shape: BoxShape.circle),
        ),
        const SizedBox(width: 10),
        Flexible(
          child: Text(
            text.toUpperCase(),
            style: const TextStyle(
              fontFamily: 'monospace',
              fontSize: 10,
              letterSpacing: 1.6,
              color: CrColors.muted,
            ),
          ),
        ),
      ],
    );
  }
}

class GradientText extends StatelessWidget {
  final String text;
  final TextStyle style;
  const GradientText(this.text, {super.key, required this.style});

  @override
  Widget build(BuildContext context) {
    return ShaderMask(
      blendMode: BlendMode.srcIn,
      shaderCallback: (bounds) => const LinearGradient(
        colors: [
          CrColors.crBlueBright,
          Colors.white,
          CrColors.crRedBright,
        ],
        stops: [0.0, 0.5, 1.0],
        begin: Alignment(-1, 0.1),
        end: Alignment(1, -0.1),
      ).createShader(bounds),
      child: Text(text, style: style),
    );
  }
}

class GlowBackground extends StatelessWidget {
  final Widget child;
  const GlowBackground({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Positioned.fill(child: CustomPaint(painter: _GlowPainter())),
        child,
      ],
    );
  }
}

class _GlowPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final blueCenter = Offset(size.width * 0.85, -size.height * 0.05);
    final blueRadius = size.width * 1.0;
    canvas.drawCircle(
      blueCenter,
      blueRadius,
      Paint()
        ..shader = RadialGradient(
          colors: const [Color(0x381D4FC4), Color(0x001D4FC4)],
          stops: const [0.0, 0.7],
        ).createShader(Rect.fromCircle(center: blueCenter, radius: blueRadius)),
    );
    final redCenter = Offset(size.width * 0.05, size.height * 1.05);
    final redRadius = size.width * 0.9;
    canvas.drawCircle(
      redCenter,
      redRadius,
      Paint()
        ..shader = RadialGradient(
          colors: const [Color(0x29CE1126), Color(0x00CE1126)],
          stops: const [0.0, 0.7],
        ).createShader(Rect.fromCircle(center: redCenter, radius: redRadius)),
    );
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}

class ApiErrorView extends StatefulWidget {
  final Object error;
  final VoidCallback onRetry;
  const ApiErrorView({super.key, required this.error, required this.onRetry});

  @override
  State<ApiErrorView> createState() => _ApiErrorViewState();
}

class _ApiErrorViewState extends State<ApiErrorView> {
  bool _showDetail = false;

  ({IconData icon, String title, String body}) _classify(String raw) {
    final lower = raw.toLowerCase();
    if (lower.contains('connection refused') ||
        lower.contains('errno = 61') ||
        lower.contains('socketexception') ||
        lower.contains('failed host lookup') ||
        lower.contains('network is unreachable')) {
      return (
        icon: Icons.wifi_off_rounded,
        title: 'No se pudo conectar al servidor',
        body:
            'El BFF no respondió. Verificá tu conexión o que el servicio esté arriba en el host correcto.',
      );
    }
    if (lower.contains('timeoutexception') || lower.contains('timed out')) {
      return (
        icon: Icons.timer_off_rounded,
        title: 'La consulta tardó demasiado',
        body: 'No recibimos respuesta a tiempo. Reintentá en un momento.',
      );
    }
    final m5xx = RegExp(r'\b5\d{2}\b').firstMatch(raw);
    if (m5xx != null) {
      return (
        icon: Icons.cloud_off_rounded,
        title: 'Servicio no disponible',
        body:
            'Algo falló del lado del servidor (${m5xx.group(0)}). Reintentá en un momento.',
      );
    }
    final m4xx = RegExp(r'\b4\d{2}\b').firstMatch(raw);
    if (m4xx != null) {
      return (
        icon: Icons.report_gmailerrorred_rounded,
        title: 'Solicitud rechazada',
        body: 'El servidor rechazó la consulta (${m4xx.group(0)}).',
      );
    }
    return (
      icon: Icons.error_outline_rounded,
      title: 'Algo salió mal',
      body: 'No pudimos cargar esta sección.',
    );
  }

  @override
  Widget build(BuildContext context) {
    final raw = widget.error.toString();
    final c = _classify(raw);
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
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: CrColors.crRedBrightDim,
              borderRadius: BorderRadius.circular(8),
            ),
            child: Icon(c.icon, color: CrColors.crRedBright, size: 22),
          ),
          const SizedBox(height: 14),
          Text(
            c.title,
            style: const TextStyle(
              fontSize: 17,
              fontWeight: FontWeight.w600,
              color: CrColors.text,
            ),
          ),
          const SizedBox(height: 6),
          Text(
            c.body,
            style: const TextStyle(
              fontSize: 13,
              color: CrColors.muted,
              height: 1.5,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              FilledButton.icon(
                onPressed: widget.onRetry,
                icon: const Icon(Icons.refresh_rounded, size: 16),
                label: const Text('Reintentar'),
                style: FilledButton.styleFrom(
                  backgroundColor: CrColors.crBlueBright,
                  foregroundColor: Colors.white,
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          InkWell(
            onTap: () => setState(() => _showDetail = !_showDetail),
            borderRadius: BorderRadius.circular(4),
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 6),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    _showDetail
                        ? Icons.keyboard_arrow_down_rounded
                        : Icons.keyboard_arrow_right_rounded,
                    size: 16,
                    color: CrColors.muted,
                  ),
                  const SizedBox(width: 4),
                  const Text(
                    'Detalle técnico',
                    style: TextStyle(
                      fontSize: 11,
                      fontFamily: 'monospace',
                      letterSpacing: 1.2,
                      color: CrColors.muted,
                    ),
                  ),
                ],
              ),
            ),
          ),
          if (_showDetail) ...[
            const SizedBox(height: 4),
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: CrColors.surface2,
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: CrColors.border),
              ),
              child: SelectableText(
                raw,
                style: const TextStyle(
                  fontSize: 11,
                  fontFamily: 'monospace',
                  color: CrColors.muted,
                  height: 1.5,
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }
}

class CrPrimaryButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  const CrPrimaryButton({super.key, required this.label, required this.onPressed});

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [CrColors.crBlue, CrColors.crBlueBright],
          begin: Alignment(-1, -0.3),
          end: Alignment(1, 0.3),
        ),
        borderRadius: BorderRadius.circular(8),
        boxShadow: const [
          BoxShadow(
            color: Color(0x99002B7F),
            blurRadius: 24,
            offset: Offset(0, 8),
            spreadRadius: -8,
          ),
        ],
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onPressed,
          borderRadius: BorderRadius.circular(8),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 14),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  label,
                  style: const TextStyle(
                    color: Colors.white,
                    fontWeight: FontWeight.w600,
                    fontSize: 14,
                  ),
                ),
                const SizedBox(width: 8),
                const Icon(Icons.arrow_forward, size: 16, color: Colors.white),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

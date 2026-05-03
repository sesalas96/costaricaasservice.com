import 'dart:async';

import 'package:flutter/material.dart';

import '../data/quick_prompts.dart';
import '../theme.dart';

sealed class ChatMessage {
  final String id;
  ChatMessage(this.id);
}

class CitizenMsg extends ChatMessage {
  final String text;
  CitizenMsg(super.id, this.text);
}

class LlmMsg extends ChatMessage {
  final String text;
  LlmMsg(super.id, this.text);
}

class SystemStepMsg extends ChatMessage {
  final PromptStep step;
  SystemStepMsg(super.id, this.step);
}

class EmailMsg extends ChatMessage {
  final EmailNotice email;
  EmailMsg(super.id, this.email);
}

class ChatTab extends StatefulWidget {
  final String cedula;
  const ChatTab({super.key, required this.cedula});

  @override
  State<ChatTab> createState() => _ChatTabState();
}

class _ChatTabState extends State<ChatTab> {
  final List<ChatMessage> _messages = [];
  final List<Timer> _timers = [];
  final TextEditingController _composer = TextEditingController();
  final ScrollController _scroll = ScrollController();

  bool _running = false;
  String? _activeCategoryId;
  bool _promptsExpanded = false;

  @override
  void dispose() {
    for (final t in _timers) {
      t.cancel();
    }
    _composer.dispose();
    _scroll.dispose();
    super.dispose();
  }

  void _cancelTimers() {
    for (final t in _timers) {
      t.cancel();
    }
    _timers.clear();
  }

  void _scrollToBottom() {
    if (!_scroll.hasClients) return;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scroll.hasClients) {
        _scroll.animateTo(
          _scroll.position.maxScrollExtent,
          duration: const Duration(milliseconds: 220),
          curve: Curves.easeOut,
        );
      }
    });
  }

  String _seq() => '${DateTime.now().microsecondsSinceEpoch}';

  void _runPrompt(QuickPrompt p) {
    if (_running) return;
    _cancelTimers();

    setState(() {
      _messages.add(CitizenMsg('${p.id}-c-${_seq()}', p.citizenAsk));
      _messages.add(LlmMsg('${p.id}-i-${_seq()}', p.llmIntro));
      _running = true;
      _promptsExpanded = false;
    });
    _scrollToBottom();

    var cumulative = 600;
    for (var i = 0; i < p.steps.length; i++) {
      final step = p.steps[i];
      _timers.add(
        Timer(Duration(milliseconds: cumulative), () {
          if (!mounted) return;
          setState(() {
            _messages.add(SystemStepMsg('${p.id}-s$i-${_seq()}', step));
          });
          _scrollToBottom();
        }),
      );
      cumulative += step.delayMs;
    }

    final summaryDelay = cumulative + 500;
    _timers.add(
      Timer(Duration(milliseconds: summaryDelay), () {
        if (!mounted) return;
        setState(() {
          _messages.add(LlmMsg('${p.id}-fs-${_seq()}', p.finalSummary));
        });
        _scrollToBottom();
      }),
    );

    var endDelay = summaryDelay + 400;
    if (p.email != null) {
      _timers.add(
        Timer(Duration(milliseconds: endDelay), () {
          if (!mounted) return;
          setState(() {
            _messages.add(EmailMsg('${p.id}-em-${_seq()}', p.email!));
          });
          _scrollToBottom();
        }),
      );
      endDelay += 400;
    }

    _timers.add(
      Timer(Duration(milliseconds: endDelay), () {
        if (!mounted) return;
        setState(() => _running = false);
      }),
    );
  }

  void _sendFreeForm() {
    final text = _composer.text.trim();
    if (text.isEmpty || _running) return;
    _composer.clear();
    _runPrompt(buildFreeFormFallback(text));
  }

  void _resetChat() {
    _cancelTimers();
    setState(() {
      _messages.clear();
      _running = false;
      _activeCategoryId = null;
      _promptsExpanded = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Expanded(
          child: _messages.isEmpty
              ? _EmptyState(
                  onPickCategory: (id) => setState(() {
                    _activeCategoryId = id;
                    _promptsExpanded = true;
                  }),
                )
              : ListView.builder(
                  controller: _scroll,
                  padding: const EdgeInsets.fromLTRB(14, 12, 14, 8),
                  itemCount: _messages.length + (_running ? 1 : 0),
                  itemBuilder: (ctx, i) {
                    if (i == _messages.length) {
                      return const _TypingDots();
                    }
                    final m = _messages[i];
                    return switch (m) {
                      CitizenMsg c => _CitizenBubble(text: c.text),
                      LlmMsg l => _LlmBubble(text: l.text),
                      SystemStepMsg s => _StepRow(step: s.step),
                      EmailMsg e => _EmailCard(email: e.email),
                    };
                  },
                ),
        ),
        _PromptsPanel(
          activeCategoryId: _activeCategoryId,
          expanded: _promptsExpanded,
          running: _running,
          onCategoryTap: (id) => setState(() {
            if (_activeCategoryId == id && _promptsExpanded) {
              _promptsExpanded = false;
            } else {
              _activeCategoryId = id;
              _promptsExpanded = true;
            }
          }),
          onPromptTap: _runPrompt,
        ),
        _Composer(
          controller: _composer,
          enabled: !_running,
          onSend: _sendFreeForm,
          onReset: _messages.isEmpty ? null : _resetChat,
        ),
      ],
    );
  }
}

class _EmptyState extends StatelessWidget {
  final void Function(String) onPickCategory;
  const _EmptyState({required this.onPickCategory});

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.fromLTRB(20, 24, 20, 16),
      children: [
        const Eyebrow('Asistente ciudadano'),
        const SizedBox(height: 12),
        const Text(
          'Resolvé cualquier consulta del Estado',
          style: TextStyle(
            color: CrColors.text,
            fontSize: 22,
            fontWeight: FontWeight.w700,
            height: 1.2,
          ),
        ),
        const SizedBox(height: 8),
        const Text(
          'Tocá una categoría para ver consultas frecuentes, o escribí tu pregunta abajo. Cada consulta se firma, se enruta por Conecta CR y queda en tu bitácora.',
          style: TextStyle(color: CrColors.muted, fontSize: 13, height: 1.5),
        ),
        const SizedBox(height: 22),
        Wrap(
          spacing: 10,
          runSpacing: 10,
          children: promptCategories
              .map(
                (c) => InkWell(
                  borderRadius: BorderRadius.circular(10),
                  onTap: () => onPickCategory(c.id),
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 12),
                    decoration: BoxDecoration(
                      color: CrColors.surface,
                      borderRadius: BorderRadius.circular(10),
                      border: Border.all(color: CrColors.border),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(c.icon, size: 16, color: c.accent),
                        const SizedBox(width: 8),
                        Text(
                          c.label,
                          style: const TextStyle(
                            color: CrColors.text,
                            fontSize: 13,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              )
              .toList(),
        ),
      ],
    );
  }
}

class _CitizenBubble extends StatelessWidget {
  final String text;
  const _CitizenBubble({required this.text});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          Flexible(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 320),
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                decoration: BoxDecoration(
                  gradient: const LinearGradient(
                    colors: [CrColors.crBlue, CrColors.crBlueBright],
                    begin: Alignment(-1, -0.2),
                    end: Alignment(1, 0.2),
                  ),
                  borderRadius: const BorderRadius.only(
                    topLeft: Radius.circular(16),
                    topRight: Radius.circular(16),
                    bottomLeft: Radius.circular(16),
                    bottomRight: Radius.circular(4),
                  ),
                ),
                child: Text(
                  text,
                  style: const TextStyle(color: Colors.white, fontSize: 14, height: 1.35),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _LlmBubble extends StatelessWidget {
  final String text;
  const _LlmBubble({required this.text});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 28,
            height: 28,
            decoration: BoxDecoration(
              gradient: const LinearGradient(
                colors: [CrColors.crBlue, CrColors.crBlueBright],
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
              ),
              borderRadius: BorderRadius.circular(14),
            ),
            alignment: Alignment.center,
            child: const Text(
              'cr',
              style: TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.w700),
            ),
          ),
          const SizedBox(width: 8),
          Flexible(
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 320),
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                decoration: BoxDecoration(
                  color: CrColors.surface2,
                  borderRadius: const BorderRadius.only(
                    topLeft: Radius.circular(4),
                    topRight: Radius.circular(16),
                    bottomLeft: Radius.circular(16),
                    bottomRight: Radius.circular(16),
                  ),
                  border: Border.all(color: CrColors.border),
                ),
                child: Text(
                  text,
                  style: const TextStyle(color: CrColors.text, fontSize: 14, height: 1.4),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _StepRow extends StatelessWidget {
  final PromptStep step;
  const _StepRow({required this.step});

  String _kindLabel(StepKind k) => switch (k) {
        StepKind.thinking => 'thinking',
        StepKind.request => 'request',
        StepKind.response => 'response',
        StepKind.audit => 'audit',
        StepKind.notify => 'notify',
      };

  @override
  Widget build(BuildContext context) {
    final color = step.color;
    return Padding(
      padding: const EdgeInsets.only(bottom: 10, left: 4, right: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            margin: const EdgeInsets.only(top: 2),
            width: 48,
            height: 28,
            alignment: Alignment.center,
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.15),
              borderRadius: BorderRadius.circular(6),
              border: Border.all(color: color.withValues(alpha: 0.4)),
            ),
            child: Text(
              step.actorTag,
              style: TextStyle(
                color: color,
                fontFamily: 'monospace',
                fontSize: 10,
                fontWeight: FontWeight.w700,
                letterSpacing: 0.6,
              ),
            ),
          ),
          const SizedBox(width: 10),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Wrap(
                  crossAxisAlignment: WrapCrossAlignment.center,
                  spacing: 8,
                  runSpacing: 4,
                  children: [
                    Text(
                      step.actor,
                      style: const TextStyle(
                        color: CrColors.text,
                        fontSize: 13,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
                      decoration: BoxDecoration(
                        color: color.withValues(alpha: 0.12),
                        borderRadius: BorderRadius.circular(3),
                        border: Border.all(color: color.withValues(alpha: 0.35)),
                      ),
                      child: Text(
                        _kindLabel(step.kind),
                        style: TextStyle(
                          color: color,
                          fontFamily: 'monospace',
                          fontSize: 9,
                          letterSpacing: 1.4,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                    if (step.endpoint != null)
                      Text(
                        step.endpoint!,
                        style: const TextStyle(
                          color: CrColors.muted,
                          fontFamily: 'monospace',
                          fontSize: 10,
                        ),
                      ),
                  ],
                ),
                const SizedBox(height: 4),
                Text(
                  step.text,
                  style: const TextStyle(color: CrColors.muted, fontSize: 12.5, height: 1.45),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _EmailCard extends StatelessWidget {
  final EmailNotice email;
  const _EmailCard({required this.email});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: Container(
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: CrColors.bg,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: CrColors.border),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.mail_outline, size: 14, color: CrColors.areaPlatform),
                const SizedBox(width: 6),
                Text(
                  'NOTIFICACIÓN ELECTRÓNICA',
                  style: TextStyle(
                    color: CrColors.areaPlatform,
                    fontFamily: 'monospace',
                    fontSize: 9.5,
                    letterSpacing: 1.4,
                    fontWeight: FontWeight.w700,
                  ),
                ),
                const Spacer(),
                const Text(
                  'tu@correo.cr',
                  style: TextStyle(
                    color: CrColors.muted,
                    fontFamily: 'monospace',
                    fontSize: 10,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              email.subject,
              style: const TextStyle(
                color: CrColors.text,
                fontSize: 14,
                fontWeight: FontWeight.w600,
                height: 1.3,
              ),
            ),
            const SizedBox(height: 4),
            Text(
              email.preview,
              style: const TextStyle(color: CrColors.muted, fontSize: 12.5, height: 1.5),
            ),
          ],
        ),
      ),
    );
  }
}

class _TypingDots extends StatefulWidget {
  const _TypingDots();
  @override
  State<_TypingDots> createState() => _TypingDotsState();
}

class _TypingDotsState extends State<_TypingDots> with SingleTickerProviderStateMixin {
  late final AnimationController _ctrl;

  @override
  void initState() {
    super.initState();
    _ctrl = AnimationController(vsync: this, duration: const Duration(milliseconds: 900))
      ..repeat();
  }

  @override
  void dispose() {
    _ctrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 8, bottom: 12),
      child: AnimatedBuilder(
        animation: _ctrl,
        builder: (ctx, _) {
          return Row(
            mainAxisSize: MainAxisSize.min,
            children: List.generate(3, (i) {
              final phase = (_ctrl.value + i * 0.18) % 1.0;
              final scale = 0.6 + (phase < 0.5 ? phase * 2 : (1 - phase) * 2) * 0.5;
              return Padding(
                padding: const EdgeInsets.symmetric(horizontal: 2),
                child: Transform.scale(
                  scale: scale,
                  child: Container(
                    width: 6,
                    height: 6,
                    decoration: const BoxDecoration(
                      color: CrColors.crBlueBright,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
              );
            }),
          );
        },
      ),
    );
  }
}

class _PromptsPanel extends StatelessWidget {
  final String? activeCategoryId;
  final bool expanded;
  final bool running;
  final void Function(String) onCategoryTap;
  final void Function(QuickPrompt) onPromptTap;

  const _PromptsPanel({
    required this.activeCategoryId,
    required this.expanded,
    required this.running,
    required this.onCategoryTap,
    required this.onPromptTap,
  });

  @override
  Widget build(BuildContext context) {
    final activeCat = activeCategoryId == null
        ? null
        : promptCategories.firstWhere(
            (c) => c.id == activeCategoryId,
            orElse: () => promptCategories.first,
          );

    return Container(
      decoration: const BoxDecoration(
        color: CrColors.surface,
        border: Border(top: BorderSide(color: CrColors.border)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            height: 44,
            child: ListView.separated(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              itemCount: promptCategories.length,
              separatorBuilder: (_, _) => const SizedBox(width: 8),
              itemBuilder: (ctx, i) {
                final c = promptCategories[i];
                final isActive = c.id == activeCategoryId && expanded;
                return _CategoryChip(
                  category: c,
                  active: isActive,
                  onTap: () => onCategoryTap(c.id),
                );
              },
            ),
          ),
          AnimatedSize(
            duration: const Duration(milliseconds: 200),
            curve: Curves.easeOut,
            child: expanded && activeCat != null
                ? Container(
                    width: double.infinity,
                    padding: const EdgeInsets.fromLTRB(12, 4, 12, 10),
                    decoration: const BoxDecoration(
                      border: Border(top: BorderSide(color: CrColors.border)),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.fromLTRB(4, 8, 4, 6),
                          child: Row(
                            children: [
                              Icon(activeCat.icon, size: 14, color: activeCat.accent),
                              const SizedBox(width: 6),
                              Text(
                                activeCat.label.toUpperCase(),
                                style: TextStyle(
                                  color: activeCat.accent,
                                  fontFamily: 'monospace',
                                  fontSize: 10,
                                  letterSpacing: 1.4,
                                  fontWeight: FontWeight.w700,
                                ),
                              ),
                            ],
                          ),
                        ),
                        ...activeCat.prompts.map(
                          (p) => Padding(
                            padding: const EdgeInsets.only(bottom: 6),
                            child: InkWell(
                              borderRadius: BorderRadius.circular(8),
                              onTap: running ? null : () => onPromptTap(p),
                              child: Container(
                                width: double.infinity,
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 12,
                                  vertical: 10,
                                ),
                                decoration: BoxDecoration(
                                  color: CrColors.surface2,
                                  borderRadius: BorderRadius.circular(8),
                                  border: Border.all(color: CrColors.border),
                                ),
                                child: Row(
                                  children: [
                                    Expanded(
                                      child: Text(
                                        p.label,
                                        style: TextStyle(
                                          color: running ? CrColors.muted : CrColors.text,
                                          fontSize: 13,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ),
                                    Icon(
                                      Icons.arrow_forward_rounded,
                                      size: 14,
                                      color: activeCat.accent,
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  )
                : const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }
}

class _CategoryChip extends StatelessWidget {
  final PromptCategory category;
  final bool active;
  final VoidCallback onTap;
  const _CategoryChip({
    required this.category,
    required this.active,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final color = category.accent;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(20),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
        decoration: BoxDecoration(
          color: active ? color.withValues(alpha: 0.18) : CrColors.surface2,
          borderRadius: BorderRadius.circular(20),
          border: Border.all(
            color: active ? color : CrColors.border,
          ),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(category.icon, size: 14, color: active ? color : CrColors.muted),
            const SizedBox(width: 6),
            Text(
              category.label,
              style: TextStyle(
                color: active ? color : CrColors.text,
                fontSize: 12.5,
                fontWeight: active ? FontWeight.w600 : FontWeight.w500,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _Composer extends StatelessWidget {
  final TextEditingController controller;
  final bool enabled;
  final VoidCallback onSend;
  final VoidCallback? onReset;

  const _Composer({
    required this.controller,
    required this.enabled,
    required this.onSend,
    required this.onReset,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.fromLTRB(12, 8, 12, 10),
      decoration: const BoxDecoration(
        color: CrColors.surface,
        border: Border(top: BorderSide(color: CrColors.border)),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          if (onReset != null)
            IconButton(
              onPressed: onReset,
              icon: const Icon(Icons.refresh_rounded, size: 18, color: CrColors.muted),
              tooltip: 'Reiniciar conversación',
              padding: const EdgeInsets.all(8),
              constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
            ),
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                color: CrColors.surface2,
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: CrColors.border),
              ),
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 4),
              child: TextField(
                controller: controller,
                enabled: enabled,
                minLines: 1,
                maxLines: 4,
                onSubmitted: (_) => onSend(),
                textInputAction: TextInputAction.send,
                style: const TextStyle(color: CrColors.text, fontSize: 14),
                decoration: const InputDecoration(
                  hintText: 'Escribí tu consulta…',
                  hintStyle: TextStyle(color: CrColors.muted, fontSize: 14),
                  border: InputBorder.none,
                  isCollapsed: true,
                  contentPadding: EdgeInsets.symmetric(vertical: 10),
                ),
              ),
            ),
          ),
          const SizedBox(width: 8),
          Material(
            color: enabled ? CrColors.crBlueBright : CrColors.surface2,
            shape: const CircleBorder(),
            child: InkWell(
              customBorder: const CircleBorder(),
              onTap: enabled ? onSend : null,
              child: SizedBox(
                width: 40,
                height: 40,
                child: Icon(
                  Icons.send_rounded,
                  size: 18,
                  color: enabled ? Colors.white : CrColors.muted,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

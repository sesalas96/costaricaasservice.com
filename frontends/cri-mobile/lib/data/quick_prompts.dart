import 'package:flutter/material.dart';

import '../l10n/app_localizations.dart';
import '../theme.dart';

enum StepKind { thinking, request, response, audit, notify }

class PromptStep {
  final String actor;
  final String actorTag;
  final Color color;
  final StepKind kind;
  final String text;
  final List<({String k, String v})>? data;
  final String? endpoint;
  final int delayMs;

  const PromptStep({
    required this.actor,
    required this.actorTag,
    required this.color,
    required this.kind,
    required this.text,
    this.data,
    this.endpoint,
    this.delayMs = 900,
  });
}

class EmailNotice {
  final String subject;
  final String preview;
  const EmailNotice({required this.subject, required this.preview});
}

class QuickPrompt {
  final String id;
  final String label;
  final String citizenAsk;
  final String llmIntro;
  final List<PromptStep> steps;
  final String finalSummary;
  final EmailNotice? email;

  const QuickPrompt({
    required this.id,
    required this.label,
    required this.citizenAsk,
    required this.llmIntro,
    required this.steps,
    required this.finalSummary,
    this.email,
  });
}

class PromptCategory {
  final String id;
  final String label;
  final IconData icon;
  final Color accent;
  final List<QuickPrompt> prompts;

  const PromptCategory({
    required this.id,
    required this.label,
    required this.icon,
    required this.accent,
    required this.prompts,
  });
}

const _categories = <PromptCategory>[
  PromptCategory(
    id: 'tributos',
    label: 'Tributos',
    icon: Icons.receipt_long_outlined,
    accent: CrColors.crBlueBright,
    prompts: [
      QuickPrompt(
        id: 'declaracion-2025',
        label: '¿Cuánto debo de renta este año?',
        citizenAsk: '¿Cuánto debo de renta este año?',
        llmIntro:
            'Voy a consultar tu declaración pre-llenada en Hacienda. Para armarla, ellos consultan a Registro Civil y CCSS por interop — vos no tenés que repetir nada.',
        steps: [
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.request,
            text: 'Consulta tu declaración 2025 firmada con tu IDUC.',
            endpoint: 'GET /tax/prefilled/{cedula}/2025',
          ),
          PromptStep(
            actor: 'Conecta CR',
            actorTag: 'X-RD',
            color: CrColors.areaInterop,
            kind: StepKind.audit,
            text: 'Registra la consulta en la bitácora hash-chained.',
            endpoint: 'cri-svc-interop-audit',
          ),
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'periodo', v: '2025'),
              (k: 'ingreso_bruto', v: '₡18.000.000'),
              (k: 'retenido', v: '₡2.700.000'),
              (k: 'deducciones', v: '₡540.000'),
              (k: 'a_pagar', v: '₡180.000'),
              (k: 'vence', v: '15-mar-2026'),
            ],
          ),
        ],
        finalSummary:
            'Te quedan ₡180.000 por pagar antes del 15 de marzo de 2026. Podés firmar el pago directo desde MiCR cuando estés listo.',
        email: EmailNotice(
          subject: 'Tu declaración 2025 está lista para firmar',
          preview:
              'Resumen: a pagar ₡180.000. Fecha límite 15-mar-2026. Abrí MiCR para firmar.',
        ),
      ),
      QuickPrompt(
        id: 'historial-renta',
        label: 'Mostrame mi historial de renta',
        citizenAsk: 'Mostrame mi historial de renta de los últimos 3 años.',
        llmIntro: 'Pido tus declaraciones cerradas de 2022, 2023 y 2024 a Hacienda.',
        steps: [
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.request,
            text: 'Lista declaraciones cerradas del ciudadano.',
            endpoint: 'GET /tax/history/{cedula}',
          ),
          PromptStep(
            actor: 'Conecta CR',
            actorTag: 'X-RD',
            color: CrColors.areaInterop,
            kind: StepKind.audit,
            text: 'Acceso firmado y registrado.',
          ),
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.response,
            text: 'Todas presentadas a tiempo, sin moras.',
            data: [
              (k: '2024', v: '₡220.000'),
              (k: '2023', v: '₡165.000'),
              (k: '2022', v: '₡98.000'),
            ],
          ),
        ],
        finalSummary:
            'Tu historial está al día — sin moras ni declaraciones pendientes en los últimos 3 años.',
      ),
      QuickPrompt(
        id: 'patentes-municipal',
        label: '¿Tengo deuda con la Muni?',
        citizenAsk: '¿Tengo alguna deuda municipal pendiente?',
        llmIntro: 'Consulto a tu municipalidad por servicios e impuestos pendientes.',
        steps: [
          PromptStep(
            actor: 'Municipalidad',
            actorTag: 'MUNI',
            color: CrColors.areaMembers,
            kind: StepKind.request,
            text: 'Consulta deuda por bienes inmuebles, basura y patente.',
            endpoint: 'GET /muni/account/{cedula}',
          ),
          PromptStep(
            actor: 'Conecta CR',
            actorTag: 'X-RD',
            color: CrColors.areaInterop,
            kind: StepKind.audit,
            text: 'Acceso registrado.',
          ),
          PromptStep(
            actor: 'Municipalidad',
            actorTag: 'MUNI',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: 'Saldo al día.',
            data: [
              (k: 'bienes_inmuebles', v: '₡0 vencido'),
              (k: 'basura', v: '₡0 vencido'),
              (k: 'patente', v: '₡0 vencido'),
              (k: 'próximo_cobro', v: '31-mar-2026 · ₡42.300'),
            ],
          ),
        ],
        finalSummary:
            'No tenés deuda con la municipalidad. Te avisaré antes del próximo cobro.',
      ),
    ],
  ),
  PromptCategory(
    id: 'salud',
    label: 'Salud',
    icon: Icons.favorite_outline_rounded,
    accent: CrColors.crRedBright,
    prompts: [
      QuickPrompt(
        id: 'recetas-activas',
        label: '¿Qué recetas tengo activas?',
        citizenAsk: '¿Qué recetas tengo activas?',
        llmIntro: 'Consulto tu expediente en CCSS.',
        steps: [
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.request,
            text: 'Lista recetas activas del ciudadano.',
            endpoint: 'GET /health/prescriptions/{cedula}',
          ),
          PromptStep(
            actor: 'Conecta CR',
            actorTag: 'X-RD',
            color: CrColors.areaInterop,
            kind: StepKind.audit,
            text: 'Acceso registrado.',
          ),
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'losartán_50mg', v: '1× día (mañana)'),
              (k: 'atorvastatina_20mg', v: '1× día (noche)'),
              (k: 'próximo_retiro', v: '03-may-2026'),
            ],
          ),
        ],
        finalSummary:
            'Recordá que el retiro mensual es la primera semana en la farmacia más cercana.',
      ),
      QuickPrompt(
        id: 'proxima-cita',
        label: '¿Cuándo es mi próxima cita?',
        citizenAsk: '¿Cuándo es mi próxima cita médica?',
        llmIntro: 'Consulto tu agenda médica en CCSS.',
        steps: [
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.request,
            text: 'Próxima cita programada del ciudadano.',
            endpoint: 'GET /health/appointments/next/{cedula}',
          ),
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'especialidad', v: 'Cardiología'),
              (k: 'hospital', v: 'Hospital México'),
              (k: 'fecha', v: '18-may-2026 · 08:30'),
              (k: 'preparativos', v: 'Ayuno 8h'),
            ],
          ),
        ],
        finalSummary:
            'Te recordaré 24h antes y te enviaré los preparativos por correo.',
        email: EmailNotice(
          subject: 'Recordatorio: cita 18-may-2026 8:30 AM',
          preview: 'Cardiología · Hospital México. Llegá 30 min antes.',
        ),
      ),
      QuickPrompt(
        id: 'incapacidad',
        label: 'Quiero solicitar una incapacidad',
        citizenAsk: 'Quiero solicitar una incapacidad por enfermedad.',
        llmIntro:
            'Para tramitar incapacidad necesito confirmar tu patrono activo (MTSS) y tu médico tratante (CCSS). Inicio el flujo.',
        steps: [
          PromptStep(
            actor: 'MTSS',
            actorTag: 'MTSS',
            color: CrColors.areaMembers,
            kind: StepKind.request,
            text: 'Verifica patrono activo del ciudadano.',
            endpoint: 'GET /labor/employer/{cedula}',
          ),
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.request,
            text: 'Reserva cita para evaluación de incapacidad.',
            endpoint: 'POST /health/incapacity/request',
          ),
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'cita_asignada', v: '04-may-2026 · 10:00'),
              (k: 'lugar', v: 'EBAIS Curridabat'),
              (k: 'evaluación', v: 'Médico tratante'),
            ],
          ),
        ],
        finalSummary:
            'Tu solicitud está abierta. Tras la consulta, la incapacidad se notifica automáticamente a tu patrono — vos no entregás papeles.',
      ),
      QuickPrompt(
        id: 'vacunas',
        label: 'Mostrame mi carné de vacunas',
        citizenAsk: 'Mostrame mi carné de vacunas.',
        llmIntro: 'Consulto MINSA y CCSS.',
        steps: [
          PromptStep(
            actor: 'MINSA',
            actorTag: 'MINSA',
            color: CrColors.crRedBright,
            kind: StepKind.request,
            text: 'Carné nacional de vacunación.',
            endpoint: 'GET /vacunas/{cedula}',
          ),
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.response,
            text: 'Esquema completo.',
            data: [
              (k: 'covid_19', v: '3 dosis · al día'),
              (k: 'influenza', v: '2025 · al día'),
              (k: 'td', v: '2022 · refuerzo en 2027'),
            ],
          ),
        ],
        finalSummary:
            'Tu carné está actualizado. Te aviso cuando toque el refuerzo Td.',
      ),
    ],
  ),
  PromptCategory(
    id: 'identidad',
    label: 'Identidad',
    icon: Icons.badge_outlined,
    accent: CrColors.areaIduc,
    prompts: [
      QuickPrompt(
        id: 'mis-datos',
        label: '¿Qué datos míos tiene el Estado?',
        citizenAsk: '¿Qué datos míos tiene el Estado hoy?',
        llmIntro:
            'Consulto Registro Civil — la fuente oficial de identidad. El resto de instituciones se conectan a este registro vía interop.',
        steps: [
          PromptStep(
            actor: 'Registro Civil',
            actorTag: 'RC',
            color: CrColors.areaIduc,
            kind: StepKind.request,
            text: 'Datos del ciudadano firmados.',
            endpoint: 'GET /persons/{cedula}',
          ),
          PromptStep(
            actor: 'Conecta CR',
            actorTag: 'X-RD',
            color: CrColors.areaInterop,
            kind: StepKind.audit,
            text: 'Acceso registrado.',
          ),
          PromptStep(
            actor: 'Registro Civil',
            actorTag: 'RC',
            color: CrColors.areaIduc,
            kind: StepKind.response,
            text: 'Datos firmados con clave X.509 del Registro Civil.',
            data: [
              (k: 'nombre_completo', v: '[Tu nombre]'),
              (k: 'cédula', v: '[Tu cédula]'),
              (k: 'fecha_nacimiento', v: '[Fecha]'),
              (k: 'estado_civil', v: '[Estado]'),
              (k: 'domicilio', v: '[Dirección]'),
            ],
          ),
        ],
        finalSummary:
            'Estos datos son la fuente única — Hacienda, CCSS y otras los consultan, no los guardan (once-only).',
      ),
      QuickPrompt(
        id: 'cambio-domicilio',
        label: 'Cambié de domicilio, ¿qué hago?',
        citizenAsk: 'Me mudé. ¿Cómo actualizo mi dirección?',
        llmIntro:
            'Cambias el domicilio una sola vez en Registro Civil. Hacienda, CCSS, INS y TSE se enteran automáticamente.',
        steps: [
          PromptStep(
            actor: 'Registro Civil',
            actorTag: 'RC',
            color: CrColors.areaIduc,
            kind: StepKind.request,
            text: 'Actualiza domicilio.',
            endpoint: 'PATCH /persons/{cedula}/address',
          ),
          PromptStep(
            actor: 'Conecta CR',
            actorTag: 'X-RD',
            color: CrColors.areaInterop,
            kind: StepKind.notify,
            text: 'Propaga evento a Hacienda, CCSS, INS, TSE.',
          ),
          PromptStep(
            actor: 'Notificaciones',
            actorTag: 'NOT',
            color: CrColors.areaPlatform,
            kind: StepKind.notify,
            text: 'Confirmación firmada al ciudadano.',
          ),
        ],
        finalSummary:
            'Listo. No tenés que ir a 4 ventanillas — esto es once-only en acción.',
        email: EmailNotice(
          subject: 'Domicilio actualizado',
          preview:
              'Tu nueva dirección quedó registrada y se notificó a 4 instituciones.',
        ),
      ),
      QuickPrompt(
        id: 'firmar-documento',
        label: 'Quiero firmar un documento',
        citizenAsk: 'Necesito firmar digitalmente un documento.',
        llmIntro: 'Inicio el flujo de firma con tu IDUC.',
        steps: [
          PromptStep(
            actor: 'IDUC',
            actorTag: 'ID',
            color: CrColors.areaIduc,
            kind: StepKind.request,
            text: 'Genera intent de firma.',
            endpoint: 'POST /signing/intents',
          ),
          PromptStep(
            actor: 'IDUC',
            actorTag: 'ID',
            color: CrColors.areaIduc,
            kind: StepKind.response,
            text:
                'Confirmá con biometría — la firma se ejecuta dentro del KMS, tu clave privada nunca sale.',
          ),
        ],
        finalSummary:
            'Cuando confirmes, recibís el documento firmado con sello de tiempo (RFC 3161).',
      ),
      QuickPrompt(
        id: 'padron-electoral',
        label: '¿En qué junta electoral voto?',
        citizenAsk: '¿En cuál junta receptora me toca votar?',
        llmIntro: 'Consulto el padrón del TSE.',
        steps: [
          PromptStep(
            actor: 'TSE',
            actorTag: 'TSE',
            color: CrColors.areaIduc,
            kind: StepKind.request,
            text: 'Consulta junta asignada.',
            endpoint: 'GET /electoral/poll-station/{cedula}',
          ),
          PromptStep(
            actor: 'TSE',
            actorTag: 'TSE',
            color: CrColors.areaIduc,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'junta', v: '4521'),
              (k: 'centro_votación', v: 'Esc. República de Argentina'),
              (k: 'cantón', v: 'San José'),
            ],
          ),
        ],
        finalSummary: 'Aplica para las próximas elecciones nacionales.',
      ),
    ],
  ),
  PromptCategory(
    id: 'migracion',
    label: 'Migración',
    icon: Icons.flight_takeoff_outlined,
    accent: CrColors.areaMembers,
    prompts: [
      QuickPrompt(
        id: 'renovar-pasaporte',
        label: 'Quiero renovar mi pasaporte',
        citizenAsk: 'Quiero renovar mi pasaporte.',
        llmIntro:
            'DGME consulta a Registro Civil y Hacienda por interop para validar identidad y solvencia.',
        steps: [
          PromptStep(
            actor: 'Migración',
            actorTag: 'DGME',
            color: CrColors.areaMembers,
            kind: StepKind.request,
            text: 'Inicia trámite de renovación.',
            endpoint: 'POST /immigration/passport/renew',
          ),
          PromptStep(
            actor: 'Registro Civil',
            actorTag: 'RC',
            color: CrColors.areaIduc,
            kind: StepKind.response,
            text: 'Identidad verificada.',
          ),
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.response,
            text: 'Solvencia OK.',
          ),
          PromptStep(
            actor: 'Migración',
            actorTag: 'DGME',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: 'Cita asignada: 22-may-2026, oficina central.',
          ),
        ],
        finalSummary:
            'Pasaporte listo en 5 días hábiles desde la cita. Te llega notificación de retiro.',
      ),
      QuickPrompt(
        id: 'estado-tramite-residencia',
        label: '¿Cómo va mi residencia?',
        citizenAsk: '¿En qué estado va mi solicitud de residencia?',
        llmIntro: 'Consulto el expediente en DGME.',
        steps: [
          PromptStep(
            actor: 'Migración',
            actorTag: 'DGME',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text:
                'Expediente RES-2026-0042: en revisión por inspector. Próximo paso: entrevista (te avisamos por correo).',
          ),
        ],
        finalSummary:
            'Sin acción de tu parte por ahora. Te aviso cuando agenden la entrevista.',
      ),
    ],
  ),
  PromptCategory(
    id: 'familia',
    label: 'Familia',
    icon: Icons.family_restroom_outlined,
    accent: CrColors.areaPlatform,
    prompts: [
      QuickPrompt(
        id: 'bono-imas',
        label: '¿Soy elegible para bono IMAS?',
        citizenAsk: '¿Soy elegible para algún bono o beca?',
        llmIntro:
            'IMAS valida con Hacienda, Registro Civil y CCSS sin que tengás que llevar papeles.',
        steps: [
          PromptStep(
            actor: 'IMAS',
            actorTag: 'IMAS',
            color: CrColors.areaMembers,
            kind: StepKind.request,
            text: 'Evalúa elegibilidad por SINIRUBE.',
            endpoint: 'GET /welfare/eligibility/{cedula}',
          ),
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.response,
            text: 'Ingreso bruto del último año.',
          ),
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.response,
            text: 'Cargas familiares.',
          ),
          PromptStep(
            actor: 'IMAS',
            actorTag: 'IMAS',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text:
                'Elegible para Bono Familiar. Monto estimado: ₡52.000/mes. Aplica beca avancemos para tu hijo/a.',
          ),
        ],
        finalSummary:
            'Te abro el formulario pre-llenado — solo confirmás y firmás.',
      ),
      QuickPrompt(
        id: 'inscribir-bebe',
        label: 'Acabo de tener un bebé, ¿qué hago?',
        citizenAsk: 'Mi hijo/a nació esta semana. ¿Qué tengo que hacer?',
        llmIntro:
            'En CR digital, el nacimiento se reporta una vez en CCSS y se propaga: Reg.Civil emite cédula menor, IMAS evalúa beca, Hacienda agrega dependiente.',
        steps: [
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.request,
            text: 'Hospital reporta nacimiento.',
            endpoint: 'POST /health/birth-event',
          ),
          PromptStep(
            actor: 'Registro Civil',
            actorTag: 'RC',
            color: CrColors.areaIduc,
            kind: StepKind.notify,
            text: 'Genera identidad del menor y cédula.',
          ),
          PromptStep(
            actor: 'PANI',
            actorTag: 'PANI',
            color: CrColors.areaMembers,
            kind: StepKind.notify,
            text: 'Activa expediente de protección.',
          ),
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.notify,
            text: 'Agrega dependiente a tu declaración.',
          ),
        ],
        finalSummary:
            'No tenés que hacer fila. Recibirás la cédula del menor y el ajuste de tu declaración por correo.',
        email: EmailNotice(
          subject: 'Cédula de tu hijo/a emitida',
          preview: 'Identidad firmada por Registro Civil. Listo para descargar.',
        ),
      ),
      QuickPrompt(
        id: 'pension-alimentaria',
        label: 'Quiero consultar mi pensión alimentaria',
        citizenAsk: '¿Cómo va mi expediente de pensión alimentaria?',
        llmIntro: 'Consulto Poder Judicial.',
        steps: [
          PromptStep(
            actor: 'Poder Judicial',
            actorTag: 'PJ',
            color: CrColors.areaGateway,
            kind: StepKind.response,
            text:
                'Expediente PA-2024-887: vigente. Próximo depósito: 5-may-2026 (₡180.000).',
          ),
        ],
        finalSummary: 'Sin alertas en el expediente.',
      ),
    ],
  ),
  PromptCategory(
    id: 'vivienda',
    label: 'Vivienda',
    icon: Icons.home_work_outlined,
    accent: CrColors.areaInterop,
    prompts: [
      QuickPrompt(
        id: 'comprar-casa',
        label: 'Quiero comprar una casa',
        citizenAsk: 'Quiero verificar una propiedad antes de comprarla.',
        llmIntro:
            'Consulto Registro Nacional para gravámenes, BCCR para tipo de cambio y Muni para impuestos pendientes.',
        steps: [
          PromptStep(
            actor: 'Reg. Nacional',
            actorTag: 'RN',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'finca', v: '1-12345'),
              (k: 'gravámenes', v: 'Ninguno'),
              (k: 'propietario', v: 'Mora & Asociados S.A.'),
              (k: 'área', v: '320 m²'),
            ],
          ),
          PromptStep(
            actor: 'Municipalidad',
            actorTag: 'MUNI',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'bienes_inmuebles', v: 'Al día'),
              (k: 'deuda_municipal', v: '₡0'),
            ],
          ),
        ],
        finalSummary:
            'Propiedad lista para escriturar. Te puedo abrir el flujo de notarial cuando estés listo.',
      ),
      QuickPrompt(
        id: 'permiso-construccion',
        label: 'Necesito permiso de construcción',
        citizenAsk: 'Quiero ampliar mi casa, ¿qué permisos necesito?',
        llmIntro:
            'Tramito en paralelo Muni, Bomberos, AyA, ICE y MINAE — antes era una vuelta de ~3 meses.',
        steps: [
          PromptStep(
            actor: 'Municipalidad',
            actorTag: 'MUNI',
            color: CrColors.areaMembers,
            kind: StepKind.request,
            text: 'Solicitud de permiso.',
          ),
          PromptStep(
            actor: 'Bomberos',
            actorTag: 'BOMB',
            color: CrColors.crRedBright,
            kind: StepKind.notify,
            text: 'Visto bueno de seguridad.',
          ),
          PromptStep(
            actor: 'AyA',
            actorTag: 'AyA',
            color: CrColors.areaInterop,
            kind: StepKind.notify,
            text: 'Disponibilidad de agua confirmada.',
          ),
          PromptStep(
            actor: 'MINAE',
            actorTag: 'MINAE',
            color: CrColors.areaIduc,
            kind: StepKind.notify,
            text: 'Sin afectación ambiental.',
          ),
        ],
        finalSummary:
            'Permiso pre-aprobado. Recibís el dictamen final firmado en 7 días hábiles.',
      ),
    ],
  ),
  PromptCategory(
    id: 'vehiculos',
    label: 'Vehículos',
    icon: Icons.directions_car_outlined,
    accent: CrColors.areaMembers,
    prompts: [
      QuickPrompt(
        id: 'marchamo',
        label: '¿Cuánto cuesta mi marchamo?',
        citizenAsk: '¿Cuánto pago de marchamo este año?',
        llmIntro: 'Consulto INS y MOPT.',
        steps: [
          PromptStep(
            actor: 'INS',
            actorTag: 'INS',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: 'Sin multas pendientes.',
            data: [
              (k: 'placa', v: 'BLP-123'),
              (k: 'soa_2026', v: '₡18.500'),
              (k: 'multas_inv', v: '0'),
            ],
          ),
          PromptStep(
            actor: 'MOPT',
            actorTag: 'MOPT',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'impuesto_circ', v: '₡126.400'),
              (k: 'total_marchamo', v: '₡144.900'),
              (k: 'vence', v: '31-dic-2026'),
            ],
          ),
        ],
        finalSummary:
            'Podés pagarlo desde MiCR. Te emito el comprobante firmado en el momento.',
      ),
      QuickPrompt(
        id: 'multas',
        label: '¿Tengo multas pendientes?',
        citizenAsk: '¿Tengo multas de tránsito pendientes?',
        llmIntro: 'Consulto Cosevi (MOPT).',
        steps: [
          PromptStep(
            actor: 'MOPT',
            actorTag: 'MOPT',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: 'Aún podés impugnar dentro de 5 días hábiles.',
            data: [
              (k: 'tipo_infracción', v: 'Exceso de velocidad'),
              (k: 'fecha', v: '12-feb-2026'),
              (k: 'monto', v: '₡52.000'),
              (k: 'estado', v: 'Pendiente'),
            ],
          ),
        ],
        finalSummary: 'Podés pagarla acá o impugnarla en línea.',
      ),
    ],
  ),
  PromptCategory(
    id: 'trabajo',
    label: 'Trabajo y empresa',
    icon: Icons.business_center_outlined,
    accent: CrColors.areaMembers,
    prompts: [
      QuickPrompt(
        id: 'planilla',
        label: '¿Estoy reportado en planilla?',
        citizenAsk: '¿Mi patrono me tiene en planilla?',
        llmIntro: 'Verifico CCSS y MTSS.',
        steps: [
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'estado', v: 'Activa'),
              (k: 'salario_reportado', v: '₡1.350.000/mes'),
              (k: 'última_planilla', v: 'abr-2026'),
            ],
          ),
          PromptStep(
            actor: 'MTSS',
            actorTag: 'MTSS',
            color: CrColors.areaMembers,
            kind: StepKind.response,
            text: '',
            data: [
              (k: 'patrono', v: 'Al día'),
              (k: 'cargas_sociales', v: '0 vencidas'),
            ],
          ),
        ],
        finalSummary: 'Todo en orden con tu patrono.',
      ),
      QuickPrompt(
        id: 'abrir-empresa',
        label: 'Quiero abrir una empresa',
        citizenAsk: 'Quiero registrar una sociedad anónima.',
        llmIntro:
            'Tramito en paralelo Reg. Nacional, Hacienda, CCSS, MTSS y SUGEF — antes era ~30 días, ahora es ~3.',
        steps: [
          PromptStep(
            actor: 'Reg. Nacional',
            actorTag: 'RN',
            color: CrColors.areaMembers,
            kind: StepKind.request,
            text: 'Inscripción de sociedad.',
          ),
          PromptStep(
            actor: 'Hacienda',
            actorTag: 'HAC',
            color: CrColors.crBlueBright,
            kind: StepKind.notify,
            text: 'Asigna número de cédula jurídica.',
          ),
          PromptStep(
            actor: 'CCSS',
            actorTag: 'CCSS',
            color: CrColors.crRedBright,
            kind: StepKind.notify,
            text: 'Activa patrono.',
          ),
          PromptStep(
            actor: 'MTSS',
            actorTag: 'MTSS',
            color: CrColors.areaMembers,
            kind: StepKind.notify,
            text: 'Registra empleador.',
          ),
        ],
        finalSummary: 'Sociedad activa. Recibís el folio firmado por correo.',
      ),
    ],
  ),
  PromptCategory(
    id: 'seguridad',
    label: 'Seguridad y justicia',
    icon: Icons.shield_outlined,
    accent: CrColors.areaGateway,
    prompts: [
      QuickPrompt(
        id: 'denuncia',
        label: 'Quiero poner una denuncia',
        citizenAsk: 'Quiero presentar una denuncia.',
        llmIntro: 'Inicio expediente con Fuerza Pública y Poder Judicial.',
        steps: [
          PromptStep(
            actor: 'Fuerza Pública',
            actorTag: 'FP',
            color: CrColors.areaGateway,
            kind: StepKind.request,
            text: 'Crea expediente preliminar.',
          ),
          PromptStep(
            actor: 'Poder Judicial',
            actorTag: 'PJ',
            color: CrColors.areaGateway,
            kind: StepKind.notify,
            text: 'Asigna número único de expediente.',
          ),
        ],
        finalSummary:
            'Expediente abierto. Vas a poder seguir el estado desde MiCR sin volver a la oficina.',
      ),
      QuickPrompt(
        id: 'antecedentes',
        label: 'Necesito mi hoja de delincuencia',
        citizenAsk: 'Necesito mi hoja de delincuencia.',
        llmIntro: 'Consulto Poder Judicial.',
        steps: [
          PromptStep(
            actor: 'Poder Judicial',
            actorTag: 'PJ',
            color: CrColors.areaGateway,
            kind: StepKind.response,
            text: 'Documento firmado con sello de tiempo RFC 3161.',
            data: [
              (k: 'estado', v: 'Limpia'),
              (k: 'antecedentes', v: '0 registros'),
              (k: 'vigencia', v: '90 días'),
            ],
          ),
        ],
        finalSummary: 'Te lo envío al correo registrado.',
        email: EmailNotice(
          subject: 'Hoja de delincuencia firmada',
          preview: 'PDF con firma CAdES y sello de tiempo RFC 3161.',
        ),
      ),
    ],
  ),
  PromptCategory(
    id: 'transparencia',
    label: 'Mi bitácora',
    icon: Icons.receipt_outlined,
    accent: CrColors.areaPlatform,
    prompts: [
      QuickPrompt(
        id: 'quien-consulto',
        label: '¿Quién consultó mis datos?',
        citizenAsk: '¿Qué instituciones consultaron mis datos esta semana?',
        llmIntro: 'Reviso la bitácora hash-chained de Conecta CR.',
        steps: [
          PromptStep(
            actor: 'Bitácora',
            actorTag: 'AUD',
            color: CrColors.areaPlatform,
            kind: StepKind.response,
            text: 'Total: 12 accesos en 7 días.',
            data: [
              (k: 'CCSS', v: '5 accesos'),
              (k: 'Hacienda', v: '4 accesos'),
              (k: 'IMAS', v: '1 acceso'),
              (k: 'Municipalidad', v: '1 acceso'),
              (k: 'Reg.Civil', v: '1 acceso'),
            ],
          ),
        ],
        finalSummary:
            'Cada acceso está firmado y encadenado. Si algo no cuadra, podés impugnar desde el detalle.',
      ),
      QuickPrompt(
        id: 'verificar-cadena',
        label: 'Verificá la integridad de mi bitácora',
        citizenAsk: 'Verificá que mi bitácora no haya sido alterada.',
        llmIntro: 'Recalculo el Merkle root del último epoch.',
        steps: [
          PromptStep(
            actor: 'Bitácora',
            actorTag: 'AUD',
            color: CrColors.areaPlatform,
            kind: StepKind.response,
            text: 'Hash chain íntegro · Merkle root coincide con el publicado.',
            data: [
              (k: 'merkle_root', v: '0x9f3a…b2c1'),
              (k: 'epoch', v: '4521'),
              (k: 'entries', v: '10.000'),
              (k: 'estado', v: 'OK'),
            ],
          ),
        ],
        finalSummary:
            'Tu historial es verificable end-to-end. Esto es la garantía estonia de transparencia.',
      ),
    ],
  ),
];

List<PromptCategory> promptCategoriesFor(AppLocalizations t) => _categories;

QuickPrompt? promptById(String id) {
  for (final c in _categories) {
    for (final p in c.prompts) {
      if (p.id == id) return p;
    }
  }
  return null;
}

final List<({List<String> kw, String promptId})> _keywordRoutes = [
  (kw: ['renta', 'declaración', 'declaracion', 'tributaria', 'impuesto sobre la renta'], promptId: 'declaracion-2025'),
  (kw: ['historial renta', 'declaraciones pasadas', 'años fiscales'], promptId: 'historial-renta'),
  (kw: ['muni', 'municipal', 'bienes inmuebles', 'patente'], promptId: 'patentes-municipal'),
  (kw: ['receta', 'medicamento', 'farmacia'], promptId: 'recetas-activas'),
  (kw: ['cita', 'consulta médica', 'consulta medica', 'ebais', 'hospital'], promptId: 'proxima-cita'),
  (kw: ['incapacidad', 'enfermedad', 'reposo'], promptId: 'incapacidad'),
  (kw: ['vacuna', 'carné', 'carne', 'inmunización'], promptId: 'vacunas'),
  (kw: ['mis datos', 'qué datos', 'que datos', 'identidad'], promptId: 'mis-datos'),
  (kw: ['domicilio', 'dirección', 'direccion', 'mudé', 'mude', 'mudanza'], promptId: 'cambio-domicilio'),
  (kw: ['firmar', 'firma digital', 'firma electrónica', 'firma electronica'], promptId: 'firmar-documento'),
  (kw: ['junta', 'voto', 'electoral', 'padrón', 'padron'], promptId: 'padron-electoral'),
  (kw: ['pasaporte'], promptId: 'renovar-pasaporte'),
  (kw: ['residencia', 'migración', 'migracion'], promptId: 'estado-tramite-residencia'),
  (kw: ['bono', 'beca', 'imas', 'sinirube', 'subsidio'], promptId: 'bono-imas'),
  (kw: ['bebé', 'bebe', 'nacimiento', 'recién nacido', 'recien nacido'], promptId: 'inscribir-bebe'),
  (kw: ['pensión alimentaria', 'pension alimentaria', 'pensión', 'pension'], promptId: 'pension-alimentaria'),
  (kw: ['casa', 'finca', 'comprar propiedad', 'comprar casa', 'gravamen'], promptId: 'comprar-casa'),
  (kw: ['construcción', 'construccion', 'ampliación', 'ampliacion', 'permiso construcción', 'permiso construccion'], promptId: 'permiso-construccion'),
  (kw: ['marchamo', 'circulación', 'circulacion', 'soa'], promptId: 'marchamo'),
  (kw: ['multa', 'parte', 'cosevi'], promptId: 'multas'),
  (kw: ['planilla', 'patrono', 'cargas sociales'], promptId: 'planilla'),
  (kw: ['empresa', 'sociedad', 'cédula jurídica', 'cedula juridica', 's.a.', 'sa.'], promptId: 'abrir-empresa'),
  (kw: ['denuncia', 'fiscalía', 'fiscalia', 'oij'], promptId: 'denuncia'),
  (kw: ['hoja de delincuencia', 'antecedentes', 'récord', 'record'], promptId: 'antecedentes'),
  (kw: ['quién consultó', 'quien consulto', 'accesos a mis datos', 'bitácora', 'bitacora'], promptId: 'quien-consulto'),
  (kw: ['integridad', 'merkle', 'hash chain', 'verificar bitácora', 'verificar bitacora'], promptId: 'verificar-cadena'),
];

QuickPrompt? matchFreeForm(String text) {
  final t = text.toLowerCase();
  for (final r in _keywordRoutes) {
    for (final k in r.kw) {
      if (t.contains(k)) {
        final p = promptById(r.promptId);
        if (p != null) {
          return QuickPrompt(
            id: 'free-${DateTime.now().microsecondsSinceEpoch}',
            label: text,
            citizenAsk: text,
            llmIntro: p.llmIntro,
            steps: p.steps,
            finalSummary: p.finalSummary,
            email: p.email,
          );
        }
      }
    }
  }
  return null;
}

QuickPrompt buildFreeFormFallback(String text) {
  final matched = matchFreeForm(text);
  if (matched != null) return matched;

  return QuickPrompt(
    id: 'free-${DateTime.now().microsecondsSinceEpoch}',
    label: text,
    citizenAsk: text,
    llmIntro:
        'Voy a enrutarla por Conecta CR para identificar qué institución es competente y traerte la respuesta firmada.',
    steps: const [
      PromptStep(
        actor: 'Conecta CR',
        actorTag: 'X-RD',
        color: CrColors.areaInterop,
        kind: StepKind.thinking,
        text: 'Analizo la consulta y selecciono al member competente del catálogo.',
        delayMs: 800,
      ),
      PromptStep(
        actor: 'Registro Civil',
        actorTag: 'RC',
        color: CrColors.areaIduc,
        kind: StepKind.request,
        text: 'Verifica identidad del ciudadano antes de enrutar.',
        endpoint: 'GET /persons/{cedula}',
        delayMs: 900,
      ),
      PromptStep(
        actor: 'Conecta CR',
        actorTag: 'X-RD',
        color: CrColors.areaInterop,
        kind: StepKind.audit,
        text: 'Consulta firmada y registrada en la bitácora hash-chained.',
        delayMs: 700,
      ),
      PromptStep(
        actor: 'Asistente',
        actorTag: 'LLM',
        color: CrColors.crBlueBright,
        kind: StepKind.response,
        text:
            'No encontré una respuesta exacta en el catálogo demo. En producción este enrutador llama al member correcto vía security-server.',
        delayMs: 800,
      ),
    ],
    finalSummary:
        'Probá reformular tu consulta o tocá una de las categorías abajo — cubren los trámites más comunes del Estado.',
  );
}

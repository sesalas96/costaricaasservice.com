export type Lang = "es" | "en";

export type Translations = {
  hero: {
    eyebrow: string;
    titleStart: string;
    titleHighlight: string;
    descBefore: string;
    descBold1: string;
    descMid1: string;
    descBold2: string;
    descMid2: string;
    descBold3: string;
    descAfter: string;
    ctaPrimary: string;
    ctaSecondary: string;
  };
  guarantees: {
    eyebrow: string;
    heading: string;
    lead: string;
    audiencesHeading: string;
    audiences: { tag: string; title: string; body: string }[];
    casesEyebrow: string;
    casesHeading: string;
    casesLead: string;
    cases: {
      tag: string;
      title: string;
      body: string;
      actors: string;
      category: string;
    }[];
    casesCtaCard: {
      tag: string;
      title: string;
      body: string;
      action: string;
    };
    invariantsEyebrow: string;
    invariantsHeading: string;
    items: { title: string; body: string }[];
  };
  casesPage: {
    eyebrow: string;
    heading: string;
    lead: string;
    backToHome: string;
    casesCount: string;
    categories: {
      gobierno: { label: string; description: string };
      identidad: { label: string; description: string };
      bienestar: { label: string; description: string };
      justicia: { label: string; description: string };
    };
  };
  diagrams: {
    eyebrow: string;
    heading: string;
    tip: string;
    audience: {
      citizen: string;
      technical: string;
      hint: string;
      activeBadge: string;
    };
    modulesHeading: string;
    modules: {
      ecosistema: { tag: string; title: string; sub: string };
      interop: { tag: string; title: string; sub: string };
      onceOnly: { tag: string; title: string; sub: string };
      identity: { tag: string; title: string; sub: string };
    };
    subs: {
      ecosistema: { citizen: string; technical: string };
      interop: { citizen: string; technical: string };
      onceOnly: { citizen: string; technical: string };
      identity: { citizen: string; technical: string };
    };
  };
  deliverables: {
    eyebrow: string;
    heading: string;
    lead: string;
    view: {
      business: string;
      technical: string;
      hint: string;
      activeBadge: string;
    };
    business: {
      intro: string;
      cards: {
        tag: string;
        title: string;
        sub: string;
        items: string[];
      }[];
    };
    technical: {
      intro: string;
      cards: {
        tag: string;
        title: string;
        sub: string;
        items: string[];
      }[];
    };
  };
  footer: {
    coreOf: string;
  };
  langSwitch: {
    aria: string;
  };
  ecosystem: {
    citizen: { title: string; sub: string };
    operator: { title: string; sub: string };
    admin: { title: string; sub: string };
    gateway: { title: string; sub: string };
    bffCitizen: { title: string; sub: string };
    bffMember: { title: string; sub: string };
    bffAdmin: { title: string; sub: string };
    iducGroup: string;
    iducIdentity: { title: string; sub: string };
    iducKeys: { title: string; sub: string };
    iducSigning: { title: string; sub: string };
    interopGroup: string;
    interopHub: { title: string; sub: string };
    securityServer: { title: string; sub: string };
    interopRouter: { title: string; sub: string };
    interopAudit: { title: string; sub: string };
    membersGroup: string;
    registroCivil: { title: string; sub: string };
    hacienda: { title: string; sub: string };
    otherMembers: { title: string; sub: string };
    platformGroup: string;
    platformAudit: { title: string; sub: string };
    notifications: { title: string; sub: string };
    files: { title: string; sub: string };
    featureFlags: { title: string; sub: string };
    edges: {
      kms: string;
      expose: string;
      exposeConsume: string;
      mtls: string;
      event: string;
    };
  };
  interop: {
    haciendaSvc: { title: string; sub: string };
    ssHac: { title: string; sub: string };
    router: { title: string; sub: string };
    ssRc: { title: string; sub: string };
    rcSvc: { title: string; sub: string };
    audit: { title: string; sub: string };
    hub: { title: string; sub: string };
    edges: {
      request: string;
      jwsSigned: string;
      routeMtls: string;
      verifyDeliver: string;
      response: string;
      kafkaEvent: string;
      catalogCa: string;
    };
  };
  onceOnly: {
    citizen: { title: string; sub: string };
    micr: { title: string; sub: string };
    bff: { title: string; sub: string };
    hacienda: { title: string; sub: string };
    ssHac: { title: string };
    router: { title: string };
    ssRc: { title: string };
    rc: { title: string; sub: string };
    callout: { title: string; sub: string };
    edges: {
      request: string;
      step2: string;
      getPrefilled: string;
      needsPerson: string;
      getPerson: string;
      citizenData: string;
      prefilled: string;
      step8: string;
      signClick: string;
    };
  };
  identity: {
    citizen: { title: string; sub: string };
    micr: { title: string; sub: string };
    iduc: { title: string; sub: string };
    keys: { title: string; sub: string };
    signing: { title: string; sub: string };
    document: { title: string; sub: string };
    member: { title: string; sub: string };
    edges: {
      login: string;
      issueJwt: string;
      sign: string;
      askKey: string;
      signedDoc: string;
      verify: string;
    };
  };
  ecosystemCitizen: {
    me: { title: string; sub: string };
    micr: { title: string; sub: string };
    conecta: { title: string; sub: string };
    institutions: { title: string; sub: string };
    audit: { title: string; sub: string };
    signing: { title: string; sub: string };
    edges: {
      ingreso: string;
      consulta: string;
      enruta: string;
      registra: string;
      veo: string;
      firmo: string;
    };
  };
  interopCitizen: {
    institutionA: { title: string; sub: string };
    bridge: { title: string; sub: string };
    institutionB: { title: string; sub: string };
    audit: { title: string; sub: string };
    edges: {
      needs: string;
      sealed: string;
      response: string;
      logged: string;
    };
  };
  onceOnlyCitizen: {
    me: { title: string; sub: string };
    micr: { title: string; sub: string };
    hacienda: { title: string; sub: string };
    bridge: { title: string; sub: string };
    rc: { title: string; sub: string };
    callout: { title: string; sub: string };
    edges: {
      request: string;
      micrAsks: string;
      hacAsksRc: string;
      rcAnswers: string;
      readyToSign: string;
      signed: string;
    };
  };
  identityCitizen: {
    me: { title: string; sub: string };
    eid: { title: string; sub: string };
    safe: { title: string; sub: string };
    document: { title: string; sub: string };
    member: { title: string; sub: string };
    callout: { title: string; sub: string };
    edges: {
      login: string;
      askSign: string;
      signs: string;
      anyoneVerifies: string;
    };
  };
  institutionPanel: {
    closeAria: string;
    badgeLive: string;
    badgePlanned: string;
    listHeading: string;
    listSub: string;
    listCount: string;
    backToList: string;
    sectionCustodia: string;
    sectionExpone: string;
    sectionConsume: string;
    emptyConsume: string;
    nodeBadgeLabel: string;
    diagramHint: string;
  };
};

export const translations: Record<Lang, Translations> = {
  es: {
    hero: {
      eyebrow: "costaricaasservice · núcleo de Conecta CR",
      titleStart: "Gobierno digital,",
      titleHighlight: "como producto.",
      descBefore:
        "Estonia tardó 25 años. Nosotros, 3 meses. ",
      descBold1: "Las instituciones se hablan entre ellas",
      descMid1: " y no te piden lo mismo dos veces. ",
      descBold2: "Firmás desde el celular",
      descMid2: " con valor legal. ",
      descBold3: "Ves quién consulta tus datos",
      descAfter: " y por qué.",
      ctaPrimary: "Ver el ecosistema",
      ctaSecondary: "Garantías estonias",
    },
    guarantees: {
      eyebrow: "Propuesta de valor",
      heading: "El sistema operativo del Estado digital.",
      lead: "costaricaasservice no es middleware. Es la plataforma sobre la que un país opera su gobierno digital — desde un trámite de cinco minutos hasta una elección verificable. La misma capa que conecta a Registro Civil con Hacienda conecta a la CCSS con una farmacia privada, con idéntica garantía de identidad, trazabilidad y firma.",
      audiencesHeading: "Quién gana con esto",
      audiences: [
        {
          tag: "Ciudadano",
          title: "Trámites en minutos, no en semanas",
          body: "Ves en MiCR quién consultó tus datos y para qué. Firmás desde el celular con valor jurídico equivalente a la firma manuscrita.",
        },
        {
          tag: "Institución",
          title: "Te conectás sin construir conectores",
          body: "Exponés un servicio una vez y todo el ecosistema lo consume con audit automático. Te ahorrás los 30 contratos bilaterales con cada par.",
        },
        {
          tag: "Jurisdicción",
          title: "Una sola gobernanza del Estado",
          body: "Catálogo único, métricas de uso entre instituciones, capacidad de revocar accesos en tiempo real. Visión de país, no de silo.",
        },
      ],
      casesEyebrow: "Casos de uso",
      casesHeading: "Lo que el Estado digital habilita.",
      casesLead: "Estos son los flujos reales que costaricaasservice permite armar — todos cumplen las 4 garantías estonias por construcción, no por diligencia.",
      cases: [
        {
          tag: "Hacienda",
          title: "Declaración pre-llenada",
          body: "Hacienda construye la declaración consultando salarios al patrono, cargas a la CCSS y composición familiar a Registro Civil — todo vía X-Road. El ciudadano abre MiCR, revisa y firma con un click. De semanas a 90 segundos.",
          actors: "Hacienda · CCSS · Reg.Civil · MiCR",
          category: "gobierno",
        },
        {
          tag: "Migración",
          title: "Pasaporte sin papeles",
          body: "Migración verifica identidad con e-ID, antecedentes con MJ y cobra el arancel a Hacienda en un solo formulario. El ciudadano agenda la entrega y firma el comprobante digitalmente.",
          actors: "Migración · MJ · Hacienda · iduc",
          category: "identidad",
        },
        {
          tag: "Salud",
          title: "Receta médica electrónica",
          body: "La CCSS emite la receta firmada por el médico. La farmacia (member privado) la valida vía X-Road antes de despachar. El ciudadano la ve en MiCR con su histórico clínico.",
          actors: "CCSS · Farmacia · MiCR",
          category: "bienestar",
        },
        {
          tag: "Once-only",
          title: "Cambio de domicilio · una sola vez",
          body: "El ciudadano actualiza su dirección en Reg.Civil. Hacienda, CCSS, INS y TSE reciben el cambio en tiempo real vía Kafka. Nunca más cargar comprobantes en cada ventanilla.",
          actors: "Reg.Civil → Hacienda · CCSS · INS · TSE",
          category: "identidad",
        },
        {
          tag: "Banca",
          title: "KYC bancario verificable",
          body: "Un banco (member) abre cuenta verificando identidad con e-ID y antecedentes vía X-Road. Sin pasaporte, sin recibos de servicios, sin esperar 48 horas para validar.",
          actors: "Banco · iduc · Reg.Civil",
          category: "identidad",
        },
        {
          tag: "Compras",
          title: "Contratación pública trazable",
          body: "Empresas firman ofertas y contratos con sellado RFC 3161. Cada paso del proceso queda en bitácora hash-chained, verificable por cualquier ciudadano sin pedir permiso.",
          actors: "Hacienda · proveedores · audit",
          category: "gobierno",
        },
        {
          tag: "Social",
          title: "Beneficios sociales sin trámite",
          body: "IMAS calcula elegibilidad consultando ingresos en Hacienda y composición familiar en Reg.Civil — el solicitante no presenta nada. Once-only en su forma más pura.",
          actors: "IMAS · Hacienda · Reg.Civil",
          category: "bienestar",
        },
        {
          tag: "Justicia",
          title: "Notificación judicial electrónica",
          body: "Los tribunales envían citaciones a MiCR con valor legal equivalente al físico. Acuse de recibo criptográfico, trazabilidad completa para el juzgado.",
          actors: "PJ · notifications · firma",
          category: "justicia",
        },
        {
          tag: "Democracia",
          title: "Voto electrónico verificable",
          body: "TSE valida identidad e-ID, registra el voto firmado y publica un resumen Merkle de la jornada. Cada votante puede verificar que su voto fue contado — sin revelar por quién votó.",
          actors: "TSE · iduc · interop-audit",
          category: "justicia",
        },
      ],
      casesCtaCard: {
        tag: "Catálogo completo",
        title: "Ver todos los casos de uso",
        body: "Catálogo organizado por dominio: gobierno, identidad, bienestar y justicia. Con flujo, actores y la garantía estonia que cumple cada uno.",
        action: "Abrir catálogo",
      },
      invariantsEyebrow: "Bajo el capó",
      invariantsHeading: "Cuatro garantías no negociables.",
      items: [
        {
          title: "Once-only",
          body: "El ciudadano da un dato al sistema una sola vez. Hacienda consulta a Registro Civil vía interop, no le pide al ciudadano.",
        },
        {
          title: "Descentralización",
          body: "Cada institución es dueña de sus datos. No hay base de datos central; los servicios se exponen y consumen vía X-Road.",
        },
        {
          title: "Trazabilidad ciudadana",
          body: "El ciudadano ve quién consultó sus datos, cuándo y con qué propósito. Bitácora hash-chained, verificable.",
        },
        {
          title: "Firma digital con valor jurídico",
          body: "Equivalente a firma manuscrita. Claves Ed25519 custodiadas en KMS, firma server-side, sellado de tiempo RFC 3161.",
        },
      ],
    },
    casesPage: {
      eyebrow: "Catálogo · casos de uso",
      heading: "Todos los flujos que costaricaasservice habilita.",
      lead: "Casos organizados por dominio. Cada uno detalla los actores que intervienen y la categoría a la que pertenece — todos cumplen las cuatro garantías estonias por construcción.",
      backToHome: "Volver al inicio",
      casesCount: "{n} casos catalogados · {c} dominios",
      categories: {
        gobierno: {
          label: "Hacienda y gobierno",
          description: "Tributación, compras públicas y trámites con valor fiscal.",
        },
        identidad: {
          label: "Identidad y movilidad",
          description: "Verificación de identidad, datos personales y acceso al sector privado.",
        },
        bienestar: {
          label: "Salud y bienestar social",
          description: "Servicios sanitarios y beneficios sociales sin trámite.",
        },
        justicia: {
          label: "Justicia y democracia",
          description: "Notificaciones legales, procesos judiciales y voto verificable.",
        },
      },
    },
    diagrams: {
      eyebrow: "Diagramas interactivos",
      heading: "Cómo funciona el ecosistema.",
      tip: "scroll para zoom · arrastrá para mover · controles abajo a la izquierda",
      audience: {
        citizen: "Vista ciudadana",
        technical: "Vista técnica",
        hint: "Cambiá entre lenguaje cotidiano y la implementación real.",
        activeBadge: "activa",
      },
      modulesHeading: "Módulos",
      modules: {
        ecosistema: {
          tag: "Vista global",
          title: "Ecosistema",
          sub: "Cómo encajan las piezas",
        },
        interop: {
          tag: "Flujo",
          title: "Conecta · X-Road",
          sub: "Cómo viaja un pedido entre instituciones",
        },
        onceOnly: {
          tag: "Caso",
          title: "Once-only en acción",
          sub: "Hacienda consulta a Registro Civil",
        },
        identity: {
          tag: "Capacidad",
          title: "Identidad y firma",
          sub: "Cómo te identificás y firmás",
        },
      },
      subs: {
        ecosistema: {
          citizen:
            "Lo que vos usás del Estado: una app, las instituciones que te atienden y un cartero digital que las conecta — sin que tengas que ir a una ventanilla.",
          technical:
            "Vista completa del monorepo: gateway, BFFs por audiencia, iduc, interop, members y platform.",
        },
        interop: {
          citizen:
            "Cuando una institución necesita un dato tuyo, no te lo pide a vos: se lo pide directo a la institución que ya lo tiene. Vos ves todo en tu bitácora.",
          technical:
            "Cómo viaja un request inter-member: JWS detached, mTLS, idempotencia ULID y evento Kafka al hash chain.",
        },
        onceOnly: {
          citizen:
            "Hacienda te arma la declaración pidiéndole los datos a Registro Civil. Vos no llenás formularios, no cargás papeles, no vas a ninguna ventanilla — solo abrís la app, revisás y firmás.",
          technical:
            "Hacienda invoca a registro-civil vía security-server + interop-router; respuesta cacheable, audit registrado.",
        },
        identity: {
          citizen:
            "Tu identidad y tu firma digital están guardadas en una caja fuerte que ni nosotros podemos abrir. Firmás con un toque desde el celular y vale igual que tu firma a mano.",
          technical:
            "iduc-identity emite JWT RS256, iduc-keys custodia Ed25519 en KMS (Vault), iduc-signing produce JWS/CAdES/PAdES con sellado RFC 3161.",
        },
      },
    },
    deliverables: {
      eyebrow: "Entregables del sistema",
      heading: "Qué se entrega cuando costaricaasservice se enciende.",
      lead: "Un mismo monorepo, dos lecturas. Lo que recibís según tu rol — y lo que hay adentro de la caja, repo por repo.",
      view: {
        business: "Vista de negocio",
        technical: "Vista técnica",
        hint: "Cambiá entre qué te llevás y qué hay adentro.",
        activeBadge: "activa",
      },
      business: {
        intro:
          "Cuatro paquetes alineados con quién consume el sistema: el Estado que lo opera, la institución que lo extiende, el ciudadano que lo usa y la plataforma transversal que sostiene el contrato legal.",
        cards: [
          {
            tag: "Realm · jurisdicción",
            title: "Para el Estado que lo opera",
            sub: "Control plane completo del realm",
            items: [
              "Control plane admin (gestión de members, claves, billing)",
              "Hub de interoperabilidad con CA del realm",
              "Gobernanza del catálogo de servicios inter-member",
              "Observabilidad end-to-end (trazas, métricas, p99 < 500ms)",
              "Runbooks y respuesta a incidentes",
            ],
          },
          {
            tag: "Member · institución",
            title: "Para cada institución",
            sub: "Conectarse al ecosistema en días, no años",
            items: [
              "Security-server desplegable (firma JWS + mTLS llave en mano)",
              "BFF Member + portal operador",
              "SDK interop-client (expone y consume servicios remotos)",
              "Plantilla de servicio (cri-templates-service) lista para producción",
              "Onboarding al catálogo y emisión de certificados X.509",
            ],
          },
          {
            tag: "Ciudadano · MiCR",
            title: "Para cada habitante",
            sub: "Una sola puerta al Estado",
            items: [
              "MiCR Web (Next.js 15) y MiCR Mobile (Flutter melos)",
              "Identidad digital única + login con WebAuthn/passkey",
              "Firma digital con valor jurídico (Ed25519 en KMS)",
              "Bitácora de accesos: quién consultó tus datos y por qué",
              "Notificaciones legales y archivos firmados (PDF + CAdES/PAdES)",
            ],
          },
          {
            tag: "Plataforma · transversal",
            title: "Capacidades comunes",
            sub: "El contrato legal y la trazabilidad que sostienen todo",
            items: [
              "Audit chain hash-encadenado (SHA-256 + Merkle por epoch)",
              "Notificaciones electrónicas con valor legal",
              "Servicio de archivos firmados (PDF + CAdES/PAdES)",
              "Feature flags por realm y por member",
              "KMS con envelope encryption (Vault Transit en dev, AWS KMS en prod)",
            ],
          },
        ],
      },
      technical: {
        intro:
          "Las seis áreas del monorepo, con los repos concretos que se despliegan en cada una. Cada servicio con su propia DB, multi-tenant desde el día uno (schema-per-realm).",
        cards: [
          {
            tag: "gateway",
            title: "Entrada al sistema",
            sub: "Gateway central + BFFs por audiencia",
            items: [
              "cri-gateway-api · JWT RS256 + RBAC + rate limit + revocación Roaring",
              "cri-bff-citizen · orquesta MiCR (web + mobile)",
              "cri-bff-member · portal operador",
              "cri-bff-admin · control plane costaricaasservice",
            ],
          },
          {
            tag: "iduc",
            title: "Identidad Digital Única",
            sub: "Autenticación, claves y firma",
            items: [
              "cri-svc-iduc-identity · login + WebAuthn + JWT RS256",
              "cri-svc-iduc-keys · Ed25519 custodiadas en KMS",
              "cri-svc-iduc-signing · JWS / CAdES / PAdES + sellado RFC 3161",
            ],
          },
          {
            tag: "interop",
            title: "Conecta · X-Road analog",
            sub: "Hub central + plano de datos federado",
            items: [
              "cri-svc-interop-hub · catálogo + CA del realm + gobernanza",
              "cri-svc-security-server · daemon que cada member instala",
              "cri-svc-interop-router · plano de datos HTTP/2 + mTLS",
              "cri-svc-interop-audit · hash chain SHA-256 + Merkle epoch",
            ],
          },
          {
            tag: "members",
            title: "Instituciones",
            sub: "Servicios mock listos para conectar reales",
            items: [
              "cri-svc-registro-civil · TSE mock con 1000 personas seed",
              "cri-svc-hacienda · declaración pre-llenada (consume Reg.Civil)",
              "+ MIT, CCSS, MOPT… (escalable por catálogo)",
            ],
          },
          {
            tag: "platform",
            title: "Plataforma transversal",
            sub: "Capacidades compartidas con valor legal",
            items: [
              "cri-svc-audit · vista ciudadana del log inter-member",
              "cri-svc-notifications · notificaciones con valor legal",
              "cri-svc-files · PDF firmados (storage S3 en prod)",
              "cri-svc-feature-flags · flags por realm/member",
              "cri-templates-service · scaffold para servicios nuevos",
            ],
          },
          {
            tag: "frontends",
            title: "Superficies cliente",
            sub: "Web + mobile + landing institucional",
            items: [
              "cri-landing · sitio público (Next.js 15)",
              "cri-web-micr · MiCR para ciudadanos (Next.js 15)",
              "cri-web-admin · control plane de members y realm",
              "cri-mobile · MiCR Mobile (Flutter melos, iOS + Android)",
            ],
          },
        ],
      },
    },
    footer: {
      coreOf: "núcleo de Conecta CR",
    },
    langSwitch: { aria: "Cambiar idioma" },
    ecosystem: {
      citizen: { title: "Ciudadano", sub: "MiCR Web · Mobile" },
      operator: { title: "Operador Member", sub: "Portal institución" },
      admin: { title: "Admin costaricaasservice", sub: "Control plane" },
      gateway: {
        title: "API Gateway",
        sub: "JWT RS256 · RBAC · rate limit · revocación Roaring",
      },
      bffCitizen: { title: "BFF Citizen", sub: "Orquesta MiCR" },
      bffMember: { title: "BFF Member", sub: "Portal operador" },
      bffAdmin: { title: "BFF Admin", sub: "Control plane costaricaasservice" },
      iducGroup: "iduc · identidad",
      iducIdentity: {
        title: "iduc-identity",
        sub: "Login · WebAuthn · JWT",
      },
      iducKeys: { title: "iduc-keys", sub: "Ed25519 en KMS" },
      iducSigning: { title: "iduc-signing", sub: "JWS · CAdES · PAdES" },
      interopGroup: "interop · Conecta",
      interopHub: { title: "interop-hub", sub: "Catálogo · CA del realm" },
      securityServer: {
        title: "security-server",
        sub: "JWS detached · mTLS",
      },
      interopRouter: {
        title: "interop-router",
        sub: "Plano de datos · idempotencia",
      },
      interopAudit: {
        title: "interop-audit",
        sub: "Hash chain · Merkle epoch",
      },
      membersGroup: "members · instituciones",
      registroCivil: {
        title: "registro-civil",
        sub: "TSE mock · 1000 personas",
      },
      hacienda: { title: "hacienda", sub: "Declaración pre-llenada" },
      otherMembers: {
        title: "+ otros members",
        sub: "MIT · CCSS · MOPT…",
      },
      platformGroup: "platform · transversal",
      platformAudit: {
        title: "audit (ciudadano)",
        sub: "Vista de quién consultó",
      },
      notifications: { title: "notifications", sub: "Valor legal" },
      files: { title: "files", sub: "PDF firmados" },
      featureFlags: {
        title: "feature-flags",
        sub: "Por realm / member",
      },
      edges: {
        kms: "KMS",
        expose: "expone",
        exposeConsume: "expone/consume",
        mtls: "mTLS",
        event: "evento",
      },
    },
    interop: {
      haciendaSvc: {
        title: "Hacienda · svc",
        sub: "Originador del request",
      },
      ssHac: {
        title: "security-server (Hacienda)",
        sub: "Firma JWS detached con clave del member",
      },
      router: {
        title: "interop-router",
        sub: "HTTP/2 + mTLS · idempotencia ULID",
      },
      ssRc: {
        title: "security-server (Reg.Civil)",
        sub: "Verifica firma · resuelve ruta interna",
      },
      rcSvc: {
        title: "Registro Civil · svc",
        sub: "Responde con datos solicitados",
      },
      audit: {
        title: "interop-audit",
        sub: "SHA-256 chain · Merkle por epoch",
      },
      hub: {
        title: "interop-hub",
        sub: "Catálogo · CA del realm · gobernanza",
      },
      edges: {
        request: "request",
        jwsSigned: "JWS firmado",
        routeMtls: "ruteo + mTLS",
        verifyDeliver: "verifica + entrega",
        response: "response",
        kafkaEvent: "evento Kafka",
        catalogCa: "catálogo + CA",
      },
    },
    onceOnly: {
      citizen: {
        title: "Ciudadano",
        sub: "Solicita declaración pre-llenada",
      },
      micr: { title: "MiCR", sub: "Web · Mobile" },
      bff: { title: "BFF Citizen", sub: "Orquesta llamadas" },
      hacienda: { title: "Hacienda · svc", sub: "Construye declaración" },
      ssHac: { title: "security-server (Hacienda)" },
      router: { title: "interop-router" },
      ssRc: { title: "security-server (Reg.Civil)" },
      rc: { title: "Registro Civil", sub: "GET /persons/{cedula}" },
      callout: {
        title: "El ciudadano nunca dio el dato a Hacienda",
        sub: "Once-only — la institución pregunta, no el ciudadano",
      },
      edges: {
        request: "1. solicita",
        step2: "2.",
        getPrefilled: "3. GET /prefilled-return",
        needsPerson: "4. necesita persona",
        getPerson: "5. GET /persons/{cedula}",
        citizenData: "6. datos del ciudadano",
        prefilled: "7. declaración pre-llenada",
        step8: "8.",
        signClick: "9. firma con un click",
      },
    },
    identity: {
      citizen: { title: "Ciudadano", sub: "Quiero firmar un documento" },
      micr: { title: "MiCR", sub: "Tu app de gobierno" },
      iduc: {
        title: "iduc-identity",
        sub: "Login · WebAuthn · JWT RS256",
      },
      keys: {
        title: "iduc-keys",
        sub: "Ed25519 en KMS · la llave nunca sale",
      },
      signing: {
        title: "iduc-signing",
        sub: "JWS · CAdES · PAdES + sellado RFC 3161",
      },
      document: {
        title: "Documento firmado",
        sub: "Valor jurídico equivalente a manuscrita",
      },
      member: {
        title: "Cualquier institución",
        sub: "Verifica firma con la clave pública",
      },
      edges: {
        login: "login + 2FA",
        issueJwt: "JWT 15 min",
        sign: "intent de firma",
        askKey: "firmar con clave",
        signedDoc: "doc + firma JWS",
        verify: "verifica",
      },
    },
    ecosystemCitizen: {
      me: { title: "Vos", sub: "Persona o empresa" },
      micr: {
        title: "Tu app del Estado",
        sub: "En el celular o en la web",
      },
      conecta: {
        title: "Conecta CR",
        sub: "El cartero digital que conecta a las instituciones",
      },
      institutions: {
        title: "Las instituciones",
        sub: "Hacienda · Registro Civil · Caja · TSE · Migración…",
      },
      audit: {
        title: "Tu bitácora",
        sub: "Mirás quién consultó tus datos y cuándo",
      },
      signing: {
        title: "Tu firma",
        sub: "Vale igual que tu firma a mano",
      },
      edges: {
        ingreso: "entrás con tu cédula digital",
        consulta: "pedís un trámite",
        enruta: "lo lleva a la institución correcta",
        registra: "queda anotado",
        veo: "vos lo podés ver",
        firmo: "firmás con un toque",
      },
    },
    interopCitizen: {
      institutionA: {
        title: "Una institución",
        sub: "Ej: Hacienda armando tu declaración",
      },
      bridge: {
        title: "Conecta CR",
        sub: "El cartero digital del Estado",
      },
      institutionB: {
        title: "Otra institución",
        sub: "Ej: Registro Civil, donde están tus datos",
      },
      audit: {
        title: "Tu bitácora",
        sub: "Cada paso queda anotado y vos lo podés ver",
      },
      edges: {
        needs: "necesita un dato tuyo",
        sealed: "lleva el pedido como correo certificado",
        response: "le entrega el dato",
        logged: "queda anotado para vos",
      },
    },
    onceOnlyCitizen: {
      me: { title: "Vos", sub: "Querés tu declaración lista" },
      micr: { title: "Tu app del Estado", sub: "Web o celular" },
      hacienda: {
        title: "Hacienda",
        sub: "Te la arma sin que vos cargues nada",
      },
      bridge: {
        title: "Conecta CR",
        sub: "El cartero digital",
      },
      rc: {
        title: "Registro Civil",
        sub: "Tiene tu nombre, edad, familia",
      },
      callout: {
        title: "Vos no llenaste un solo formulario",
        sub: "Hacienda le preguntó a Registro Civil — no te preguntó a vos",
      },
      edges: {
        request: "1. abrís la app y la pedís",
        micrAsks: "2. la app le avisa a Hacienda",
        hacAsksRc: "3. Hacienda le pide tus datos a Registro Civil",
        rcAnswers: "4. Registro Civil los entrega",
        readyToSign: "5. Hacienda te la arma",
        signed: "6. firmás con un toque",
      },
    },
    identityCitizen: {
      me: { title: "Vos", sub: "Persona o empresa" },
      eid: {
        title: "Tu cédula digital",
        sub: "Una identidad oficial para todo el Estado",
      },
      safe: {
        title: "Tu firma guardada",
        sub: "En una caja fuerte que nadie puede abrir — ni nosotros",
      },
      document: {
        title: "Tu documento firmado",
        sub: "Vale igual que tu firma a mano",
      },
      member: {
        title: "Cualquier institución",
        sub: "Puede confirmar que fuiste vos quien firmó",
      },
      callout: {
        title: "Tu firma digital vale como tu firma a mano",
        sub: "Y nadie puede falsificarla — ni nosotros",
      },
      edges: {
        login: "entrás una sola vez",
        askSign: "pedís firmar un documento",
        signs: "se firma desde tu caja fuerte",
        anyoneVerifies: "cualquiera puede confirmar que fuiste vos",
      },
    },
    institutionPanel: {
      closeAria: "Cerrar panel",
      badgeLive: "live",
      badgePlanned: "próximamente",
      listHeading: "Catálogo de instituciones",
      listSub:
        "Cualquier institución del Estado se conecta vía Conecta CR con las mismas garantías de identidad, trazabilidad y firma. Esto es lo que cada una expone — o expondrá — al ecosistema.",
      listCount: "{n} instituciones · {l} live · {p} próximamente",
      backToList: "← Volver al catálogo",
      sectionCustodia: "Datos que custodia",
      sectionExpone: "Servicios expuestos vía Conecta",
      sectionConsume: "Consume de otras instituciones",
      emptyConsume: "No consume datos de otros — fuente de oro.",
      nodeBadgeLabel: "ver",
      diagramHint:
        "tip: hacé click en una institución para ver qué expone vía Conecta",
    },
  },
  en: {
    hero: {
      eyebrow: "costaricaasservice · core of Conecta CR",
      titleStart: "Digital government,",
      titleHighlight: "as a product.",
      descBefore:
        "Estonia took 25 years. We're aiming for 3 months. ",
      descBold1: "Institutions talk to each other",
      descMid1: " — no one asks you the same thing twice. ",
      descBold2: "Sign from your phone",
      descMid2: " with full legal weight. ",
      descBold3: "See who looks at your data",
      descAfter: " and why.",
      ctaPrimary: "Explore the ecosystem",
      ctaSecondary: "Estonian guarantees",
    },
    guarantees: {
      eyebrow: "Value proposition",
      heading: "The operating system of the digital state.",
      lead: "costaricaasservice is not middleware. It is the platform on which a country runs its digital government — from a five-minute errand to a verifiable election. The same layer that connects the civil registry to the tax authority connects healthcare to a private pharmacy, with identical guarantees of identity, traceability and signature.",
      audiencesHeading: "Who wins",
      audiences: [
        {
          tag: "Citizen",
          title: "Errands in minutes, not weeks",
          body: "You see in MiCR who queried your data and why. You sign from your phone with the same legal weight as a handwritten signature.",
        },
        {
          tag: "Institution",
          title: "You join without building connectors",
          body: "You expose a service once and the whole ecosystem consumes it with automatic audit. No more 30 bilateral contracts with every counterpart.",
        },
        {
          tag: "Jurisdiction",
          title: "One governance for the state",
          body: "Single catalog, inter-institutional usage metrics, real-time access revocation. A country-wide view, not a per-silo view.",
        },
      ],
      casesEyebrow: "Use cases",
      casesHeading: "What the digital state unlocks.",
      casesLead: "These are real flows costaricaasservice enables — all four Estonian guarantees come by construction, not by diligence.",
      cases: [
        {
          tag: "Tax",
          title: "Prefilled tax return",
          body: "The tax authority builds the return by querying the employer for wages, the healthcare fund for contributions and the civil registry for household composition — all over X-Road. The citizen opens MiCR, reviews and signs in one click. From weeks to 90 seconds.",
          actors: "Tax · Healthcare · Civil Reg. · MiCR",
          category: "gobierno",
        },
        {
          tag: "Immigration",
          title: "Paperless passport",
          body: "Immigration verifies identity with e-ID, background with the justice ministry and collects the fee from the tax authority in a single form. The citizen books pickup and signs the receipt digitally.",
          actors: "Immigration · Justice · Tax · iduc",
          category: "identidad",
        },
        {
          tag: "Health",
          title: "Electronic prescription",
          body: "The healthcare fund issues the prescription signed by the physician. The pharmacy (private member) validates it over X-Road before dispensing. The citizen sees it in MiCR alongside their full clinical history.",
          actors: "Healthcare · Pharmacy · MiCR",
          category: "bienestar",
        },
        {
          tag: "Once-only",
          title: "Address change · once and done",
          body: "The citizen updates their address in the civil registry. Tax, healthcare, insurance and the electoral body receive the change in real time over Kafka. No more uploaded utility bills at every counter.",
          actors: "Civil Reg. → Tax · Healthcare · Insurance · Electoral",
          category: "identidad",
        },
        {
          tag: "Banking",
          title: "Verifiable bank KYC",
          body: "A bank (member) opens an account by verifying identity with e-ID and background over X-Road. No passport scans, no utility bills, no 48-hour wait for validation.",
          actors: "Bank · iduc · Civil Reg.",
          category: "identidad",
        },
        {
          tag: "Procurement",
          title: "Traceable public procurement",
          body: "Companies sign bids and contracts with RFC 3161 timestamps. Every step lives in a hash-chained log, verifiable by any citizen without asking permission.",
          actors: "Tax · suppliers · audit",
          category: "gobierno",
        },
        {
          tag: "Welfare",
          title: "Social benefits, no application",
          body: "The welfare agency computes eligibility by querying incomes from tax and household composition from civil registry — the applicant submits nothing. Once-only at its purest.",
          actors: "Welfare · Tax · Civil Reg.",
          category: "bienestar",
        },
        {
          tag: "Justice",
          title: "Electronic court service",
          body: "Courts deliver subpoenas to MiCR with the same legal weight as physical service. Cryptographic delivery receipt, full traceability for the court.",
          actors: "Judiciary · notifications · signing",
          category: "justicia",
        },
        {
          tag: "Democracy",
          title: "Verifiable electronic vote",
          body: "The electoral body validates identity with e-ID, records the signed ballot and publishes a Merkle summary of the day. Each voter can verify their ballot was counted — without revealing how they voted.",
          actors: "Electoral · iduc · interop-audit",
          category: "justicia",
        },
      ],
      casesCtaCard: {
        tag: "Full catalog",
        title: "See every use case",
        body: "Catalog grouped by domain: government, identity, welfare and justice. With the flow, the actors and which Estonian guarantee each one upholds.",
        action: "Open catalog",
      },
      invariantsEyebrow: "Under the hood",
      invariantsHeading: "Four non-negotiable guarantees.",
      items: [
        {
          title: "Once-only",
          body: "The citizen provides a piece of data to the system only once. Tax authority queries Civil Registry via interop — never asks the citizen again.",
        },
        {
          title: "Decentralization",
          body: "Each institution owns its data. No central database; services are exposed and consumed over X-Road.",
        },
        {
          title: "Citizen-facing traceability",
          body: "The citizen sees who queried their data, when and for what purpose. Hash-chained, verifiable audit log.",
        },
        {
          title: "Legally binding digital signature",
          body: "Equivalent to a handwritten signature. Ed25519 keys held in a KMS, server-side signing, RFC 3161 timestamping.",
        },
      ],
    },
    casesPage: {
      eyebrow: "Catalog · use cases",
      heading: "Every flow costaricaasservice enables.",
      lead: "Use cases grouped by domain. Each one shows the actors involved and the category it belongs to — all four Estonian guarantees come by construction.",
      backToHome: "Back to home",
      casesCount: "{n} catalogued cases · {c} domains",
      categories: {
        gobierno: {
          label: "Tax & government",
          description: "Taxation, public procurement and fiscal-grade procedures.",
        },
        identidad: {
          label: "Identity & mobility",
          description: "Identity verification, personal data flows and access to the private sector.",
        },
        bienestar: {
          label: "Health & welfare",
          description: "Healthcare services and social benefits with no paperwork.",
        },
        justicia: {
          label: "Justice & democracy",
          description: "Legal notifications, judicial proceedings and verifiable voting.",
        },
      },
    },
    diagrams: {
      eyebrow: "Interactive diagrams",
      heading: "How the ecosystem works.",
      tip: "scroll to zoom · drag to pan · controls in the bottom-left",
      audience: {
        citizen: "Citizen view",
        technical: "Technical view",
        hint: "Switch between everyday language and the actual implementation.",
        activeBadge: "active",
      },
      modulesHeading: "Modules",
      modules: {
        ecosistema: {
          tag: "Big picture",
          title: "Ecosystem",
          sub: "How the pieces fit together",
        },
        interop: {
          tag: "Flow",
          title: "Conecta · X-Road",
          sub: "How a request travels between institutions",
        },
        onceOnly: {
          tag: "Scenario",
          title: "Once-only in action",
          sub: "Tax queries the Civil Registry",
        },
        identity: {
          tag: "Capability",
          title: "Identity & signing",
          sub: "How you identify yourself and sign",
        },
      },
      subs: {
        ecosistema: {
          citizen:
            "What you actually use of the state: an app, the institutions that serve you, and a digital postman that connects them — so you never have to walk into an office again.",
          technical:
            "Full view of the monorepo: gateway, audience-specific BFFs, iduc, interop, members and platform.",
        },
        interop: {
          citizen:
            "When an institution needs a piece of your data, it doesn't ask you — it asks the institution that already has it. You see every step in your log.",
          technical:
            "How an inter-member request travels: detached JWS, mTLS, ULID idempotency, and a Kafka event into the hash chain.",
        },
        onceOnly: {
          citizen:
            "The tax office builds your return by asking the civil registry for your data. You fill out no forms, upload no papers, walk to no counter — you just open the app, review and sign.",
          technical:
            "Tax service invokes civil-registry via security-server + interop-router; cacheable response, audit recorded.",
        },
        identity: {
          citizen:
            "Your identity and your digital signature live in a vault that not even we can open. You sign with one tap on your phone, and it carries the same weight as your handwritten signature on paper.",
          technical:
            "iduc-identity issues RS256 JWTs, iduc-keys custodies Ed25519 in a KMS (Vault), iduc-signing produces JWS/CAdES/PAdES with RFC 3161 timestamping.",
        },
      },
    },
    deliverables: {
      eyebrow: "System deliverables",
      heading: "What ships when costaricaasservice goes live.",
      lead: "One monorepo, two readings. What you receive based on your role — and what is inside the box, repo by repo.",
      view: {
        business: "Business view",
        technical: "Technical view",
        hint: "Switch between what you get and what is inside.",
        activeBadge: "active",
      },
      business: {
        intro:
          "Four packages aligned with who consumes the system: the State that operates it, the institution that extends it, the citizen who uses it, and the platform layer that holds the legal contract.",
        cards: [
          {
            tag: "Realm · jurisdiction",
            title: "For the State that operates it",
            sub: "Full realm control plane",
            items: [
              "Admin control plane (members, keys, billing)",
              "Interoperability hub with the realm CA",
              "Inter-member service catalog governance",
              "End-to-end observability (traces, metrics, p99 < 500ms)",
              "Runbooks and incident response",
            ],
          },
          {
            tag: "Member · institution",
            title: "For each institution",
            sub: "Connect to the ecosystem in days, not years",
            items: [
              "Deployable security-server (JWS signing + mTLS turnkey)",
              "Member BFF + operator portal",
              "interop-client SDK (expose and consume remote services)",
              "Service template (cri-templates-service) production-ready",
              "Catalog onboarding and X.509 certificate issuance",
            ],
          },
          {
            tag: "Citizen · MiCR",
            title: "For every resident",
            sub: "A single front door to the State",
            items: [
              "MiCR Web (Next.js 15) and MiCR Mobile (Flutter melos)",
              "Single digital identity + WebAuthn/passkey login",
              "Legally-binding digital signature (Ed25519 in KMS)",
              "Access log: who looked at your data and why",
              "Legal notifications and signed files (PDF + CAdES/PAdES)",
            ],
          },
          {
            tag: "Platform · cross-cutting",
            title: "Shared capabilities",
            sub: "The legal contract and traceability that hold it together",
            items: [
              "Hash-chained audit (SHA-256 + Merkle per epoch)",
              "Electronic notifications with legal weight",
              "Signed-file service (PDF + CAdES/PAdES)",
              "Feature flags per realm and per member",
              "KMS with envelope encryption (Vault Transit in dev, AWS KMS in prod)",
            ],
          },
        ],
      },
      technical: {
        intro:
          "The six monorepo areas, with the concrete repos deployed in each one. Every service owns its database, multi-tenant from day one (schema-per-realm).",
        cards: [
          {
            tag: "gateway",
            title: "System ingress",
            sub: "Central gateway + audience-specific BFFs",
            items: [
              "cri-gateway-api · JWT RS256 + RBAC + rate limit + Roaring revocation",
              "cri-bff-citizen · orchestrates MiCR (web + mobile)",
              "cri-bff-member · operator portal",
              "cri-bff-admin · costaricaasservice control plane",
            ],
          },
          {
            tag: "iduc",
            title: "Single Digital Identity",
            sub: "Authentication, keys and signing",
            items: [
              "cri-svc-iduc-identity · login + WebAuthn + JWT RS256",
              "cri-svc-iduc-keys · Ed25519 custodied in KMS",
              "cri-svc-iduc-signing · JWS / CAdES / PAdES + RFC 3161 timestamping",
            ],
          },
          {
            tag: "interop",
            title: "Conecta · X-Road analog",
            sub: "Central hub + federated data plane",
            items: [
              "cri-svc-interop-hub · catalog + realm CA + governance",
              "cri-svc-security-server · daemon each member installs",
              "cri-svc-interop-router · HTTP/2 + mTLS data plane",
              "cri-svc-interop-audit · SHA-256 hash chain + Merkle epoch",
            ],
          },
          {
            tag: "members",
            title: "Institutions",
            sub: "Mock services ready to swap for real ones",
            items: [
              "cri-svc-registro-civil · TSE mock with 1000 seeded persons",
              "cri-svc-hacienda · pre-filled tax return (consumes Reg.Civil)",
              "+ MIT, CCSS, MOPT… (catalog-scalable)",
            ],
          },
          {
            tag: "platform",
            title: "Cross-cutting platform",
            sub: "Shared capabilities with legal value",
            items: [
              "cri-svc-audit · citizen view of the inter-member log",
              "cri-svc-notifications · notifications with legal weight",
              "cri-svc-files · signed PDFs (S3 storage in prod)",
              "cri-svc-feature-flags · per-realm/per-member flags",
              "cri-templates-service · scaffold for new services",
            ],
          },
          {
            tag: "frontends",
            title: "Client surfaces",
            sub: "Web + mobile + institutional landing",
            items: [
              "cri-landing · public site (Next.js 15)",
              "cri-web-micr · MiCR for citizens (Next.js 15)",
              "cri-web-admin · members and realm control plane",
              "cri-mobile · MiCR Mobile (Flutter melos, iOS + Android)",
            ],
          },
        ],
      },
    },
    footer: {
      coreOf: "core of Conecta CR",
    },
    langSwitch: { aria: "Change language" },
    ecosystem: {
      citizen: { title: "Citizen", sub: "MiCR Web · Mobile" },
      operator: { title: "Member operator", sub: "Institution portal" },
      admin: { title: "costaricaasservice admin", sub: "Control plane" },
      gateway: {
        title: "API Gateway",
        sub: "JWT RS256 · RBAC · rate limit · Roaring revocation",
      },
      bffCitizen: { title: "BFF Citizen", sub: "Orchestrates MiCR" },
      bffMember: { title: "BFF Member", sub: "Operator portal" },
      bffAdmin: { title: "BFF Admin", sub: "costaricaasservice control plane" },
      iducGroup: "iduc · identity",
      iducIdentity: {
        title: "iduc-identity",
        sub: "Login · WebAuthn · JWT",
      },
      iducKeys: { title: "iduc-keys", sub: "Ed25519 in KMS" },
      iducSigning: { title: "iduc-signing", sub: "JWS · CAdES · PAdES" },
      interopGroup: "interop · Conecta",
      interopHub: { title: "interop-hub", sub: "Catalog · realm CA" },
      securityServer: {
        title: "security-server",
        sub: "Detached JWS · mTLS",
      },
      interopRouter: {
        title: "interop-router",
        sub: "Data plane · idempotency",
      },
      interopAudit: {
        title: "interop-audit",
        sub: "Hash chain · Merkle epoch",
      },
      membersGroup: "members · institutions",
      registroCivil: {
        title: "civil-registry",
        sub: "TSE mock · 1000 people",
      },
      hacienda: { title: "tax-authority", sub: "Prefilled return" },
      otherMembers: {
        title: "+ other members",
        sub: "MIT · CCSS · MOPT…",
      },
      platformGroup: "platform · cross-cutting",
      platformAudit: {
        title: "audit (citizen)",
        sub: "View of who queried",
      },
      notifications: { title: "notifications", sub: "Legally binding" },
      files: { title: "files", sub: "Signed PDFs" },
      featureFlags: {
        title: "feature-flags",
        sub: "Per realm / member",
      },
      edges: {
        kms: "KMS",
        expose: "exposes",
        exposeConsume: "exposes/consumes",
        mtls: "mTLS",
        event: "event",
      },
    },
    interop: {
      haciendaSvc: {
        title: "Tax authority · svc",
        sub: "Originator of the request",
      },
      ssHac: {
        title: "security-server (Tax)",
        sub: "Signs detached JWS with member key",
      },
      router: {
        title: "interop-router",
        sub: "HTTP/2 + mTLS · ULID idempotency",
      },
      ssRc: {
        title: "security-server (Civil Reg.)",
        sub: "Verifies signature · resolves internal route",
      },
      rcSvc: {
        title: "Civil Registry · svc",
        sub: "Returns the requested data",
      },
      audit: {
        title: "interop-audit",
        sub: "SHA-256 chain · Merkle per epoch",
      },
      hub: {
        title: "interop-hub",
        sub: "Catalog · realm CA · governance",
      },
      edges: {
        request: "request",
        jwsSigned: "signed JWS",
        routeMtls: "route + mTLS",
        verifyDeliver: "verify + deliver",
        response: "response",
        kafkaEvent: "Kafka event",
        catalogCa: "catalog + CA",
      },
    },
    onceOnly: {
      citizen: {
        title: "Citizen",
        sub: "Requests prefilled return",
      },
      micr: { title: "MiCR", sub: "Web · Mobile" },
      bff: { title: "BFF Citizen", sub: "Orchestrates calls" },
      hacienda: { title: "Tax authority · svc", sub: "Builds the return" },
      ssHac: { title: "security-server (Tax)" },
      router: { title: "interop-router" },
      ssRc: { title: "security-server (Civil Reg.)" },
      rc: { title: "Civil Registry", sub: "GET /persons/{id}" },
      callout: {
        title: "The citizen never gave this data to the tax authority",
        sub: "Once-only — the institution asks, not the citizen",
      },
      edges: {
        request: "1. requests",
        step2: "2.",
        getPrefilled: "3. GET /prefilled-return",
        needsPerson: "4. needs person",
        getPerson: "5. GET /persons/{id}",
        citizenData: "6. citizen data",
        prefilled: "7. prefilled return",
        step8: "8.",
        signClick: "9. one-click signature",
      },
    },
    identity: {
      citizen: { title: "Citizen", sub: "I want to sign a document" },
      micr: { title: "MiCR", sub: "Your government app" },
      iduc: {
        title: "iduc-identity",
        sub: "Login · WebAuthn · RS256 JWT",
      },
      keys: {
        title: "iduc-keys",
        sub: "Ed25519 in KMS · the key never leaves",
      },
      signing: {
        title: "iduc-signing",
        sub: "JWS · CAdES · PAdES + RFC 3161 timestamp",
      },
      document: {
        title: "Signed document",
        sub: "Same legal weight as a handwritten signature",
      },
      member: {
        title: "Any institution",
        sub: "Verifies the signature with the public key",
      },
      edges: {
        login: "login + 2FA",
        issueJwt: "JWT 15 min",
        sign: "sign intent",
        askKey: "sign with key",
        signedDoc: "doc + JWS signature",
        verify: "verify",
      },
    },
    ecosystemCitizen: {
      me: { title: "You", sub: "Person or business" },
      micr: {
        title: "Your state app",
        sub: "On your phone or in a browser",
      },
      conecta: {
        title: "Conecta CR",
        sub: "The digital postman that connects every institution",
      },
      institutions: {
        title: "The institutions",
        sub: "Tax · Civil Registry · Healthcare · Electoral · Immigration…",
      },
      audit: {
        title: "Your log",
        sub: "You see who looked at your data and when",
      },
      signing: {
        title: "Your signature",
        sub: "Same weight as your handwritten signature",
      },
      edges: {
        ingreso: "you log in with your digital ID",
        consulta: "you request something",
        enruta: "it's delivered to the right institution",
        registra: "it gets logged",
        veo: "you can see it",
        firmo: "you sign with a tap",
      },
    },
    interopCitizen: {
      institutionA: {
        title: "One institution",
        sub: "e.g. Tax putting your return together",
      },
      bridge: {
        title: "Conecta CR",
        sub: "The digital postman of the state",
      },
      institutionB: {
        title: "Another institution",
        sub: "e.g. Civil Registry, where your data already lives",
      },
      audit: {
        title: "Your log",
        sub: "Every step is recorded and you can see it",
      },
      edges: {
        needs: "needs a piece of your data",
        sealed: "carries the request like certified mail",
        response: "hands over the data",
        logged: "gets logged for you",
      },
    },
    onceOnlyCitizen: {
      me: { title: "You", sub: "You want your return done" },
      micr: { title: "Your state app", sub: "Web or mobile" },
      hacienda: {
        title: "Tax authority",
        sub: "Builds your return without you uploading anything",
      },
      bridge: {
        title: "Conecta CR",
        sub: "The digital postman",
      },
      rc: {
        title: "Civil Registry",
        sub: "Has your name, age, household",
      },
      callout: {
        title: "You filled in zero forms",
        sub: "Tax asked the Civil Registry — not you",
      },
      edges: {
        request: "1. you open the app and ask for it",
        micrAsks: "2. the app tells the tax office",
        hacAsksRc: "3. Tax asks the Civil Registry for your data",
        rcAnswers: "4. Civil Registry hands it over",
        readyToSign: "5. Tax puts your return together",
        signed: "6. you sign with a tap",
      },
    },
    identityCitizen: {
      me: { title: "You", sub: "Person or business" },
      eid: {
        title: "Your digital ID",
        sub: "One official identity for the whole state",
      },
      safe: {
        title: "Your signature, locked away",
        sub: "In a vault no one can open — not even us",
      },
      document: {
        title: "Your signed document",
        sub: "Same weight as your handwritten signature",
      },
      member: {
        title: "Any institution",
        sub: "Can confirm that you were the one who signed",
      },
      callout: {
        title: "Your digital signature is as binding as your handwritten one",
        sub: "And no one can forge it — not even us",
      },
      edges: {
        login: "you log in just once",
        askSign: "you ask to sign a document",
        signs: "it's signed from your vault",
        anyoneVerifies: "anyone can confirm it was you",
      },
    },
    institutionPanel: {
      closeAria: "Close panel",
      badgeLive: "live",
      badgePlanned: "coming soon",
      listHeading: "Institutions catalog",
      listSub:
        "Any institution of the state plugs into Conecta CR with the same guarantees of identity, traceability and signature. Here is what each one exposes — or will expose — to the ecosystem.",
      listCount: "{n} institutions · {l} live · {p} coming soon",
      backToList: "← Back to catalog",
      sectionCustodia: "Data it custodies",
      sectionExpone: "Services exposed via Conecta",
      sectionConsume: "Consumes from other institutions",
      emptyConsume: "Does not consume from others — it's a golden source.",
      nodeBadgeLabel: "view",
      diagramHint: "tip: click an institution to see what it exposes via Conecta",
    },
  },
};

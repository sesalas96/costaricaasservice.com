import { ImageResponse } from "next/og";

export const alt = "costaricaasservice — Gobierno digital, as a product";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default async function OpenGraphImage() {
  return new ImageResponse(
    (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          backgroundColor: "#07090d",
          backgroundImage:
            "radial-gradient(900px 500px at 85% -10%, rgba(29,79,196,0.35), transparent 60%), radial-gradient(700px 400px at 5% 110%, rgba(206,17,38,0.25), transparent 60%)",
          color: "#e7eaf0",
          fontFamily: "system-ui, sans-serif",
          position: "relative",
        }}
      >
        {/* Bandera CR — franjas verticales arriba */}
        <div style={{ display: "flex", height: 14, width: "100%" }}>
          <div style={{ flex: 1, background: "#002b7f" }} />
          <div style={{ flex: 1, background: "#f7f5ef" }} />
          <div style={{ flex: 2, background: "#ce1126" }} />
          <div style={{ flex: 1, background: "#f7f5ef" }} />
          <div style={{ flex: 1, background: "#002b7f" }} />
        </div>

        <div
          style={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "space-between",
            flex: 1,
            padding: "60px 72px 56px",
          }}
        >
          {/* Eyebrow */}
          <div
            style={{
              display: "flex",
              alignItems: "center",
              gap: 14,
              fontSize: 18,
              fontFamily: "monospace",
              letterSpacing: "0.22em",
              textTransform: "uppercase",
              color: "#8b93a3",
            }}
          >
            <div
              style={{
                width: 12,
                height: 12,
                borderRadius: 999,
                background: "#1d4fc4",
              }}
            />
            costaricaasservice · propuesta independiente inspirada en e-Estonia
          </div>

          {/* Título principal */}
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              gap: 28,
            }}
          >
            <div
              style={{
                fontSize: 96,
                fontWeight: 700,
                lineHeight: 1.02,
                letterSpacing: "-0.03em",
                color: "#e7eaf0",
                display: "flex",
                flexWrap: "wrap",
                gap: 18,
              }}
            >
              <span>Gobierno digital,</span>
              <span
                style={{
                  backgroundImage:
                    "linear-gradient(100deg, #1d4fc4 0%, #ffffff 52%, #ef3145 100%)",
                  backgroundClip: "text",
                  color: "transparent",
                }}
              >
                as a product.
              </span>
            </div>

            <div
              style={{
                fontSize: 28,
                lineHeight: 1.4,
                color: "#8b93a3",
                maxWidth: 980,
              }}
            >
              Identidad ciudadana, interoperabilidad entre instituciones y firma
              digital — pensados para que el ciudadano dé un dato una sola vez.
            </div>
          </div>

          {/* Footer */}
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              fontSize: 18,
              fontFamily: "monospace",
              color: "#8b93a3",
              letterSpacing: "0.18em",
              textTransform: "uppercase",
              borderTop: "1px solid #232936",
              paddingTop: 24,
            }}
          >
            <div style={{ display: "flex", alignItems: "center", gap: 14 }}>
              <div style={{ display: "flex", gap: 6 }}>
                <div
                  style={{
                    width: 10,
                    height: 10,
                    borderRadius: 999,
                    background: "#10b981",
                  }}
                />
                <div
                  style={{
                    width: 10,
                    height: 10,
                    borderRadius: 999,
                    background: "#38bdf8",
                  }}
                />
                <div
                  style={{
                    width: 10,
                    height: 10,
                    borderRadius: 999,
                    background: "#f59e0b",
                  }}
                />
                <div
                  style={{
                    width: 10,
                    height: 10,
                    borderRadius: 999,
                    background: "#a78bfa",
                  }}
                />
              </div>
              identidad · interop · members · platform
            </div>
            <div>costaricaasservice.com</div>
          </div>
        </div>
      </div>
    ),
    size,
  );
}

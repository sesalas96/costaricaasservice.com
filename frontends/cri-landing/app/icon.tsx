import { ImageResponse } from "next/og";

export const size = { width: 32, height: 32 };
export const contentType = "image/png";

export default function Icon() {
  return new ImageResponse(
    (
      <div
        style={{
          width: "100%",
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          background:
            "linear-gradient(135deg, #002b7f 0%, #1d4fc4 50%, #ef3145 100%)",
          color: "#fff",
          fontSize: 22,
          fontWeight: 700,
          fontFamily: "system-ui",
          borderRadius: 6,
          letterSpacing: "-0.02em",
        }}
      >
        cr
      </div>
    ),
    size,
  );
}

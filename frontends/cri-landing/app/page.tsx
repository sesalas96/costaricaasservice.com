import Hero from "./components/Hero";
import Guarantees from "./components/Guarantees";
import DiagramSection from "./components/DiagramSection";
import Footer from "./components/Footer";

export default function Page() {
  return (
    <main className="min-h-screen">
      <Hero />
      <Guarantees />
      <DiagramSection />
      <Footer />
    </main>
  );
}

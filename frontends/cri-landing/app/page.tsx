import Hero from "./components/Hero";
import Guarantees from "./components/Guarantees";
import DiagramSection from "./components/DiagramSection";
import Deliverables from "./components/Deliverables";
import Footer from "./components/Footer";

export default function Page() {
  return (
    <main className="min-h-screen">
      <Hero />
      <Guarantees />
      <DiagramSection />
      <Deliverables />
      <Footer />
    </main>
  );
}

import Hero from "./components/Hero";
import Guarantees from "./components/Guarantees";
import MembersCarousel from "./components/MembersCarousel";
import DiagramSection from "./components/DiagramSection";
import DemoChat from "./components/DemoChat";
import Deliverables from "./components/Deliverables";
import Footer from "./components/Footer";
import { InstitutionProvider } from "./components/InstitutionContext";
import InstitutionPanel from "./components/InstitutionPanel";

export default function Page() {
  return (
    <main className="min-h-screen">
      <Hero />
      <Guarantees />
      <InstitutionProvider>
        <InstitutionPanel />
        <MembersCarousel />
        <DiagramSection />
      </InstitutionProvider>
      <DemoChat />
      <Deliverables />
      <Footer />
    </main>
  );
}

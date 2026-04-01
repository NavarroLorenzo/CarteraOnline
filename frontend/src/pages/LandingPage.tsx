import { LandingCTA } from "../components/landing/LandingCTA";
import { LandingFeatures } from "../components/landing/LandingFeatures";
import { LandingHeader } from "../components/landing/LandingHeader";
import { LandingHero } from "../components/landing/LandingHero";
import { LandingInfo } from "../components/landing/LandingInfo";

export function LandingPage() {
  return (
    <div className="landing-shell">
      <div className="landing-shell__glow landing-shell__glow--top" />
      <div className="landing-shell__glow landing-shell__glow--bottom" />

      <div className="landing-container">
        <LandingHeader />
        <LandingHero />
        <LandingFeatures />
        <LandingInfo />
        <LandingCTA />
      </div>
    </div>
  );
}

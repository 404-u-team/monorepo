import type { JSX } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { Hero, Benefits, HowItWorks, TargetAudience, CallToAction } from "@/widgets/Landing";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index(): JSX.Element {
  return (
    <>
      <Hero />
      <Benefits />
      <HowItWorks />
      <TargetAudience />
      <CallToAction />
    </>
  );
}

import type { JSX } from "react";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index(): JSX.Element {
  return (
    <div className="p-2">
      <h3>Welcome! This is the home page.</h3>
    </div>
  );
}

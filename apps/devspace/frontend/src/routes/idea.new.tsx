import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import { CreateIdeaForm } from "@/features/idea/create";

export const Route = createFileRoute("/idea/new")({
  beforeLoad: () => {
    // We might want to check auth here if context has it,
    // but UserStore is usually accessed via useStore
  },
  component: CreateIdeaPage,
});

function CreateIdeaPage(): JSX.Element {
  return (
    <div style={{ padding: "0 24px" }}>
      <CreateIdeaForm />
    </div>
  );
}

import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import { IdeaDetail } from "@/widgets/IdeaDetail";

export const Route = createFileRoute("/idea/$ideaId")({
  component: IdeaDetailPage,
});

function IdeaDetailPage(): JSX.Element {
  return (
    <div style={{ padding: "0 24px" }}>
      <IdeaDetail />
    </div>
  );
}

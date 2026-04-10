import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import { AdminPanel } from "@/widgets/AdminPanel";

export const Route = createFileRoute("/admin")({
  component: AdminPage,
});

function AdminPage(): JSX.Element {
  return <AdminPanel />;
}

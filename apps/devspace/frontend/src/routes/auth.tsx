import { createFileRoute } from "@tanstack/react-router";

import { AuthForm } from "@/features/auth";

export interface AuthSearch {
  redirect?: string | undefined;
}

export const Route = createFileRoute("/auth")({
  validateSearch: (search: Record<string, unknown>): AuthSearch => ({
    redirect:
      typeof search.redirect === "string" && search.redirect !== "" ? search.redirect : undefined,
  }),
  component: AuthPage,
});

function AuthPage(): React.JSX.Element {
  return (
    <div
      style={{
        display: "flex",
        minHeight: "100vh",
        padding: "2rem 1rem",
        backgroundColor: "var(--bg--main)",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <AuthForm />
    </div>
  );
}

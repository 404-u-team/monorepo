import { createRootRoute, Outlet, useRouterState } from "@tanstack/react-router";
import { clsx } from "clsx";

import { PageLoader } from "@/shared/ui";
import { Footer } from "@/widgets/Footer";
import { Navbar } from "@/widgets/Navbar";

import styles from "./__root.module.scss";

function RootComponent(): React.JSX.Element {
  const isLoading = useRouterState({ select: (s) => s.isLoading });

  return (
    <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
      <PageLoader />
      <Navbar />
      <main className={clsx(styles.main, isLoading && styles.blurred)}>
        <Outlet />
      </main>
      <Footer />
    </div>
  );
}

export const Route = createRootRoute({
  component: RootComponent,
});

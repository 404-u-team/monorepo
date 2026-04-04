import { RouterProvider, createRouter } from "@tanstack/react-router";
import { configure } from "mobx";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "@/app/styles/index.scss";
// Import the generated route tree
import { routeTree } from "@/app/generated/routeTree.gen";
import { verifyInterceptors } from "@/app/providers/apiInterceptors";
import { rootStore } from "@/app/providers/store";
import { ThemeProvider } from "@/app/providers/ThemeProvider";
import { StoreContext } from "@/shared/lib/store";

// Create a new router instance
const router = createRouter({ routeTree });

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

// MobX configuration
configure({
  enforceActions: "always",
  computedRequiresReaction: true,
  reactionRequiresObservable: true,
});

// Initial auth check
void rootStore.userStore.checkAuth();

// Configure Axios with FSD-compliant rules
verifyInterceptors();

// Render the app
const rootElement = document.getElementById("root");
if (rootElement !== null && !rootElement.innerHTML) {
  const root = createRoot(rootElement);
  root.render(
    <StrictMode>
      <ThemeProvider>
        <StoreContext.Provider value={rootStore}>
          <RouterProvider router={router} />
        </StoreContext.Provider>
      </ThemeProvider>
    </StrictMode>,
  );
}

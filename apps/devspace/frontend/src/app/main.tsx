import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { configure } from 'mobx';
import '@/app/styles/index.scss'
import { StoreContext } from '@/shared/lib/store';
import { rootStore } from '@/app/providers/store';

import { verifyInterceptors } from '@/app/providers/apiInterceptors';

// Import the generated route tree
import { routeTree } from '@/app/generated/routeTree.gen'

// Create a new router instance
const router = createRouter({ routeTree })

// Register the router instance for type safety
declare module '@tanstack/react-router' {
    interface Register {
        router: typeof router
    }
}

// MobX configuration
configure({
    enforceActions: 'always',
    computedRequiresReaction: true,
    reactionRequiresObservable: true,
});

// Initial auth check
void rootStore.userStore.checkAuth();

// Configure Axios with FSD-compliant rules
verifyInterceptors();

// Render the app
const rootElement = document.getElementById('root')
if (rootElement !== null && !rootElement.innerHTML) {
    const root = createRoot(rootElement)
    root.render(
        <StrictMode>
            <StoreContext.Provider value={rootStore}>
                <RouterProvider router={router} />
            </StoreContext.Provider>
        </StrictMode>,
    )
}

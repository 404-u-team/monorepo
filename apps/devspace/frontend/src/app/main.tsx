import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { configure } from 'mobx';
import '@/app/styles/index.scss'
import { RootStore, StoreContext } from '@/app/providers/store';

// Import the generated route tree
// eslint-disable-next-line import-x/no-unresolved
import { routeTree } from '@/app/generated/routeTree.gen'

// Create a new router instance
// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
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

// Create the global store
const store = new RootStore();

// Render the app
const rootElement = document.getElementById('root')
if (rootElement !== null && !rootElement.innerHTML) {
    const root = createRoot(rootElement)
    root.render(
        <StrictMode>
            <StoreContext.Provider value={store}>
                <RouterProvider router={router} />
            </StoreContext.Provider>
        </StrictMode>,
    )
}

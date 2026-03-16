import { createRootRoute, Outlet } from '@tanstack/react-router'
import { Navbar } from '@/widgets/Navbar'
import { Footer } from '@/widgets/Footer'

export const Route = createRootRoute({
    component: () => (
        <div style={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            <Navbar />
            <main style={{ flex: 1 }}>
                <Outlet />
            </main>
            <Footer />
        </div>
    ),
})

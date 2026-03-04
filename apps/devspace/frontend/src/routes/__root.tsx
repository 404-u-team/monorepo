import { createRootRoute, Outlet } from '@tanstack/react-router'
import { Navbar } from '@/widgets/Navbar'

export const Route = createRootRoute({
    component: () => (
        <>
            <Navbar />
            <Outlet />
        </>
    ),
})

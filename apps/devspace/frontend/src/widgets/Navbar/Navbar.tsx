import { useState, type JSX } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'
import { Bell, ChevronDown, Menu, X } from 'lucide-react'
import { observer } from 'mobx-react-lite'
import { useStore } from '@/shared/lib/store'
import { Button } from '@/shared/ui'
import { clsx } from 'clsx'
import styles from './Navbar.module.scss'

export const Navbar = observer(function Navbar(): JSX.Element {
    const { userStore } = useStore()
    const navigate = useNavigate()
    const [isMenuOpen, setIsMenuOpen] = useState(false)

    const toggleMenu = (): void => { setIsMenuOpen(!isMenuOpen) }
    const closeMenu = (): void => { setIsMenuOpen(false) }

    const navLinks = [
        { to: '/', label: 'Сообщество' },
        { to: '/projects', label: 'Проекты' },
        { to: '/ideas', label: 'Идеи' },
    ]

    return (
        <nav className={styles.navbar}>
            <Link to="/" className={styles.logo} onClick={closeMenu}>
                <img
                    src="/DevSpaceLogo-removebg.png"
                    alt="DevSpace"
                    className={styles.logoImage}
                />
            </Link>

            <div className={styles.actions}>
                <ul className={styles.nav}>
                    {navLinks.map((link) => (
                        <li key={link.to}>
                            <Link to={link.to} className={styles.link}>
                                {link.label}
                            </Link>
                        </li>
                    ))}
                </ul>

                {userStore.isAuthenticated ? (
                    <>
                        <Button
                            className={styles.createButton}
                            onClick={() => { void navigate({ to: '/' }) }}
                        >
                            Мои проекты
                        </Button>

                        <div className={styles.userControls}>
                            <button className={styles.iconButton} aria-label="Уведомления">
                                <Bell size={24} />
                            </button>

                            <div className={styles.avatarWrapper}>
                                <div className={styles.avatar} />
                                <ChevronDown />
                            </div>
                        </div>
                    </>
                ) : (
                    <div className={styles.userControls}>
                        <Button onClick={() => { void navigate({ to: '/auth' }) }}>
                            Войти
                        </Button>
                    </div>
                )}
                <div className={styles.burgerWrapper}>
                    <button
                        className={styles.burgerButton}
                        onClick={toggleMenu}
                        aria-label={isMenuOpen ? 'Закрыть меню' : 'Открыть меню'}
                    >
                        {isMenuOpen ? <X size={24} /> : <Menu size={24} />}
                    </button>
                </div>
            </div>

            {/* Mobile Menu */}
            <div className={clsx(styles.mobileMenu, isMenuOpen && styles.opened)}>
                <div className={styles.mobileMenuContent}>
                    <ul className={styles.mobileNav}>
                        {navLinks.map((link) => (
                            <li key={link.to}>
                                <Link
                                    to={link.to}
                                    className={styles.mobileLink}
                                    onClick={closeMenu}
                                >
                                    {link.label}
                                </Link>
                            </li>
                        ))}
                    </ul>

                    <div className={styles.mobileActions}>
                        {userStore.isAuthenticated ? (
                            <Button
                                className={styles.mobileCreateButton}
                                onClick={() => {
                                    void navigate({ to: '/' })
                                    closeMenu()
                                }}
                                fullWidth
                            >
                                Мои проекты
                            </Button>
                        ) : (
                            <Button
                                onClick={() => {
                                    void navigate({ to: '/auth' })
                                    closeMenu()
                                }}
                                fullWidth
                            >
                                Войти
                            </Button>
                        )}
                    </div>
                </div>
                <div className={styles.overlay} onClick={closeMenu} />
            </div>
        </nav>
    )
})
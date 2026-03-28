import { useState, useRef, useEffect, type JSX } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'
import { Bell, ChevronDown, Menu, X } from 'lucide-react'
import { observer } from 'mobx-react-lite'
import { useUserStore } from '@/entities/user'
import { Button, Logo, UserAvatar } from '@/shared/ui'
import { ThemeToggle } from '@/features/theme'
import { clsx } from 'clsx'
import styles from './Navbar.module.scss'

export const Navbar = observer(function Navbar(): JSX.Element {
    const userStore = useUserStore()
    const navigate = useNavigate()
    const [isMenuOpen, setIsMenuOpen] = useState(false)
    const [isDropdownOpen, setIsDropdownOpen] = useState(false)
    const dropdownReference = useRef<HTMLDivElement>(null)

    const toggleMenu = (): void => { setIsMenuOpen(!isMenuOpen) }
    const closeMenu = (): void => { setIsMenuOpen(false) }

    const handleLogout = (): void => {
        userStore.invalidateToken()
        userStore.invalidateUser()
        setIsDropdownOpen(false)
        void navigate({ to: '/' })
    }

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent): void => {
            if (dropdownReference.current !== null && !dropdownReference.current.contains(event.target as Node)) {
                setIsDropdownOpen(false)
            }
        }
        document.addEventListener('mousedown', handleClickOutside)
        return (): void => { document.removeEventListener('mousedown', handleClickOutside) }
    }, [])

    const navLinks = [
        { to: '/community', label: 'Сообщество' },
        { to: '/projects', label: 'Проекты' },
        { to: '/ideas', label: 'Идеи' },
    ]

    const currentUser = userStore.user

    return (
        <nav className={styles.navbar}>
            <Link to="/" className={styles.logo} onClick={closeMenu}>
                <Logo className={styles.logoImage} />
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
                            onClick={() => { void navigate({ to: '/profile' }) }}
                        >
                            Мои проекты
                        </Button>

                        <div className={styles.userControls}>
                            <ThemeToggle />

                            <button className={styles.iconButton} aria-label="Уведомления">
                                <Bell size={24} />
                            </button>

                            <div ref={dropdownReference} className={styles.avatarWrapper}>
                                <button
                                    className={styles.avatarButton}
                                    onClick={() => { setIsDropdownOpen(!isDropdownOpen) }}
                                    aria-label="Меню пользователя"
                                    aria-expanded={isDropdownOpen}
                                >
                                    <UserAvatar
                                        avatarUrl={currentUser?.avatarUrl}
                                        nickname={currentUser?.nickname}
                                        size={32}
                                        className={styles.avatar}
                                    />
                                    <ChevronDown
                                        size={16}
                                        className={clsx(styles.chevron, isDropdownOpen && styles.chevronOpen)}
                                    />
                                </button>

                                {isDropdownOpen && (
                                    <div className={styles.dropdown}>
                                        {currentUser !== undefined && (
                                            <div className={styles.dropdownHeader}>
                                                <span className={styles.dropdownNickname}>
                                                    {currentUser.nickname}
                                                </span>
                                                <span className={styles.dropdownEmail}>
                                                    {currentUser.email}
                                                </span>
                                            </div>
                                        )}
                                        <div className={styles.dropdownDivider} />
                                        <button
                                            className={styles.dropdownItem}
                                            onClick={() => {
                                                setIsDropdownOpen(false)
                                                void navigate({ to: '/profile' })
                                            }}
                                        >
                                            Личный кабинет
                                        </button>
                                        {currentUser !== undefined && (
                                            <button
                                                className={styles.dropdownItem}
                                                onClick={() => {
                                                    setIsDropdownOpen(false)
                                                    void navigate({ to: '/users/$userId', params: { userId: currentUser.id } })
                                                }}
                                            >
                                                Мой профиль
                                            </button>
                                        )}
                                        <div className={styles.dropdownDivider} />
                                        <button
                                            className={clsx(styles.dropdownItem, styles.dropdownItemDanger)}
                                            onClick={handleLogout}
                                        >
                                            Выйти
                                        </button>
                                    </div>
                                )}
                            </div>
                        </div>
                    </>
                ) : (
                    <div className={styles.userControls}>
                        <ThemeToggle />
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
                            <>
                                <Button
                                    className={styles.mobileCreateButton}
                                    onClick={() => {
                                        void navigate({ to: '/profile' })
                                        closeMenu()
                                    }}
                                    fullWidth
                                >
                                    Мои проекты
                                </Button>
                                <Button
                                    variant="outline"
                                    onClick={handleLogout}
                                    fullWidth
                                >
                                    Выйти
                                </Button>
                            </>
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

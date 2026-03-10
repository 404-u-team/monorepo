import type { JSX } from 'react'
import { Link, useNavigate } from '@tanstack/react-router'
import { Bell, ChevronDown } from 'lucide-react'
import { observer } from 'mobx-react-lite'
import { useStore } from '@/shared/lib/store'
import { Button } from '@/shared/ui'
import styles from './Navbar.module.scss'

export const Navbar = observer(function Navbar(): JSX.Element {
    const { userStore } = useStore()
    const navigate = useNavigate()

    return (
        <nav className={styles.navbar}>
            <Link to="/" className={styles.logo}>
                <img
                    src="/DevSpaceLogo-removebg.png"
                    alt="DevSpace"
                    className={styles.logoImage}
                />
            </Link>

            <div className={styles.actions}>
                <ul className={styles.nav}>
                    <li>
                        <Link to="/" className={styles.link}>
                            Сообщество
                        </Link>
                    </li>
                    <li>
                        <Link to="/" className={styles.link}>
                            Проекты
                        </Link>
                    </li>
                    <li>
                        <Link to="/" className={styles.link}>
                            Идеи
                        </Link>
                    </li>
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
            </div>
        </nav>
    )
})
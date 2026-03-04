import type { JSX } from 'react'
import { Link } from '@tanstack/react-router'
import { Bell, ChevronDown } from 'lucide-react'
import styles from './Navbar.module.scss'

export function Navbar(): JSX.Element {
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

                <Link to="/" className={styles.createButton}>
                    Мои проекты
                </Link>

                <div className={styles.userControls}>
                    <button className={styles.iconButton} aria-label="Уведомления">
                        <Bell size={24} />
                    </button>

                    <div className={styles.avatarWrapper}>
                        <div className={styles.avatar} />
                        <ChevronDown />
                    </div>
                </div>
            </div>
        </nav>
    )
}
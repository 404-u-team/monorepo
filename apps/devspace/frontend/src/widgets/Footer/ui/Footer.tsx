import { Link } from "@tanstack/react-router";
import type { JSX } from "react";

import { Logo } from "@/shared/ui";

import styles from "./Footer.module.scss";

export function Footer(): JSX.Element {
  return (
    <footer className={styles.footer}>
      <div className={styles.container}>
        <div className={styles.top}>
          <div className={styles.brand}>
            <Link to="/" className={styles.logo}>
              <Logo className={styles.logoImage} />
            </Link>
            <p className={styles.description}>
              Платформа для совместной разработки ИТ-продуктов и создания сильных команд.
            </p>
          </div>

          <div className={styles.linksBlock}>
            <div className={styles.column}>
              <h4 className={styles.columnTitle}>Продукт</h4>
              <ul className={styles.list}>
                <li>
                  <Link to="/projects" className={styles.link}>
                    Проекты
                  </Link>
                </li>
                <li>
                  <Link to="/ideas" className={styles.link}>
                    Идеи
                  </Link>
                </li>
                <li>
                  <Link to="/community" className={styles.link}>
                    Люди
                  </Link>
                </li>
              </ul>
            </div>
            <div className={styles.column}>
              <h4 className={styles.columnTitle}>Ресурсы</h4>
              <ul className={styles.list}>
                <li>
                  <Link to="/" className={styles.link}>
                    Документация
                  </Link>
                </li>
                <li>
                  <Link to="/" className={styles.link}>
                    Блог
                  </Link>
                </li>
                <li>
                  <Link to="/" className={styles.link}>
                    Правила
                  </Link>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <div className={styles.bottom}>
          <p className={styles.copyright}>© 2026 DevSpace. Все права защищены.</p>
          <div className={styles.legal}>
            <Link to="/" className={styles.link}>
              Политика конфиденциальности
            </Link>
            <Link to="/" className={styles.link}>
              Условия использования
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}

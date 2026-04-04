import { Award, Briefcase, Globe } from "lucide-react";
import type { JSX } from "react";

import styles from "./Benefits.module.scss";

export function Benefits(): JSX.Element {
  return (
    <section className={styles.benefits}>
      <div className={styles.container}>
        <div className={styles.header}>
          <h2 className={styles.title}>Почему DevSpace?</h2>
          <p className={styles.subtitle}>
            Мы решаем главную проблему начинающих специалистов - отсутствие реального опыта работы в
            команде над настоящими продуктами.
          </p>
        </div>

        <div className={styles.grid}>
          <div className={styles.card}>
            <div className={styles.iconWrapper}>
              <Award size={24} />
            </div>
            <h3 className={styles.cardTitle}>Реальный опыт</h3>
            <p className={styles.cardText}>
              Получите практический опыт разработки в команде, который высоко ценят работодатели при
              найме.
            </p>
          </div>

          <div className={styles.card}>
            <div className={styles.iconWrapper}>
              <Briefcase size={24} />
            </div>
            <h3 className={styles.cardTitle}>Сильное портфолио</h3>
            <p className={styles.cardText}>
              Превратите одиночные, скучные пет-проекты в полноценные совместные кейсы для вашего
              резюме.
            </p>
          </div>

          <div className={styles.card}>
            <div className={styles.iconWrapper}>
              <Globe size={24} />
            </div>
            <h3 className={styles.cardTitle}>Полезный нетворкинг</h3>
            <p className={styles.cardText}>
              Находите крутых разработчиков, дизайнеров и продактов. Создавайте связи для будущих
              стартапов или работы.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}

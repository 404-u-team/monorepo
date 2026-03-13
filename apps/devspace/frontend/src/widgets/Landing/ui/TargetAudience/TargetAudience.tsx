import styles from './TargetAudience.module.scss';
import type { JSX } from 'react';
import { clsx } from 'clsx';

export function TargetAudience(): JSX.Element {
  return (
    <section className={styles.targetAudience}>
      <div className={styles.container}>
        <div className={styles.header}>
          <h2 className={styles.title}>Платформа для каждого</h2>
          <p className={styles.subtitle}>
            DevSpace объединяет людей с разными целями и опытом для создания крутых продуктов.
          </p>
        </div>

        <div className={styles.grid}>
          {/* Card 1: Beginners */}
          <div className={styles.card}>
            <div className={clsx(styles.cardHeader, styles.blue)}>
              <span className={styles.tag}>Для джунов и стажеров</span>
              <h3 className={styles.cardTitle}>Начинающим специалистам</h3>
            </div>
            <div className={styles.cardBody}>
              <p className={styles.cardText}>
                Получите первую практику работы в распределенной команде, переймите опыт и соберите крутое портфолио для уверенного старта карьеры. Забудьте про скучные туду-листы или магазины.
              </p>
            </div>
          </div>

          {/* Card 2: Pro */}
          <div className={styles.card}>
            <div className={clsx(styles.cardHeader, styles.orange)}>
              <span className={styles.tag}>Для мидлов и сеньоров</span>
              <h3 className={styles.cardTitle}>Опытным профессионалам</h3>
            </div>
            <div className={styles.cardBody}>
              <p className={styles.cardText}>
                Реализуйте свои амбициозные идеи, соберите команду мечты или просто попробуйте новые для себя технологии на реальных, интересных задачах вне работы.
              </p>
            </div>
          </div>

          {/* Card 3: Enthusiasts */}
          <div className={styles.card}>
            <div className={clsx(styles.cardHeader, styles.green)}>
              <span className={styles.tag}>Для энтузиастов</span>
              <h3 className={styles.cardTitle}>Энтузиастам и фаундерам</h3>
            </div>
            <div className={styles.cardBody}>
              <p className={styles.cardText}>
                Найдите тех, кто так же горит вашей идеей, и воплотите ее в жизнь. Отличный старт для будущих фаундеров, продактов.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

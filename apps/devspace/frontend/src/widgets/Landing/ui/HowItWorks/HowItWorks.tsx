import type { JSX } from "react";

import styles from "./HowItWorks.module.scss";

export function HowItWorks(): JSX.Element {
  return (
    <section className={styles.howItWorks}>
      <div className={styles.container}>
        <div className={styles.header}>
          <h2 className={styles.title}>Как это работает</h2>
          <p className={styles.subtitle}>Весь путь от поиска команды до успешного релиза</p>
        </div>

        <div className={styles.steps}>
          <div className={styles.step}>
            <div className={styles.stepHeader}>
              <div className={styles.stepNumber}>01</div>
              <div className={styles.stepLine} />
            </div>
            <div className={styles.stepContent}>
              <h3 className={styles.stepTitle}>Расскажите о себе</h3>
              <p className={styles.stepText}>
                Создайте профиль, укажите свой стек технологий, навыки и проекты, над которыми вы
                хотели бы поработать.
              </p>
            </div>
          </div>

          <div className={styles.step}>
            <div className={styles.stepHeader}>
              <div className={styles.stepNumber}>02</div>
              <div className={styles.stepLine} />
            </div>
            <div className={styles.stepContent}>
              <h3 className={styles.stepTitle}>Найдите или предложите</h3>
              <p className={styles.stepText}>
                Опубликуйте свою идею проекта или изучите уже существующие, откликнувшись на
                открытые для вас слоты.
              </p>
            </div>
          </div>

          <div className={styles.step}>
            <div className={styles.stepHeader}>
              <div className={styles.stepNumber}>03</div>
            </div>
            <div className={styles.stepContent}>
              <h3 className={styles.stepTitle}>Создавайте вместе</h3>
              <p className={styles.stepText}>
                Объединяйтесь с другими участниками, распределяйте задачи, общайтесь и доводите
                продукт до финального релиза.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

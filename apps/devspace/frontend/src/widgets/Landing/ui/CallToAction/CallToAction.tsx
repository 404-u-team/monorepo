import type { JSX } from "react";

import { Button } from "@/shared/ui";

import styles from "./CallToAction.module.scss";

export function CallToAction(): JSX.Element {
  return (
    <section className={styles.cta}>
      <div className={styles.container}>
        <div className={styles.content}>
          <h2 className={styles.title}>Готовы создать что-то выдающееся?</h2>
          <p className={styles.subtitle}>
            Присоединяйтесь к тысячам специалистов. Найдите команду мечты, начните свой проект или
            внесите вклад в уже существующий.
          </p>
          <div className={styles.actions}>
            <Button variant="primary" className={styles.button}>
              Присоединиться сейчас
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
}

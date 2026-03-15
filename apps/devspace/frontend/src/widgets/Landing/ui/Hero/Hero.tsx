import styles from './Hero.module.scss';
import type { JSX } from 'react';
import { Button } from '@/shared/ui';
import { heroTechnologies } from '../../config/technologies';

export function Hero(): JSX.Element {
  return (
    <section className={styles.hero}>
      <div className={styles.visual}>
        {heroTechnologies.map((tech) => (
          <div 
            key={tech.name}
            className={styles.tag} 
            style={{
              top: tech.top,
              left: tech.left,
              color: tech.color,
              animationDelay: tech.delay,
              animationDuration: tech.duration,
              '--tx1': tech.tx1,
              '--ty1': tech.ty1,
              '--tx2': tech.tx2,
              '--ty2': tech.ty2,
            } as React.CSSProperties}
            title={tech.name}
          >
            <svg role="img" viewBox="0 0 24 24" fill="currentColor" width="28" height="28">
              <path d={tech.icon.path} />
            </svg>
          </div>
        ))}
      </div>
      
      <div className={styles.container}>
        <div className={styles.content}>
          <h1 className={styles.title}>
            DevSpace:<br />
            Ваше пространство<br />
            для совместной разработки
          </h1>
          <p className={styles.subtitle}>
            Объединяем специалистов разного уровня для создания реальных ИТ-продуктов. Нарабатывайте опыт в команде, расширяйте портфолио и находите единомышленников от первой идеи до готового релиза.
          </p>
          <div className={styles.actions}>
            <Button variant="primary">Создать профиль</Button>
            <Button variant="outline">Изучить проекты и идеи</Button>
          </div>
        </div>
      </div>
    </section>
  );
}

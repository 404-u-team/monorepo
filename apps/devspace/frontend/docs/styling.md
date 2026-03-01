# Стилизация — Frontend

---

## Обязательное требование: только CSS-переменные

**Запрещено** использовать хардкоженные цвета, отступы или любые другие значения дизайн-системы напрямую в коде.
Любое значение из дизайн-системы должно быть вынесено в CSS-переменную и использоваться через неё.

```scss
// Правильно
color: var(--text--primary);
background-color: var(--bg--main);
border-color: var(--border--main);

// Неправильно
color: #08101A;
background-color: #FFFFFF;
border: 1px solid #D0D5DD;
```

Это обеспечивает единый источник правды для всех значений

---

## Расположение стилей

```
src/app/styles/
├── index.scss          # Основные стили
└── colors.scss         # CSS-переменные цветовой палитры
```

Глобальные стили подключаются **один раз** в `src/app/main.tsx`:\
Добавляя новый scss-файл в `src/app/styles/`, подключать его через `@use` в `index.scss`.

---

## CSS Modules для компонентов

**Все стили компонентов** пишутся в файлах `*.module.scss` — рядом с tsx-файлом компонента.
Глобальные стили (`src/app/styles/`) — только для CSS-переменных и reset.

```
Button/
├── Button.tsx
└── Button.module.scss   ← стили только этого компонента
```

Использование:

```tsx
import styles from './Button.module.scss'

export function Button() {
    return <button className={styles.root}>Click</button>
}
```

```scss
// Button.module.scss
.root {
    color: var(--text--onPrimary);
    background-color: var(--brand);

    &:hover {
        background-color: var(--brand--hover);
    }
}
```

**Запрещено:**
- Писать стили компонентов в глобальных файлах
- Использовать `className="my-button"` без CSS Modules (кроме случаев с внешними библиотеками)
- Хардкодить цвета — только через `var(--...)`

---

## CSS-переменные: именование

### Правила

- Семантические части имени разделяются `--` (двойное тире)
- Предлоги (`on`, `for` и т.д.) пишутся слитно с camelCase

```scss
// Правильно
--text--primary
--text--onSurface
--brand--hover
--status--error

// Неправильно
--text-primary
--text-on-surface
--brand-hover
--status-error
```

---

## CSS-переменные: доступные токены

### Brand

| Переменная | Значение | Использование |
|-----------|----------|--------------|
| `--brand` | `#005ED9` | Основной акцентный цвет |
| `--brand--hover` | `#1A75E0` | Hover-состояние бренд-элементов |
| `--brand--active` | `#0047A3` | Active/pressed состояние |
| `--brand--subtle` | `#E5F0FA` | Мягкий бренд-фон (бейджи, тэги) |

### Text

| Переменная | Значение | Использование |
|-----------|----------|--------------|
| `--text--primary` | `#08101A` | Основной текст |
| `--text--secondary` | `#434445` | Вспомогательный текст |
| `--text--disabled` | `#98A2B3` | Неактивный текст |
| `--text--onError` | `#FFFFFF` | Текст поверх красного фона |
| `--text--onSuccess` | `#FFFFFF` | Текст поверх зелёного фона |
| `--text--onPrimary` | `#FFFFFF` | Текст поверх бренд-цвета |
| `--text--onWarning` | `#08101A` | Текст поверх жёлтого фона |
| `--text--onSurface` | `#08101A` | Текст поверх нейтрального фона |

### Background

| Переменная | Значение | Использование |
|-----------|----------|--------------|
| `--bg--main` | `#FFFFFF` | Основной фон страницы |
| `--bg--surface` | `#EDEDED` | Фон карточек, секций |
| `--bg--disabled` | `#F2F4F7` | Фон неактивных элементов |

### Status

| Переменная | Значение | Использование |
|-----------|----------|--------------|
| `--status--error` | `#D92D20` | Ошибка |
| `--status--warning` | `#F79009` | Предупреждение |
| `--status--success` | `#039855` | Успех |

### Border

| Переменная | Значение | Использование |
|-----------|----------|--------------|
| `--border--main` | `#D0D5DD` | Основная граница |
| `--border--hover` | `#98A2B3` | Граница при hover |
| `--border--divider` | `#EAECF0` | Разделители |
| `--border--active` | `#005CFF` | Граница активного элемента |

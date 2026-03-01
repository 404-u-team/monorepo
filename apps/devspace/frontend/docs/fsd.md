# Feature-Sliced Design — Frontend Architecture

Архитектура проекта основана на методологии [Feature-Sliced Design (FSD)](https://feature-sliced.design/).

![Feature-Sliced Design Schema](https://feature-sliced.design/assets/images/visual_schema-e826067f573946613dcdc76e3f585082.jpg)

---

## Структура директорий

```
src/
├── app/                        # Инициализация приложения
│   ├── main.tsx                # Точка входа (React root, router, StrictMode)
│   ├── generated/              # Автогенерация (не редактировать вручную)
│   ├── providers/              # React Context, MobX RootStore Provider и т.д.
│   └── styles/                 # Глобальные стили
│       ├── index.scss          # Точка входа стилей (импортирует всё остальное)
│       └── colors.scss         # CSS-переменные цветовой палитры
│
├── routes/                     # Роутинг (TanStack Router)
│   ├── __root.tsx              # Корневой лейаут приложения
│   ├── index.tsx               # Главная страница
│   └── profile.tsx             # Пример страницы
│
├── widgets/                    # Крупные самостоятельные блоки страниц
│   ├── Header/                 # Шапка сайта
|    <...>
│
├── features/                   # Конкретные действия пользователя
│   ├── auth/                   # Авторизация
|    <...>
│
├── entities/                   # Бизнес-сущности
│   ├── user/                   # Модель пользователя, MobX UserStore, типы
|    <...>
│
└── shared/                     # Переиспользуемый код без бизнес-логики
    ├── ui/                     # UI-кит: кнопки, инпуты, модалки
    ├── api/                    # Настроенный fetch / axios
    ├── lib/                    # Утилиты, хелперы, хуки
    └── config/                 # Константы, env-переменные
```

---

## Иерархия слоёв

Слои упорядочены от **нижнего** (наименее зависимый) к **верхнему** (наиболее зависимый):

```
app  ▲
routes
widgets
features
entities
shared  ▼
```

**Правило:** модуль может импортировать **только из слоёв ниже себя**.

| Слой | Может импортировать из |
|------|----------------------|
| `app` | всё |
| `routes` | `widgets`, `features`, `entities`, `shared` |
| `widgets` | `features`, `entities`, `shared` |
| `features` | `entities`, `shared` |
| `entities` | `shared` |
| `shared` | ничего |

---

## Требования к импортам

### 1. Только абсолютные импорты через `@/`

Алиас `@` указывает на `src/`. Все межслойные импорты — **только через `@/`**.

```ts
// Правильно
import { Button } from '@/shared/ui/Button'
import { UserStore } from '@/entities/user'

// Неправильно
import { Button } from '../../shared/ui/Button'
import { UserStore } from '../user'
```

### 2. Запрет импорта из вышестоящих слоёв

Нарушение иерархии слоёв вызывает ошибку ESLint (`import-x/no-restricted-paths`):

```ts
// Ошибка: [FSD] "shared" cannot import from layers above it.
// в файле src/shared/lib/utils.ts:
import { useAuth } from '@/features/auth'

// Ошибка: [FSD] "entities" cannot import from layers above it.
// в файле src/entities/user/model.ts:
import { Header } from '@/widgets/Header'
```

### 3. Запрет относительных импортов за пределы своего модуля

Правило `import-x/no-relative-parent-imports` запрещает `../` за границы слайса:

```ts
// Неправильно — выход из своего слайса через ../
import { api } from '../../api/client'

// Правильно
import { api } from '@/shared/api/client'
```

> Соглашения по стилизации, CSS-переменные и цветовые токены — см. [docs/styling.md](./styling.md).

---

## Структура слайса и публичное API

Каждый слайс (папка внутри слоя) имеет **единственную точку входа** — `index.ts`.
Импортировать напрямую из внутренних папок слайса **запрещено**.

```
entities/user/
├── index.ts              # Публичное API — единственный способ импортировать снаружи
├── model/
│   ├── UserStore.ts
│   └── user.types.ts
└── ui/
    ├── UserCard.tsx
    └── UserCard.module.scss
```

```ts
// Правильно
import { UserStore } from '@/entities/user'

// Неправильно — нарушение публичного API
import { UserStore } from '@/entities/user/model/UserStore'
```

Внутри самого слайса (между его файлами) — относительные импорты допустимы:

```ts
// внутри entities/user/ui/UserCard.tsx — ок
import type { User } from '../model/user.types'
```

---

## Автогенерируемые файлы

| Файл | Генератор | Редактировать |
|------|-----------|--------------|
| `src/app/generated/routeTree.gen.ts` | TanStack Router (Vite plugin) | Никогда |

Файл генерируется автоматически при запуске `bun run dev` или `bun run build`.
Он находится в `.gitignore` и исключён из ESLint.

---

## Инструменты контроля архитектуры

| Правило | Плагин | Что проверяет |
|---------|--------|---------------|
| `import-x/no-restricted-paths` | `eslint-plugin-import-x` | Запрет импортов из вышестоящих FSD-слоёв |
| `import-x/no-relative-parent-imports` | `eslint-plugin-import-x` | Запрет `../` за пределы слайса |
| `import-x/no-cycle` | `eslint-plugin-import-x` | Запрет циклических зависимостей |

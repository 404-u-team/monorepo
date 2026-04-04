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

| Слой       | Может импортировать из                      |
| ---------- | ------------------------------------------- |
| `app`      | всё                                         |
| `routes`   | `widgets`, `features`, `entities`, `shared` |
| `widgets`  | `features`, `entities`, `shared`            |
| `features` | `entities`, `shared`                        |
| `entities` | `shared`                                    |
| `shared`   | ничего                                      |

---

## Требования к импортам

### 1. Только абсолютные импорты через `@/`

Алиас `@` указывает на `src/`. Все межслойные импорты — **только через `@/`**.

```ts
// Правильно
import { Button } from "@/shared/ui/Button";
import { UserStore } from "@/entities/user";

// Неправильно
import { Button } from "../../shared/ui/Button";
import { UserStore } from "../user";
```

### 2. Запрет импорта из вышестоящих слоёв

Нарушение иерархии слоёв вызывает ошибку ESLint (`import-x/no-restricted-paths`):

```ts
// Ошибка: [FSD] "shared" cannot import from layers above it.
// в файле src/shared/lib/utils.ts:
import { useAuth } from "@/features/auth";

// Ошибка: [FSD] "entities" cannot import from layers above it.
// в файле src/entities/user/model.ts:
import { Header } from "@/widgets/Header";
```

### 3. Запрет относительных импортов за пределы своего модуля

Правило `import-x/no-relative-parent-imports` запрещает `../` за границы слайса:

```ts
// Неправильно — выход из своего слайса через ../
import { api } from "../../api/client";

// Правильно
import { api } from "@/shared/api/client";
```

> Соглашения по стилизации, CSS-переменные и цветовые токены — см. [docs/styling.md](./styling.md).

---

## Структура слайса и публичное API

Каждый слайс (папка внутри слоя) имеет **единственную точку входа** — `index.ts`.
Импортировать напрямую из внутренних папок слайса **запрещено**.

### Эталонная структура слайса

```
entities/idea/                          # Слайс
├── index.ts                            # Публичное API — единственный способ импортировать снаружи
│
├── model/                              # Типы, интерфейсы, сторы, бизнес-логика
│   ├── IIdea.ts                        # Интерфейс сущности
│   └── IdeaStore.ts                    # MobX observable store (при необходимости)
│
├── api/                                # Запросы к бэкенду, относящиеся к слайсу
│   └── ideaApi.ts                      # fetchIdeaById, toggleFavorite и т.д.
│
├── lib/                                # Хуки, утилиты, хелперы слайса
│   ├── useIdea.ts                      # Кастомный React-хук
│   └── formatIdea.ts                   # Функции-утилиты (форматирование, валидация)
│
├── config/                             # Константы и конфигурация слайса
│   └── ideaConfig.ts                   # Магические числа, маппинги, enum-значения
│
└── ui/                                 # React-компоненты сущности
    ├── IdeaCard/                       # Один компонент — одна папка
    │   ├── IdeaCard.tsx
    │   ├── IdeaCard.module.scss        # CSS Modules стили
    │   └── IdeaCard.stories.tsx        # Storybook stories
    └── IdeaCardSkeleton/
        ├── IdeaCardSkeleton.tsx
        └── IdeaCardSkeleton.module.scss
```

### Сегменты слайса

| Сегмент   | Назначение                                                                 | Обязательный     |
| --------- | -------------------------------------------------------------------------- | ---------------- |
| `model/`  | Типы, интерфейсы (`I*.ts`), MobX-сторы, бизнес-логика                      | Да               |
| `api/`    | Функции запросов к бэкенду (используют `apiClient` из shared)              | По необходимости |
| `ui/`     | React-компоненты с `.module.scss` и `.stories.tsx`                         | По необходимости |
| `lib/`    | Кастомные React-хуки (`use*.ts`), утилиты, хелперы, специфичные для слайса | По необходимости |
| `config/` | Константы, enum-маппинги, конфигурация (не env)                            | По необходимости |

> **Когда добавлять `lib/` и `config/`?** Если хук или утилита нужны **только внутри одного слайса** — кладите в `lib/`. Если они могут быть полезны глобально — выносите в `shared/lib/`. Аналогично: `config/` для констант привязанных к слайсу, `shared/config/` для глобальных.

### Публичное API (`index.ts`)

```ts
// entities/idea/index.ts — реэкспортируем только то, что нужно снаружи
export { IdeaCard } from "./ui/IdeaCard/IdeaCard";
export type { IIdea } from "./model/IIdea";
export { fetchIdeaById } from "./api/ideaApi";
```

```ts
// ✅ Правильно
import { IdeaCard } from "@/entities/idea";

// ❌ Неправильно — нарушение публичного API
import { IdeaCard } from "@/entities/idea/ui/IdeaCard/IdeaCard";
```

### Внутренние импорты

Внутри самого слайса (между его файлами) — относительные импорты допустимы:

```ts
// внутри entities/idea/ui/IdeaCard/IdeaCard.tsx — ок
import type { IIdea } from "../../model/IIdea";
import { fetchIdeaById } from "../../api/ideaApi";
```

---

## Автогенерируемые файлы

| Файл                                 | Генератор                     | Редактировать |
| ------------------------------------ | ----------------------------- | ------------- |
| `src/app/generated/routeTree.gen.ts` | TanStack Router (Vite plugin) | Никогда       |

Файл генерируется автоматически при запуске `bun run dev` или `bun run build`.
Он находится в `.gitignore` и исключён из ESLint.

---

## Инструменты контроля архитектуры

| Правило                               | Плагин                   | Что проверяет                            |
| ------------------------------------- | ------------------------ | ---------------------------------------- |
| `import-x/no-restricted-paths`        | `eslint-plugin-import-x` | Запрет импортов из вышестоящих FSD-слоёв |
| `import-x/no-relative-parent-imports` | `eslint-plugin-import-x` | Запрет `../` за пределы слайса           |
| `import-x/no-cycle`                   | `eslint-plugin-import-x` | Запрет циклических зависимостей          |

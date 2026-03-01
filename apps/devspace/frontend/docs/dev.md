# Разработка — Frontend

Общий гайд по разработке фронтенда. Здесь — всё что нужно знать перед тем как начать писать код.

---

## Требования к окружению

| Инструмент | Версия | Установка |
|-----------|--------|-----------|
| [Bun](https://bun.sh) | ≥ 1.x | bun.com |
| [Node.js](https://nodejs.org) | ≥ 20 LTS | nodejs.org |
| [Git](https://git-scm.com) | любая | git-scm.com |

> Bun используется как пакетный менеджер и раннер скриптов вместо npm/yarn.

---

## Быстрый старт

```bash
bun install       # установить зависимости
bun run dev       # запустить dev-сервер (http://localhost:5173)
bun run lint      # проверить линтер
bun run storybook # запустить Storybook (http://localhost:6006)
```

---

## Технологический стек

| Область | Инструмент | Версия |
|---------|-----------|--------|
| Язык | TypeScript | ~5.9 |
| UI | React | ^19 |
| Роутинг | TanStack Router | ^1.163 |
| Состояние | MobX | — |
| Стили | SCSS + CSS Modules | — |
| Сборка | Vite | ^8 |
| Пакетный менеджер | Bun | — |
| Компонентная документация | Storybook | ^10 |
| Линтер | ESLint | ^10 |

---

## Архитектура и стилизация

| Документ | Описание |
|----------|---------|
| [docs/fsd.md](./fsd.md) | Структура проекта, FSD-слои, правила импортов |
| [docs/styling.md](./styling.md) | CSS-переменные, CSS Modules, правила стилизации |

---

## Соглашения по коду

### Именование файлов и компонентов

| Сущность | Пример |
|---------|--------|
| React-компонент | `Button.tsx`, `UserCard.tsx` |
| CSS Module | `Button.module.scss` |
| MobX Store | `UserStore.ts` |
| Хук | `useUser.ts` |
| Утилита | `formatDate.ts` |
| Тип/интерфейс | `user.types.ts` или `types.ts` внутри слайса |
| Storybook Story | `Button.stories.tsx` |

---

## MobX

MobX настроен в `src/app/main.tsx` в строгом режиме

- Всё изменение стейта — через `action` или `runInAction`
- Нельзя читать `computed` вне `observer`/`reaction`

---

## Роутинг (TanStack Router)

- Роуты находятся в `src/routes/` — файловый роутинг
- Дерево роутов автогенерируется в `src/app/generated/routeTree.gen.ts` — **не редактировать вручную**
- Файл генерируется при `bun run dev` или `bun run build`

---

## Линтер

Запуск:
```bash
bun run lint
```

Подробнее — в `eslint.config.js`.

---

## Storybook

Используется для документирования UI-компонентов из `shared/ui/` и визуального тестирования.

```bash
bun run storybook          # dev-режим на :6006
```

### CSF Next (Component Story Format)

Проект использует **CSF Next** — современный формат с максимальной типобезопасностью через `satisfies`.

```tsx
// Button.stories.tsx
import type { Meta, StoryObj } from '@storybook/react'
import { Button } from './Button'

const meta = {
    component: Button,
    tags: ['autodocs'],          // автодокументация из prop-types
    args: {
        label: 'Click me',
    },
} satisfies Meta<typeof Button>

export default meta
type Story = StoryObj<typeof meta>

export const Primary: Story = {
    args: {
        variant: 'primary',
    },
}

export const Disabled: Story = {
    args: {
        disabled: true,
    },
}
```

**Ключевые практики:**
- `satisfies Meta<typeof Component>` вместо аннотации типом — даёт вывод типов в `args`
- `tags: ['autodocs']` — автогенерация страницы документации с prop-таблицей
- Базовые `args` в `meta`, перегрузки — в конкретных stories
- Файл stories кладётся рядом с компонентом: `Button.stories.tsx`

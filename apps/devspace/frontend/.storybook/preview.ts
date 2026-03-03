import type { Preview } from '@storybook/react-vite'

const preview = {
  // Все stories получают страницу Autodocs автоматически.
  // Можно переопределить на уровне конкретного файла: tags: []
  tags: ['autodocs'],

  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },

    a11y: {
      // 'todo'  — показывать нарушения только в Test UI
      // 'error' — падать в CI при нарушениях
      // 'off'   — пропускать проверку
      test: 'todo',
    },
  },
} satisfies Preview

export default preview
import { definePreview } from '@storybook/react-vite';
// @ts-expect-error - SCSS типы определены в tsconfig.app.json (vite/client), а Storybook использует tsconfig.node.json
// eslint-disable-next-line import-x/no-relative-parent-imports
import '@/app/styles/index.scss'; // Подключение глобальных стилей для всех историй

const preview = definePreview({
  addons: [],
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
})

export default preview
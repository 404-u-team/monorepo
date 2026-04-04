import { definePreview } from "@storybook/react-vite";
import { initialize, mswLoader } from "msw-storybook-addon";
import { createElement } from "react";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import "@/app/styles/index.scss";
import { StoreContext, type IRootStore } from "@/shared/lib/store";

// Initialize MSW
initialize();

/**
 * Мок-стор по умолчанию для Storybook.
 * Можно переопределить через декоратор конкретного story.
 */
const defaultMockStore: IRootStore = {
  userStore: { isAuthenticated: false },
};

const preview = definePreview({
  addons: [],
  tags: ["autodocs"],

  loaders: [mswLoader],

  decorators: [
    // eslint-disable-next-line @typescript-eslint/explicit-function-return-type
    (Story) =>
      createElement(StoreContext.Provider, { value: defaultMockStore }, createElement(Story)),
  ],

  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },

    a11y: {
      test: "todo",
    },
  },
});

export default preview;

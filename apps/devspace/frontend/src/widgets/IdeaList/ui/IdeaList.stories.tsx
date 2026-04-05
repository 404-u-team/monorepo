import type { Meta, StoryObj } from "@storybook/react";
import {
  RouterProvider,
  createMemoryHistory,
  createRootRoute,
  createRouter,
} from "@tanstack/react-router";
import { http, HttpResponse } from "msw";
import { createElement } from "react";

import { StoreContext } from "@/shared/lib/store";

import { IdeaList } from "./IdeaList";

const memoryHistory = createMemoryHistory({
  initialEntries: ["/ideas"],
});

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const withRouter = (Story: any) => {
  const rootRoute = createRootRoute({
    component: Story,
  });

  const router = createRouter({
    routeTree: rootRoute,
    history: memoryHistory,
  });

  return <RouterProvider router={router} />;
};

const meta = {
  title: "widgets/IdeaList",
  component: IdeaList,
  parameters: {
    layout: "padded",
    msw: {
      handlers: [
        http.get("*/ideas/:ideaId", ({ params }) => {
          return HttpResponse.json({
            id: params.ideaId,
            title: `Идея ${params.ideaId}`,
            description: "Описание идеи из мока",
            category: "MVP",
            author_id: "user-1",
            created_at: "2026-03-01T12:00:00Z",
            updated_at: "2026-03-10T15:30:00Z",
            views_count: 100,
            favorites_count: 5,
          });
        }),
        http.get("*/users/:userId", () => {
          return HttpResponse.json({
            id: "user-1",
            nickname: "AuthorName",
            avatar_url: "",
          });
        }),
      ],
    },
  },
  decorators: [
    withRouter,
    (Story) =>
      createElement(
        StoreContext.Provider,
        { value: { userStore: { isAuthenticated: true } } as any },
        createElement(Story),
      ),
  ],
  tags: ["autodocs"],
} satisfies Meta<typeof IdeaList>;

export default meta;
type Story = StoryObj<typeof meta>;

// Моковые идеи, чтобы заполнить список
const MOCK_IDEAS = Array.from({ length: 6 }).map(
  (_, i) =>
    ({
      id: `idea-${i}`,
    }) as any,
);

export const Default: Story = {
  args: {
    ideas: MOCK_IDEAS,
    totalPages: 5,
  },
};

export const Empty: Story = {
  args: {
    ideas: [],
    totalPages: 0,
  },
};

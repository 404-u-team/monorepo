import type { Meta, StoryObj } from "@storybook/react";
import {
  RouterProvider,
  createMemoryHistory,
  createRootRoute,
  createRouter,
} from "@tanstack/react-router";
import { http, HttpResponse } from "msw";

import type { IProject } from "@/entities/project/model/IProject";

import { ProjectList } from "./ProjectList";

const memoryHistory = createMemoryHistory({
  initialEntries: ["/"],
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

const MOCK_PROJECT = {
  id: "1",
  title: "Project 1",
  description: "First project description here",
  leader_id: "leader1",
  status: "open",
  idea_id: null,
  created_at: "2025-01-01",
  updated_at: "2025-01-01",
  slots: [
    {
      id: "slot-1",
      project_id: "1",
      skill_category_id: "cat-1",
      skill: { id: "skill-1", name: "React", color: "3B82F6", icon: "react" },
      title: "Frontend Developer",
      status: "open",
      user_id: null,
      created_at: "2025-01-01",
    },
  ],
};

const MOCK_PROJECT_2 = {
  ...MOCK_PROJECT,
  id: "2",
  title: "Project 2",
  description: "Second project description here",
  status: "closed",
  slots: [],
};

const meta = {
  title: "widgets/ProjectList",
  component: ProjectList,
  parameters: {
    layout: "padded",
    msw: {
      handlers: [
        http.get("*/projects/:projectId", ({ params }) => {
          const id = params.projectId as string;
          if (id === "1") return HttpResponse.json(MOCK_PROJECT);
          if (id === "2") return HttpResponse.json(MOCK_PROJECT_2);
          return new HttpResponse(null, { status: 404 });
        }),
        http.get("*/users/:userId", () => {
          // empty users response
          return new HttpResponse(null, { status: 404 });
        }),
      ],
    },
  },
  decorators: [withRouter],
  tags: ["autodocs"],
} satisfies Meta<typeof ProjectList>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    projects: [{ id: "1" } as IProject, { id: "2" } as IProject],
    totalPages: 5,
    total: 10,
  },
};

export const Empty: Story = {
  args: {
    projects: [],
    totalPages: 0,
    total: 0,
  },
};

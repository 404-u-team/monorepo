import { createElement } from "react";
import { http, HttpResponse } from "msw";
import type { Meta, StoryObj } from "@storybook/react";
import { UserCard } from "./UserCard";

const MOCK_USER = {
  id: "user-1",
  avatar_uri: "https://i.pravatar.cc/150?u=user-1",
  main_role: "Senior Frontend Developer",
  nickname: "alice_wonder",
  bio: "Frontend developer with 5+ years of experience. Love React, TypeScript and clean code. Passionate about UI/UX and performance optimization.",
  skills: [
    "React",
    "TypeScript",
    "Next.js",
    "Redux",
    "Jest",
    "Webpack",
    "CSS-in-JS",
    "Storybook",
    "Figma",
  ],
};

const MOCK_USERS: Record<string, typeof MOCK_USER> = {
  "user-1": MOCK_USER,
  "user-2": {
    id: "user-2",
    avatar_uri: "https://i.pravatar.cc/150?u=user-2",
    main_role: "Backend Developer",
    nickname: "bob_builder",
    bio: "Backend engineer specializing in Python and microservices. Experienced with FastAPI, Django, and PostgreSQL.",
    skills: [
      "Python",
      "FastAPI",
      "Django",
      "PostgreSQL",
      "Redis",
      "Docker",
      "Kafka",
      "REST API",
    ],
  },
  "user-3": {
    id: "user-3",
    avatar_uri: "https://i.pravatar.cc/150?u=user-3",
    main_role: "DevOps Engineer",
    nickname: "charlie_ops",
    bio: "Infrastructure and automation enthusiast. Kubernetes certified. Love building scalable and reliable systems.",
    skills: [
      "Kubernetes",
      "Docker",
      "AWS",
      "Terraform",
      "CI/CD",
      "Prometheus",
      "Grafana",
      "Linux",
    ],
  },
  "user-4": {
    id: "user-4",
    avatar_uri: "https://i.pravatar.cc/150?u=user-4",
    main_role: "Fullstack Developer",
    nickname: "diana_fullstack",
    bio: "Fullstack developer who loves building complete products. Experience with React, Node.js, and MongoDB.",
    skills: [
      "React",
      "Node.js",
      "MongoDB",
      "Express",
      "TypeScript",
      "GraphQL",
      "TailwindCSS",
      "Prisma",
    ],
  },
  "user-5": {
    id: "user-5",
    avatar_uri: "https://i.pravatar.cc/150?u=user-5",
    main_role: "Mobile Developer",
    nickname: "eve_mobile",
    bio: "Mobile app developer specializing in React Native and Flutter. Love creating smooth and beautiful mobile experiences.",
    skills: [
      "React Native",
      "Flutter",
      "TypeScript",
      "Redux",
      "Firebase",
      "iOS",
      "Android",
      "Expo",
    ],
  },
  "user-6": {
    id: "user-6",
    avatar_uri: "https://i.pravatar.cc/150?u=user-6",
    main_role: "Data Scientist",
    nickname: "frank_ml",
    bio: "Machine learning engineer with focus on NLP and computer vision. Experienced with PyTorch and TensorFlow.",
    skills: [
      "Python",
      "PyTorch",
      "TensorFlow",
      "Pandas",
      "NumPy",
      "Scikit-learn",
      "Jupyter",
      "SQL",
    ],
  },
};

const meta = {
  title: "entities/UserCard",
  component: UserCard,
  parameters: {
    layout: "centered",
    msw: {
      handlers: [
        http.get("*/users/:id", ({ params }) => {
          const user = MOCK_USERS[params.id as string];
          if (user !== undefined) {
            return HttpResponse.json(user);
          }
          return new HttpResponse(null, { status: 404 });
        }),
      ],
    },
  },
  tags: ["autodocs"],
  decorators: [
    (Story) =>
      createElement("div", { style: { width: "380px" } }, createElement(Story)),
  ],
} satisfies Meta<typeof UserCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    id: "user-1",
  },
};

export const WithInviteButton: Story = {
  args: {
    id: "user-1",
    project_id: "project-123",
    slot_id: "slot-456",
  },
};

export const LongBio: Story = {
  args: {
    id: "user-long-bio",
  },
  parameters: {
    msw: {
      handlers: [
        http.get("*/users/:id", () => {
          return HttpResponse.json({
            id: "user-long-bio",
            avatar_uri: "https://i.pravatar.cc/150?u=user-long",
            main_role: "Principal Software Architect",
            nickname: "dr_architect",
            bio: "Visionary software architect with over 15 years of experience designing and building large-scale distributed systems for Fortune 500 companies. Passionate about system design, microservices architecture, event-driven systems, and cloud-native technologies. Mentored hundreds of engineers and contributed to open-source projects. Speaker at international conferences including AWS re:Invent, KubeCon, and React Summit. Author of several technical books on software architecture and design patterns. Believer in clean code, SOLID principles, and domain-driven design. Always eager to learn new technologies and share knowledge with the community.",
            skills: [
              "System Design",
              "Microservices",
              "Event-Driven Architecture",
              "DDD",
              "Cloud Native",
              "Kubernetes",
              "AWS",
              "Go",
              "Rust",
              "Kafka",
              "gRPC",
              "CQRS",
              "Event Sourcing",
            ],
          });
        }),
      ],
    },
  },
};

export const ManySkills: Story = {
  args: {
    id: "user-many-skills",
  },
  parameters: {
    msw: {
      handlers: [
        http.get("*/users/:id", () => {
          return HttpResponse.json({
            id: "user-many-skills",
            avatar_uri: "https://i.pravatar.cc/150?u=user-skills",
            main_role: "Polyglot Developer",
            nickname: "polyglot",
            bio: "Developer who loves learning new programming languages and technologies.",
            skills: [
              "JavaScript",
              "TypeScript",
              "Python",
              "Java",
              "Go",
              "Rust",
              "C++",
              "C#",
              "PHP",
              "Ruby",
              "React",
              "Vue",
              "Angular",
              "Svelte",
              "Next.js",
              "Nuxt",
              "Gatsby",
              "Node.js",
              "Deno",
              "Bun",
              "Express",
              "Fastify",
              "NestJS",
              "Django",
              "Flask",
              "FastAPI",
              "Spring Boot",
              "Laravel",
              "Rails",
              "PostgreSQL",
              "MySQL",
              "MongoDB",
              "Redis",
              "Elasticsearch",
              "Cassandra",
              "Docker",
              "Kubernetes",
              "Terraform",
              "AWS",
              "Azure",
              "GCP",
              "GraphQL",
              "REST",
              "gRPC",
              "WebSocket",
              "MQTT",
              "Jest",
              "Mocha",
              "Cypress",
              "Playwright",
              "Selenium",
            ],
          });
        }),
      ],
    },
  },
};

export const Loading: Story = {
  args: {
    id: "user-loading",
  },
  parameters: {
    msw: {
      handlers: [
        http.get("*/users/:id", async () => {
          await new Promise(() => {});
          return HttpResponse.json(MOCK_USER);
        }),
      ],
    },
  },
};

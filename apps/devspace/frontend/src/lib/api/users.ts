// import { apiClient } from "../api-client";
// import type {
//   PrivateUserProfile,
//   PublicUserProfile,
//   UUID,
// } from "@/types/api.types";

// export const usersApi = {
//   getMyProfile: (token: string) =>
//     apiClient.request<PrivateUserProfile>("/users/me", {
//       method: "GET",
//       token,
//     }),

//   updateMyProfile: (
//     token: string,
//     data: {
//       nickname?: string;
//       avatar_uri?: string;
//       bio?: string;
//     },
//   ) =>
//     apiClient.request<PrivateUserProfile>("/users/me", {
//       method: "PUT",
//       token,
//       body: JSON.stringify(data),
//     }),

//   addSkill: (token: string, skill_id: UUID) =>
//     apiClient.request<void>("/users/me/skills", {
//       method: "POST",
//       token,
//       body: JSON.stringify({ skill_id }),
//     }),

//   removeSkill: (token: string, skill_id: UUID) =>
//     apiClient.request<void>(`/users/me/skills/${skill_id}`, {
//       method: "DELETE",
//       token,
//     }),

//   getUserProfile: (user_id: UUID) =>
//     apiClient.request<PublicUserProfile>(`/users/${user_id}`, {
//       method: "GET",
//     }),
// };

// import type { ApiError, AuthResponse, UUID } from "@/types/api.types";

// interface FetchOptions extends RequestInit {
//   token?: string;
// }

// class ApiClient {
//   private baseUrl: string;
//   private defaultHeaders: HeadersInit;

//   constructor(baseUrl = "https://api.example.com/v1") {
//     this.baseUrl = baseUrl;
//     this.defaultHeaders = {
//       "Content-Type": "application/json",
//     };
//   }

//   private async request<T>(
//     endpoint: string,
//     options: FetchOptions = {},
//   ): Promise<T> {
//     const url = `${this.baseUrl}${endpoint}`;
//     const headers = {
//       ...this.defaultHeaders,
//       ...(options.token && { Authorization: `Bearer ${options.token}` }),
//       ...options.headers,
//     };

//     const response = await fetch(url, {
//       ...options,
//       headers,
//     });

//     if (!response.ok) {
//       const error = await response.json().catch(() => ({}));
//       throw new ApiError(
//         error.statusCode || response.status,
//         error.message || response.statusText,
//       );
//     }

//     if (response.status === 204) {
//       return {} as T;
//     }

//     return response.json();
//   }

//   setAuthToken(token: string) {
//     this.defaultHeaders = {
//       ...this.defaultHeaders,
//       Authorization: `Bearer ${token}`,
//     };
//   }

//   clearAuthToken(): void {
//     const { Authorization, ...rest } = this.defaultHeaders;
//     this.defaultHeaders = rest;
//   }
// }

// export const apiClient = new ApiClient();

// export class ApiError extends Error {
//   constructor(
//     public statusCode: number,
//     message: string,
//   ) {
//     super(message);
//     this.name = "ApiError";
//   }
// }

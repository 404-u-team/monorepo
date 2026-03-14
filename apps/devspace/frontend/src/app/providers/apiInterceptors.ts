import axios, { type InternalAxiosRequestConfig } from 'axios';
import { apiClient } from '@/shared/api/client';
import { rootStore } from '@/app/providers/store';

interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
    _isRetry?: boolean;
}

export const verifyInterceptors = (): void => {
    // Add Access token to every outgoing request and handle /auth routes baseURL
    apiClient.interceptors.request.use((config) => {
        if (config.url?.startsWith('/auth')) {
            const currentBaseURL = config.baseURL ?? apiClient.defaults.baseURL;
            if (currentBaseURL?.endsWith('/api')) {
                config.baseURL = currentBaseURL.slice(0, -4);
            }
        }

        const token = rootStore.userStore.accessToken;
        if (typeof token === 'string' && token.length > 0) {
            config.headers.set('Authorization', `Bearer ${token}`);
        }
        return config;
    });

    // Handle 401 and refresh token
    apiClient.interceptors.response.use(
        (response) => response,
        async (error: unknown) => {
            if (axios.isAxiosError(error) && error.response?.status === 401 && error.config !== undefined) {
                const originalRequest = error.config as CustomAxiosRequestConfig;

                if (originalRequest._isRetry !== true) {
                    originalRequest._isRetry = true;

                    try {
                        const baseURL = apiClient.defaults.baseURL ?? '/api';
                        const authBaseURL = baseURL.endsWith('/api') ? baseURL.slice(0, -4) : baseURL;

                        const refreshResponse = await axios.post<{ accessToken: string }>(
                            `${authBaseURL}/auth/refresh`,
                            {},
                            { withCredentials: true }
                        );

                        const newAccessToken = refreshResponse.data.accessToken;
                        rootStore.userStore.setAccessToken(newAccessToken);

                        originalRequest.headers.set('Authorization', `Bearer ${newAccessToken}`);
                        return await apiClient(originalRequest);
                    } catch (refreshError) {
                        rootStore.userStore.invalidateToken();
                        rootStore.userStore.invalidateUser();

                        if (refreshError instanceof Error) {
                            throw refreshError;
                        }
                        throw new Error(String(refreshError), { cause: refreshError });
                    }
                }
            }

            if (error instanceof Error) {
                throw error;
            }
            throw new Error(String(error), { cause: error });
        }
    );
};

import axios from 'axios';

const baseURL = import.meta.env.VITE_API_URL as string | undefined;

if (baseURL === undefined || baseURL === '') {
    throw new Error('VITE_API_URL environment variable is not defined. Please set it in your Vite environment configuration.');
}

export const apiClient = axios.create({
    baseURL,
    withCredentials: true,
});

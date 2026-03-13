import axios from 'axios';

const baseURL = import.meta.env.VITE_API_URL;

if (!baseURL) {
    throw new Error('VITE_API_URL environment variable is not defined. Please set it in your Vite environment configuration.');
}

export const apiClient = axios.create({
    baseURL,
    withCredentials: true,
});

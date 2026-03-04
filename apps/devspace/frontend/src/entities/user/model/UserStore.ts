import { makeAutoObservable } from "mobx";
import type { IUser } from "./IUser";

export class UserStore {
    user: IUser | null = null;
    accessToken: string | null = null;

    constructor() {
        makeAutoObservable(this);
    }

    setUser(user: IUser) {
        this.user = user;
    }

    setAccessToken(token: string) {
        this.accessToken = token;
    }

    invalidateToken() {
        this.accessToken = null;
    }

    invalidateUser() {
        this.user = null;
    }

    isAuthenticated(): boolean {
        return this.user !== null && this.accessToken !== null;
    }
}
import { makeAutoObservable } from "mobx";
import type { IUser } from "./IUser";

export class UserStore {
    user: IUser | undefined = undefined;
    accessToken: string | undefined = undefined;

    constructor() {
        makeAutoObservable(this);
    }

    setUser(user: IUser): void {
        this.user = user;
    }

    setAccessToken(token: string): void {
        this.accessToken = token;
    }

    invalidateToken(): void {
        this.accessToken = undefined;
    }

    invalidateUser(): void {
        this.user = undefined;
    }

    isAuthenticated(): boolean {
        return this.user !== undefined && this.accessToken !== undefined;
    }
}
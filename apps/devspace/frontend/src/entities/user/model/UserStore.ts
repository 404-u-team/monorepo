import { makeAutoObservable, runInAction } from "mobx";

import type { IUser } from "./IUser";

export class UserStore {
  user: IUser | undefined = undefined;
  accessToken: string | undefined = localStorage.getItem("accessToken") ?? undefined;
  isAuthLoading = false;

  constructor() {
    makeAutoObservable(this);
  }

  setUser(user: IUser): void {
    this.user = user;
  }

  setAccessToken(token: string): void {
    this.accessToken = token;
    localStorage.setItem("accessToken", token);
  }

  invalidateToken(): void {
    this.accessToken = undefined;
    localStorage.removeItem("accessToken");
  }

  invalidateUser(): void {
    this.user = undefined;
  }

  async checkAuth(): Promise<void> {
    console.warn("checking auth");
    if (this.accessToken === undefined) return;

    runInAction(() => {
      this.isAuthLoading = true;
    });

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL as string}users/me`, {
        headers: {
          Authorization: `Bearer ${this.accessToken}`,
        },
      });

      if (response.ok) {
        const data = (await response.json()) as Record<string, unknown>;
        const user: IUser = {
          id: data.id as string,
          nickname: data.nickname as string,
          email: data.email as string,
          avatarUrl: (data.avatar_url ?? data.avatarUrl) as string | undefined,
          bio: data.bio as string | undefined,
        };
        runInAction(() => {
          this.user = user;
        });
      } else {
        runInAction(() => {
          this.invalidateToken();
          this.invalidateUser();
        });
      }
    } catch (error) {
      console.error("Check auth error:", error);
      runInAction(() => {
        this.invalidateToken();
        this.invalidateUser();
      });
    } finally {
      runInAction(() => {
        this.isAuthLoading = false;
      });
    }
  }

  get isAuthenticated(): boolean {
    return this.accessToken !== undefined;
  }
}

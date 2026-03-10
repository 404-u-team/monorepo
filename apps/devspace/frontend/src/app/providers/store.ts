import { UserStore } from "@/entities/user";
import type { IRootStore } from "@/shared/lib/store";

export class RootStore implements IRootStore {
    readonly userStore: UserStore;

    constructor() {
        this.userStore = new UserStore();
    }
}

export const rootStore = new RootStore();

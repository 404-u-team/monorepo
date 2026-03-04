import { createContext, useContext, type Context } from "react";
// eslint-disable-next-line import-x/no-relative-parent-imports
import { UserStore } from "@/entities/user";

export class RootStore {
    readonly userStore: UserStore;

    constructor() {
        this.userStore = new UserStore();
    }
}

const rootStore = new RootStore();
export const StoreContext: Context<RootStore> = createContext<RootStore>(rootStore);
export const useStore = (): RootStore => useContext(StoreContext);
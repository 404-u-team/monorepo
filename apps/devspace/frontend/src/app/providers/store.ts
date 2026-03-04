import { createContext, useContext } from "react";
import { UserStore } from "@/entities/user";

export class RootStore {
    readonly userStore: UserStore;

    constructor() {
        this.userStore = new UserStore();
    }
}

const rootStore = new RootStore();
export const StoreContext = createContext<RootStore>(rootStore);
export const useStore = () => useContext(StoreContext);
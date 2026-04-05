import { useContext } from "react";

import { StoreContext } from "@/shared/lib/store";

import type { UserStore } from "../model/UserStore";

interface UserStoreContainer {
  readonly userStore: UserStore;
}

export function useUserStore(): UserStore {
  const store = useContext(StoreContext) as UserStoreContainer | undefined;
  if (store === undefined) {
    throw new Error("useUserStore must be used within a StoreContext.Provider");
  }
  return store.userStore;
}

import { createContext, useContext, type Context } from "react";

/**
 * Интерфейс корневого стора — описываем минимальный контракт,
 * чтобы shared не зависел от конкретной реализации (app/entities).
 * Потребители из вышестоящих слоёв, которые знают более конкретный тип,
 * могут передать его как generic-параметр в useStore<T>().
 */
export interface IRootStore {
  readonly userStore: {
    readonly isAuthenticated: boolean;
  };
}

// Context создаётся с undefined, реальное значение подставляется провайдером в app/
export const StoreContext: Context<IRootStore | undefined> = createContext<IRootStore | undefined>(
  undefined,
);

export function useStore(): IRootStore {
  const store = useContext(StoreContext);
  if (store === undefined) {
    throw new Error("useStore must be used within a StoreContext.Provider");
  }
  return store;
}

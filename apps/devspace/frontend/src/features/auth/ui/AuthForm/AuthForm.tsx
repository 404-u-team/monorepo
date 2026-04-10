import { useNavigate, useSearch } from "@tanstack/react-router";
import { isAxiosError } from "axios";
import { Mail, Lock, Eye, EyeOff, User } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useState, type JSX, type SyntheticEvent } from "react";

import type { IUser, UserStore } from "@/entities/user";
import { apiClient } from "@/shared/api/client";
import { useStore } from "@/shared/lib/store";
import { Button, Input, Logo } from "@/shared/ui";

import styles from "./AuthForm.module.scss";

const GoogleIcon = (): JSX.Element => (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path
      d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
      fill="#4285F4"
    />
    <path
      d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
      fill="#34A853"
    />
    <path
      d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
      fill="#FBBC05"
    />
    <path
      d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
      fill="#EA4335"
    />
  </svg>
);

const GithubIcon = (): JSX.Element => (
  <svg
    width="20"
    height="20"
    viewBox="0 0 24 24"
    fill="currentColor"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" />
  </svg>
);

type AuthMode = "login" | "register";

export const AuthForm = observer((): JSX.Element => {
  const { userStore } = useStore() as { userStore: UserStore };
  const navigate = useNavigate();
  const search = useSearch({ strict: false });
  const redirectTo = (search as { redirect?: string }).redirect ?? "/";

  const [mode, setMode] = useState<AuthMode>("login");
  const [nickname, setNickname] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const toggleMode = (): void => {
    setMode(mode === "login" ? "register" : "login");
    setError("");
  };

  const handleSubmit = async (event: SyntheticEvent): Promise<void> => {
    event.preventDefault();
    setError("");
    setIsLoading(true);
    try {
      const endpoint = mode === "login" ? "/auth/login" : "/auth/register";
      const payload = mode === "login" ? { login: email, password } : { email, password, nickname };
      const response = await apiClient.post<{ access_token: string }>(endpoint, payload);
      const accessToken = response.data.access_token;
      userStore.setAccessToken(accessToken);
      const meResponse = await apiClient.get<Record<string, unknown>>("/users/me");
      const data = meResponse.data;
      const user: IUser = {
        id: data.id as string,
        nickname: data.nickname as string,
        email: data.email as string,
        avatarUrl: data.avatar_url as string | undefined,
        bio: data.bio as string | undefined,
      };
      userStore.setUser(user);

      void navigate({ to: redirectTo as "/" });
    } catch (error_: unknown) {
      if (isAxiosError(error_)) {
        const data = error_.response?.data as { message?: string } | undefined;
        setError(
          data?.message ??
            (mode === "login" ? "Произошла ошибка при входе" : "Произошла ошибка при регистрации"),
        );
      } else {
        setError(
          mode === "login" ? "Произошла ошибка при входе" : "Произошла ошибка при регистрации",
        );
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.logo}>
        <Logo className={styles.logoImage} height={60} />
      </div>

      <h1 className={styles.title}>{mode === "login" ? "Вход" : "Регистрация"}</h1>
      <p className={styles.subtitle}>
        {mode === "login"
          ? "Войдите в свой аккаунт, чтобы продолжить"
          : "Создайте аккаунт для начала работы"}
      </p>

      <form
        className={styles.form}
        onSubmit={(event) => {
          void handleSubmit(event);
        }}
      >
        {mode === "register" && (
          <div className={styles.field}>
            <span className={styles.label}>Имя пользователя</span>
            <Input
              placeholder="Например: JohnDoe"
              iconLeft={<User size={20} />}
              value={nickname}
              onChange={(event) => {
                setNickname(event.target.value);
              }}
              required
            />
          </div>
        )}
        <div className={styles.field}>
          <span className={styles.label}>
            {mode === "login" ? "Email или имя пользователя" : "Электронная почта"}
          </span>
          <Input
            placeholder={
              mode === "login" ? "your.email@example.com или username" : "your.email@example.com"
            }
            iconLeft={<Mail size={20} />}
            value={email}
            onChange={(event) => {
              setEmail(event.target.value);
            }}
            required
            type={mode === "login" ? "text" : "email"}
          />
        </div>

        <div className={styles.field}>
          <span className={styles.label}>Пароль</span>
          <Input
            placeholder={mode === "register" ? "Минимум 8 символов" : "Введите пароль"}
            type={showPassword ? "text" : "password"}
            iconLeft={<Lock size={20} />}
            iconRight={
              <button
                type="button"
                className={styles.iconButton}
                onClick={() => {
                  setShowPassword(!showPassword);
                }}
              >
                {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
              </button>
            }
            value={password}
            onChange={(event) => {
              setPassword(event.target.value);
            }}
            required
            minLength={mode === "register" ? 8 : undefined}
          />
        </div>

        {mode === "login" && (
          <div className={styles.options}>
            <label className={styles.checkbox}>
              <input type="checkbox" />
              <span>Запомнить меня</span>
            </label>
            <a href="#" className={styles.link}>
              Забыли пароль?
            </a>
          </div>
        )}

        <Button type="submit" fullWidth disabled={isLoading}>
          {((): string => {
            if (isLoading) return "Загрузка...";
            if (mode === "login") return "Войти";
            return "Зарегистрироваться";
          })()}
        </Button>

        {error && <div className={styles.errorMessage}>{error}</div>}
      </form>

      <div className={styles.divider}>
        <span>Или войдите через</span>
      </div>

      <div className={styles.socialButtons}>
        <Button variant="outline" type="button" className={styles.socialBtn} disabled>
          <GoogleIcon />
          Войти через Google (Скоро)
        </Button>
        <Button variant="outline" type="button" className={styles.socialBtn} disabled>
          <GithubIcon />
          Войти через GitHub (Скоро)
        </Button>
      </div>

      <div className={styles.register}>
        {mode === "login" ? (
          <>
            Нет аккаунта?{" "}
            <button type="button" onClick={toggleMode} className={styles.linkBtn}>
              Зарегистрироваться
            </button>
          </>
        ) : (
          <>
            Уже есть аккаунт?{" "}
            <button type="button" onClick={toggleMode} className={styles.linkBtn}>
              Войти
            </button>
          </>
        )}
      </div>
    </div>
  );
});

import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from "react";
import { authApi, type LoginPayload, type RegisterPayload } from "../api/auth";
import { clearStoredAuth, loadStoredAuth, saveStoredAuth } from "../lib/storage";
import type { AuthResponse, User } from "../types/api";

type AuthStatus = "loading" | "authenticated" | "unauthenticated";

type AuthContextValue = {
  status: AuthStatus;
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  login: (payload: LoginPayload) => Promise<void>;
  register: (payload: RegisterPayload) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

function persistAuth(response: AuthResponse): void {
  saveStoredAuth({
    token: response.token,
    user: response.user,
  });
}

export function AuthProvider({ children }: PropsWithChildren) {
  const [status, setStatus] = useState<AuthStatus>("loading");
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<User | null>(null);

  const applyAuth = (response: AuthResponse) => {
    persistAuth(response);
    setToken(response.token);
    setUser(response.user);
    setStatus("authenticated");
  };

  const clearAuthState = () => {
    clearStoredAuth();
    setToken(null);
    setUser(null);
    setStatus("unauthenticated");
  };

  const refreshUser = async () => {
    const currentToken = loadStoredAuth()?.token;
    if (!currentToken) {
      clearAuthState();
      return;
    }

    const currentUser = await authApi.me();
    saveStoredAuth({
      token: currentToken,
      user: currentUser,
    });
    setToken(currentToken);
    setUser(currentUser);
    setStatus("authenticated");
  };

  useEffect(() => {
    let cancelled = false;

    const hydrate = async () => {
      const stored = loadStoredAuth();
      if (!stored?.token) {
        if (!cancelled) {
          setStatus("unauthenticated");
        }
        return;
      }

      if (!cancelled) {
        setToken(stored.token);
        setUser(stored.user);
      }

      try {
        const freshUser = await authApi.me();
        if (cancelled) {
          return;
        }

        saveStoredAuth({
          token: stored.token,
          user: freshUser,
        });
        setToken(stored.token);
        setUser(freshUser);
        setStatus("authenticated");
      } catch {
        if (!cancelled) {
          clearAuthState();
        }
      }
    };

    void hydrate();

    return () => {
      cancelled = true;
    };
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      status,
      token,
      user,
      isAuthenticated: status === "authenticated" && Boolean(token),
      login: async (payload) => {
        const response = await authApi.login(payload);
        applyAuth(response);
      },
      register: async (payload) => {
        const response = await authApi.register(payload);
        applyAuth(response);
      },
      logout: clearAuthState,
      refreshUser,
    }),
    [status, token, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth debe usarse dentro de AuthProvider");
  }

  return context;
}

"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  getStoredUser,
  removeStoredUser,
  type StoredUser,
} from "@/lib/auth-storage";

export type AuthUser = {
  user_id: number;
  name: string;
  email: string;
};

type UseAuthResult = {
  user: AuthUser | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  logout: () => void;
};

export function useAuth(): UseAuthResult {
  const router = useRouter();
  const [user, setUser] = useState<StoredUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      setUser(getStoredUser());
      setIsLoading(false);
    }, 0);

    return () => window.clearTimeout(timeoutId);
  }, []);

  const logout = useCallback(() => {
    removeStoredUser();
    setUser(null);
    router.replace("/login");
  }, [router]);

  return {
    user,
    isLoading,
    isAuthenticated: Boolean(user),
    logout,
  };
}

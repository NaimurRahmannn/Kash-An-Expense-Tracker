export type StoredUser = {
  user_id: number;
  name: string;
  email: string;
};

const STORAGE_KEY = "expense_tracker_user";

function canUseStorage() {
  return typeof window !== "undefined" && Boolean(window.localStorage);
}

function isStoredUser(value: unknown): value is StoredUser {
  if (!value || typeof value !== "object") {
    return false;
  }

  const user = value as Partial<StoredUser>;

  return (
    typeof user.user_id === "number" &&
    typeof user.name === "string" &&
    typeof user.email === "string"
  );
}

export function saveUser(user: StoredUser): void {
  if (!canUseStorage()) {
    return;
  }

  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(user));
}

export function getStoredUser(): StoredUser | null {
  if (!canUseStorage()) {
    return null;
  }

  const rawUser = window.localStorage.getItem(STORAGE_KEY);

  if (!rawUser) {
    return null;
  }

  try {
    const parsedUser = JSON.parse(rawUser);

    return isStoredUser(parsedUser) ? parsedUser : null;
  } catch {
    return null;
  }
}

export function removeStoredUser(): void {
  if (!canUseStorage()) {
    return;
  }

  window.localStorage.removeItem(STORAGE_KEY);
}

export function getStoredUserID(): number | null {
  return getStoredUser()?.user_id ?? null;
}

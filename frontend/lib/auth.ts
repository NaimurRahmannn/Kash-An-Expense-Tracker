import { apiRequest } from "@/lib/api";
import type { ApiResponse } from "@/types/api";
import type { LoginInput, LoginResponseData, RegisterInput } from "@/types/auth";

export function registerUser(input: RegisterInput): Promise<ApiResponse> {
  return apiRequest<ApiResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function loginUser(
  input: LoginInput,
): Promise<ApiResponse<LoginResponseData>> {
  return apiRequest<ApiResponse<LoginResponseData>>("/auth/login", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

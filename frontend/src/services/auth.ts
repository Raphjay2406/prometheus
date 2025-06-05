// frontend/src/services/auth.ts
import axios, { AxiosError } from 'axios';

// 1. Define Type Aliases/Interfaces (mirroring backend structures)
// These should align with your Go backend's request and response structs.

export interface LoginPayload {
  username?: string; // Backend uses 'username' which can be username or email
  email?: string;    // Allow either for frontend clarity, backend handles logic
  password?: string;
}

export interface RegisterPayload {
  username?: string;
  email?: string;
  password?: string;
  role_id?: number; // Optional, backend defaults to 'staff' if not provided
}

// Compact user information returned on login
export interface User {
  id: number;
  username: string;
  email: string;
  role_name: string;
  is_active: boolean;
}

// Response from a successful login
export interface AuthSuccessResponse {
  user: User;
  access_token: string;
  refresh_token?: string; // Optional, if you implement refresh tokens
}

// Response from a successful registration (matches UserResponse in Go handler)
export interface RegisteredUserResponse {
  id: number;
  username: string;
  email: string;
  is_active: boolean;
  role_id: number;
  role_name?: string;
  created_at: string; // Dates are typically strings in JSON
  updated_at: string;
}

// Standard error structure from the backend's utils.SendErrorResponse
export interface ApiErrorResponse {
  status: string; // e.g., "error"
  message: string;
  errors?: Record<string, string> | string[]; // Optional detailed errors
}

// 2. Create an Axios instance
const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // Your backend API base URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// 3. Define Authentication Service Functions

/**
 * Logs in a user.
 * @param credentials - The user's login credentials (username/email and password).
 * @returns A promise that resolves with the authentication response (user data and token).
 */
export const loginUser = async (
  credentials: LoginPayload
): Promise<AuthSuccessResponse> => {
  try {
    const payload: Record<string, string | undefined> = {};
    // Backend expects 'username' field for either username or email
    payload.username = credentials.username || credentials.email;
    payload.password = credentials.password;

    if (!payload.username || !payload.password) {
      throw new Error('Username/Email and password are required.');
    }

    const response = await apiClient.post<AuthSuccessResponse>(
      '/auth/login',
      payload
    );
    return response.data;
  } catch (error) {
    const axiosError = error as AxiosError<ApiErrorResponse>;
    if (axiosError.response && axiosError.response.data) {
      throw new Error(axiosError.response.data.message || 'Login failed. Please try again.');
    }
    throw new Error('Login failed due to a network or server error. Please try again.');
  }
};

/**
 * Registers a new user.
 * @param userData - The data for the new user.
 * @returns A promise that resolves with the registered user's information.
 */
export const registerUser = async (
  userData: RegisterPayload
): Promise<RegisteredUserResponse> => {
  try {
    // Basic frontend validation (can be more extensive)
    if (!userData.username || !userData.email || !userData.password) {
      throw new Error('Username, email, and password are required for registration.');
    }
    if (userData.password.length < 6) {
        throw new Error('Password must be at least 6 characters long.');
    }

    const response = await apiClient.post<RegisteredUserResponse>(
      '/auth/register',
      userData
    );
    return response.data;
  } catch (error) {
    const axiosError = error as AxiosError<ApiErrorResponse>;
    if (axiosError.response && axiosError.response.data) {
      throw new Error(axiosError.response.data.message || 'Registration failed. Please try again.');
    }
    throw new Error('Registration failed due to a network or server error. Please try again.');
  }
};

// You can add other auth-related API calls here later, e.g.:
// - logoutUser (might just be local token removal, or could hit a backend endpoint)
// - refreshToken
// - forgotPassword
// - resetPassword

export default apiClient; // Exporting the configured Axios instance can be useful
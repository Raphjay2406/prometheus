// frontend/src/contexts/AuthContext.tsx
import React, {
  createContext,
  useState,
  useContext,
  useEffect,
  ReactNode,
  useCallback,
} from 'react';
import {
  loginUser as apiLoginUser,
  registerUser as apiRegisterUser,
  LoginPayload,
  RegisterPayload,
  User,
  AuthSuccessResponse,
} from '../services/auth'; // Assuming auth.ts is in ../services

// 1. Define the shape of the context data
interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  login: (credentials: LoginPayload) => Promise<void>;
  register: (userData: RegisterPayload) => Promise<void>; // Or Promise<User> if auto-login
  logout: () => void;
  clearError: () => void;
}

// 2. Create the context with a default undefined value (or a default state)
// Throw an error if used outside a provider
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// 3. Define the Props for the AuthProvider
interface AuthProviderProps {
  children: ReactNode;
}

// 4. Create the AuthProvider component
export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(true); // Start true for initial load check
  const [error, setError] = useState<string | null>(null);

  // Function to clear errors
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Check localStorage for token on initial load
  useEffect(() => {
    setIsLoading(true);
    try {
      const storedToken = localStorage.getItem('authToken');
      const storedUserString = localStorage.getItem('authUser');

      if (storedToken && storedUserString) {
        const storedUser: User = JSON.parse(storedUserString);
        setToken(storedToken);
        setUser(storedUser);
        setIsAuthenticated(true);
      }
    } catch (e) {
      // If parsing fails or any error, ensure clean state
      localStorage.removeItem('authToken');
      localStorage.removeItem('authUser');
      console.error("Error loading auth state from localStorage:", e);
    } finally {
      setIsLoading(false);
    }
  }, []);

  const handleAuthSuccess = (data: AuthSuccessResponse) => {
    setUser(data.user);
    setToken(data.access_token);
    setIsAuthenticated(true);
    localStorage.setItem('authToken', data.access_token);
    localStorage.setItem('authUser', JSON.stringify(data.user)); // Store user object
    setError(null); // Clear any previous errors
  };

  const login = async (credentials: LoginPayload) => {
    setIsLoading(true);
    clearError();
    try {
      const data = await apiLoginUser(credentials);
      handleAuthSuccess(data);
    } catch (err: any) {
      setError(err.message || 'Login failed. Please check your credentials.');
      setIsAuthenticated(false);
      setUser(null);
      setToken(null);
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (userData: RegisterPayload) => {
    setIsLoading(true);
    clearError();
    try {
      // The registerUser service returns RegisteredUserResponse.
      // For a seamless experience, we can attempt to log the user in
      // or use a separate login step. Here, we'll try to log them in
      // by calling the login service with their credentials if registration is successful.
      // This assumes the registration backend doesn't return a token directly.
      // If it did, this flow would be simpler.
      await apiRegisterUser(userData);
      // After successful registration, automatically log the user in
      // Ensure userData has username/email and password for login
      const loginCredentials: LoginPayload = {
        username: userData.username || userData.email, // Use username or email
        password: userData.password,
      };
      // We need to call the actual login function to get a token and set auth state
      await login(loginCredentials); // This will call apiLoginUser and handleAuthSuccess
      // If login after register is successful, `handleAuthSuccess` within `login` function will manage state.
    } catch (err: any) {
      setError(err.message || 'Registration failed.');
      // Ensure auth state is cleared if registration or subsequent login fails
      setIsAuthenticated(false);
      setUser(null);
      setToken(null);
      // Re-throw the error if you want calling components to also handle it
      // throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = useCallback(() => {
    setUser(null);
    setToken(null);
    setIsAuthenticated(false);
    localStorage.removeItem('authToken');
    localStorage.removeItem('authUser');
    setError(null);
    // Here you might also want to redirect to login page using useNavigate() from react-router-dom
    // or notify other parts of the app.
  }, []);


  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated,
        isLoading,
        error,
        login,
        register,
        logout,
        clearError,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

// 5. Create a custom hook to use the AuthContext
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
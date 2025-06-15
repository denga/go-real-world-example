"use client";

import React, { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { getCurrentUser, storeUserData, logout as authLogout } from "@/lib/auth";

// Define the User interface
interface User {
  username: string;
  email: string;
  token: string;
  bio: string;
  image: string;
}

// Define the context type
interface AuthContextType {
  user: User | null;
  login: (user: User) => void;
  logout: () => void;
  isLoading: boolean;
}

// Create the context with a default value
const AuthContext = createContext<AuthContextType>({
  user: null,
  login: () => {},
  logout: () => {},
  isLoading: true,
});

// Create a hook to use the auth context
export const useAuth = () => useContext(AuthContext);

// Create the auth provider component
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Load user from localStorage on initial render
  useEffect(() => {
    const loadUser = () => {
      const currentUser = getCurrentUser();
      setUser(currentUser);
      setIsLoading(false);
    };

    loadUser();
  }, []);

  // Function to login a user
  const login = (userData: User) => {
    storeUserData(userData);
    setUser(userData);
  };

  // Function to logout a user
  const logout = () => {
    authLogout();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}
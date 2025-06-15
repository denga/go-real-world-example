// Authentication utilities

// Define the User interface
interface User {
  email: string;
  token: string;
  username: string;
  bio: string;
  image: string;
}

/**
 * Get the current user from localStorage
 */
export function getCurrentUser(): User | null {
  if (typeof window === 'undefined') {
    return null;
  }

  const storedUser = localStorage.getItem("user");
  if (!storedUser) {
    return null;
  }

  try {
    return JSON.parse(storedUser);
  } catch (error) {
    console.error("Failed to parse user data:", error);
    localStorage.removeItem("user");
    localStorage.removeItem("token");
    return null;
  }
}

/**
 * Get the JWT token from localStorage
 */
export function getToken() {
  if (typeof window === 'undefined') {
    return null;
  }

  return localStorage.getItem("token");
}

/**
 * Check if a user is logged in
 */
export function isLoggedIn() {
  return !!getToken();
}

/**
 * Log out the current user
 */
export function logout() {
  if (typeof window === 'undefined') {
    return;
  }

  localStorage.removeItem("user");
  localStorage.removeItem("token");
}

/**
 * Store user data and token in localStorage
 */
export function storeUserData(user: User) {
  if (typeof window === 'undefined') {
    return;
  }

  localStorage.setItem("token", user.token);
  localStorage.setItem("user", JSON.stringify(user));
}

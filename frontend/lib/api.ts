// API client for connecting to the backend API

// Base URL for the API
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Generic fetch function with error handling
async function fetchAPI<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'An error occurred while fetching the data.');
  }

  return response.json();
}

// API functions for different endpoints
export const api = {
  // Articles
  getArticles: (params?: { tag?: string; author?: string; favorited?: string; limit?: number; offset?: number }) => 
    fetchAPI<any>('/articles' + (params ? `?${new URLSearchParams(params as any).toString()}` : '')),
  
  getArticle: (slug: string) => 
    fetchAPI<any>(`/articles/${slug}`),
  
  // Tags
  getTags: () => 
    fetchAPI<any>('/tags'),
  
  // User
  getCurrentUser: (token: string) => 
    fetchAPI<any>('/user', {
      headers: {
        Authorization: `Token ${token}`,
      },
    }),
  
  // Authentication
  login: (email: string, password: string) => 
    fetchAPI<any>('/users/login', {
      method: 'POST',
      body: JSON.stringify({
        user: {
          email,
          password,
        },
      }),
    }),
  
  register: (username: string, email: string, password: string) => 
    fetchAPI<any>('/users', {
      method: 'POST',
      body: JSON.stringify({
        user: {
          username,
          email,
          password,
        },
      }),
    }),
};
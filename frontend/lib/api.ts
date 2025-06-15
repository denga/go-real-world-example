// API client for connecting to the backend API

// Base URL for the API
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

// Define interfaces for API responses based on the OpenAPI spec
interface Profile {
  username: string;
  bio: string;
  image: string;
  following: boolean;
}

interface Article {
  slug: string;
  title: string;
  description: string;
  body: string;
  tagList: string[];
  createdAt: string;
  updatedAt: string;
  favorited: boolean;
  favoritesCount: number;
  author: Profile;
}

interface MultipleArticlesResponse {
  articles: Article[];
  articlesCount: number;
}

interface SingleArticleResponse {
  article: Article;
}

interface TagsResponse {
  tags: string[];
}

interface User {
  email: string;
  token: string;
  username: string;
  bio: string;
  image: string;
}

interface UserResponse {
  user: User;
}

// Type for URL search params to avoid 'as any' cast
type SearchParamsObject = Record<string, string | number | undefined>;

// Helper function to convert SearchParamsObject to a format acceptable by URLSearchParams
function convertToURLSearchParams(params: SearchParamsObject): Record<string, string> {
  const result: Record<string, string> = {};
  for (const key in params) {
    const value = params[key];
    if (value !== undefined) {
      result[key] = String(value);
    }
  }
  return result;
}

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
    fetchAPI<MultipleArticlesResponse>('/articles' + (params ? `?${new URLSearchParams(convertToURLSearchParams(params)).toString()}` : '')),

  getArticle: (slug: string) => 
    fetchAPI<SingleArticleResponse>(`/articles/${slug}`),

  // Tags
  getTags: () => 
    fetchAPI<TagsResponse>('/tags'),

  // User
  getCurrentUser: (token: string) => 
    fetchAPI<UserResponse>('/user', {
      headers: {
        Authorization: `Token ${token}`,
      },
    }),

  // Authentication
  login: (email: string, password: string) => 
    fetchAPI<UserResponse>('/users/login', {
      method: 'POST',
      body: JSON.stringify({
        user: {
          email,
          password,
        },
      }),
    }),

  register: (username: string, email: string, password: string) => 
    fetchAPI<UserResponse>('/users', {
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

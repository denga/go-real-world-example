"use client";

import Link from "next/link";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { useAuth } from "@/contexts/auth-context";

export default function RegisterPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const response = await api.register(username, email, password);

      // Update auth context with user data
      login(response.user);

      // Redirect to home page
      router.push("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-md">
      <h1 className="text-2xl font-bold text-center mb-8">Sign Up</h1>

      <div className="text-center mb-4">
        <Link href="/login" className="text-primary hover:underline">
          Have an account?
        </Link>
      </div>

      {error && (
        <div className="bg-destructive/10 text-destructive p-3 rounded-md mb-4">
          {error}
        </div>
      )}

      <form className="space-y-4" onSubmit={handleSubmit}>
        <div>
          <input
            type="text"
            placeholder="Username"
            className="w-full p-2 border rounded-md"
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            disabled={isLoading}
          />
        </div>

        <div>
          <input
            type="email"
            placeholder="Email"
            className="w-full p-2 border rounded-md"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            disabled={isLoading}
          />
        </div>

        <div>
          <input
            type="password"
            placeholder="Password"
            className="w-full p-2 border rounded-md"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            disabled={isLoading}
          />
        </div>

        <div>
          <button
            type="submit"
            className="w-full bg-primary text-primary-foreground p-2 rounded-md hover:bg-primary/90 disabled:opacity-50"
            disabled={isLoading}
          >
            {isLoading ? "Signing up..." : "Sign up"}
          </button>
        </div>
      </form>
    </div>
  );
}

"use client";

import Link from "next/link";
import { ThemeToggle } from "./theme-toggle";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/auth-context";

export function Nav() {
  const router = useRouter();
  const { user, logout } = useAuth();

  const handleLogout = () => {
    // Clear user data and token from localStorage
    logout();
    router.push("/");
  };

  return (
    <nav className="w-full flex justify-between items-center py-4 px-4 md:px-8 border-b">
      <div className="flex items-center space-x-4">
        <Link href="/" className="text-xl font-bold">
          Real World Example
        </Link>
        <div className="hidden md:flex space-x-4">
          <Link href="/" className="text-sm hover:text-primary">
            Home
          </Link>
          {!user ? (
            <>
              <Link href="/login" className="text-sm hover:text-primary">
                Sign in
              </Link>
              <Link href="/register" className="text-sm hover:text-primary">
                Sign up
              </Link>
            </>
          ) : (
            <>
              <Link href="/editor" className="text-sm hover:text-primary">
                New Article
              </Link>
              <Link href="/settings" className="text-sm hover:text-primary">
                Settings
              </Link>
              <Link href={`/profile/${user.username}`} className="text-sm hover:text-primary">
                {user.username}
              </Link>
              <button 
                onClick={handleLogout}
                className="text-sm hover:text-primary"
              >
                Logout
              </button>
            </>
          )}
        </div>
      </div>
      <div className="flex items-center space-x-4">
        {!user && (
          <div className="md:hidden flex space-x-2">
            <Link href="/login" className="text-sm hover:text-primary">
              Sign in
            </Link>
            <Link href="/register" className="text-sm hover:text-primary">
              Sign up
            </Link>
          </div>
        )}
        {user && (
          <div className="md:hidden flex space-x-2">
            <Link href="/editor" className="text-sm hover:text-primary">
              New Article
            </Link>
            <button 
              onClick={handleLogout}
              className="text-sm hover:text-primary"
            >
              Logout
            </button>
          </div>
        )}
        <ThemeToggle />
      </div>
    </nav>
  );
}

import Link from "next/link";
import { ThemeToggle } from "./theme-toggle";

export function Nav() {
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
          <Link href="/login" className="text-sm hover:text-primary">
            Sign in
          </Link>
          <Link href="/register" className="text-sm hover:text-primary">
            Sign up
          </Link>
        </div>
      </div>
      <div className="flex items-center space-x-4">
        <Link href="/editor" className="text-sm hover:text-primary hidden md:inline-block">
          New Article
        </Link>
        <Link href="/settings" className="text-sm hover:text-primary hidden md:inline-block">
          Settings
        </Link>
        <ThemeToggle />
      </div>
    </nav>
  );
}
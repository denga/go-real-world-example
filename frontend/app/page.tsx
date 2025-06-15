import Image from "next/image";
import { ThemeToggle } from "@/components/theme-toggle";

export default function Home() {
  return (
    <div className="grid grid-rows-[auto_1fr_auto] items-center justify-items-center min-h-screen p-4 md:p-8 font-[family-name:var(--font-geist-sans)]">
      <header className="w-full flex justify-between items-center py-4 px-4 md:px-8">
        <h1 className="text-xl font-bold">Real World Example</h1>
        <ThemeToggle />
      </header>

      <main className="flex flex-col gap-6 items-center text-center max-w-4xl mx-auto py-8">
        <h2 className="text-3xl md:text-4xl font-bold">Welcome to the Real World Example</h2>
        <p className="text-lg text-muted-foreground">
          A Next.js application with chadcn UI, dark mode, and mobile responsiveness
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8 w-full">
          <div className="bg-card p-6 rounded-lg shadow-sm border">
            <h3 className="text-xl font-semibold mb-2">Dark Mode</h3>
            <p className="text-muted-foreground">
              Toggle between light and dark mode using the button in the header.
            </p>
          </div>

          <div className="bg-card p-6 rounded-lg shadow-sm border">
            <h3 className="text-xl font-semibold mb-2">Mobile Ready</h3>
            <p className="text-muted-foreground">
              The UI is fully responsive and works well on all device sizes.
            </p>
          </div>

          <div className="bg-card p-6 rounded-lg shadow-sm border">
            <h3 className="text-xl font-semibold mb-2">Backend API</h3>
            <p className="text-muted-foreground">
              Connected to the backend API defined in the OpenAPI specs.
            </p>
          </div>

          <div className="bg-card p-6 rounded-lg shadow-sm border">
            <h3 className="text-xl font-semibold mb-2">chadcn UI</h3>
            <p className="text-muted-foreground">
              Built with chadcn UI components for a consistent and beautiful design.
            </p>
          </div>
        </div>
      </main>

      <footer className="w-full py-6 px-4 md:px-8 text-center border-t">
        <p className="text-sm text-muted-foreground">
          © {new Date().getFullYear()} Real World Example. All rights reserved.
        </p>
      </footer>
    </div>
  );
}

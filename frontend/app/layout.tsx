import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/theme-provider";
import { Nav } from "@/components/nav";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Real World Example",
  description: "A real world example using Next.js and chadcn UI",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <div className="min-h-screen flex flex-col">
            <Nav />
            <main className="flex-1">
              {children}
            </main>
            <footer className="py-6 px-4 md:px-8 text-center border-t">
              <p className="text-sm text-muted-foreground">
                © {new Date().getFullYear()} Real World Example. All rights reserved.
              </p>
            </footer>
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}

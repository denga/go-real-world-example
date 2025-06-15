import Link from "next/link";

export default function LoginPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-md">
      <h1 className="text-2xl font-bold text-center mb-8">Sign In</h1>
      
      <div className="text-center mb-4">
        <Link href="/register" className="text-primary hover:underline">
          Need an account?
        </Link>
      </div>
      
      <form className="space-y-4">
        <div>
          <input
            type="email"
            placeholder="Email"
            className="w-full p-2 border rounded-md"
            required
          />
        </div>
        
        <div>
          <input
            type="password"
            placeholder="Password"
            className="w-full p-2 border rounded-md"
            required
          />
        </div>
        
        <div>
          <button
            type="submit"
            className="w-full bg-primary text-primary-foreground p-2 rounded-md hover:bg-primary/90"
          >
            Sign in
          </button>
        </div>
      </form>
    </div>
  );
}
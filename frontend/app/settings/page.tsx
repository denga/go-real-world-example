export default function SettingsPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-lg">
      <h1 className="text-2xl font-bold text-center mb-8">Your Settings</h1>
      
      <form className="space-y-6">
        <div className="space-y-2">
          <label className="text-sm font-medium">Profile Picture URL</label>
          <input
            type="url"
            placeholder="URL of profile picture"
            className="w-full p-2 border rounded-md"
          />
        </div>
        
        <div className="space-y-2">
          <label className="text-sm font-medium">Username</label>
          <input
            type="text"
            placeholder="Username"
            className="w-full p-2 border rounded-md"
          />
        </div>
        
        <div className="space-y-2">
          <label className="text-sm font-medium">Bio</label>
          <textarea
            placeholder="Short bio about you"
            className="w-full p-2 border rounded-md h-24"
          ></textarea>
        </div>
        
        <div className="space-y-2">
          <label className="text-sm font-medium">Email</label>
          <input
            type="email"
            placeholder="Email"
            className="w-full p-2 border rounded-md"
          />
        </div>
        
        <div className="space-y-2">
          <label className="text-sm font-medium">New Password</label>
          <input
            type="password"
            placeholder="New Password"
            className="w-full p-2 border rounded-md"
          />
        </div>
        
        <div>
          <button
            type="submit"
            className="w-full bg-primary text-primary-foreground p-2 rounded-md hover:bg-primary/90"
          >
            Update Settings
          </button>
        </div>
      </form>
      
      <hr className="my-8" />
      
      <button
        className="w-full bg-destructive text-destructive-foreground p-2 rounded-md hover:bg-destructive/90"
      >
        Logout
      </button>
    </div>
  );
}
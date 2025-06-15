export default function EditorPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <h1 className="text-2xl font-bold mb-8">Create New Article</h1>
      
      <form className="space-y-6">
        <div className="space-y-2">
          <input
            type="text"
            placeholder="Article Title"
            className="w-full p-2 border rounded-md text-xl"
            required
          />
        </div>
        
        <div className="space-y-2">
          <input
            type="text"
            placeholder="What's this article about?"
            className="w-full p-2 border rounded-md"
            required
          />
        </div>
        
        <div className="space-y-2">
          <textarea
            placeholder="Write your article (in markdown)"
            className="w-full p-2 border rounded-md h-64 font-mono"
            required
          ></textarea>
        </div>
        
        <div className="space-y-2">
          <input
            type="text"
            placeholder="Enter tags (comma separated)"
            className="w-full p-2 border rounded-md"
          />
        </div>
        
        <div>
          <button
            type="submit"
            className="float-right bg-primary text-primary-foreground px-4 py-2 rounded-md hover:bg-primary/90"
          >
            Publish Article
          </button>
        </div>
      </form>
    </div>
  );
}
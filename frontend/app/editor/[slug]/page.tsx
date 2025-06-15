export default function EditArticlePage({ params }: { params: { slug: string } }) {
  // In a real app, we would fetch the article data based on the slug
  const mockArticle = {
    title: "Example Article",
    description: "This is an example article",
    body: "# Example Article\n\nThis is the body of the example article in markdown format.\n\n## Section 1\n\nSome content here.\n\n## Section 2\n\nMore content here.",
    tagList: ["example", "article", "markdown"]
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <h1 className="text-2xl font-bold mb-8">Edit Article</h1>
      
      <form className="space-y-6">
        <div className="space-y-2">
          <input
            type="text"
            placeholder="Article Title"
            className="w-full p-2 border rounded-md text-xl"
            defaultValue={mockArticle.title}
            required
          />
        </div>
        
        <div className="space-y-2">
          <input
            type="text"
            placeholder="What's this article about?"
            className="w-full p-2 border rounded-md"
            defaultValue={mockArticle.description}
            required
          />
        </div>
        
        <div className="space-y-2">
          <textarea
            placeholder="Write your article (in markdown)"
            className="w-full p-2 border rounded-md h-64 font-mono"
            defaultValue={mockArticle.body}
            required
          ></textarea>
        </div>
        
        <div className="space-y-2">
          <input
            type="text"
            placeholder="Enter tags (comma separated)"
            className="w-full p-2 border rounded-md"
            defaultValue={mockArticle.tagList.join(', ')}
          />
        </div>
        
        <div className="flex justify-end space-x-4">
          <button
            type="button"
            className="bg-destructive text-destructive-foreground px-4 py-2 rounded-md hover:bg-destructive/90"
          >
            Delete Article
          </button>
          <button
            type="submit"
            className="bg-primary text-primary-foreground px-4 py-2 rounded-md hover:bg-primary/90"
          >
            Update Article
          </button>
        </div>
      </form>
    </div>
  );
}
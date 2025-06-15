import Link from "next/link";
import Image from "next/image";

// This function is required for static site generation with dynamic routes
export async function generateStaticParams() {
  // In a real app, we would fetch all article slugs from an API or database
  // For now, we'll return a few mock slugs
  return [
    { slug: 'how-to-use-react-hooks' },
    { slug: 'getting-started-with-golang' },
    { slug: 'nextjs-app-router' }
  ];
}

export default async function ArticlePage({ params }: { params: Promise<{ slug: string }> }) {
  // In a real app, we would fetch the article data based on the slug
  const resolvedParams = await params;
  const mockArticle = {
    slug: resolvedParams.slug,
    title: "Example Article",
    description: "This is an example article",
    body: "# Example Article\n\nThis is the body of the example article in markdown format.\n\n## Section 1\n\nSome content here.\n\n## Section 2\n\nMore content here.",
    tagList: ["example", "article", "markdown"],
    createdAt: "2023-06-15",
    favorited: false,
    favoritesCount: 12,
    author: {
      username: "johndoe",
      bio: "I work at statefarm",
      image: "https://picsum.photos/id/40/300/200",
      following: false
    }
  };

  // Mock comments
  const mockComments = [
    {
      id: 1,
      createdAt: "2023-06-16",
      body: "This is a great article!",
      author: {
        username: "janedoe",
        image: "https://picsum.photos/id/40/300/200",
      }
    },
    {
      id: 2,
      createdAt: "2023-06-17",
      body: "I learned a lot from this, thanks for sharing!",
      author: {
        username: "alexsmith",
        image: "https://picsum.photos/id/40/300/200",
      }
    }
  ];

  // Simple markdown renderer (in a real app, use a proper markdown library)
  const renderMarkdown = (markdown: string) => {
    const html = markdown
      .replace(/^# (.*$)/gm, '<h1 class="text-3xl font-bold my-4">$1</h1>')
      .replace(/^## (.*$)/gm, '<h2 class="text-2xl font-bold my-3">$1</h2>')
      .replace(/^### (.*$)/gm, '<h3 class="text-xl font-bold my-2">$1</h3>')
      .replace(/\n/g, '<br />');

    return <div dangerouslySetInnerHTML={{ __html: html }} />;
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Article header */}
      <div className="bg-muted py-8 px-4 rounded-lg mb-8">
        <h1 className="text-3xl font-bold mb-4">{mockArticle.title}</h1>
        <div className="flex items-center mb-4">
          <Image 
            src={mockArticle.author.image} 
            alt={mockArticle.author.username}
            width={40}
            height={40}
            className="w-10 h-10 rounded-full mr-2"
          />
          <div>
            <Link 
              href={`/profile/${mockArticle.author.username}`}
              className="font-medium hover:underline"
            >
              {mockArticle.author.username}
            </Link>
            <p className="text-xs text-muted-foreground">{mockArticle.createdAt}</p>
          </div>
        </div>
        <div className="flex space-x-2">
          <button className="text-sm px-3 py-1 border rounded-md hover:bg-muted-foreground/10">
            + Follow {mockArticle.author.username}
          </button>
          <button className="text-sm px-3 py-1 border rounded-md hover:bg-muted-foreground/10">
            ♥ Favorite Article ({mockArticle.favoritesCount})
          </button>
          {/* Edit/Delete buttons (only shown to article's author) */}
          <Link 
            href={`/editor/${mockArticle.slug}`}
            className="text-sm px-3 py-1 border rounded-md hover:bg-muted-foreground/10"
          >
            Edit Article
          </Link>
          <button className="text-sm px-3 py-1 border rounded-md bg-destructive text-destructive-foreground hover:bg-destructive/90">
            Delete Article
          </button>
        </div>
      </div>

      {/* Article body */}
      <div className="prose dark:prose-invert max-w-none mb-8">
        {renderMarkdown(mockArticle.body)}
      </div>

      {/* Tags */}
      <div className="flex flex-wrap gap-2 mb-8">
        {mockArticle.tagList.map((tag) => (
          <span 
            key={tag} 
            className="text-xs px-2 py-1 bg-muted rounded-full"
          >
            {tag}
          </span>
        ))}
      </div>

      <hr className="my-8" />

      {/* Comments section */}
      <div className="max-w-3xl mx-auto">
        <h3 className="text-xl font-bold mb-4">Comments</h3>

        {/* Comment form */}
        <div className="border rounded-lg p-4 mb-6">
          <textarea
            placeholder="Write a comment..."
            className="w-full p-2 border rounded-md h-24 mb-2"
          ></textarea>
          <div className="flex justify-end">
            <button className="bg-primary text-primary-foreground px-4 py-2 rounded-md hover:bg-primary/90">
              Post Comment
            </button>
          </div>
        </div>

        {/* Comments list */}
        <div className="space-y-4">
          {mockComments.map((comment) => (
            <div key={comment.id} className="border rounded-lg p-4">
              <p className="mb-4">{comment.body}</p>
              <div className="flex justify-between items-center">
                <div className="flex items-center">
                  <Image 
                    src={comment.author.image} 
                    alt={comment.author.username}
                    width={32}
                    height={32}
                    className="w-8 h-8 rounded-full mr-2"
                  />
                  <div>
                    <Link 
                      href={`/profile/${comment.author.username}`}
                      className="text-sm font-medium hover:underline"
                    >
                      {comment.author.username}
                    </Link>
                    <p className="text-xs text-muted-foreground">{comment.createdAt}</p>
                  </div>
                </div>
                {/* Delete button (only shown to comment's author) */}
                <button className="text-xs text-destructive hover:text-destructive/80">
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

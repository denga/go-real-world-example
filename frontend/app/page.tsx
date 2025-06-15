import Link from "next/link";
import Image from "next/image";

export default function Home() {
  // Mock data for demonstration
  const tags = ["programming", "javascript", "webdev", "react", "nextjs", "golang", "tutorial"];
  const articles = [
    {
      slug: "how-to-use-react-hooks",
      title: "How to Use React Hooks",
      description: "A comprehensive guide to React Hooks",
      author: {
        username: "johndoe",
        image: "https://via.placeholder.com/40"
      },
      createdAt: "2023-06-15",
      favorited: false,
      favoritesCount: 12,
      tagList: ["react", "javascript", "webdev"]
    },
    {
      slug: "getting-started-with-golang",
      title: "Getting Started with Golang",
      description: "Learn the basics of Go programming language",
      author: {
        username: "janedoe",
        image: "https://via.placeholder.com/40"
      },
      createdAt: "2023-06-10",
      favorited: true,
      favoritesCount: 8,
      tagList: ["golang", "programming", "tutorial"]
    },
    {
      slug: "nextjs-app-router",
      title: "Understanding Next.js App Router",
      description: "A deep dive into the Next.js App Router",
      author: {
        username: "alexsmith",
        image: "https://via.placeholder.com/40"
      },
      createdAt: "2023-06-05",
      favorited: false,
      favoritesCount: 5,
      tagList: ["nextjs", "react", "javascript"]
    }
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col md:flex-row gap-8">
        {/* Main content */}
        <div className="flex-1">
          {/* Feed tabs */}
          <div className="border-b mb-6">
            <div className="flex space-x-6">
              <button className="py-2 border-b-2 border-primary font-medium">
                Global Feed
              </button>
              <button className="py-2 text-muted-foreground">
                Your Feed
              </button>
            </div>
          </div>

          {/* Articles list */}
          <div className="space-y-6">
            {articles.map((article) => (
              <div key={article.slug} className="border rounded-lg p-6">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex items-center">
                    <Image 
                      src={article.author.image} 
                      alt={article.author.username}
                      width={40}
                      height={40}
                      className="w-10 h-10 rounded-full mr-2"
                    />
                    <div>
                      <Link 
                        href={`/profile/${article.author.username}`}
                        className="font-medium hover:underline"
                      >
                        {article.author.username}
                      </Link>
                      <p className="text-xs text-muted-foreground">{article.createdAt}</p>
                    </div>
                  </div>
                  <button className="text-sm px-2 py-1 border rounded-md hover:bg-muted">
                    ♥ {article.favoritesCount}
                  </button>
                </div>
                <Link href={`/article/${article.slug}`}>
                  <h2 className="text-xl font-bold mb-2 hover:underline">{article.title}</h2>
                  <p className="text-muted-foreground mb-4">{article.description}</p>
                </Link>
                <div className="flex justify-between items-center">
                  <Link 
                    href={`/article/${article.slug}`}
                    className="text-sm text-muted-foreground hover:text-foreground"
                  >
                    Read more...
                  </Link>
                  <div className="flex flex-wrap gap-2">
                    {article.tagList.map((tag) => (
                      <span 
                        key={tag} 
                        className="text-xs px-2 py-1 bg-muted rounded-full"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Pagination */}
          <div className="flex justify-center mt-8">
            <div className="flex space-x-2">
              <button className="px-3 py-1 border rounded-md hover:bg-muted">1</button>
              <button className="px-3 py-1 border rounded-md hover:bg-muted">2</button>
              <button className="px-3 py-1 border rounded-md hover:bg-muted">3</button>
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="w-full md:w-80 shrink-0">
          <div className="bg-muted p-4 rounded-lg">
            <h3 className="font-semibold mb-2">Popular Tags</h3>
            <div className="flex flex-wrap gap-2">
              {tags.map((tag) => (
                <Link 
                  key={tag}
                  href={`/?tag=${tag}`}
                  className="text-xs px-2 py-1 bg-secondary text-secondary-foreground rounded-full hover:bg-secondary/80"
                >
                  {tag}
                </Link>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

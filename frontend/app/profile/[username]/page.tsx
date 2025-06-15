import Link from "next/link";

export default function ProfilePage({ params }: { params: { username: string } }) {
  // In a real app, we would fetch the profile data based on the username
  const mockProfile = {
    username: params.username,
    bio: "I work at statefarm",
    image: "https://via.placeholder.com/100",
    following: false
  };

  // Mock articles
  const mockArticles = [
    {
      slug: "how-to-use-react-hooks",
      title: "How to Use React Hooks",
      description: "A comprehensive guide to React Hooks",
      author: {
        username: params.username,
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
        username: params.username,
        image: "https://via.placeholder.com/40"
      },
      createdAt: "2023-06-10",
      favorited: true,
      favoritesCount: 8,
      tagList: ["golang", "programming", "tutorial"]
    }
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Profile header */}
      <div className="bg-muted py-8 px-4 rounded-lg mb-8 text-center">
        <img 
          src={mockProfile.image} 
          alt={mockProfile.username}
          className="w-24 h-24 rounded-full mx-auto mb-4"
        />
        <h1 className="text-2xl font-bold mb-2">{mockProfile.username}</h1>
        <p className="text-muted-foreground mb-4">{mockProfile.bio}</p>
        <button className="text-sm px-3 py-1 border rounded-md hover:bg-muted-foreground/10">
          {mockProfile.following ? 'Unfollow' : 'Follow'} {mockProfile.username}
        </button>
      </div>

      {/* Articles tabs */}
      <div className="border-b mb-6">
        <div className="flex space-x-6">
          <Link 
            href={`/profile/${params.username}`}
            className="py-2 border-b-2 border-primary font-medium"
          >
            My Articles
          </Link>
          <Link 
            href={`/profile/${params.username}/favorites`}
            className="py-2 text-muted-foreground"
          >
            Favorited Articles
          </Link>
        </div>
      </div>

      {/* Articles list */}
      <div className="space-y-6">
        {mockArticles.map((article) => (
          <div key={article.slug} className="border rounded-lg p-6">
            <div className="flex justify-between items-start mb-4">
              <div className="flex items-center">
                <img 
                  src={article.author.image} 
                  alt={article.author.username}
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
    </div>
  );
}
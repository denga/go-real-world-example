import Link from "next/link";
import Image from "next/image";

// This function is required for static site generation with dynamic routes
export async function generateStaticParams() {
  // In a real app, we would fetch all usernames from an API or database
  // For now, we'll return a few mock usernames
  return [
    { username: 'johndoe' },
    { username: 'janedoe' },
    { username: 'alexsmith' }
  ];
}

export default async function ProfileFavoritesPage({ params }: { params: Promise<{ username: string }> }) {
  // In a real app, we would fetch the profile data based on the username
  const resolvedParams = await params;
  const mockProfile = {
    username: resolvedParams.username,
    bio: "I work at statefarm",
    image: "https://picsum.photos/id/100/300/200",
    following: false
  };

  // Mock favorited articles
  const mockFavoritedArticles = [
    {
      slug: "nextjs-app-router",
      title: "Understanding Next.js App Router",
      description: "A deep dive into the Next.js App Router",
      author: {
        username: "alexsmith",
        image: "https://picsum.photos/id/40/300/200"
      },
      createdAt: "2023-06-05",
      favorited: true,
      favoritesCount: 5,
      tagList: ["nextjs", "react", "javascript"]
    },
    {
      slug: "tailwind-css-tips",
      title: "Tailwind CSS Tips and Tricks",
      description: "Learn how to use Tailwind CSS effectively",
      author: {
        username: "janedoe",
        image: "https://picsum.photos/id/40/300/200"
      },
      createdAt: "2023-06-01",
      favorited: true,
      favoritesCount: 10,
      tagList: ["css", "tailwind", "webdev"]
    }
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Profile header */}
      <div className="bg-muted py-8 px-4 rounded-lg mb-8 text-center">
        <Image 
          src={mockProfile.image} 
          alt={mockProfile.username}
          width={96}
          height={96}
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
            href={`/profile/${resolvedParams.username}`}
            className="py-2 text-muted-foreground"
          >
            My Articles
          </Link>
          <Link 
            href={`/profile/${resolvedParams.username}/favorites`}
            className="py-2 border-b-2 border-primary font-medium"
          >
            Favorited Articles
          </Link>
        </div>
      </div>

      {/* Articles list */}
      <div className="space-y-6">
        {mockFavoritedArticles.map((article) => (
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
    </div>
  );
}

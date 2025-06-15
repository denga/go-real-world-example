# Requirements

## Backend

- We want to implement a real world example in go with the chi framework and oapi-codegen
- The open api definition can be found under https://github.com/gothinkster/realworld/blob/main/api/openapi.yml and should be placed in the local project folder
- The openapi spec should be embedded in the binary 
- As a datasource we use an simple in memory database

## Frontend

- We want to implement a web frontend that use nextjs framework and chadcn ui
- Use npx create-next-app@latest to create the frontend so we have a good structure to start
- The frontend is placed in the frontend directory
- The frontend calls the backend api defined in the openapi specs
- The frontend should have a dark mode and is mobile ready
- The build frontend is embedded in the backend binary and we can access it through the server 

### Routing
- Home page (URL: / )
  - List of tags
  - List of articles pulled from either Feed, Global, or by Tag
  - Pagination for list of articles
- Sign in/Sign up pages (URL: /login, /register )
  - Uses JWT (store the token in localStorage)
  - Authentication can be easily switched to session/cookie based
  - Settings page (URL: /settings )
- Editor page to create/edit articles (URL: /editor, /editor/article-slug-here )
- Article page (URL: /article/article-slug-here )
  - Delete article button (only shown to article’s author)
  - Render markdown from server client side
  - Comments section at bottom of page
  - Delete comment button (only shown to comment’s author)
- Profile page (URL: /profile/:username, /profile/:username/favorites )
  - Show basic user info
  - List of articles populated from author’s created articles or author’s favorited articles
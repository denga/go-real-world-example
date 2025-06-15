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
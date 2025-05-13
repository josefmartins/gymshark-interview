## How to use
- Access `http://ec2-18-212-206-26.compute-1.amazonaws.com/` and use the UI.
- If you want to see the backend API, visit the page `localhost:8080/docs` - requires local build. 

#### Calculate Number of Packages
![calculate packaging](docs/calculate_packaging.png "Calculate Packaging")

#### Create a New Product
![create product](docs/create_product_1.png "Create Product 1")
![create product 2](docs/create_product_2.png "Create Product 2")

#### Delete a Product
![delete product](docs/delete_product_1.png "Delete Product")

## How to build
#### Requirements
- Go and Node should be installed locally. If go is not installed, there is a Dockerfile available under `/build`.
#### Backend 
- in a shell: change the working directory to `/backend` 
- run `make` . That'll start the application and you'll see a log entry similar to `2025/09/32 21:37:04 server started. listening on port :8080` 
#### Frontend
- in a shell: change the working directory to `/backend/product-app`
- run `npm start`. It might take some seconds to minutes but eventually the app will be available in a browser via `localhost:3000`

## Notes on implementation
- Using SQLite for ease of deployment. In production, this could be a SQL/NoSQL database with persistence, such as PostgreSQL or DynamoDB.
- Using github.com/rubenv/sql-migrate for setting the database scheme and seed test examples.
- List products endpoint should be paginated.   
- Split the backend in 3 different layers to keep domains segregated: server, service (actual business logic) and storage. models package is common to the logical layers and makes mapping easier.
- REST API can be split into 2: CRUD for products and specific add/remove package size to product and calculate package units. API docs can be consulted in `/docs` HTTP endpoint.
- Calculation Algorithm was implemented using a dynamic programming logic that looks for the least amount of packages to be sent and then looks for the solution that ships less units, ie. smaller package sizes. In addition it looks up the least common multiple to create an offset of biggest package in size. 
- I spent much more time on the backend than in the frontend. Frontend was quickly built using React and Typescript since those are the technologies I'm more comfortable with. 
- Disclaimer: I've used AI (ie. chatgpt) to create boilerplate code. This task took me some hours and using AI made it a bit faster and less tedious.

## Notes on testing 
- Wrote unit tests for service layer. 
- Wrote Integration tests (`/backend/tests`) to validate storage and server layer comply with the requirements.
- My opinion on unit tests is: they're good for complex use cases but they miss the integration layers (which in a service like this one are more than 50% of the logic). That being said I think they should be complemented with integration or E2E tests that actually validate the API is working as expected and won't break on any changes.

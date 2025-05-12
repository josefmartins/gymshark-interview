## How to build

## How to use

## Notes on implementation
- Using SQLite for ease of deployment. In production, this could be a SQL/NoSQL database with persistence, such as PostgreSQL or DynamoDB.
- Using github.com/rubenv/sql-migrate for setting the database scheme and seed test examples.
- List products endpoint should be paginated.   
- Split the backend in 3 different layers to keep domains segregated: server, service (actual business logic) and storage. models package is common to the logical layers and makes mapping easier.
- REST API can be split into 2: CRUD for products and specific add/remove package size to product and calculate package units. API docs can be consulted in: TODO 
- Calculation Algorithm was implemented using a dynamic programming logic that looks for the least amount of packages to be sent and then looks for the solution that ships less units, ie. smaller package sizes. 

## TODOS 
- tests
- how to build
- how to use
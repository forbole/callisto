# Database
Since BDJuno relies on a PostgreSQL database in order to store the parsed data, one of the most important things is to create such database. To do this the first thing you need to do is install [PostgreSQL](https://www.postgresql.org/). 

Once installed you need to create a new database, and a new user that is going to read and write data inside it.  
Then, once that's one, you need to run the SQL queries that you can find inside the [`database/schema` folder](../database/schema).  

Once that's done, you are ready to [continue the setup](setup.md).
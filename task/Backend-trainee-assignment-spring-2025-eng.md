# Test assignment for a Backend trainee AVITO (spring wave 2025)

## Service for working with PVZ

New goods that were ordered through Avito are delivered to the PVZ several times a day. Before giving them to the customer, you must first check and enter the information into the database. Due to the fact that there are many PVZS and even more goods, it is necessary to implement a mechanism that allows, in the context of each PVZ, to see how many times a day goods arrived at them for acceptance and which goods were received.

## Task description

Develop a backend service for PVZ employees that will allow you to enter information on orders as part of the acceptance of goods.

1. User authorization:
* Using the /dummyLogin handle and passing the desired user type (employee, moderator) to it,
  In response, the service will return a token with the appropriate access level — a regular user or moderator.
  This token must be transferred to all endpoints that require authorization.

2. Registration and authorization of users by mail and password:
    * endpoint/register is used during registration.
      A new user of the desired type is created and saved in the database: a regular user (employee) or moderator (moderator).
      The created user gets the endpoint /login token.
      Upon successful authorization by mail and password, a token is returned for the user with the appropriate access level.

3. Setting up a PVZ:
* Only a user with the role of "moderator" can set up a PVZ in the system.
    * If the request is successful, the full information about the created IDP is returned. The establishment of a health insurance company is possible only in three cities: Moscow, St. Petersburg and Kazan. In other cities, it is impossible to start an air defense system at first, in which case the error must be returned.
    * A new record in the data warehouse should be the result of adding the PVZ.

4. Adding information about the acceptance of goods:
* Only an authorized user of the system with the role of "PVZ employee" can initiate the acceptance of goods.
    * A new record in the data warehouse should be the result of product acceptance initiation.
    * If the previous acceptance of goods has not been closed, then the operation to create a new acceptance of goods is not possible.

5. Adding goods within a single acceptance:
* Only an authorized user of the system with the role of "PVZ employee" can add goods after his inspection.
    * At the same time, the product must be linked to the last unclosed receipt of goods within the framework of the current PVZ.
    * If there is no new unclosed acceptance of goods, then an error should be returned, and the product should not be added to the system.
    * If the last acceptance of the product has not yet been closed, the result should be that the product is linked to the current PVZ and the current acceptance, with the subsequent addition of data to the repository.

6. Deletion of goods within the framework of non-closed acceptance:
* Only an authorized user of the system with the role of "PVZ employee" can delete goods that have been added
  as part of the current acceptance for the PVZ.
    * The removal of goods is possible only before the acceptance is closed, after that it is no longer possible to change the composition of the goods that
      they were accepted to the PVZ.
    * Items are deleted according to the LIFO principle, i.e. it is possible to delete items only in the order in which
      they were added as part of the current acceptance.

7. Closing the acceptance:
* Only an authorized user of the system with the role of "PVZ employee" can close the acceptance of goods.
    * If the acceptance of goods has already been closed (or there has been no acceptance of goods in this PVZ yet),
      then an error should be returned.
    * In all other cases, it is necessary to update the data in the storage and register the goods.,
      which were within the scope of this acceptance.

8. Receiving data:
    * Only an authorized user of the system with the role of "PVZ employee" or "moderator" can receive this data.
    * It is necessary to get a list of PVZ and all the information on them using pagination.
    * At the same time, add a filter by the date of acceptance of goods, i.e. display only those PVZ and all information on them that are within the specified time range
      We conducted product receptions.

## General introductory

The entity "Order Acceptance Point (PVZ)" has:
* Unique identifier
* Date of registration in the system
* City

The item "Acceptance of goods" has:
* Unique identifier
* Date and time of acceptance
* The PVZ in which the acceptance was carried out
* Products that were accepted as part of this acceptance
* Status (in_progress, close)

The "Product" entity has:
* Unique identifier
* Date and time of receipt of the goods (the date and time when the goods were added to the system as part of the acceptance of goods)
* Type (we work with three types of goods: electronics, clothing, shoes)

## Conditions

1. Use this [API](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2025/swagger.yaml).
2. Implement all the requirements specified in the terms of the assignment.
3. The server must be running on port 8080.
4. The implementation of user authorizations is not a prerequisite. In this case, the authorization token can be obtained from the /dummyLogin method.

   In the request parameters, you can select the user role: moderator or regular user.
   Depending on the role, a token with a certain access level will be generated.
6. Non—functional requirements:
* RPS - 1000
    * SLI response time — 100 ms
    * The response success rate SLI is 99.99%
7. The code must be covered by unit tests. The test coverage is at least 75%.
8. One integration test should be developed, which:
* First of all creates a new PVZ
    * Adds new order acceptance
    * Adds 50 products as part of the current order acceptance
    * Closes order acceptance

## Additional tasks

They are optional, but they will give you an advantage over other candidates.

1. Implement user authorization using the /register and /login methods
   (while the /dummyLogin method must still be implemented)
2. Implement a gRPC method that will simply return all the PVZ added to the system. For him, it's not
   authorization verification and validation of user roles are required. The server for gRPC must be running on port 3000.
   Please note that in the file [pvz.proto](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-spring-2025/pvz.proto ) it is necessary to register go_package for your project structure
3. Add prometheus to the project and collect the following metrics:
    * Technical:
* Number of requests
    * Response time
    * Business:
        * The number of PVZ created
        * The number of order acceptances created
        * The number of added items
          The server for prometheus must be raised on port 9000 and send data via the handle /metrics.
4. Set up logging in the project
5. Configure the DTO endpoint code generation according to the openapi scheme

## Stack requirements

* **The language of the service:** Go is preferred, but the following languages are also acceptable: PHP, Java, Python, C#.
* **Database:** preferably PostgreSQL, but you can choose another one that is convenient for you. You can't use an ORM to interact with the database.
* It is acceptable to use **builders for queries**, for example, like this: https://github.com/Masterminds/squirrel
* To deploy dependencies and the service itself, you need to use Docker or Docker & DockerCompose

## Additions to the solution
* If you have any questions that are not answered in the terms, then you can make decisions on your own.
* In this case, attach a Readme file with a list of questions and explanations of your decisions to the project.

## Making and sending the solution

Create a public git repository on any host (GitHub, GitLab, and others) containing the master/main branch:

1. Service code
2. Docker or Docker & DockerCompose or described in Readme.md launch instructions
3. Described in Readme.md questions or problems that you have encountered and a description of your solutions
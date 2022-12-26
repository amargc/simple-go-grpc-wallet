# simple-go-grpc-wallet

### Objective
As a test you are required to design and create a user wallet service from scratch, which is to record customer balance, and all transactions occurred through the system.

### Requirements
1. Design a data structure for storing customer ledgers, balance details etc.
2. By making use of the created data structure, create a gRPC service to support below features
  a. Create user wallet – assuming user profile is created already in other systems and a unique user id (uuid) is created
  b. Record keeping for user's transaction (credit to or debit from user's wallet) in ledgers
  c. Retrieve user wallet summary 
  d. Get user transaction history (pagination required)
3. The service should also support
  a. Multiple currency (e.g ETC, BTC, USD etc.)
  b. Multiple users
4. The service should have proper unit testing coverage
5. Proper logging in the service
6. Documentation
7. Anything you think it can make the service better

* Please try not to over complicate the task.
* You can use any third-party framework or libraries that you think it can help on you task.	
* Please provide a repository URL (e.g Github) to send it after completion.


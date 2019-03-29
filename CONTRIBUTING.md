## Design decisions
Before contributing, please read through the following design decisions:

 1. *Target architecture:* Docker Container run in our Kubernetes cluster
 2. *Trigger* The application should generate output based on timed triggers (e.g., once a day) and respond to manual requests.
 3. *Outputs* The application should be able to generate with CSV and be able to integrate with SAP as closely as possible, but not fully automate the billing process.
 4. *Persistence* If possible, the application should not use any persistent volumes, but may use Amazon S3 for persistent storage
 5. *Authentication* The container will be authenticated with AWS roles in our cluster and should not store any secrets or credentials.

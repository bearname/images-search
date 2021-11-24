# Simple face recognizer
[Available at https://damp-lake-99927.herokuapp.com/](https://damp-lake-99927.herokuapp.com/)
[Available at https://thawing-sea-83431.herokuapp.com/](https://thawing-sea-83431.herokuapp.com/)
## Task
[https://freelance.habr.com/tasks/376077](https://freelance.habr.com/tasks/376077)

## Setting Credentials
Setting your credentials for use by the AWS SDK for Java can be done in a number of ways,  
but here are the recommended approaches:

Set credentials in the AWS credentials profile file on your local system, located at:

`~/.aws/credentials` on Linux, macOS, or Unix

`C:\Users\USERNAME\.aws\credentials` on Windows

This file should contain lines in the following format:

`+
[default]`  
`aws_access_key_id = your_access_key_id` 
`aws_secret_access_key = your_secret_access_key`
## First building run command
`make download`
## Build
`make build`
## Run
`make run`

Open http://localhost:3000
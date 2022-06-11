# Minecraft Server Backup Tool



## Description
This is a tool that will watch the change of the backup folder of the map. When there has a new backup of the map in the backup folder, this tool upload the backup file to the AWS S3 service.

## Installation
- Clone this project first
- Excute `go mod download` command under the root folder of this project

## Usage
This tool will fetch settings from environment variable.

You need to create a .env file in project root for environment variable settings.

The .env file should look like as below:
```
WATCH_PATH=<Fill in the target path>
AWS_REGION=<Fill in the aws region>
AWS_ACCESS_KEY_ID=<Fill in the aws access key ID>
AWS_SECRET_ACCESS_KEY=<Fill in the aws access key>
S3_BUCKET_NAME=<Fill in the aws S3 bucket name>
```

Then, you can use `go build cmd\minecraft-server-backup-tool\main.go` to build the executable and run it.

Make sure the .env file and the assets folder are in the same directory as the executable you just build.

## Co-op Workflow
You shold follow the step below while you are working on this project:
1. Get the tasks you want to contribute from jira
2. Open an issue for it
3. Start working. At the same time, remember to update your status of your work on jira
4. When you are finish, open an merge request for it
5. Notify your teammates that you are already done with the task
6. (Optional) Let the reviewer of the MR review the code
7. Let the assignee of the MR review the code
8. If it all looks good, the assignee can merge the code into the branch

## Authors
The following are the developers who mainly maintain this project.
- [GUAN-YU CHEN](https://gitlab.guanyu.dev/ares30841167)
- [TING-HSUN SHEN](https://gitlab.guanyu.dev/ls3165)

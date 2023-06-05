# Minecraft Server Backup Tool



## Description
This is a tool that will watch the change of the backup folder of the map. When there has a new backup of the map in the backup folder, this tool upload the backup file to the AWS S3 service.

## Usage
This tool will fetch settings from environment variable.

You need to create a .env file in project root for environment variable settings.

The .env file should look like as below:
```
BACKUP_FOLDER_PATH=<Fill in the path of the folder where the backup file in>
BACKUP_FILE_REGEXP=<Fill in the regular expression of the target backup file>
AWS_REGION=<Fill in the aws region>
AWS_ACCESS_KEY_ID=<Fill in the aws access key ID>
AWS_SECRET_ACCESS_KEY=<Fill in the aws access key>
S3_BUCKET_NAME=<Fill in the aws S3 bucket name>
```

Then, you can use `go build cmd\minecraft-server-backup-tool\main.go` to build the executable and run it.

Make sure the .env file and the assets folder are in the same directory as the executable you just build.

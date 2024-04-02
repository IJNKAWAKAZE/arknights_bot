# Development Documentation

Environment requirements:

- [Go lang](https://go.dev) 1.10+
- [MySQL](https://www.mysql.com/) 5.7+ ([MariaDB](https://mariadb.org) is also supported)
- [Redis](https://redis.io) 5.0+ ([Redis stack](https://redis.io/docs/about/about-stack/) is also supported)

## Build from Source

If you need to update your default path for Go, you can use the following command:

```shell
export GOPATH="${workspaceFolder}/src/.go/"
```

You might need to write this command in your shell script or shell configuration file.

Then you can build the project:

```shell
git clone https://github.com/IJNKAWAKAZE/arknights_bot
cd arknights_bot/src/
go build -v
```

parameter `-o` can be used to specify the output file name. It will affect the last running step.

## Run the project

1. Copy the config file:

   ```shell
   cp ./arknights.example.yaml ./arknights.yaml
   ```

   Edit arknights.yaml, modify the configuration to your own.

2. Import the database schema(arknights.sql) to your SQL database.

   ```sql
   source ./arknights.sql
   ```

4. Start the bot

   ```shell
   ./src/arknights_bot
   ```

If things gonna all right, you will see the output like this:

```log
2024-04-02T10:34:38.015+08:00 INF 2024/04/02 10:34:38 init_db.go:25: 数据库连接成功
2024-04-02T10:34:38.017+08:00 INF 2024/04/02 10:34:38 init_redis.go:19: redis连接成功
2024-04-02T10:34:38.982+08:00 INF 2024/04/02 10:34:38 init_bot.go:22: 机器人初始化完成
2024-04-02T10:34:38.983+08:00 INF 2024/04/02 10:34:38 runner.go:48: 定时任务已启动
2024-04-02T10:34:38.984+08:00 INF 2024/04/02 10:34:38 bot.go:21: 机器人启动成功
```

## Local development

### Visual Studio Code

1. Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) for Visual Studio Code.

2. Install requirements:

    ```shell
    go install -v golang.org/x/tools/gopls@latest
    ```

    This help you connect to go language server, which helps you to check your cgo code grammar and format your code.
    
    ```shell
    go install -v github.com/go-delve/delve/cmd/dlv@latest
    ```

    This help you debug your go code.

3. Open the project in Visual Studio Code.

4. Config the `.vscode/launch.json` file.

    ```json
    {
        // Use IntelliSense to learn about possible attributes.
        // Hover to view descriptions of existing attributes.
        // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
        "version": "0.2.0",
        "configurations": [
            {
                "name": "Build and Run Package",
                "type": "go",
                "request": "launch",
                "mode": "exec",
                "cwd": "${workspaceFolder}",
                "program": "${workspaceFolder}/src/arknights_bot",
                "args": [],
                "env": {
                    "GOPATH": "${workspaceFolder}/src/.go/"
                }
            }
        ]
    }
    ```

5. Config the `.vscode/tasks.json` file.

    ```json
    {
        "version": "2.0.0",
        "tasks": [
            {
                "label": "Build Package",
                "type": "shell",
                "command": "go build -v",
                "group": {
                    "kind": "build",
                    "isDefault": true
                },
                "options": {
                    "cwd": "${workspaceFolder}/src"
                },
                "problemMatcher": []
            }
        ],
        "options": {
            "env": {
                "GOPATH": "${workspaceFolder}/src/.go/"
            }
        }
    }
    ```

6. Run the project in debug mode. Breakpoints can be set in the code.

## GoLand

1. Install [GoLand](https://www.jetbrains.com/go/).
2. Open the project in GoLand.
3. Run the project.

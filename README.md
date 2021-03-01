[![Build Status](https://travis-ci.com/fr123k/confluence-slackbot.svg?branch=main)](https://travis-ci.com/fr123k/confluence-slackbot)

# confluence-slackbot

> **ⓘ Note: The natural language processing works only for english at the moment.**
> [Natural Language Processing](#natural-language-processing)

See this bot in action.

![demo](./image/demo.gif)

# Configuration

## Slack Application configuration

### OAuth
*This slack application has to installed as a legacy application because it needs
the bot scope in order to use the RTM (Real Time Messages) feature of slack.*

The following slack OAuth scopes are needed.
* bot
* channels:read
![slack scope](./image/slack_scopes.png)

### Interactive Buttons

Following this [Preparing your Request URL](https://api.slack.com/legacy/interactive-messages#making-messages-interactive__readying-your-application-for-interactive-messages) guide of slack.

## Confluence Bot Application Server
The configuration can be passed as an yaml file and/or environment variables.
For duplicated values the following value overwrite each other.
* config files
* environment variables
  
Example yaml configuration.
```
  confluence:
    url: https://xxxxxxx.atlassian.net/wiki
    username: jane.doe@example.com
    token: jira-api-token

  slack:
    token: xoxb-slack-bot-token

  server:
    actionurl: /actions
    port: 3000

  debug: true
```
Environment Variables
```
  CONFLUENCE_URL=https://xxxxxxx.atlassian.net/wiki
  CONFLUENCE_USERNAME=jane.doe@example.com
  CONFLUENCE_TOKEN=jira-api-token
  SLACK_TOKEN=xoxb-slack-bot-token
  SERVER_ACTION_URL,ACTION_URL=/actions
  SERVER_PORT,PORT=3000
  SERVER_DEBUG,DEBUG=true
```

## Local Development

Prerequisites:
* slack workspace to add this confluence slack bot application
* for interactive buttons a public reachable instance

### Interactive Buttons

Start the localhost server with the following command
```
make run
```
Have it public reachable for slack.

#### ngrok

Follow this wonderful instraction to setup a public reachable localhost instance with [ngrok](https://ngrok.com/).
[tunneling-with-ngrok](https://api.slack.com/tutorials/tunneling-with-ngrok)

Start tunnel with the following command. (3000 port is the default server port for this app)
```
ngrok http 3000
```

## Targets

### Build

The following command will build the golang binary and run the unit tests.
The result of this build step is the standalone binary in the `./build/` folder. 

```
make build
```

Example Output:
```
  go build -o build/main cmd/main.go
  go test -v --cover ./...
  ?       github.com/fr123k/confluence-slackbot/cmd       [no test files]
  ?       github.com/fr123k/confluence-slackbot/pkg/confluence-cli        [no test files]
```

### Run

The following make target will first build and then execute the golang binary.
```
  make run
```

Example Output:
```
  go build -o build/main cmd/main.go
  go test -v --cover ./...
  ?       github.com/fr123k/confluence-slackbot/cmd       [no test files]
  ?       github.com/fr123k/confluence-slackbot/pkg/confluence-cli        [no test files]
  ./build/main
  Event Received: Event Received:
```

### Clean

The following make target will remove the `./build/` folder.
**No confirmation needed**
```
  make clean
```

Example Output:
```
  rm -rfv ./build
  ./build/main
  ./build
```

# Natural Language Processing

The slack message language processing is using the following github repository [jdkato/prose](https://github.com/jdkato/prose)

It provides a full natural language processing model for at the moment only *english* language.

# Changelog

* setup travis build
* add unit test for natural language processing

# Todos

* add unit test for confluence query building
* verify slack identity for the interactive button support. [verifying-requests-from-slack](https://api.slack.com/authentication/verifying-requests-from-slack)
* make the confluence query configurable
* configure on which messages to respond with search results (direct messages, non thread, thread messages, ...)
* add commands to perform confluence search results directly
* add infrastructure deployments examples for aws, hetzner, ...

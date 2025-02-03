# Gator CLI

Gator CLI is a command-line application designed to interact with RSS feeds. It allows users to add, follow, and browse feeds, among other functionalities.

## Prerequisites

Before you can run the Gator CLI, ensure you have the following installed on your system:

-  **PostgreSQL**: The application uses PostgreSQL as its database. You can download and install it from [PostgreSQL's official site](https://www.postgresql.org/download/).
-  **Go**: The application is written in Go. You can download and install it from [Go's official site](https://golang.org/dl/).

## Installation

To install the Gator CLI, use the `go install` command:

```bash
go install github.com/yourusername/gator@latest
```

Make sure your `GOPATH/bin` directory is in your system's `PATH` so you can run the `gator` command from anywhere.

## Configuration

Before running the application, you need to set up a configuration file. Create a `.gatorconfig.json` file in your home directory with the following content:

```json
{
  "DbUrl": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "CurrentUserName": "your-username"
}
```

-  Replace `username` and `password` with your PostgreSQL credentials.
-  Replace `your-username` with the username you will use to log in to the application.

## Running the Program

Once installed and configured, you can run the Gator CLI using various commands. Here are a few examples:

-  **Register a User**: Register a new user with the application.
  ```bash
  gator register <username>
  ```

-  **Add a Feed**: Add a new RSS feed to your account.
  ```bash
  gator addfeed "Feed Name" "https://example.com/rss"
  ```

-  **Follow a Feed**: Follow an existing feed by URL.
  ```bash
  gator follow "https:// example.com/rss"
  ```

-  **Browse Feeds**: View posts from feeds you are following, with an optional limit on the number of posts.
  ```bash
  gator browse [limit]
  ```

-  **Aggregate Feeds**: Continuously fetch and print posts from your feeds.
  ```bash
  gator agg <time_between_reqs>
  ```

## Commands Overview

-  **register**: Register a new user with the application.
-  **addfeed**: Add a new RSS feed to your account.
-  **follow**: Follow an existing feed by URL.
-  **browse**: View posts from feeds you are following.
-  **agg**: Continuously fetch and print posts from your feeds.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

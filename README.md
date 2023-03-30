# Wholesome Living

Backend for Wholesome Living

## Prerequisites

-   go
-   task
-   Wholesome Living-backend (repo)

#### 1. Step: Clone the repository

```bash
git clone git@github.com:Wholesome Living-social/Wholesome Living-backend.git
```

#### 2. Step: Install task

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

## Dependencies

To install the dependencies, please run the following command:

```bash
task install
```

## Docs

For creating the initial documentation.

```bash
task docs
```

## Starting the Server

### Dev environment (hot-reloading)

```bash
task dev
```

### Production environment

#### 1. Step: Build the server

```bash
task build
```

#### 2. Step: Start the server

```bash
task start
```

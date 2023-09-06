# Wholesome Living

Backend for Wholesome Living

### Prerequisites

-   go
-   task
-   Wholesome-Living-backend (repo)

1.Step: **Clone the repository**

```bash
git clone https://github.com/Wholesome-Living-Project/wholesome-living-backend.git 
```

2.Step: **Install task**

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### Dependencies and Docs

To install the dependencies, please run the following command:

```bash
task install
```

For creating the initial documentation:
```bash
task docs
```
---

## Starting the Server

### Dev environment (hot-reloading)

```bash
task dev
```

### Production environment

 1. Step: Build the server

```bash
task build
```
2. Step: Start the server

```bash
task start
```

---

## Testing

Running all tests:
```bash
task test
```
Check for coverage:

```bash
task testcov
```

with prettier output:
```bash
task guicov
```

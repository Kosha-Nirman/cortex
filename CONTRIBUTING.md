# Contributing to Cortex

Thank you for considering contributing to **Cortex – Subdomain Discovery & Reconnaissance Tool**! 🚀
Your contributions help improve the project and make it more powerful for security researchers and developers. Please follow the guidelines below to contribute effectively.

---

## 📌 Getting Started

1. **Fork the Repository**
   Click the **Fork** button on the top-right corner of the [repository](https://github.com/rivetron/cortex) to create your own copy.

2. **Clone the Repository**

   ```bash
   git clone https://github.com/<your-username>/cortex.git
   cd cortex
   ```

3. **Create a New Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Install Dependencies**

   ```bash
   go mod tidy
   ```

---

### Project Structure 📂

```bash
cortex/
├── scripts/                  # Install/uninstall scripts
│   ├── install.sh
│   ├── uninstall.sh
├── src/                      # Source code
│   ├── cmd/                  # CLI commands
│   │   ├── root.go
│   │   ├── version.go
│   ├── config/               # Configuration handling
│   │   └── config.go
│   ├── domain/               # Core domain logic
│   │   ├── subdomain.go
│   │   └── subdomain_test.go
│   ├── orchestrator/         # Resolver orchestration
│   │   └── orchestrator.go
│   ├── report/               # Report generation
│   │   └── markdown.go
│   ├── resolver/             # Subdomain resolvers
│   │   ├── brute.go
│   │   ├── certificate.go
│   │   ├── dns.go
│   │   ├── dns_test.go
│   │   └── passive.go
│   ├── utils/                # Utility helpers
│   │   ├── banner.go
│   │   ├── downloads.go
│   │   ├── downloads_test.go
│   │   └── result.go
│   └── main.go               # Application entry point
├── go.mod                    # Go module definition
├── go.sum                    # Go dependencies lockfile
├── LICENSE.md                # Project license
├── README.md                 # Documentation
└── Makefile                  # Build/test automation
```

---

## 🚀 Making Changes

1. **Implement Your Changes**

   * Follow best practices and keep the code consistent.
   * Add comments where necessary for clarity.
   * Ensure new resolvers or report features are modular.

2. **Run and Test the Application**

   ```bash
   go run src/main.go
   ```

   Example:

   ```bash
   go run src/main.go example.com
   ```

3. **Build the Project**

   ```bash
   go build -o cortex ./src
   ```

4. **Run the CLI Tool**

   ```bash
   ./cortex example.com
   ```

5. **Commit Your Changes**

   ```bash
   git add .
   git commit -m "feat: add DNS resolver with timeout support"
   ```

6. **Push the Changes**

   ```bash
   git push origin feature/your-feature-name
   ```

---

## ✅ Submitting a Pull Request

1. Navigate to the original repository: [Cortex](https://github.com/rivetron/cortex).
2. Click on the **New Pull Request** button.
3. Select your fork and branch, and compare it with the `main` branch.
4. Add a meaningful **title and description** for your changes.
5. Submit the pull request and wait for review.

---

## 🛠 Code Guidelines

* Write clear, descriptive commit messages.

* Follow Go best practices and existing project structure.

* Use **table-driven tests** for functions where applicable.

* Run linters before pushing changes:

  ```bash
  golangci-lint run
  ```

* Run all tests to ensure nothing breaks:

  ```bash
  go test ./...
  ```

---

## 💬 Need Help?

If you have any questions or need clarification, feel free to:

* Open an issue in the [Issues tab](https://github.com/rivetron/cortex/issues)
* Start a discussion in your pull request

Happy hacking! 🎯

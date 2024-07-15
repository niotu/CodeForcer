# ChooseTwoOption

---

## Table of Contents

---

- [Used Stack](#used-stack)
   - [Backend](#backend)
   - [Frontend](#frontend)
   - [Deployment](#deployment)
   - [Summary](#summary)
- [Prerequisites](#prerequisites)
   - [Installing Docker](#installing-docker)
      - [On Windows](#on-windows)
      - [On macOS](#on-macos)
      - [On Linux](#on-linux)
- [Quick Start](#quick-start)

## Used stack

---

This project utilizes the following technologies:

### Backend

- **Go**: The backend is built using the Go programming language. We utilized the standard `net/http` package to implement the HTTP server.

### Frontend

- **React**: The frontend is developed using React, a popular JavaScript library for building user interfaces.

### Deployment

- **Docker**: A `Dockerfile` is provided to facilitate easy deployment of the application. Docker allows for consistent and reproducible environments, ensuring that the application runs smoothly across different systems.

### Summary

- **Backend Language**: Go
- **Backend Framework**: `net/http` (standard library)
- **Frontend Library**: React
- **Deployment**: Docker

## Prerequisites

---

Make sure you have Docker installed on your machine. If Docker is not installed, follow these steps:

### Installing Docker


#### On Windows:

1. Download and install Docker Desktop from [Docker Hub](https://www.docker.com/products/docker-desktop).
2. Follow the installation instructions.
3. After installation, Docker Desktop should start automatically. If not, start it manually.

#### On macOS:

1. Download and install Docker Desktop from [Docker Hub](https://www.docker.com/products/docker-desktop).
2. Follow the installation instructions.
3. After installation, Docker Desktop should start automatically. If not, start it manually.

#### On Linux:

1. Update your package database.

    ```bash
    sudo apt-get update
    ```

2. Install Docker.

    ```bash
    sudo apt-get install -y docker.io
    ```

3. Start Docker.

    ```bash
    sudo systemctl start docker
    ```

4. Enable Docker to start at boot.

    ```bash
    sudo systemctl enable docker
    ```

## Quick start

---

To start the whole application you need to build and run the docker image.
Follow these steps to quickly get the project up and running.


Clone the repository to your local machine using the following command:

```bash
git clone https://gitlab.pg.innopolis.university/n.solomennikov/choosetwooption.git
```

Move to cloned project directory and build docker image.

```bash
cd choosetwooption
docker build -t codeforcer .
docker run -d -p 8080:8080 -p 3000:3000 codeforcer
```

Now the web application runs on http://localhost:3000.

# Notes API Tester 📝

A lightweight frontend interface built to test authentication endpoints and core functionality for a work-in-progress Notes API. 

## 🚀 Overview

This project is an experiment created to stress-test token-based authentication workflows. It provides a visual interface to verify how the backend handles secure note creation and token based authentication.

## 🛠️ Tech Stack

* **Frontend:** [Vite/React]
* **Styling:** [CSS]
* **HTTP Client:** [Fetch API]
* **HTTP Server:** [Gin]
* **Backend:** [Go]
* **Tokens:** [JWT]
* **Database:** [sqlc/Postgresql]

## 📦 Getting Started

### Prerequisites
* Go (version 1.26)
* Docker Desktop [typical docker install steps](https://www.docker.com/get-started/)
* A running instance of the [Notes API](https://github.com/meads/firstly-api)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/meads/firstly-api.git
    cd firstly-api
    ```

2. Configure environment variables:
   Create a `.env` file in the root directory:
    ```env
    # api database connection string
    DATABASE_URL=postgresql://username:password@db:5432/databasename?sslmode=disable

    # used for jwt signing
    SECRET=

    # used for configuring gin router cors middleware
    ALLOW_ORIGINS=

    # variables for db service
    POSTGRES_USER=
    POSTGRES_PASSWORD=
    POSTGRES_DB=
    ```

3. Build and start Docker containers:
   ```bash
   docker compose build 
   docker compose up
   ```


### Running the App

Open `http://localhost:3000` in your browser to view the UI.

## 🧪 How to Test the Auth Flow

Use the curl commands to inspect the http api or launch the web application and witness the
token expiration. 
~~1. **Register/Login:** Create an account to receive your initial access and refresh tokens.~~
~~2. **Inspect Tokens:** Check the UI dashboard to see active token strings and lifetimes.~~
~~3. **Simulate Expiration:** Shorten token lifespans on your API to witness the UI use the refresh token automatically.~~

## 🚧 Project Status

**Current Status:** Active Experiment / Work in Progress
* This is a utility sandbox, not a production-ready application.
* Features may change frequently as the backend API evolves.

## 📄 License

This project is open-source and available under the **MIT License**.

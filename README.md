# Book It - Web Application for Book Exchange, Reading Progress Tracking, and Recommendations

## Project Description

The "Book It" project is a web application developed as part of an engineering thesis. The aim of the application is to enable users to exchange books, track their reading progress, and receive recommendations for new books to read.

The application is designed for a community of book enthusiasts who want to easily organize their library, find new readers for book exchanges, and effectively plan their reading goals.

## Technologies

The project was developed using:

- Golang
- Gorm
- HTMX
- SQLite
- SSR (Server-Side Rendering) Architecture

## Features

- Create a virtual book library
- Schedule your reading
- Book exchange matching system
- Book recommendations based on preferences
- Geolocation support for better exchange matching

## Running the Project

### Requirements

- Installed [Golang](https://go.dev/dl/)
- An account and configured project on [Google Cloud Platform](https://console.cloud.google.com/) for login and registration
- Configured environment variables in the `.env` file:

```
PORT=3000
APP_ENV=local
DB_URL=./test.db

GOOGLE_CLIENT_ID="example-google-client-id"
GOOGLE_CLIENT_SECRET="example-google-secret"
GOOGLE_API_KEY="example-google-api-key"
GEOAPIFY_KEY="example-geoapify-key"
```

### Installation

Clone the repository to your local environment:

```bash
git clone https://github.com/FilipBudzynski/book-it
cd book_it
```

### Running the Application

To run the application locally:

```bash
make build
make run
```

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

# Template System

A robust template system that allows users to upload, generate, and download document templates. This project is built using Go, Gin for the web framework, and integrates with Amazon S3 for file storage.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Endpoints](#endpoints)
- [Configuration](#configuration)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

- User Authentication (Register, Login)
- Upload DOCX Templates
- Generate Documents from Templates
- Download Generated Documents
- Password Reset Functionality

## Getting Started

### Prerequisites

- Go 1.16+ installed on your machine
- PostgreSQL database
- AWS S3 bucket for storing templates and documents

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/geekychandan/template-system.git
    cd template-system
    ```

2. Set up your environment variables. Create a `.env` file with the following content:

    ```plaintext
    DB_HOST=localhost
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    DB_PORT=5432

    S3_REGION=your_aws_region
    S3_BUCKET=your_s3_bucket
    AWS_ACCESS_KEY_ID=your_aws_access_key_id
    AWS_SECRET_ACCESS_KEY=your_aws_secret_access_key

    JWT_SECRET=your_jwt_secret

    EMAIL_ADDRESS=your_email@example.com
    EMAIL_PASSWORD=your_email_password
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    ```

3. Install dependencies:

    ```sh
    go mod tidy
    ```

4. Run the application:

    ```sh
    go run main.go
    ```

## Usage

### Endpoints

#### Authentication

- **Register**
  - `POST api/auth/register`
  - Request Body:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```

- **Login**
  - `POST api/auth/login`
  - Request Body:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```

- **Request Password Reset**
  - `POST api/auth/reset-password`
  - Request Body:
    ```json
    {
      "email": "user@example.com"
    }
    ```

- **Reset Password**
  - `POST api/auth/reset-password/confirm`
  - Request Body:
    ```json
    {
      "token": "reset_token",
      "new_password": "newpassword123"
    }
    ```

#### Templates

- **Upload Template**
  - `POST api/templates/upload`
  - Headers: `Authorization: Bearer <your_jwt_token>`
  - Form-data: `template` (File)

- **Get User Templates**
  - `GET api/templates`
  - Headers: `Authorization: Bearer <your_jwt_token>`

- **Get Placeholders**
  - `GET api/templates/:id/placeholders`
  - Headers: `Authorization: Bearer <your_jwt_token>`

#### Documents

- **Generate Document**
  - `POST api/documents/:id/generate`
  - Headers: `Authorization: Bearer <your_jwt_token>`
  - Request Body:
    ```json
    {
      "placeholder1": "value1",
      "placeholder2": "value2"
    }
    ```

- **Get User Documents**
  - `GET api/documents`
  - Headers: `Authorization: Bearer <your_jwt_token>`

- **Download Document**
  - `GET api/documents/:id/download`
  - Headers: `Authorization: Bearer <your_jwt_token>`

## Configuration

Configuration is done using environment variables. Ensure that the `.env` file is created with the necessary values as mentioned in the [Installation](#installation) section.

## Project Structure

```plaintext
template-system/
├── controllers
│   ├── authController.go
│   ├── documentController.go
│   └── templateController.go
├── middleware
│   └── authMiddleware.go
├── models
│   └── models.go
├── routes
│   ├── authRoutes.go
│   ├── documentRoutes.go
│   └── templateRoutes.go
├── services
│   ├── authService.go
│   ├── documentService.go
│   └── templateService.go
├── utils
│   ├── database.go
│   ├── email.go
│   └── s3.go
│   └── cache.go
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md

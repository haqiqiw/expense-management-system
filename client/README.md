# Frontend

This directory contains frontend for Expense Management System.

## Prerequisites

Before you begin, make sure you have the following installed on your machine:

- [Node.js](https://nodejs.org/id) (version 20.19.4 or later)

## Development Setup

The following commands should be run from the `client/` directory.

### Configure Environment Variables

Copy the sample environment file to create your local configuration.

```bash
cp env.sample .env.local
```

Adjust the variables in `.env.local` if needed.

### Install Dependencies

Install the required depedencies:

```bash
npm install
```

### Start the Development Server

Run the application in development mode with hot-reloading:

```bash
npm run dev
```

### Testing

To run the test:

```bash
npm run test
```

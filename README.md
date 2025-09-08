# Expense Management System

## Overview

An expense management system where employee can submit expenses, manager can approve them, and approved expenses get processed for payment.

## Getting Started

### Prerequisites

Before running the project locally, make sure you have the following installed:

- [Git](https://git-scm.com/)
- [Go](https://go.dev/) (version 1.24.3 or later)
- [Node.js](https://nodejs.org/) (version 20.19.4 or later)
- [NPM](https://nodejs.org/) (version 10.8.2 or later)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [golang-migrate](https://github.com/golang-migrate/migrate) – follow installation instructions: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Optional tools for development and debugging:

- **PostgreSQL client** (`psql`) or GUI tools (e.g., [DBeaver](https://dbeaver.io/)) for inspecting database, running queries, and checking migrations
- **Redis client** (`redis-cli`) – for inspecting keys and debugging

### Installation & Running the App

Clone the repository:

```bash
git clone https://github.com/haqiqiw/expense-management-system.git

cd expense-management-system
```

Start the entire App and dependencies using Docker Compose:

```bash
docker-compose up --build -d
```

> Ensure all containers are running and healthy, before continuing to the next step

### Database Migration

Go to `server/` directory:

```bash
cd server/
```

Copy the sample environment file:

```bash
cp env.sample .env
```

Run database migration:

```bash
make migrate-up
```

Run database seeding:

```bash
go run dev/seeder/main.go
```

> For simplicity in this demo, each seeded user’s password is set to match the prefix of their email address (the part before the @).
>
> Example: `john@mail.com` → `john`

### Access the App:

- Web: http://localhost:5173/
- API: http://localhost:8500/
- Swagger UI: http://localhost:8500/swagger/

<details>
<summary>Web preview</summary>

![Expense List](/docs/expense-sample.png)

</details>

## Development

For instructions on how to run the frontend and backend services locally for development, please see the `README.md` files in the [`client`](https://github.com/haqiqiw/expense-management-system/blob/main/client/README.md) and [`server`](https://github.com/haqiqiw/expense-management-system/blob/main/server/README.md) directories respectively.

## Business Rules & Assumptions

### How Expense Approval Works

For this project, I kept things simple, any user with a `manager` role can approve or reject expenses. It doesn’t matter who submitted the expense, `employee` or `manager`, as long as the approver is a `manager`.

The only restriction is that managers cannot approve their own expenses, if a `manager` submits an expense, it must be approved by `another manager`. I skipped implementing a full reporting-line hierarchy due to time constraints.

### Changing or Rolling Back Expenses

Once an expense reaches a final state (`approved`, `rejected`, or `completed`), it can't be rolled back. So if a manager accidentally rejects an expense, the employee needs to create a new expense to get it approved.

### How Users Are Created

The requirements didn’t mention how users are registered. For now, we can add users directly into the database or use a script. I did build an endpoint for registration, but I didn’t integrate it into a UI because it wasn’t a core requirement.

### Receipt Upload

Receipt upload is mocked. If a user chooses an image, it’s temporarily stored on client side. When the expense is submitted, the receipt URL is set to a dummy image (https://placehold.co/500x700).

## Architecture Decisions

### Abstraction with Interfaces

The backend is built around interfaces to decouple business logic from the database. This allows for easy mocking of dependencies, making the codebase testable and easier to maintain.

### Database Constraints

While the backend validates all incoming data, `CHECK` (prevent negative amount) and `FOREIGN KEY` constraints are also used in the database. This provides a second layer of protection to guarantee data integrity.

### Unique Indexes

`UNIQUE` indexes are used to enforce critical business rules at the database level, such as preventing duplicate emails for users and ensuring only one approval record can be linked to an expense.

### Composite Indexes

A composite index was implemented on the expenses table for performance:

- An index on (`user_id`, `status`) is used for personal expense queries, `user_id` is placed first due to its higher cardinality, which filters the data more effectively. This also allows for efficient lookups by `user_id` alone
- A separate index on (`status`, `user_id`) is used for the manager's approval queue. This allows the database to first find all expenses with an `awaiting_approval` status before checking the user

### Asynchronous Processing with Kafka

Kafka is used for background payment processing to keep API requests fast and non-blocking. It was chosen for its ability to handle high throughput and reliably decouple the API from the payment worker.

### Distributed Locking with Redis

Redis is used for distributed locking when processing payments. This ensures that a single payment job is not processed by multiple workers at the same time. Redis was chosen for its atomic operations.

### Rate Limiting

Rate limiting is implemented at the app / backend level as a middleware. For a larger-scale system, this would ideally be handled by an API Gateway or WAF to avoid burdening the backend service.

## Things I Would Improve With More Time

### Approval Hierarchy

I’d make the approval flow more realistic by connecting employees to specific managers. That way, expenses could only be approved by the manager who the employee reports to, and managers’ expenses would need approval from their upper managers.

We could have a `user_reporting_lines` table that maps each employee to their manager and a manager to their upper manager.

### Payment Records for Audit

I’d record every payment attempt for each expense. This would give us proof of what was attempted, what the partner returned, and help troubleshoot errors when payments fail.

We could create an `expense_payment_attempts` table that records each payment attempt for an expense. It would store the partner ID returned by the partner API, the response, and the status. This way, we can track failures, and have an audit trail for all payment interactions.

### Notifications

I’d implement notifications to alert managers about submitted expenses. If an expense is submitted outside business hours, the notification could wait until business hours.

We could create an `expense_notifications` table that records each notification to be sent. It would include the recipient, related expense ID, and status (e.g., pending, sent, failed). A cron job would periodically check this table to send pending notifications

### Refresh Token Mechanism

Currently, when access tokens expire, users are forced to log in again. I’d implement a refresh token system so clients could get new tokens without logging out.

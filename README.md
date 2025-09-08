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

For this project, I kept things simple: any user with a manager role can approve or reject expenses. It doesn’t matter who submitted the expense, employee or manager, as long as the approver is a manager. The only restriction is that managers cannot approve their own expenses. I skipped implementing a full reporting-line hierarchy due to time constraints.

### Changing or Rolling Back Expenses

Once an expense reaches a final state (approved, rejected, or completed), it cannot be changed or rolled back. So if a manager accidentally rejects an expense, the employee needs to create a new expense to get it approved.

### How Users Are Created

The requirements didn’t mention how users are registered. For now, we can add users directly into the database or use a script. I did build an endpoint for registration, but I didn’t integrate it into a UI because it wasn’t a core requirement.

### Receipt Upload

Receipt upload is mocked. If a user chooses an image, it’s temporarily stored on client side. When the expense is submitted, the receipt URL is set to a dummy image (https://placehold.co/500x700).

## Things I Would Improve With More Time

### Approval Hierarchy

I’d make the approval flow more realistic by connecting employees to specific managers. That way, expenses could only be approved by the manager who the employee reports to, and managers’ expenses would need approval from their upper managers.

We could have a `user_reporting_lines` table that maps each employee to their manager and a manager to their upper manager.

### Payment Records for Audit

I’d record every payment attempt for each expense. This would give us proof of what was attempted, what the partner returned, and help troubleshoot errors when payments fail.

We could create an `expense_payment_attempts` table that records each payment attempt for an expense. It would store the partner ID returned by the API, the response, and the status. This way, we can track failures, and have an audit trail for all payment interactions.

### Notifications

I’d implement notifications to alert managers about submitted expenses. If an expense is submitted outside business hours, the notification could wait until business hours. I’d also allow employees to manually trigger reminder notifications to managers.

We could create an `expense_notifications` table that records each notification to be sent. It would include the recipient, related expense ID, and status (e.g., pending, sent, failed). A cron job would periodically check this table to send pending notifications

### Dead Letter Queue for Failed Events

I’d add a separate queue for events that fail, like invalid data or messages exceeding retry limits. Right now, failed events just stop retrying after a limit, which might lose important information.

### Refresh Token Mechanism

Currently, when access tokens expire, users are forced to log in again. I’d implement a refresh token system so clients could get new tokens without logging out.

### Monitoring

If I had more time, I would improve the monitoring and metrics in the system. For example, I would record detailed metrics for outgoing requests, such as latency and throughput, so we can track how the system interacts with external services.

I’d also add metrics for consumers, to monitor their processing throughput, success/failure status, which would help quickly detect bottlenecks or failures in message processing.

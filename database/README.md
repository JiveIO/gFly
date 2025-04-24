# Database

This directory contains database-related files, including migration scripts and database setup instructions.

## Purpose

The database directory is used for:
- Database migration scripts
- Database setup and configuration instructions
- Database schema definitions
- Database seeding scripts

## Structure

- **migrations/**: Contains database migration scripts
  - **mysql/**: Migration scripts for MySQL
  - **postgresql/**: Migration scripts for PostgreSQL
- **MySQL.md**: Instructions for setting up and configuring MySQL
- **PostgreSQL.md**: Instructions for setting up and configuring PostgreSQL

## Usage

### Migrations

Database migrations are managed using the `migrate` tool. Migration files follow the naming convention:

```
{version}_{description}.{up|down}.sql
```

- **up.sql**: Contains SQL statements for applying the migration
- **down.sql**: Contains SQL statements for reverting the migration

To run migrations:

```bash
# Apply all migrations
make migrate.up

# Revert last migration
make migrate.down

# Create a new migration
make migrate.create name=create_new_table
```

### Database Setup

For detailed instructions on setting up databases:

- See **MySQL.md** for MySQL setup instructions
- See **PostgreSQL.md** for PostgreSQL setup instructions

## Best Practices

- Always create both up and down migrations
- Keep migrations small and focused on a single change
- Test migrations before applying them to production
- Document any manual steps required for migrations
- Use descriptive names for migration files
- Maintain backward compatibility when possible
- Consider data migration needs, not just schema changes

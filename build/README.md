# Build Directory

This directory contains the compiled binaries and runtime files for the application.

## Purpose

The build directory is used for:
- Storing compiled application binaries
- Storing compiled CLI tool binaries
- Containing runtime configuration files
- Providing a clean, isolated environment for running the application

## Structure

After building the application, this directory will contain:
- **app**: The main application binary
- **artisan**: The CLI tool for running commands, schedules, and queues
- **.env**: Environment configuration file (copied from the project root)

## Usage

The build directory is populated by running the build command:

```bash
make build
```

The compiled binaries can be executed directly:

```bash
# Run the main application
./build/app

# Run a CLI command
./build/artisan cmd:run <command-name>

# Run scheduled tasks
./build/artisan schedule:run

# Run queue workers
./build/artisan queue:run
```

## Best Practices

- Do not commit compiled binaries to version control
- Always use the latest build for testing and deployment
- Keep the build directory clean by using `make clean` when needed
- Use the same build process for development and production to ensure consistency
- Consider using a CI/CD pipeline for automated builds

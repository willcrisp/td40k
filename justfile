# Counter App Justfile

# List all available commands (runs by default if you just type `just`)
default:
    @just --list

# Build containers from scratch and start them in the background
# We include --build here so it automatically captures any new changes you make!
up:
    docker compose up --build -d

# Build and run the stack in Development Mode (HMR for Frontend)
dev:
    docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build


# Stop all containers
down:
    docker compose down

# Follow the live logs for all services (Useful for observing our WebSocket prints!)
logs:
    docker compose logs -f

# Stop containers and delete the PostgreSQL volume
kill:
    @echo "☢️ Deleting PostgreSQL volume..."
    docker compose down -v

# Hard Reset: Destroy containers, WIPE THE DATABASE VOLUME, and rebuild
reset:
    @echo "☢️ Wiping database and rebuilding..."
    docker compose down -v
    docker compose up --build -d

## Relevant Files

- `Dockerfile` - Backend multi-stage dockerfile.
- `docker-compose.yml` - Docker compose configuration for the backend service.
- `.dockerignore` - Files to exclude from the docker context.
- `.env.example` - Example environment variables for the docker container.

### Notes

- Use multi-stage builds in the `Dockerfile` to separate the build environment from the runtime environment, reducing the final image size.
- Ensure only necessary artifacts are copied to the final stage.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [x]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Create and checkout a new branch for this feature (e.g., `git checkout -b feature/backend-docker`)
- [x] 1.0 Setup Backend Dockerfile with Multi-Stage Build
  - [x] 1.1 Create `Dockerfile` in the root backend directory
  - [x] 1.2 Define the `builder` stage: use appropriate base image, copy dependency files, and install dependencies
  - [x] 1.3 Copy source code and build the backend application in the `builder` stage
  - [x] 1.4 Define the `production` (final) stage: use a lightweight base image (e.g., alpine or slim)
  - [x] 1.5 Copy only the compiled artifacts and production dependencies from the `builder` stage to the `production` stage
  - [x] 1.6 Set the default command (`CMD` or `ENTRYPOINT`) to start the backend server
- [x] 2.0 Configure Docker Ignore
  - [x] 2.1 Create a `.dockerignore` file
  - [x] 2.2 Add unnecessary files to `.dockerignore` to optimize build context (e.g., `node_modules`, `.git`, `.env`)
- [x] 3.0 Configure Docker Compose for Backend
  - [x] 3.1 Create or update `docker-compose.yml`
  - [x] 3.2 Define the backend service, pointing `build.context` to the backend directory
  - [x] 3.3 Map necessary ports
  - [x] 3.4 Define volume mounts if needed for local development
  - [x] 3.5 Configure network and dependencies
- [x] 4.0 Setup Environment Configuration for Docker
  - [x] 4.1 Create a `.env.example` file for required environment variables
  - [x] 4.2 Map environment variables from `.env` to the backend container in `docker-compose.yml`
- [x] 5.0 Verify and Test Docker Image
  - [x] 5.1 Build the docker image locally using `docker build` or `docker compose build`
  - [x] 5.2 Run the container using `docker run` or `docker compose up`
  - [x] 5.3 Verify the backend API is accessible on the exposed port
  - [x] 5.4 Check container logs for any runtime errors

## Relevant Files

- `Dockerfile` - Backend multi-stage dockerfile.
- `docker-compose.yml` - Docker compose configuration for the backend service.
- `.dockerignore` - Files to exclude from the docker context.
- `.env.example` - Example environment variables for the docker container.

### Notes

- Use multi-stage builds in the `Dockerfile` to separate the build environment from the runtime environment, reducing the final image size.
- Ensure only necessary artifacts are copied to the final stage.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [ ] 0.0 Create feature branch
- [ ] 1.0 Setup Backend Dockerfile with Multi-Stage Build
- [ ] 2.0 Configure Docker Ignore
- [ ] 3.0 Configure Docker Compose for Backend
- [ ] 4.0 Setup Environment Configuration for Docker
- [ ] 5.0 Verify and Test Docker Image

# easymirror-backend
Source code to EasyMirror's backend

## Docker containers
- This repo is containerized using Docker.
- All containers start with the name `easymirror-backend`
- It follows advice from the following sources:
    - https://www.youtube.com/watch?v=C5y-14YFs_8
    - https://www.youtube.com/watch?v=vIfS9bZVBaw
    - https://www.docker.com/blog/developing-go-apps-docker/
    - https://laurent-bel.medium.com/running-go-on-docker-comparing-debian-vs-alpine-vs-distroless-vs-busybox-vs-scratch-18b8c835d9b8

### Building Docker images
- To build this docker image, run the following command:  `$ docker build -t easymirror-backend:TAG .`

### Dockerhub
- As we will not be paying for an organization Dockerhub, all containers will be stored in a personal docker hub.

## TODOs
- [ ] Endpoint to allow uploads
- [ ] Account endpoint
    - [ ] Endpoint that returns account info
    - [ ] Endpoint that allows updating account info
- [ ] History
    - [ ] Endpoint that returns upload history
    - [ ] Endpoint to allow renaming of history item
    - [ ] Endpoint that gives list of files in history link
    - [ ] Endpoint to delete history item
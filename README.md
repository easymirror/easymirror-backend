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
- To build a docker image with env variables: `$ docker-compose up`

### Dockerhub
- As we will not be paying for an organization Dockerhub, all containers will be stored in a personal docker hub.


## CI/CD
### Process
- In order to have a proper CI/CD workflow, updates will be done in stages.
1. Create new branch with goal feature
    - Do all necessary commits to make the feature work
2. Make a pull request to merge into the `staging` branch
3. If all tests pass and nothing breaks, make a pull request into the `v*.*.*` branch.
### Branches
| Name | Description | Example
| - | - | - |
| `\*Feature name*` | Will be used as a development branch | `add_user_endpoint` branch
| `staging` | Nearly exact replica of a production environment for testing. | `staging`
| `v*.*.*` | Quality Assurance branch. | `v1.0.2`
| `main` | Production branch. Will be used by clients. | `main`


## Mirroring Flow
1. User makes a request to get a presigned URL to upload the files to AWS S3
2. User makes a request telling server which hosts to mirror to
3. Server mirrors to other hosts
    1. Using the mirror ID (UUID), lookup the folder in the S3 bucket
    2. For each file in folder:
        1. Create a private presigned URL
        2. Download the contents from the presigned URL and upload to other host

## TODOs
- [x] Integrate postgresSQL
- [x] Upload endpoints
    - [x] Accept multiple files
    - [x] For each file uploaded:
        - [x] Add file data to `files` table in database
        - [x] Upload to AWS S3 bucket
            - [x] Create a new folder for uploads that are uploaded together
        - [x] Upload to other hosts
        - [x] After uploading to other hosts, delete from S3 bucket
- [x] Account endpoint
    - [x] Endpoint that returns account info
    - [x] Endpoint that allows updating account info
- [x] History
    - [x] Endpoint that returns upload history
    - [x] Endpoint to allow renaming of history item
    - [x] Endpoint that gives list of files in history link
    - [x] Endpoint to delete history item
- Authentication
    - [x] When a new user joins, set a JWT
    - [x] JWT refresh every 12 hours
- [x] When deleting mirror links, cascade delete all relevant files too
- [x] When creating a new user, set the `member_since` column
- [x] Refactor getting JWT token
    - [x] 1. in the `user` package, create a `FromEcho` function to convert JWT to user
    - [x] 2. Refactor all code to get JWT token
- [x] own package for `mirrorlinks`?
- [x] Endpoint for shareable mirror link

## Benchmarks
### Download benchmarks
Benchmark to see which downloading method will be most efficient. Benchmark(s) can be found [here](/tests/download_test.go).
```MD
goos: darwin
goarch: arm64
pkg: github.com/easymirror/easymirror-backend/tests
BenchmarkDownloadPresigned
BenchmarkDownloadPresigned-14                  1        1515022625 ns/op          301816 B/op       3332 allocs/op
BenchmarkDownloadS3Manager
BenchmarkDownloadS3Manager-14                  1        13914433208 ns/op       84075549712 B/op            7245 allocs/op
PASS
ok      github.com/easymirror/easymirror-backend/tests  15.567s
```
### Upload Benchmarks
Benchmark to see which uploading method will be most efficient. Benchmark(s) can be found [here](/tests/download_test.go).
```
goos: darwin
goarch: arm64
pkg: github.com/easymirror/easymirror-backend/tests
BenchmarkUploadWithPipe-14    	1000000000	         0.08321 ns/op	       0 B/op	       0 allocs/op
BenchmarkUploadWithBuf-14     	       1	1275108626 ns/op	   28520 B/op	     157 allocs/op
BenchmarkRealWorldPipe-14     	       1	2171659833 ns/op	  373672 B/op	    2127 allocs/op
PASS
ok  	github.com/easymirror/easymirror-backend/tests	4.652s
```
## YouTube Video Fetcher API

### Documentation
For detailed information on how to use the project's API, please refer to the [Postman Documentation](https://documenter.getpostman.com/view/19816367/2s9Y5R1Rh4).

The YouTube Video Fetcher API is a GoLang application that fetches the latest videos from YouTube based on a specified search query. It stores video data in a PostgreSQL database and provides endpoints to retrieve video information in a paginated response. This project is built using GoFiber for handling HTTP requests and GORM for database operations.


### Getting Started

Clone the repository:

``git clone https://github.com/souvik150/Fampay-Backend-Assignment``

Create a .env file with all the details.
Install Make, Docker and Docker Compose if not already installed.

Run the following commands in the project root directory:

````

make build         # Build Docker containers
make up            # Start Docker containers in the background

````

Your API will be running on http://localhost:80

Test whether API is up at http://localhost/api/v1/healthcheck

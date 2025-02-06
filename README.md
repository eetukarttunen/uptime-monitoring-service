# Uptime Monitoring Service

A fully Dockerized uptime monitoring service built using **Golang**, **PostgreSQL**, and **Grafana**. This project allows user to monitor websites' uptime and performance (status code and response time) at regular intervals.

I personally use this service to monitor my own websites and experiments. 😊

## Tech Stack

- **Golang** ensures that the website monitoring runs with high performance and simplifies the maintenance of the backend.
- **PostgreSQL** for storing monitoring results.
- **Grafana** for visualizing the data.
- **Docker** to run everything in isolated containers.

## Requirements

**Docker** installed on your machine.

## Setup

### 1. Set up environment variables
In the root folder of the project (`uptime-monitoring-service`), run the following commands to create your `.env` file with your preferred configurations:

```
chmod +x setup_env.sh
./setup_env.sh 
```

### 2. Start the services:
Once the environment variables are set, user can run the following command to build and start all services (PostgreSQL, Grafana, and the Golang-based monitoring service):

```
docker-compose up -d --build
```

### 3. Access the services
**PostgreSQL**: The database will be available on port 5432 for any client that connects with the correct credentials. You can also access the database with this command:
```
docker exec -it <container_name> psql -U <username> -d <database_name>
```

**Grafana**: Open http://localhost:3000 in your browser. The username and password is defined in the .env file.

In grafana you can add new data source from here:

![image](https://github.com/user-attachments/assets/12f74b2c-0175-4fa6-add5-e5072fed26ef)

You can find the host URL (default is uptime_postgres) of your running containers by using the following command:
```
docker ps
```

In dashboards view you can create dashboard and visualization to visualize data:

![image](https://github.com/user-attachments/assets/a8857a86-5a3c-4478-9b83-3fe050d75371)

Example:

![image](https://github.com/user-attachments/assets/e68d3842-61e5-4f9b-a6e9-198f70b54e78)




# Auto Deploy

## Overview
This is a tool that allows a user to spin up a simple HTTPS server that will listen for incoming POST requests and execute a git update on a specified repo in a specified location on the host.

## Usage
See the `start-ad-docker.sh` script in `server-scripts` and configure environment variables as needed.


## tracking changes:
Here is the npm version of the dockerfile 
```
# Use an official Node runtime as a parent image
FROM node:18

# Set the working directory in the container
WORKDIR /

# Copy package.json and package-lock.json to the working directory
COPY public/ ./public
COPY src/ ./src
COPY package*.json ./

# Install app dependencies
RUN npm install

# Define the command to run your app
CMD ["npm", "start"]
```

# Status

Not working as of 2022/05/09 because domain look up service is broken. See [issue](https://github.com/Cgboal/SonarSearch/issues/47).

# About

Demo: One of these should work: [cloudsandbox.dev](https://cloudsandbox.dev/), [ipmapper.cloudsandbox.dev](https://ipmapper.cloudsandbox.dev/)

This app will:

  - Search for all the sub-domains of the domain entered
  - Find the IPs of the domain and each sub-domain found
  - Fetch location data including data of any nearby urban area of each IP found
  - Store the data fetched from the different services into the database
  - Show the locations ~~and associated data~~ of those IPs on the map

I have no idea if this app will be of use to anybody. I made this app for the purpose of using it as a demo of my skills in looking for software engineering positions. The app's features were basically decided on which services out there are free and do not require any sign ups. There's one service that require a sign up though. This service provides the location of an IP address. However, this is optional (since multiple IP location services are used in the app). If a particular API key is not found, this service won't be used.

# Current Issues

Because of garbage data from the [urban area data service](api.teleport.org), I've disabled showing data from this service.

# Backend
The backend is written in Go. I decided to use my own project structure instead of using something like [golang-standards project layout](https://github.com/golang-standards/project-layout). In some frameworks like Ruby on Rails, the code for the models and the functionality to fetch data from APIs are put in separate directories like "models" and "services". I decided to combine the models and services code into one location under the "entities" directory. Adding features and debugging seems easier not having to jump between multiple levels of directories.

The "entities" are basically the meat of the backend. In it, the app fetches data from different services, combines the data, transforms the data, stores the transformed data into the database and sends the data to the client. Getting data from all API services takes time. By storing the data fetched from those APIs in the database, if a client enters a domain that is already found in the database, the response will be a lot quicker.

# Frontend
The frontend is a React app. I decided to keep it simple. My aim was to quickly get a demo up and running. Also, my design skills are horrendous. I need to improve in this area.

# Deployment
The app is deployed to a single AWS EC2 instance using nginx as a reverse proxy. The app uses Terraform to provision the VPC, EC2 instance, RDS, and other components. TODO: use auto scaling group.

Configuration management is handled by Ansible. Among other things, Ansible installs the OS packages, sets up the language environments for Go and JavaScript, builds the backend, builds the frontend, installs and enables a systemd init script, initializes the database, and sets up the certificates using Let's Encrypt's CertBot tool.

There's a Makefile task that pretty much builds the app from zero to launch if all the tools like the Terraform and Ansible executables, etc. and all the required secret files and AWS keys are present.

# TODO

This will be an ongoing project for the foreseeable future. There are lots of things to do:

- Clean up. Among other things to clean up, there are a lot of Makefile tasks in different Makefiles that either need to be improved or removed.
- CI/CD
- Use a CDN service to serve the frontend build files
- Restructure the code so it would be easier for anybody to build and deploy the app. One option is to create a file that lists all the paths of the secret files, key files, etc.
- More tests
- Make the tests self-contained. Some end-to-end tests use external APIs. Use a test server instead. The challenge is injecting the test server information into the code under test.
- Separate deployment code. Make the deployment code more generalized. Probably just give it the git repository and list of secret files to deploy instead of hard-coding them.

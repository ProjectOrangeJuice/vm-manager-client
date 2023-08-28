# VM Manager Client

The VM Manager Client is a lightweight service designed to run on virtual machines (VMs) and Containers. It connects securely to the `vm-manager-server` using Mutual TLS (MTLS) authentication. Once connected, it waits for commands from the server and sends back results, allowing for remote monitoring and management of VMs.

## Features

- **Secure Connection**: Uses MTLS for a secure connection to the server, ensuring both client and server can trust each other.
- **System Monitoring**: Reports back vital system information such as:
  - CPU usage
  - RAM utilization
  - Storage capacity and usage
- **Remote Service Management**: Allows specific services to be restarted remotely.

## Prerequisites

- A running instance of `vm-manager-server`.
- A copy of the servers cert

## Installation


## Configuration

After running for the first time a filed called `config.json` will be generated. You should edit the file:

- `Name`: This clients (vm) name.
- `KeyLocation`: Path to the certificate (servers, and the ones this client will generate).
- `ServerAddress`: Address for the server.

## Usage


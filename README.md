# Movie-hub-server

## Prerequisites

To run this project, you need to have the following tools installed on your system:

- **make**: 
  ```sh
  sudo apt install make
  ```
- **Go**: 
  ```sh
  sudo snap install go --classic
  ```
- **sqlc**: 
  ```sh
  sudo snap install sqlc
  ```
- **Docker**

## Getting Started

Follow these steps to start the server:

1. **Initialize MySQL**:
   ```sh
   make mysql
   ```

2. **Create Tables**:
   ```sh
   make create_table
   ```

3. **Start the Server**:
   ```sh
   make server
   ```


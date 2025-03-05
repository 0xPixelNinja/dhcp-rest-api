# Usage

This document provides instructions on how to run the DHCP REST API and interact with its endpoints.

## Running the Server

1. **Install Dependencies:**  
   Ensure you have Python 3.7+ installed. Then install the dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. **Set Up Environment Variables:**  
   Create a `.env` file at the root of the project with the appropriate configuration. See [Configuration](CONFIGURATION.md) for details.

3. **Start the Server:**  
   Run the FastAPI server using uvicorn:
   ```bash
   uvicorn app.main:app --reload
   ```
   The `--reload` flag enables auto-reload during development.

4. **Access the API:**  
   Open your browser and navigate to [http://127.0.0.1:8000/docs](http://127.0.0.1:8000/docs) for the interactive API documentation.

## API Endpoints

### DHCP Hosts Endpoints

- **GET /hosts/**  
  Retrieves a list of all DHCP host entries.

- **POST /hosts/**  
  Adds a new DHCP host entry.  
  *Example JSON:*
  ```json
  {
      "name": "vm1006",
      "hardware_ethernet": "00:16:3e:aa:bb:cc",
      "option_routers": "192.168.1.1",
      "option_subnet_mask": "255.255.255.0",
      "fixed_address": "192.168.1.10",
      "option_domain_name_servers": "8.8.8.8,8.8.4.4"
  }
  ```

- **PUT /hosts/{name}**  
  Updates an existing DHCP host entry. Supports partial updates—only include fields to be updated.  
  *Example JSON:*
  ```json
  {
      "option_routers": "192.168.1.254",
      "fixed_address": "192.168.1.20"
  }
  ```

- **DELETE /hosts/{name}**  
  Deletes the specified DHCP host entry.

### Interfaces Endpoints

- **GET /interfaces/**  
  Returns the current interface configuration.

- **POST /interfaces/**  
  Adds a new interface to the configuration.  
  *Example JSON:*
  ```json
  {
      "type": "v4",
      "interface": "eth2"
  }
  ```

- **DELETE /interfaces/**  
  Removes an interface from the configuration.  
  *Example JSON:*
  ```json
  {
      "type": "v4",
      "interface": "vmbr0"
  }
  ```

**Note:** All API calls require an `Authorization` header with a Bearer token, for example:  
`Authorization: Bearer your-secret-token`

## Testing the API

You can test the endpoints using tools such as cURL, Postman, or directly through the Swagger UI available at `/docs`.

Happy managing your DHCP configurations!
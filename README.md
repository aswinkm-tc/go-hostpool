# go-hostpool

# Host Pool Implementation

This project demonstrates a host pool implementation using the Beta distribution to dynamically select the best-performing host based on response times. The logic is implemented in Go and simulates a scenario where multiple servers respond with varying delays.

## Key Components

### Host Pool
The `hostpool` struct maintains a list of server addresses and their respective alpha and beta values for the Beta distribution. These values are updated dynamically based on the performance of each host.

- **Alpha and Beta Values**: These parameters are used to calculate the Beta distribution for each host. Higher alpha values indicate better performance, while higher beta values indicate worse performance.
- **Host Selection**: The `selectHost` method uses the Beta distribution to probabilistically select the best-performing host.
- **Reward Update**: The `update` method adjusts the alpha and beta values based on the response time of the selected host.

### Server Simulation
Each server is simulated with a unique response time. The `startServer` function creates an HTTP server that responds with a delay to simulate varying performance.

### Main Logic
1. **Server Initialization**: Multiple servers are started with different response times.
2. **Host Selection**: The `selectHost` method is called to choose a host based on its performance.
3. **Request Handling**: An HTTP GET request is sent to the selected host, and the response is processed.
4. **Reward Calculation**: The reward is calculated based on the response time, and the host's alpha and beta values are updated accordingly.

## How It Works
1. Servers are started on different ports with simulated response times.
2. The host pool dynamically selects the best-performing host using the Beta distribution.
3. The performance of each host is continuously evaluated and updated based on response times.

## Example Output
```
Starting server on :8083
Starting server on :8081
Starting server on :8084
Starting server on :8082
Final counts: map[:8081:959 :8082:8 :8083:7 :8084:26]
...
```

## Usage
1. Clone the repository.
2. Run the program using `go run main.go`.
3. Observe the dynamic host selection and response handling in the console output.

## Dependencies
- [gonum/stat](https://pkg.go.dev/github.com/gonum/stat): Used for Beta distribution calculations.

## License
This project is licensed under the MIT License. See the LICENSE file for details.
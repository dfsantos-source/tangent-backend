# tangent-backend

## Overview
Tangent is an iOS application designed to help users find places of interest along their navigation routes. Users will enter their destination using the search bar on top and they will receive a list of possible detours that can be chosen to be included in the route. These include restaurants, stores, events, scenic areas, and more. After users choose which places they would like to visit, they will see a route traced on the map that will bring them through each one ending in their destination.

This API is built in Go and is responsible for fetching data from two dependencies: Mapbox Directions API and Yelp Fusion API. By leveraging Go, concurrent requests are optimized to improve API response time. Additionally, this project has focused on clean and scalable architecture, focused on dependency injection, separation of concerns, and best design principles. Synchronization techniques are also used alongside cancellation of context to propagate errors for common tasks. 

## Visual
<img width="661" alt="Screen Shot 2023-05-16 at 2 01 03 PM" src="https://github.com/dfsantos-source/tangent-backend/assets/64881219/64d4bf10-c489-4865-87b8-080fd3ecf148">

## How to build
1. Create a local.env file and insert it into the `/root` directory. The reason we need this as an `.env` file is because we need to keep our API keys a secret, for both Yelp and Mapbox. To retrieve these tokens, visit: https://docs.mapbox.com/api/navigation/directions/ and https://fusion.yelp.com/ 
2. This `.env` file should have the following API keys and be configured as follows:
```
  TOKEN_MAPBOX=[insert your own key here]
  TOKEN_YELP=[insert your own key here]
```
3. Change directory to /cmd
```
  ~ cd /cmd
```
4. Run the server. The server will default to port 3000. 
```
  ~ go run main.go 
```

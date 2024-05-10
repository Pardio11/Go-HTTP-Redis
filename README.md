# Go CRUD HTTP with Redis

A simple CRUD HTTP API built with Go's standard library `net/http` and using Redis as the database to store car information.

## Endpoints

- `GET /cars/{id}`: Retrieves information about a specific car.
  - Example: `localhost:8080/cars/mustang-1992`
  - Response:
    ```json
    {
        "brand": "Ford",
        "model": "Mustang",
        "year": 1992,
        "motor": {
            "size": 2.8,
            "horsepower": 543,
            "torque": 34.5,
            "max_rpm": 10000
        }
    }
    ```

- `GET /cars`: Retrieves information about all cars in the database.
  - Response:
    ```json
    {
        "mustang-1992": {
            "brand": "Ford",
            "model": "Mustang",
            "year": 1992,
            "motor": {
                "size": 2.8,
                "horsepower": 543,
                "torque": 34.5,
                "max_rpm": 10000
            }
        },
        "mustang-1993": {
            "brand": "Ford",
            "model": "Mustang",
            "year": 1993,
            "motor": {
                "size": 2.8,
                "horsepower": 543,
                "torque": 34.5,
                "max_rpm": 10000
            }
        }
    }
    ```

- `POST /cars`: Creates a new car entry.
  - Request Body:
    ```json
    {
        "brand": "Ford",
        "model": "Mustang",
        "year": 1999,
        "motor": {
            "size": 2.8,
            "horsepower": 543,
            "torque": 34.5,
            "max_rpm": 10000
        }
    }
    ```
  - Response:
    ```json
    {
        "key": "mustang-1999"
    }
    ```

- `DELETE /cars/{id}`: Deletes a car entry.
  - Example: `localhost:8080/cars/mustang-1999`
  - Response: `HTTP Status OK`

## Installation

[Add instructions on how to install and run your project here.]

## Usage

[Add instructions on how to use your API and interact with its endpoints.]

## Contributing

[Explain how others can contribute to your project, including guidelines for submitting pull requests and reporting issues.]

## License

[Specify the license for your project.]

 
 

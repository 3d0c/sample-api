## Sample RESTful API

### Prerequisites

You have to had ready to use PostgreSQL installed. Now create sample database:

```sql
CREATE DATABASE sampleapi;
```
### Installation

```sh
go env -w GO111MODULE=auto
go get github.com/3d0c/sample-api
```

No you should have an `sample-api` binay inside your `$GOPATH/bin/`.

Run it by 

```sh
DBUSER=validuser $GOPATH/bin/sample-api
```
Also there are another available environment variables:

- `DBHOST` Database hostname, default `127.0.0.1`
- `DBPORT` Database port, default `5432`
- `DBUSER` Valid database user, default `postgres`
- `DBNAME` Database name, default `sampleapi`

### Endpoints and API

#### User registration

```
POST /users
```

Endpoint expects valid JSON object. Required fields:

- `name`
- `password`

Example:

```sh
curl  -H "Content-Type: application/json" \
--data '{"name":"test4","password":"xyz"}' \
-XPOST http://localhost:5560/users
```

Expected result:

```javascript
{
    "ID": 6,
    "name": "test4"
}
```

#### User login

```
POST /users/login
```

Endpoint expects valid JSON object. Required fields:

- `name`
- `password`

Example:

```sh
curl  -H "Content-Type: application/json" \
--data '{"name":"test4","password":"xyz"}' \
-XPOST http://localhost:5560/users/login
```

Expected result is a JWT token:

```javascript
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTYyNDYwMjIsImlkIjo0LCJuYW1lIjoidGVzdDMifQ.LNMR-KIHe79l7rb68f40FRrZ2KdzzCgztzsWenCKUt4"
}
```

#### Add a flight

```
POST /flights
```

Endpoint expectes valid JSON object. Required fields:

- `name` Flight name
- `number` Flight number
- `scheduled` Flight scheduled time
- `arrival` Flight arrival time
- `departure` Flight departure time
- `destination` Flight destination
- `fare` Flight fare
- `duration` Exsimated flight duration

All time fields should be passed in format `2021-01-01T09:09:09Z`

Example

```sh
curl \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTYyNDYwMjIsImlkIjo0LCJuYW1lIjoidGVzdDMifQ.LNMR-KIHe79l7rb68f40FRrZ2KdzzCgztzsWenCKUt4" \
--data '{"name": "Test flight", "number": "AB555", "scheduled": "2021-04-01T09:00:00Z", "arrival": "2021-04-01T10:00:00Z", "departure": "2021-04-01T09:00:00Z", "destination": "Moscow", "fare": 140, "duration": 60}' \
-XPOST http://localhost:5560/flights
```

Expected result:

```javascript
{
    "ID": 3,
    "name": "Test flight",
    "number": "AB555",
    "Scheduled": "2021-04-01T09:00:00Z",
    "Arrival": "2021-04-01T10:00:00Z",
    "Departure": "2021-04-01T09:00:00Z",
    "destination": "Moscow",
    "Fare": 140,
    "Duration": 60
}
```

#### Update flight

```
PUT /flights/:id
```

Example

```sh
curl \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTYyNTA0MTEsImlkIjo0LCJuYW1lIjoidGVzdDMifQ.e9byTsyeX5FUw-e1uTmjDuzoGYIztqIm780K5yRTSNc" \
--data '{"fare": 145, "number": "AB556"}' \
-XPUT http://localhost:5560/flights/3
```

To avoid doing extra SELECT request, the expected result is 200 OK or error.

#### Search for a flight

```
GET /flights
```

By default all flights has been returned. To filter the output there are few parameters available:

- `flight_name`
- `scheduled_date`
- `departure`
- `destination`

Example

```sh
curl \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTYyNTQxODgsImlkIjo2LCJuYW1lIjoidGVzdDQifQ.j4RjviXHe9y4K7D_ZDVMo5Ut1MunqjMvG8AoPMTNHMk" \
-XGET http://localhost:5560/flights\?destination\=Moscow
```

Expected result:

```javascript
[
    {
        "ID": 4,
        "name": "Test flight",
        "number": "AB551",
        "Scheduled": "2021-04-01T12:00:00+03:00",
        "Arrival": "2021-04-01T13:00:00+03:00",
        "Departure": "2021-04-01T12:00:00+03:00",
        "destination": "Moscow",
        "Fare": 140,
        "Duration": 60
    }
]
```

#### Remove flight

```
DELETE /flights/:id
```

Example 

```sh
curl \
-H "Content-Type: application/json" \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTYyNTQxODgsImlkIjo2LCJuYW1lIjoidGVzdDQifQ.j4RjviXHe9y4K7D_ZDVMo5Ut1MunqjMvG8AoPMTNHMk" \
-XDELETE http://localhost:5560/flights/4
```

Expected result `200 OK` or error.
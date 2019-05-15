# go4graphql

Simple example using [gqlgen](https://github.com/99designs/gqlgen)

## Run

```sh
docker-compose up
```

Connect to [http://localhost:3000](http://localhost:3000)

Data is kept in memory. Try:

```graphql
query {
  vendors {
    name
    urls {
      url
    }
  }
}
```

```graphql
mutation {
  createVendor(input: {name: "FLV", code: "Flavorah"}) {
    name
    code
  }
}
```

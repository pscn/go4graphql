# go4graphql

Simple example using [gqlgen](https://github.com/99designs/gqlgen)

## Run

```sh
docker-compose up
```

Then connect to [http://localhost:3000](http://localhost:3000).

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

query {
  concentrates {
    name
    vendor {
      code
    }
  }
}

mutation {
  createVendor(input: {name: "FLV", code: "Flavorah"}) {
    name
    code
  }
}
```

## Develop

Install gqlgen:

```sh
go get github.com/99designs/gqlgen
```

Fetch dependencies:

```sh
go get ./...
```

If you modified the schema run:

```sh
gqlgen
```

Modifiy `resolver.go`.

# Bling 

Bling is a Go HTTP client library for creating and sending API requests written in Golang/Go.

![bling Logo](bling-logo.png)

[![Build Status](https://travis-ci.org/blinghttp/bling.svg "Travis CI status")](https://travis-ci.org/blinghttp/bling)    

Bling has cool methods which simplify setting HTTP Request properties and sending request. Check [usage](#usage) or the [examples](examples) to learn how to use Bling in your code.


We're still under active development, adding support methods for retries, callbacks, etc.

### Features

* Method Setters: Get/Post/Put/Patch/Delete/Head
* Adding Request Headers
* Encode structs into URL query parameters

## Install

    go get github.com/blinghttp/bling


## Usage

Use Bling client to set request properties and easily send http request
```
blingClient := bling.New() //use's default http client if .Client(customClient) is not used
resp, err := blingClient.Get("https://github.com/blinghttp/bling").DoRaw() //DoRaw method returns raw http response
if err != nil {
    fmt.Println("Some err", err)
    return
}
defer resp.Body.Close()
responseBody, err := ioutil.ReadAll(resp.Body)
if err != nil {
    fmt.Println("Some err", err)
    return
}
fmt.Println(string(responseBody))
```

## Contributing

We love PR's and highly encourage to contribute.

## License

[Apache License](LICENSE)
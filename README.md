# backbone
A basic library for golang.

## Install http api generator

  ```
  go install github.com/RivenZoo/backbone/cmd/http_codegen
  ```
  
## Use example
  
  ```
  http_codegen -input ./examples/demo_server/controllers/short_url.go
  ```

## Install project creator

  ```
  go install github.com/RivenZoo/backbone/cmd/projcreator
  ```

## Use example

  ```
  projcreator -project test -gitRepo https://github.com/RivenZoo/injectgo.git -output /tmp/injectgo -gitVer v0.4.0
  ```
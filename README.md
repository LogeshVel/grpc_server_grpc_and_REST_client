# grpc_server_grpc_and_REST_client
Simple CRUD gRPC API Server in Golang which can serve both gRPC and REST API clients using gRPC-Gateway.


This repo is the extenstion of the [Golang_grpc_server_Python_grpc_client](https://github.com/LogeshVel/Golang_grpc_server_Python_grpc_client) repo (Golang gRPC Server with Python gRPC Client) here inaddition to gRPC client it also supports the REST API clients


For Swagger I have trying a lot.


Need to clone 2 repo.

    1) for the annotation proto -  git clone https://github.com/grpc-ecosystem/grpc-gateway.git and then copy the protoc-gen-openapiv2 folder to our Project folder
    2) for the UI codes - git clone https://github.com/swagger-api/swagger-ui.git and copy the contents of the dist folder into the OpenAPi folder and edit the swagger-intializer.js file to load our json file

And install **protoc-gen-swagger** plugin to generate the swagger.json file from the proto file by annotation

     go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest



The JSON encoding for the Mask is different

[see here](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#json-encoding-of-field-masks)


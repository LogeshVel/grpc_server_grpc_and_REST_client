# grpc_server_grpc_and_REST_client
Simple CRUD gRPC API Server in Golang which can serve both gRPC and REST API clients using gRPC-Gateway.


This repo is the extenstion of the [Golang_grpc_server_Python_grpc_client](https://github.com/LogeshVel/Golang_grpc_server_Python_grpc_client) repo - Golang gRPC Server with Python gRPC Client here inaddition to gRPC client it also supports the REST API clients


Pre-requisites:

- Plugin for **_protoc_**. In addition to **protoc-gen-go** plugin to generate pb.go file and **protoc-gen-go-grpc** plugin to generate grpc.pb.go file, we need to have **protoc-gen-grpc-gateway** plugin to generate the pb.gw.go file (proxy server file)

    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

    $ go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

- Get the Protobuf package to work with Protobuf

For Golang

    $ go get google.golang.org/protobuf

For Python

    $ pip install protobuf

- Get the grpc package to work with gRPC

For Golang

    $ go get -u google.golang.org/grpc

For Python

    $ pip install grpcio-tools


### To compile the proto file to Go files we can use protoc and to Python files we can use protoc for only to generate pb files for grpc we can use grpcio-tools.

#### Using protoc

    $ protoc -I<path/to/protofile> --go_out=<outdir> --go-grpc_out=<outdir> --grpc-gateway_out=<outdir> --python_out=<outdir> protofilename.proto

here for the 3 go outputs we need the protoc plugin as we mentioned.

#### Using grpcio-tools for Python

    $ python -m grpc_tools.protoc -I<path/to/protofile> --python_out=outdir> --grpc_python_out=outdir> proto_filename.proto

##### For grpc-gateway output file for Python, I haven't tried out will update that once I tried that.


## grpc-gateway

To make the gRPC Server to serve both the gRPC and the REST clients we have do some annotations in the proto file for the rpc and compile that file to pb.gw.go file and make the server to serve the REST API Clients as well.

#### Step - 1

- Making annotations

To make the annotations for the rpc in the proto file we need to import the **"google/api/annotations.proto"** file. Unfortunately we have to manually import this file.

Clone or Download the codes of the repo [googleapis](https://github.com/googleapis/googleapis) and copy the **google** folder from this repo and paste it in your Project's root dir.

Note: Please avoid cloning the repo to the Project dir itself. Clone somewhere else and only copy the google folder from that repo.

Now the error should be resolved.

We can now annotate the rpc of the Services in the proto file.

Sample annotations.

```
import "google/api/annotations.proto";

service SomeService{
    rpc someGETRPC(RequestMessageType) returns (ResponseMessageType){
        option (google.api.http) = {
            get: "/v1/emp/{id}" // HTTP GET with this path(id here represents the id field of the RequestMessageType) returns ResponseMessageType
          };
    }
    rpc somePOSTRPC(RequestMessageType) returns (ResponseMessageType){
        option (google.api.http) = {
            post: "/v1/path" // HTTP POST with this path and the body creates the resources
            body: "*" // * in the body represents all the fields in the RequestMessageType. So all the fileds on the RequestMessageType will be provided in the body. (as a JSON since it is REST API)
        };
    }
}

```

#### Step - 2

Once we are done with the proto, now we move forward to co,pile the proto file.

    $ protoc -I<path/to/protofile> --go_out=<outdir> --go-grpc_out=<outdir> --grpc-gateway_out=<outdir> protofilename.proto

#### Step - 3

Now we need to configure the server to serve the REST API Clients also.

For that we need to make the gRPC server implementation, since the grpc-gateway connects to our gRPC server as the gRPC client and acts as proxy to our REST client.

Once the gRPC server implementation is done, make a gRPC client connection(our proxy) to connect with the gRPC server that we have implemented and create New Mux serve from the grpc-gateway runtime pkg.

Register our gRPC client(our proxy) with our Service handler function (this register function is present in the pb.gw.go file generated using the protoc-gateway plugin) by providing the client connection object and also the MuxServe object.

Now create the Gateway server by passing the socket address to serve and the Mux hanlder.

Make the Gateway server to Listen and Serve.

Our grpc-proxy server is up and running. It will convert the gRPC to REST and vice versa.

### Swagger

**Will revisit to complete the README for this**

For Swagger I have trying a lot.


Need to clone 2 repo.

    1) for the annotation proto -  git clone https://github.com/grpc-ecosystem/grpc-gateway.git and then copy the protoc-gen-openapiv2 folder to our Project folder
    2) for the UI codes - git clone https://github.com/swagger-api/swagger-ui.git and copy the contents of the dist folder into the OpenAPi folder and edit the swagger-intializer.js file to load our json file

And install **protoc-gen-swagger** plugin to generate the swagger.json file from the proto file by annotation

     go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest



The JSON encoding for the Mask is different

[see here](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#json-encoding-of-field-masks)



## This repo usage

Need to document

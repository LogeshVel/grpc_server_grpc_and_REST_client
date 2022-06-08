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

Making annotations

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

Once we are done with the proto, now we move forward to compile the proto file.

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

OpenAPI Swagger for the REST implementation can be created with some more little work.

To achieve Swagger implementation we have 2 main steps.

- annotation in proto

- UI codes

#### Step - 1

To make the OpenAPI related annotations in the proto file we have to import the **protoc-gen-openapiv2/options/annotations.proto**
Unfortunately we have to manually clone the [grpc-gateway repo](https://github.com/grpc-ecosystem/grpc-gateway.git) in any location and just copy the folder **protoc-gen-openapiv2** from this repo to our Project's root directory. And then do import the **_protoc-gen-openapiv2/options/annotations.proto_** and make the annotations.

Sample initilization in the proto file

```
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
    info: {
      version: "1.0";
    };
    external_docs: {
      url: "https://github.com/LogeshVel/grpc_server_grpc_and_REST_client";
      description: "gRPC Server - gRPC and REST clients";
    }
    schemes: HTTPS;
  };
```
Sample annotation in the rpc

```
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
          summary: "Get an employee"
          description: "Get an employee for the given ID"
          tags: "Employee"
        };
```


#### Step - 2

To generated the swagger JSON file from our proto file we need to install the **protoc-gen-swagger** plugin.

```
$ go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
```
Now compile the proto file with this plugin.

Example:
```
$ protoc -I<path/to/protofile> --go_out=<outdir> --go-grpc_out=<outdir> --grpc-gateway_out=<outdir> --swagger_out=<openAPIDir> protofilename.proto
```

To make the UI codes we have to clone the [swagger-ui repo](https://github.com/swagger-api/swagger-ui.git) and copy the contents of the _dist_ folder into the OpenAPi folder under swagger folder in the root directory of our Project(Project_Root_dir/swagger/OpenAPI for openAPI related stuffs) and edit the **_swagger-intializer.js_** file to load our json file. (The JSON file generated by compiling the proto file using the protoc-gen-swagger plugin)

Under the **swagger** folder we have _embed.go_ file that has the FileSystem code for our OpenAPI.

```
package OpenAPI

import (
	"embed"
)

//go:embed OpenAPI/*
var OpenAPI embed.FS

```

### JSON encoding for the FiledMask

The JSON encoding for the FiledMask is different than the encoding done with the Programming languages

[Reference doc](https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#json-encoding-of-field-masks)


## Folders not in my control

- google (has the proto file for the HTTP annotations)

- protoc-gen-openapiv2 (has the proto file for the OpenAPI annotations)


## This repo usage

_Need to document_

- run the Server

- then the gRPC server serves the gRPC client and the REST clients (has some code in the client/rest_api_client.py, also can use the Postman to test the endpoints)

**Add Medium article link**
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"

	oAPI "github.com/LogeshVel/grpc_server_grpc_and_REST_client/swagger"

	pb "github.com/LogeshVel/grpc_server_grpc_and_REST_client/proto/emp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const socket string = "localhost:50051"
const gwSocket string = "localhost:51151"
const swaggerSocket string = "localhost:50050"

var employeeList []*pb.Employee

type Server struct {
	pb.EmployeeManagementServer
}

func main() {
	lisn, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatalln("Errored while Listen to : ", socket, err)
	}
	log.Println("gRPC Server Listening at ", socket)
	s := grpc.NewServer()
	pb.RegisterEmployeeManagementServer(s, &Server{})
	go s.Serve(lisn)

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		socket,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Employee
	err = pb.RegisterEmployeeManagementHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    gwSocket,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on ", gwSocket)
	// go gwServer.ListenAndServe()
	gwServer.ListenAndServe()
	// oa := getOpenAPIHandler()
	// swaggerServer := &http.Server{
	// 	Addr: swaggerSocket,
	// 	Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if strings.HasPrefix(r.URL.Path, "/api") {
	// 			gwmux.ServeHTTP(w, r)
	// 			return
	// 		}
	// 		oa.ServeHTTP(w, r)
	// 	}),
	// }
	// log.Println("Serving OpenAPI Doc at", swaggerSocket)

	// swaggerServer.ListenAndServe()

}

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(oAPI.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

func (s *Server) GetEmployee(ctx context.Context, empID *pb.EmployeeID) (*pb.Employee, error) {
	log.Println("Hitted GetEmployee with the employee ID", empID.Id)
	for _, e := range employeeList {
		if e.Id == empID.Id {
			return e, nil
		}
	}
	return nil, status.Errorf(
		codes.NotFound,
		"Given employee ID is not found",
	)
}

func (s *Server) GetEmployeeByQP(ctx context.Context, empID *pb.EmployeeID) (*pb.Employee, error) {
	log.Println("Hitted GetEmployeeByQP with the employee ID", empID.Id)
	for _, e := range employeeList {
		if e.Id == empID.Id {
			return e, nil
		}
	}
	return nil, status.Errorf(
		codes.NotFound,
		"Given employee ID is not found",
	)
}

func (s *Server) CreateEmployee(ctx context.Context, emp *pb.Employee) (*pb.EmployeeID, error) {
	log.Println("Hitted CreateEmployee with the emp ID", emp.Id)
	log.Println("Checking whether the given emp id is already there")
	for _, e := range employeeList {
		if e.Id == emp.Id {
			log.Println("The Given employee ID already present. So skipping the create process")
			return nil, status.Errorf(
				codes.AlreadyExists,
				"Given employee ID already exists. use UpdateEmployee to update",
			)
		}
	}
	employeeList = append(employeeList, emp)
	empID := pb.EmployeeID{Id: string(emp.Id)}
	log.Println("Given employee ID doesn't exists. Succesfully created.")
	return &empID, nil

}

func (s *Server) ListEmployees(emt *emptypb.Empty, stream pb.EmployeeManagement_ListEmployeesServer) error {
	log.Println("Hitted ListEmployees")
	for _, e := range employeeList {
		stream.Send(e)
	}
	return nil
}

func (s *Server) UpdateEmployee(ctx context.Context, emp *pb.Employee) (*emptypb.Empty, error) {
	log.Println("Hitted UpdateEmployee to update the employee ID", emp.Id)
	for i, e := range employeeList {
		if e.Id == emp.Id {
			employeeList[i] = emp
			log.Println("Updated employee\n", emp)
			return &emptypb.Empty{}, nil
		}
	}

	log.Println("Employee not found to update")
	return nil, status.Errorf(
		codes.NotFound,
		"Given employee ID not found to update",
	)
}

func (s *Server) DeleteEmployee(ctx context.Context, empID *pb.EmployeeID) (*emptypb.Empty, error) {
	log.Println("Hitted DeleteEmployee to Delete the emp ", empID)
	for i, e := range employeeList {
		if e.Id == empID.Id {
			// slicing won't raise Index out of range error.
			// at this point employeeList[i+1:]  if i is the last index the i+1 return empty slice in slicing.
			// no error
			employeeList = append(employeeList[:i], employeeList[i+1:]...)
			log.Println("Deleted the Employee")
			fmt.Println(employeeList)
			return &emptypb.Empty{}, nil
		}
	}

	log.Println("Employee not found to delete")
	return nil, status.Errorf(
		codes.NotFound,
		"Given employee ID is not found to delete",
	)
}

func (s *Server) PatchEmployee(ctx context.Context, req *pb.UpdateEmpRequest) (*emptypb.Empty, error) {
	log.Println("Hitted PatchEmployee to update the emp ", req.Emp.GetId())
	log.Printf("Given request message\n%v", req)
	var matchEmp *pb.Employee
	for _, e := range employeeList {
		if e.Id == req.Emp.GetId() {
			matchEmp = e
			break
		}
	}
	if matchEmp == nil {
		n := fmt.Sprintf("Employee with ID %q could not be found", req.GetEmp().GetId())
		log.Println(n)
		return &emptypb.Empty{}, status.Errorf(codes.NotFound, n)
	}

	for _, path := range req.GetUpdateMask().GetPaths() {
		log.Printf("Requested to update %q", path)
		switch path {
		case "emp.first_name":
			matchEmp.FirstName = req.GetEmp().GetFirstName()
		case "emp.last_name":
			matchEmp.LastName = req.GetEmp().GetLastName()
		case "emp.role":
			matchEmp.Role = req.GetEmp().GetRole()
		case "emp.contact.home_addr":
			matchEmp.Contact.HomeAddr = req.GetEmp().GetContact().GetHomeAddr()
		case "emp.contact.mob_num":
			matchEmp.Contact.MobNum = req.GetEmp().GetContact().GetMobNum()
		case "emp.contact.mail_id":
			matchEmp.Contact.MailId = req.GetEmp().GetContact().GetMailId()
		default:
			invalid := fmt.Sprintf("cannot update field %q on Employee", path)
			log.Println(invalid)
			return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, invalid)
		}
	}
	log.Println("Partial Update done")

	return &emptypb.Empty{}, nil
}

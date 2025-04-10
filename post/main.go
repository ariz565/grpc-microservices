package main

import (
	"context"
	"log"
	"net"
	"sync"

	postpb "grpc-microservices/post/postpb"
	userpb "grpc-microservices/user/userpb"

	"google.golang.org/grpc"
)

type Post struct {
	Id      int32
	Title   string
	Content string
	UserId  int32
}

var (
	posts         = make(map[int32]Post)
	counter int32 = 1
	mu      sync.Mutex
)

type postServer struct {
	postpb.UnimplementedPostServiceServer
	userClient userpb.UserServiceClient
}

func (s *postServer) CreatePost(ctx context.Context, req *postpb.PostRequest) (*postpb.PostResponse, error) {
	mu.Lock()
	defer mu.Unlock()

	id := counter
	counter++

	posts[id] = Post{
		Id:      id,
		Title:   req.Title,
		Content: req.Content,
		UserId:  req.UserId,
	}

	return &postpb.PostResponse{
		Id:      id,
		Title:   req.Title,
		Content: req.Content,
		UserId:  req.UserId,
	}, nil
}

func (s *postServer) GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.PostResponse, error) {
	post, ok := posts[req.Id]
	if !ok {
		return nil, grpc.Errorf(404, "Post not found")
	}

	userResp, err := s.userClient.GetUser(ctx, &userpb.UserRequest{Id: post.UserId})
	if err != nil {
		return nil, err
	}

	return &postpb.PostResponse{
		Id:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		UserId:    post.UserId,
		UserName:  userResp.Name,
		UserEmail: userResp.Email,
	}, nil
}

func main() {
	// Connect to UserService
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer conn.Close()

	userClient := userpb.NewUserServiceClient(conn)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	postpb.RegisterPostServiceServer(server, &postServer{userClient: userClient})
	log.Println("PostService running on :50052")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

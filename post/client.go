package main

import (
	"context"
	"log"
	"time"

	pb "grpc-microservices/post/postpb"

	"google.golang.org/grpc"
)

func main() {
	// Connect to PostService
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPostServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// ✅ Create a post
	createResp, err := client.CreatePost(ctx, &pb.PostRequest{
		Title:   "gRPC is Awesome",
		Content: "This is a demo of gRPC microservices",
		UserId:  1,
	})
	if err != nil {
		log.Fatalf("CreatePost failed: %v", err)
	}
	log.Printf("Post Created: %+v\n", createResp)

	// ✅ Get the post
	getResp, err := client.GetPost(ctx, &pb.GetPostRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("GetPost failed: %v", err)
	}
	log.Printf("Post Fetched: %+v\n", getResp)
}

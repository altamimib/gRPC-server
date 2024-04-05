package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/grpc"

	textpb "github.com/altamimib/gRPC-server/textpb"
)

const (
	port     = ":9000"
	filePath = `./textpb/text_message.proto` // Updated
)

type textServer struct {
	textpb.UnimplementedTextServiceServer
}

func (s *textServer) SendText(ctx context.Context, in *textpb.TextLine) (*textpb.TextLine, error) {
	// This method is not used in this example but can be extended for future functionalities
	return in, nil
}

func readFileLines() ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	textpb.RegisterTextServiceServer(s, &textServer{})
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		for {
			lines, err := readFileLines()
			if err != nil {
				log.Printf("error reading file: %v", err)
			}
			for _, line := range lines {
				fmt.Println(line) // Simulate sending lines to client (replace with actual gRPC call)
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

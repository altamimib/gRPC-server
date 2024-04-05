package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

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

func (s *textServer) SendTextLines(stream textpb.TextService_SendTextLinesServer) error {
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			return nil // Client finished sending lines
		}
		if err != nil {
			return err // Handle errors during receiving lines
		}
		fmt.Println(line.GetContent()) // Process/log the received line
	}
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
				continue // Skip sending lines if reading fails
			}
			conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure()) // Update address if client runs elsewhere
			if err != nil {
				log.Printf("failed to connect to client: %v", err)
				continue // Skip sending lines if connection fails
			}
			defer conn.Close()
			client := textpb.NewTextServiceClient(conn)
			stream, err := client.SendTextLines(context.Background()) // Initiate stream with client
			if err != nil {
				log.Printf("failed to create stream: %v", err)
				continue // Skip sending lines if stream creation fails
			}
			for _, line := range lines {
				if err := stream.Send(&textpb.TextLine{Content: line}); err != nil {
					log.Printf("failed to send line: %v", err)
					// Consider handling individual line sending errors or closing the stream
				}
			}
			stream.CloseSend() // Close the sending stream after sending all lines
		}
	}()
}

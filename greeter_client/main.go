/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")

	DefaultMaxRetries = uint(5)
	DefaultTimeout    = 1 * time.Second
	DefaultBackoff    = 50 * time.Millisecond
)

func main() {
	flag.Parse()
	// Set up a connection to the server.

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithTimeout(DefaultTimeout),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			//grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(clientName)),
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithBackoff(grpc_retry.BackoffExponential(DefaultBackoff)),
				grpc_retry.WithCodes(codes.Aborted, codes.Unavailable),
				grpc_retry.WithMax(DefaultMaxRetries),
				grpc_retry.WithPerRetryTimeout(DefaultTimeout), // use the same timeout for dealing and the call
			),
		)),
	}
	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

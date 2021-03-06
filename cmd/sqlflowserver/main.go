// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// To run this program:
//	go generate .. && go run main.go
//
package main

import (
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sql-machine-learning/sqlflow/server"
	"github.com/sql-machine-learning/sqlflow/server/proto"
	"github.com/sql-machine-learning/sqlflow/sql"
)

const (
	port = ":50051"
)

func start(datasource string) {
	db, err := sql.Open(datasource)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterSQLFlowServer(s, server.NewServer(sql.Run, db))
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Println("Server Started at", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	ds := flag.String("datasource", "", "database connect string")
	flag.Parse()
	start(*ds)
}

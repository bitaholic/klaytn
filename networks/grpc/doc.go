// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

/*
Package grpc implements the gRPC protocol for Klaytn.

This package allows you to use klaytn's RPC API using gRPC.
See below for gRPC: https://grpc.io/docs/quickstart/go/

Source files

Each file provides the following features
 - gClient.go : gRPC client implementation.
 - gServer.go : gRPC server implementation.
 - klaytn.proto : Define interface and message to use in gRPC server and clieklaytnnt.
 - klaytn.pb.go : Go file created by compiling klaytn.proto.
 */
package grpc

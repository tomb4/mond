package cmd

import "text/template"

var protoTemplate, _ = template.New("").Parse(`

syntax = "proto3";
package {{.AppId}};

option go_package = "./{{.AppId}}";

// import "mond/wind/proto/plugin/plugin.proto";

message Nil{}

message PingReq {
  string name = 1;
  int32 age=2;
  string x=3;
  double dx=6;
}

message PingResp {
  string message = 1;
  int32 x=2;
  string y=3;
  float z=4;
  double a=5;
}

service {{.AppId}}Service {
  rpc Ping (PingReq) returns (PingResp) {}
}
`)

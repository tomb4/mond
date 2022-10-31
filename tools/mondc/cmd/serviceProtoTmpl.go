package cmd

import "text/template"

var protoTemplate, _ = template.New("").Parse(`

syntax = "proto3";
package {{.AppId}};

import "mond/wind/proto/plugin/plugin.proto";

message Nil{}

message PingReq {
    string name = 1 [(plugin.doc)={
        desc: "name desc";
    },(plugin.validator)={
        notEmpty: true;
        eq: 3.1;
        gte: 2.2;
        gt: 1.3;
        lte: 5.4;
        lt: 6.5;
  }];
  int32 age=2[(plugin.doc)={
        desc: "age desc";
    },(plugin.validator)={
        notEmpty: true;
  }];
  string x=3;
  double dx=6[(plugin.doc)={
       desc: "dx desc";
    },(plugin.validator)={
       notEmpty: true;
       eq: 3.1;
       gte: 2.2;
       gt: 1.3;
       lte: 5.4;
       lt: 6.123456789;
  }];
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

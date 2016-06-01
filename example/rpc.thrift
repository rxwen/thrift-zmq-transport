namespace go demo.rpc
namespace java demo.rpc

// 测试服务
service RpcService {
    string foo(1:i64 index, 2:string code, 3:string owner),
}

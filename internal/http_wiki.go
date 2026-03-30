package internal

// http - internet application communication protocol, which has it's own set of standards
// REST/HTTP communication involves 2 entity(software application/servers), one acts as the
// client(the one who sends the request to the server) another one acts as the server, whic responds to the
// the request meant for him from the client,
// request & responds messages are structured  under the standards defined in HTTP/1.1
// 1. request message - should consist of
//         -------------------------------------------
//         GET[METHOD- tells the intent of the request] /coffee[requests target] HTTP/1.1[HTTP version] \r\n[CRLF]  :request line
//         HOST: localhost:8080\r\n                     Header --> consist of key: value pairs, which includes meta data about the request
//         Content-Type: application/json\r\n
//         Accept: */*\r\n
//         \r\n
//         {Body >}

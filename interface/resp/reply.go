package resp

//简单字符串（Simple Strings）: 以 “+” 开头，例如 “+OK\r\n” 表示一个成功的响应。
//错误（Errors）: 以 “-” 开头，例如 “-ERR unknown command\r\n” 表示一个错误响应。
//整数（Integers）: 以 “:” 开头，例如 “:1000\r\n” 表示整数1000。
//批量字符串（Bulk Strings）: 以 “$” 开头，例如 “$6\r\nfoobar\r\n” 表示一个长度为6的字符串 “foobar”。
//数组（Arrays）: 以 “*” 开头，例如 “*3\r\n:1\r\n:2\r\n:3\r\n” 表示包含3个整数的数组 [1, 2, 3]。

type Reply interface {
	ToBytes() []byte
}

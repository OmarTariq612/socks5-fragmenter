
# socks5 handshake fragmenter
```
socks5 client <--> socks5-fragmenter <--> socks5 server
```
socks5 protocol specifies a handshake that can be easily detected by censorship firewalls. (`[0x05 0x01 0x00]` for no-auth-required method).

`socks5-fragmenter` is a relayer that relays connections to a specified socks server.

Its job is to fragment the handshake into (`[0x05 auth_methods_num]` and `[auth_methods...]`) instead of sending the handshake as a whole message (`[0x05 auth_methods_num auth_methods...]`) then it starts relaying.
## Build

```
go build .
```
## Usage
```
usage: socks5-fragmenter [bind_addr:bind_port] server_addr:server_port
```
`bind_addr:bind_port` is `:5555` by default.

ex: to provide `192.168.1.50:3456` as the socks server address:
```
./socks5-fragmenter 192.168.1.50:3456
```
```
2022/08/11 04:59:42 Serving on :5555
2022/08/11 04:59:42 the provided server address is 192.168.1.50:3456




```
Now the program is ready to accept socks5 client connections on `:5555` and relays them to the specified socks5 server `192.168.1.50:3456`
## REF
* socks 5 (rfc 1928) : https://datatracker.ietf.org/doc/html/rfc1928
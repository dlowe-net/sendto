# sendto - receive files from a web browser

I wrote this in a fit of pique after someone wanted to send me some
large multi-gigabyte files and asked which service I wanted to use.
"This is ridiculous," I thought.  "Sending stuff from one computer to
another is what the internet is _for_."

This is currently only tested on Ubuntu Linux 20.02.

## Fastest way to use

1. Configure your router to forward port 80 to another port (called
   `$PORT`) on your computer.  `$PORT` must be greater than 1024). Note
   your external address on the router (called `$ADDR`)
2. At an sh command line, run `mkdir -p ~/sendto/files`
2. Run `sendto -dir ~/sendto/files/`
3. Give the sender the url `http\://$ADDR:$PORT/`
4. The sender can now send you a file using that url!

## More convenient, more secure setup

1. Configure your router to forward port 80 and port 443 to your
   computer. Note your external address on the router (called $ADDR).
   The ports on your computer must also be 80 and 443 for certbot to
   work.
2. At a sh command line, run `mkdir -p ~/sendto/{files,certs}`
3. Set up a DNS entry on a domain name you own, like
   sendto.example.com
4. Use certbot from https://letsencrypt.com/ to generate an SSL
   private key and public certificate for your domain and put it in
   the certs directory.
```sh
   sudo certbot certonly --standalone
   sudo chown $USER {privkey,cert}.pem
   mv {privkey,cert}.pem ~/sendto/certs
```
5. Under Linux, give the binary the low-port binding capability:
```sh
sudo setcap cap_net_bind_service=+ep ~/go/bin/sendto
```
6. Run sendto:
```sh
sendto -dir ~/sendto/files/ \
       -key ~/sendto/certs/privkey.pem \
       -cert ~/sendto/certs/cert.pem \
       -hostname sendto.example.com
```
7. The sender can now send you a file using `https\://sendto.example.com/`!

## Argument reference

* `-dir` : (required) Path to output directory. `sendto` must have
write access.
* `-port` : Port to use for insecure web connections.  Defaults to 80.
* `-sport` : Port to use for secure web connections.  Defaults to 443.
* `-key` : Path to private key for secure web connections.
* `-cert` : Path to SSL certificate for secure web connections.
* `-hostname` : External hostname.  This is used only for a redirect
  when someone connects with an insecure http connection.

## Build instructions

1. Download the sendto source tarball.
2. You should have the Go language installed, which you can get from
[https://golang.com/].
3. Put the source where the Go build tools expect:
```sh
mkdir -p ~/go/src/
tar -C ~/go/src/ sendto.tgz
```
4. Build and install:
```sh
go install sendto
```
The sendto program should be at `~/go/bin/sendto`

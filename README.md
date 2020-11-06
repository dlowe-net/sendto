# sendto - receive files from a web browser

I wrote this in a fit of pique after someone wanted to send me some
large multi-gigabyte files and asked which service I wanted to use.
"This is ridiculous," I thought.  "Sending stuff from one computer to
another is what the internet is _for_."

This is currently only tested on Ubuntu Linux 20.02.

## Fastest way to use

1. Configure your router to forward a port (called $PORT) to port 80
   on your computer. Note your external address on the router (called
   $ADDR)
2. At an sh command line, run `mkdir -p ~/sendto/files`
2. Run `sendto -dir ~/sendto/files/`
3. Give the sender the url http\://$ADDR:$PORT/
4. The sender can now send you a file using that url!

## More convenient, more secure setup

1. Configure your router to forward port 80 and port 443 to your
   computer. Note your external address on the router (called $ADDR)
2. At a sh command line, run `mkdir -p ~/sendto/{files,certs}`
2. Set up a DNS entry on a domain name you own, like
   sendto.example.com
3. Use certbot from https://letsencrypt.com/ to generate an SSL
   private key and public certificate for your domain and put it in
   the certs directory.
   ```sh
   sudo certbot certonly --standalone
   sudo chown $USER {privkey,cert}.pem
   mv {privkey,cert}.pem ~/sendto/certs
   ```
4. Run sendto:
```sh
sendto -dir ~/sendto/files/ \
       -key ~/sendto/certs/privkey.pem \
       -cert ~/sendto/certs/cert.pem \
       -hostname sendto.example.com
       ```
5. The sender can now send you a file using http\://sendto.example.com/!

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
5. (optional) If you use the standard 80/443 ports, you won't be able
   to use them unless running as root.  Under Linux, give the binary
   the low-port binding capability:
```sh
sudo setcap cap_net_bind_service=+ep ~/go/bin/sendto
```

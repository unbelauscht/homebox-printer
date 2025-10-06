# homebox-printer

a simple api designed to print labels, but really prints anything you throw at it

# run

```bash
PRINTER=QL-700 BIND=:8080 SECRET=your-password go run main.go
```

# install
> you will need a working install of golang: https://go.dev/doc/install#install
> 
> you will need a configured label printer and its name

```bash
git clone https://github.com/unbelauscht/homebox-printer
cd homebox-printer
go build -o homebox-printer main.go
install -o root -g root -m 0755 homebox-printer /usr/local/bin/

adduser homebox-printer --disabled-password --disabled-login --no-create-home --system --verbose
adduser homebox-printer lpadmin
install -o root -g root -m 0644 homebox-printer.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable homebox-printer.service
systemctl start homebox-printer.service
```

> now expose it using something like tailscale

```bash
tailscale funnel --bg 8080
```

Now you should have an internet connected label printer.

To use it in homebox, add an environment variable to your Homebox container like

```yaml
HBOX_LABEL_MAKER_PRINT_COMMAND: "wget --post-file={{.FileName}} -O- https://labelprinter.your-tailscalenet.ts.net/print-your_secret_value"
```

# Wake on LAN (wol)

Utility written with Go. It is my implementation of WakeOnLan and it's magic packet


## Building

```bash
git clone https://github.com/KopyTKG/wol.git
cd wol
go build
./wol --version
```

## Usage

```bash
wol -m xx:xx:xx:xx:xx:xx
```

or

```bash
wol --mac xx:xx:xx:xx:xx:xx
```



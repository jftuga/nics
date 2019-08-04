# nics
Display information about Network Inferface Cards (NICs). This is easier to read compared to `ipconfig`, `ip`, or `ifconfig`.

Binaries for Windows, MacOS, Linux and FreeBSD can be found on the [Releases Page](https://github.com/jftuga/nics/releases) -- expand the `Assets` to see the downloads.

## Usage

```
nics: Display information about Network Inferface Cards (NICs)
usage: nics [options]
  -a	show all details on ALL interfaces
  -d	show debug information
  -v	show program version
```

## Install From Source

```
go get github.com/jftuga/nics
```

## Examples

```
C:\GitHub\nics>nics
+----------+----------------+-------------------+------+------------------------+
|   NAME   |      IPV4      |    MAC ADDRESS    | MTU  |         FLAGS          |
+----------+----------------+-------------------+------+------------------------+
| Ethernet | 172.16.7.89/24 | 98:ff:aa:cb:24:a0 | 1500 | up|broadcast|multicast |
+----------+----------------+-------------------+------+------------------------+
```

```
pi@pi8:~ $ ./nics -a
+---------+---------------+------------------------------+-------------------+-------+------------------------+
|  NAME   |     IPV4      |             IPV6             |    MAC ADDRESS    |  MTU  |         FLAGS          |
+---------+---------------+------------------------------+-------------------+-------+------------------------+
| lo      | 127.0.0.1/8   | ::1/128                      |                   | 65536 | up|loopback            |
| eth0    | 172.16.7.6/24 | fe80::51d3:4fc2:5a11:3abc/64 | b8:27:eb:b2:ea:11 |  1500 | up|broadcast|multicast |
| wlan0   |               |                              | b8:27:eb:c4:4e:2a |  1500 | up|broadcast|multicast |
| docker0 | 172.17.0.1/16 |                              | 02:42:60:1b:aa:30 |  1500 | up|broadcast|multicast |
+---------+---------------+------------------------------+-------------------+-------+------------------------+
```

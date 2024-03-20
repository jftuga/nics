# nics
Display information about Network Interface Cards (NICs). This program's output is easier to read compared to `ipconfig`, `ip`, or `ifconfig`.

## Usage

```
nics: Display information about Network Interface Cards (NICs)
usage: nics [options]
  -a	show all details on ALL interfaces
  -d	show debug information
  -v	show program version
```

## Installation

* macOS: `brew tap jftuga/homebrew-tap; brew update; brew install jftuga/tap/nics`
* Binaries for Linux, macOS and Windows are provided in the [releases](https://github.com/jftuga/nics/releases) section.

## Examples

```
C:\GitHub\nics>nics
+----------+----------------+-------------------+------+------------------------+
|   NAME   |      IPV4      |    MAC ADDRESS    | MTU  |         FLAGS          |
+----------+----------------+-------------------+------+------------------------+
| Ethernet | 172.16.7.89/24 | 98:ff:aa:cb:24:a0 | 1500 | up|broadcast|multicast |
+----------+----------------+-------------------+------+------------------------+

+------------+------------+------------+
|  GATEWAY   |    DNS1    |    DNS2    |
+------------+------------+------------+
| 172.22.7.1 | 172.16.7.2 | 172.16.7.3 |
+------------+------------+------------+
```

```
pi@raspberrypi:~ $ nics -a
+---------+---------------+------------------------------+-------------------+-------+-----------+
|  NAME   |     IPV4      |             IPV6             |    MAC ADDRESS    |  MTU  |   FLAGS   |
+---------+---------------+------------------------------+-------------------+-------+-----------+
| lo      | 127.0.0.1/8   | ::1/128                      |                   | 65536 | up        |
|         |               |                              |                   |       | loopback  |
+---------+---------------+------------------------------+-------------------+-------+-----------+
| eth0    | 172.16.7.6/24 | fe80::51d3:4fc2:5a11:3abc/64 | b8:27:eb:b2:ea:11 |  1500 | up        |
|         |               |                              |                   |       | broadcast |
|         |               |                              |                   |       | multicast |
+---------+---------------+------------------------------+-------------------+-------+-----------+
| wlan0   |               |                              | b8:27:eb:c4:4e:2a |  1500 | up        |
|         |               |                              |                   |       | broadcast |
|         |               |                              |                   |       | multicast |
+---------+---------------+------------------------------+-------------------+-------+-----------+
| docker0 | 172.17.0.1/16 |                              | 02:42:60:1b:aa:30 |  1500 | up        |
|         |               |                              |                   |       | broadcast |
|         |               |                              |                   |       | multicast |
+---------+---------------+------------------------------+-------------------+-------+-----------+

+------------+-----------+-------+
|  GATEWAY   |   DNS 1   | DNS 2 |
+------------+-----------+-------+
| 172.16.7.1 | 127.0.0.1 |       |
+------------+-----------+-------+
```

```
jftuga@debian:~$ nics -a

+---------+---------------+------------------------------+-------------------+-------+-----------+
|  NAME   |     IPV4      |             IPV6             |    MAC ADDRESS    |  MTU  |   FLAGS   |
+---------+---------------+------------------------------+-------------------+-------+-----------+
| lo      | 127.0.0.1/8   | ::1/128                      |                   | 65536 | up        |
|         |               |                              |                   |       | loopback  |
+---------+---------------+------------------------------+-------------------+-------+-----------+
| enp3s0  | 172.22.2.6/24 | fe80::51d3:4fc2:face:6b4c/64 | d4:b4:e7:aa:73:c2 |  1500 | up        |
|         |               |                              |                   |       | broadcast |
|         |               |                              |                   |       | multicast |
+---------+---------------+------------------------------+-------------------+-------+-----------+
| docker0 | 172.17.0.1/16 |                              | 02:42:60:42:af:a3 |  1500 | up        |
|         |               |                              |                   |       | broadcast |
|         |               |                              |                   |       | multicast |
+---------+---------------+------------------------------+-------------------+-------+-----------+

+------------+------------+------------+
|  GATEWAY   |   DNS 1    |   DNS 2    |
+------------+------------+------------+
| 172.22.2.1 | 172.22.2.2 | 172.22.2.3 |
+------------+------------+------------+
```


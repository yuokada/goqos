# goqos

original => [matsumotory/qos-control: qos-control.pl](https://github.com/matsumotory/qos-control "matsumotory/qos-control: qos-control.pl")
```
Traffic control tool using cbq, tc and iproute for CentOS and Ubuntu.
```

## Usage


    $ ./goqos
    Usage: goqos [global flags] <set|view|clear|version> [command flags]

    global flags:

## set method

     $ ./goqos set -h
     Usage of set:
      -direction string
          Traffic direction(in or out)
      -ip string
          Ipv4 Address (required)
      -protocol string
          Protocol(Ex. HTTP, DNS) (default "all")
      -src string
          Traffic Source
      -traffic uint
          traffic volume


## view method

    $ ./goqos view -h
    Usage of view:
     -ip string
         Ipv4 Address


## clear method

    $ ./goqos clear -h
    Usage of clear:
      -clsid string
        	Class ID (required)
      -ip string
        	Ipv4 Address (required)

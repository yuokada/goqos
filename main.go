package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

const (
	QOS_INIT     = "/etc/init.d/qos.init"
	QOS_MODULE   = "/usr/local/sbin/qos.sh"
	QOS_CONF_DIR = "/etc/sysconfig/qos"
)

type Direction struct {
	id        int
	direction string
	burstRate int
}

var DirectionIn = Direction{
	id:        1,
	direction: "in",
	burstRate: 5,
}

var DirectionOut = Direction{
	id:        5001,
	direction: "out",
	burstRate: 10,
}

var DirectionList = map[string]Direction{
	"in":  DirectionIn,
	"out": DirectionOut,
}

type Protocol struct {
	id   int
	port int
}

var PROTOCOL_LIST = map[string]Protocol{
	"all":   Protocol{id: 0, port: 0},
	"ftp":   Protocol{id: 1, port: 20},
	"ssh":   Protocol{id: 2, port: 22},
	"smtp":  Protocol{id: 3, port: 25},
	"http":  Protocol{id: 4, port: 80},
	"pop3":  Protocol{id: 5, port: 110},
	"imap":  Protocol{id: 6, port: 143},
	"https": Protocol{id: 7, port: 443},
	"imaps": Protocol{id: 8, port: 993},
	"pop3s": Protocol{id: 9, port: 995},
	"dns":   Protocol{id: 10, port: 53},
}

type ConfigVariables struct {
	Direction   string
	NWInterface string
	Traffic     int
	Weight      int
	ServerIP    string
	ServerPort  int
	ClientIP    string
}

/*
sudo -E ../qos-control/qos-control.pl -i `hostname -i` --method set --src 172.21.254.229 --direction=in --traffic=256
*** old [172.17.36.15] settings ***
--------------------------
CLSID->0016     interface->eth0       direction->in         bandwidth->64Mbps     server_ip_port>172.17.36.15         src_ip->nothing
--------------------------

Setting ... OK

*** current [172.17.36.15] settings ***
------------------------------
CLSID->0016     interface->eth0       direction->in         bandwidth->256Mbps    server_ip_port>172.17.36.15         src_ip->172.21.254.229
------------------------------
*/
func setCmd(args []string) {
	var ip_s, protocol string
	var source, direction string
	var traffic uint
	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	fs.StringVar(&ip_s, "ip", "", "Ipv4 Address (required)")
	fs.StringVar(&source, "src", "", "Traffic Source")
	fs.StringVar(&protocol, "protocol", "all", "Protocol(Ex. HTTP, DNS)")
	fs.StringVar(&direction, "direction", "", "Traffic direction(in or out)")
	fs.UintVar(&traffic, "traffic", 0, "traffic volume")
	err := fs.Parse(args[1:])
	if err != nil {
		os.Exit(1)
	}
	ip := net.ParseIP(ip_s)
	if ip == nil {
		fmt.Printf("ip: %s is invalid IPv4 Address", ip_s)
		os.Exit(1)
	}
	if source != "" {
		if sip := net.ParseIP(source); sip == nil {
			fmt.Printf("src: %s is invalid IPv4 Address", sip)
			os.Exit(1)
		}
	}

	fmt.Printf("*** old [%s] settings ***\n", ip_s)
	fmt.Println("--------------------------")
	view(ip)
	fmt.Println("--------------------------")

	cls_id, err := getClassId(ip, protocol, direction)
	if err != nil {
		log.Println(errors.Wrap(err, ""))
		return
	}
	filename := fmt.Sprintf("qos-%s.%s_%s",
		cls_id, direction, protocol)
	config_file := getFilepath(QOS_CONF_DIR, filename)
	weight := calcTrafficWeight(traffic, direction)
	r := ConfigVariables{
		Direction:   direction,
		NWInterface: "eth0",
		Traffic:     int(traffic),
		Weight:      int(weight),
		ServerIP:    ip.String(),
		ClientIP:    source,
	}
	writeConfigfile(config_file, r)

	fmt.Println("\nSetting ... \nOK\n")

	time.Sleep(1 * time.Second)
	cmd := exec.Command(QOS_INIT, "restart")
	output, err := cmd.Output()
	fmt.Print(string(output))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("*** current [%s] settings ***\n", ip_s)
	fmt.Println("--------------------------")
	view(ip)
	fmt.Println("--------------------------")
}

func writeConfigfile(config_file string, p ConfigVariables) {
	tpl := ""
	tpl += "DIRECTION={{.Direction}}\n"
	tpl += "DEVICE={{.NWInterface}},100Mbit,10Mbit\n"
	tpl += "RATE={{.Traffic}}Mbit\n"
	tpl += "WEIGHT={{.Weight}}Mbit\n"
	tpl += "PRIO=5\n"
	tpl += "RULE={{ .ServerIP }}{{ .ServerPort }},{{ .ClientIP }}\n"

	t := template.Must(template.New("configfile").Parse(tpl))
	f, err := os.OpenFile(config_file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println("os.Create: ", errors.WithStack(err))
		return
	}
	defer f.Close()
	err = t.Execute(f, p)
	if err != nil {
		log.Println("executing template:", errors.WithStack(err))
	}
}

func calcTrafficWeight(traffic uint, direction string) uint {
	burst_rate := DirectionList[direction].burstRate
	return traffic / uint(burst_rate)
}

/*
$ ./qos-control.pl --method view --ip `hostname -i`
CLSID->0016 interface->eth0       direction->in         bandwidth->64Mbps     server_ip_port>172.17.36.15         src_ip->nothing
*/
func viewCmd(args []string) {
	var ip_s string
	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	fs.StringVar(&ip_s, "ip", "", "Ipv4 Address")
	err := fs.Parse(args[1:])
	if err != nil {
		os.Exit(1)
	}
	if ip_s == "" {
		fmt.Println("ip is empty.")
		// TODO: う〜ん、不親切なUsage
		fs.Usage()
		os.Exit(1)
	}
	ip := net.ParseIP(ip_s)
	if ip == nil {
		fmt.Printf("%s is invalid IPv4 Address\n", ip_s)
		os.Exit(1)
	}
	view(ip)
	os.Exit(0)
}

func view(ip net.IP) {
	outputs := []string{}
	for protocol, _ := range PROTOCOL_LIST {
		for direction, _ := range DirectionList {
			cls_id, err := getClassId(ip, protocol, direction)
			if err != nil {
				// TODO: Write this later.
				continue
			}
			filename := fmt.Sprintf("qos-%s.%s_%s",
				cls_id,
				direction,
				protocol,
			)
			config_file := getFilepath(QOS_CONF_DIR, filename)
			_, err = os.Stat(config_file)
			if err != nil && !os.IsExist(err) {
				continue
			}
			contents, err := ioutil.ReadFile(config_file)
			if err != nil {
				continue
			}
			m, err := parseQosConfigfile(contents)
			if err != nil {
				log.Println(config_file, err)
				continue
			}
			output := formatOutput(cls_id, m)
			outputs = append(outputs, output)
		}
	}
	if len(outputs) != 0 {
		fmt.Printf(strings.Join(outputs, ""))
	} else {
		fmt.Println("nothing")
	}
}

func clearCmd(args []string) {
	var ip_s, cls_id string
	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	fs.StringVar(&ip_s, "ip", "", "Ipv4 Address (required)")
	fs.StringVar(&cls_id, "clsid", "", "Class ID (required)")
	err := fs.Parse(args[1:])
	if err != nil {
		os.Exit(1)
	}
	if ip_s == "" {
		fmt.Println("ip is empty")
		fs.Usage()
		os.Exit(1)
	}
	ip := net.ParseIP(ip_s)
	if ip == nil {
		fmt.Printf("%s is invalid IPv4 Address", ip_s)
		os.Exit(1)
	}
	if cls_id == "" {
		fmt.Println("clsid is empty")
		os.Exit(1)
	}

	fmt.Printf("*** old [%s] settings ***\n", ip_s)
	fmt.Println("--------------------------")
	view(ip)
	fmt.Println("--------------------------")

	filename := fmt.Sprintf("qos-%s.*", cls_id)
	config_files, err := filepath.Glob(getFilepath(QOS_CONF_DIR, filename))
	for _, f := range config_files {
		err = os.Remove(f)
		if err != nil {
			log.Println(err)
			fmt.Println("Cant clear settings")
			os.Exit(1)
		}
	}
	cmd := exec.Command(QOS_INIT, "restart")
	output, err := cmd.Output()
	fmt.Print(string(output))
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("\nSetting Clear ... \nOK\n")
	fmt.Printf("*** current [%s] settings ***\n", ip_s)
	fmt.Println("--------------------------")
	view(ip)
	fmt.Println("--------------------------")
}

func formatOutput(cls_id string, m map[string]interface{}) string {
	return fmt.Sprintf("CLSID->%-10s interface->%-10s direction->%-10s bandwidth->%-10s server_ip_port>%-20s src_ip->%-10s\n",
		cls_id,
		m["interface"],
		m["direction"],
		m["bandwidth"],
		m["server_ip_port"],
		m["src_ip"])
}
// getClassId is convert class(number) from ip , protocol and direct
func getClassId(ip net.IP, protocol, direct string) (string, error) {
	d_num := DirectionList[direct].id
	p_num := PROTOCOL_LIST[protocol].id
	i_num, err := strconv.Atoi(strings.Split(ip.String(), ".")[3])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d", i_num+(255*p_num)+d_num), nil
}

// getFilepath is a function like a File::Spec::catfile
func getFilepath(args ...string) string {
	sep := os.PathSeparator
	configfile := strings.Join(args, string(sep))
	return configfile
}

func parseQosConfigfile(contents []byte) (map[string]interface{}, error) {
	// TODO: べた書き修正したい
	//m := map[string]interface{}
	m := map[string]interface{}{}
	for _, ln := range strings.Split(string(contents), "\n") {
		splits := strings.SplitN(ln, "=", 2)
		if len(splits) != 2 {
			continue
		}
		//log.Printf("%#v\n", splits)
		k, v := splits[0], splits[1]
		switch k {
		case "DEVICE":
			m["interface"] = strings.Split(v, ",")[0]
		case "RATE":
			m["bandwidth"] = strings.Replace(v, "Mbit", "Mbps", -1)
		case "RULE":
			v_list := strings.SplitN(v, ",", 2)
			dst, src := v_list[0], v_list[1]
			if src == "" {
				src = "nothing"
			}
			m["server_ip_port"] = dst
			m["src_ip"] = src
		case "DIRECTION":
			m["direction"] = v
		}
	}
	return m, nil
}

func versionCmd(args []string) {
	fs := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: goqos [global flags] <TODO> [command flags]")
	}
	fs.Usage()
	os.Exit(0)
}

func init() {
	debug_s := os.Getenv("DEBUG")
	if debug_s == "true" || debug_s == "1" {
		debug = true
	} else {
		debug = false
	}
	flag.Usage = func() {
		fmt.Println("Usage: goqos [global flags] <set|view|clear|version> [command flags]")
		fmt.Printf("\nglobal flags:\n")
	}
	flag.Parse()
}

func main() {
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	switch flag.Args()[0] {
	case "set":
		setCmd(flag.Args())
	case "view":
		viewCmd(flag.Args())
	case "clear":
		clearCmd(flag.Args())
	case "version":
		versionCmd(flag.Args())
	default:
		flag.Usage()
		os.Exit(1)
	}
}

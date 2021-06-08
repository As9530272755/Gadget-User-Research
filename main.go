package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var (
	Uname          = `uname -i`
	OsVersion      = `docker info -f '{{.OperatingSystem}}'`
	kernel         = `uname -r`
	DcVersion      = `docker info | grep -i "server version" |awk -F: '{print $2}'`
	K8SNode        = `kubectl get node  | sed -n '1!p' | wc -l`
	LinkNetr       = `ping baidu.com -w1`
	ContainerdV    = ` containerd -v`
	CRIv           = `crio -v `
	DockerComposeV = `docker-compose version`
	K8Sversion     = `kubectl version --short `
	OKDVersion     = `openshift version`
	K8Sdeployed    = `kubectl get pod -n kube-system | grep kube`
)

// 调用 Linux 命令
func Cmd(CMD string) string {
	command := exec.Command("/bin/bash", "-c", CMD)
	stdout, err := command.StdoutPipe()
	if err != nil {
		err1 := fmt.Sprintln(err)
		return err1
	}

	if err := command.Start(); err != nil {
		err1 := fmt.Sprintln("Error:The command is err,", err)
		return err1
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		err1 := fmt.Sprintln("ReadAll Stdout:", err.Error())
		return err1
	}

	fmt := string(bytes)
	return fmt
}

func ExcelAndComd() {
	// 传入 Linux 命令变量
	Uname := Cmd(Uname)
	OsVersion := Cmd(OsVersion)
	kernel := Cmd(kernel)
	DcVersion := Cmd(DcVersion)
	K8SNode := Cmd(K8SNode)
	LinkNetr := Cmd(LinkNetr)
	ContainerdV := Cmd(ContainerdV)
	CRIv := Cmd(CRIv)
	DockerComposeV := Cmd(DockerComposeV)
	K8Sversion := Cmd(K8Sversion)
	OKDVersion := Cmd(OKDVersion)
	K8Sdeployed := Cmd(K8Sdeployed)

	// 第一列获取信息列表
	categories := map[string]string{
		"A1": "主机环境：", "B2": "CPU 架构：", "B3": "操作系统类型及版本", "B4": "内核版本",

		"A6": "集群容器环境：", "B7": "集群规模节点数", "B8": "是否连互联网", "B9": "Docker 版本",
		"B10": "Containerd 版本", "B11": "CRI-O 版本信息", "B12": "Docker-Compose 版本",
		"B13": "Kubernetes 版本", "B14": "OpenShift 版本", "B15": "Rancher 版本", "B16": "其他集群编排软件",
		"B17": "网络插件名称及版本", "B18": "Docker 存储方式", "B19": "集群部署方式", "B20": "运行的业务",
		"A22": "环境连接方式：", "B23": "连接方式", "B24": "Kubernetes 主节点 IP",
		"A26": "镜像仓库环境：", "B27": "镜像仓库软件及版本", "B28": "仓库内镜像数量",
	}

	// 获取信息
	info := map[string]string{
		"E2": Uname, "E3": OsVersion, "E4": kernel, "E7": K8SNode, "E8": LinkNetr, "E9": DcVersion,
		"E10": ContainerdV, "E11": CRIv, "E12": DockerComposeV, "E13": K8Sversion, "E14": OKDVersion,
		"E15": "如有 Rancher 请填写对应版本", "E16": "如有其他集群编排软件请填写对应版本", "E17": "请填写对应的网络插件名称及版本",
		"E18": "填写对应的存储方式如: GlusterFS, CephFS, NFS, Local", "E19": "",
		"E20": "请填写如: 登录服务, 注册服务", "E23": "请填写对应的登录方式如：直接 SSH  通过 VPN 或者 4A 验证等等", "E24": "如：需要部署集群主节点IP（master IP）",
		"E27": "填写使用的镜像仓库的类型及版本（如 Harbor 1.8.0）", "E28": "仓库内镜像的数量",
	}

	if K8Sdeployed != "" {
		info["E19"] = "容器化部署"
	} else {
		info["E19"] = "二进制部署"
	}

	f := excelize.NewFile()
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}

	for k, v := range info {
		f.SetCellValue("Sheet1", k, v)
	}

	if err := f.SaveAs("用户环境信息调研.xlsx"); err != nil {
		fmt.Println(err)
	}
}

func Flag() {
	var help bool
	var start bool
	var quit bool
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&start, "s", false, "Start the program")
	flag.BoolVar(&quit, "q", false, "quit the program")
	flag.Parse()
	flag.Usage = func() {
		fmt.Println("执行程序：Research -s")
		flag.PrintDefaults()
	}

	if start {
		ExcelAndComd()
		return
	}

	if help {
		flag.Usage()
		return
	}

	if quit {
		fmt.Println("退出程序中...")
		return
	}

	flag.PrintDefaults()
}

func main() {
	Flag()
}

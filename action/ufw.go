package action

import "os/exec"

func Ban(ip string, port string) error {
	//调用系统UFW禁止IP入站22端口
	exec.Command("ufw", "deny", "from", ip, "to", "any", "port", port).Run()
	return nil
}

func Unban(ip string, port string) error {
	//调用系统UFW解除禁止IP入站22端口
	exec.Command("ufw", "delete", "deny", "from", ip, "to", "any", "port", port).Run()
	return nil
}

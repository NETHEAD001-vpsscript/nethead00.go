package main

import (
    "bufio"
    "crypto/rand"
    "fmt"
    "log"
    "os"
    "strings"

    "golang.org/x/crypto/ssh"
)

var host = "YOUR_VPS_IP"
var port = "22"
var user = "root"

func sshRun(cmd string) string {
    keyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
    key, err := os.ReadFile(keyPath)
    if err != nil {
        log.Fatal(err)
    }
    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        log.Fatal(err)
    }
    config := &ssh.ClientConfig{
        User:            user,
        Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    conn, err := ssh.Dial("tcp", host+":"+port, config)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    session, err := conn.NewSession()
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    out, err := session.CombinedOutput(cmd)
    if err != nil {
        fmt.Println("Command error:", err)
    }
    return string(out)
}

func uuid() string {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }
    return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func main() {
    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Println("\n===== NETHEAD VPN PANEL =====")
        fmt.Println("1. VPS Status")
        fmt.Println("2. Update VPS")
        fmt.Println("3. Restart Xray")
        fmt.Println("4. Install Tools (jq/curl)")
        fmt.Println("5. Custom Command")
        fmt.Println("6. Create VLESS User")
        fmt.Println("7. Create VMess User")
        fmt.Println("8. Create Trojan User")
        fmt.Println("9. Install Hysteria2")
        fmt.Println("10. Enable UDP Note")
        fmt.Println("11. Install DNS (dnsmasq)")
        fmt.Println("0. Exit")
        fmt.Print("Select: ")

        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        switch input {
        case "1":
            fmt.Println(sshRun("uptime && df -h"))
        case "2":
            fmt.Println(sshRun("apt update -y && apt upgrade -y"))
        case "3":
            fmt.Println(sshRun("systemctl restart xray"))
        case "4":
            fmt.Println(sshRun("apt install jq curl unzip -y"))
        case "5":
            fmt.Print("Command: ")
            cmd, _ := reader.ReadString('\n')
            cmd = strings.TrimSpace(cmd)
            fmt.Println(sshRun(cmd))
        case "6":
            id := uuid()
            cmd := fmt.Sprintf(`jq '.inbounds[0].settings.clients += [{"id":"%s","email":"vless-user"}]' /etc/xray/config.json > /tmp/v.json && mv /tmp/v.json /etc/xray/config.json && systemctl restart xray && echo "VLESS CREATED: %s"`, id, id)
            fmt.Println(sshRun(cmd))
        case "7":
            id := uuid()
            cmd := fmt.Sprintf(`jq '.inbounds[0].settings.clients += [{"id":"%s","email":"vmess-user"}]' /etc/xray/config.json > /tmp/m.json && mv /tmp/m.json /etc/xray/config.json && systemctl restart xray && echo "VMESS CREATED: %s"`, id, id)
            fmt.Println(sshRun(cmd))
        case "8":
            fmt.Print("Enter password: ")
            pass, _ := reader.ReadString('\n')
            pass = strings.TrimSpace(pass)
            cmd := fmt.Sprintf(`jq '.inbounds[0].settings.clients += [{"password":"%s","email":"trojan-user"}]' /etc/xray/config.json > /tmp/t.json && mv /tmp/t.json /etc/xray/config.json && systemctl restart xray && echo "TROJAN CREATED"`, pass)
            fmt.Println(sshRun(cmd))
        case "9":
            fmt.Println(sshRun("curl -fsSL https://get.hy2.sh | bash"))
        case "10":
            fmt.Println(sshRun("echo 'UDP note: requires badvpn or udp-relay setup'"))
        case "11":
            fmt.Println(sshRun("apt install dnsmasq -y && systemctl enable dnsmasq && systemctl restart dnsmasq"))
        case "0":
            fmt.Println("Exiting...")
            return
        default:
            fmt.Println("Invalid option")
        }
    }
}

package main

import (
    "bufio"
    "crypto/rand"
    "fmt"
    "log"
    "os"
    "strings"

    "golang.org/x/crypto/ssh"
    "golang.org/x/term"
)

var host = "172.104.233.204" // Your remote server IP
var port = "22"
var user = "root" // SSH login user

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

func readPassword(prompt string) string {
    fmt.Print(prompt)
    bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println()
    return strings.TrimSpace(string(bytePassword))
}

func generatePassword() string {
    b := make([]byte, 12)
    _, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }
    return fmt.Sprintf("%x", b)[:12]
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
        fmt.Println("12. Create Server")
        fmt.Println("0. Exit")
        fmt.Print("Select: ")

        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        switch input {
        case "1":
            fmt.Println(sshRun("uptime && df -h"))

        case "2":
            fmt.Println("Updating VPS...")
            fmt.Println(sshRun("apt update -y && apt upgrade -y"))

        case "3":
            fmt.Println("Restarting Xray...")
            fmt.Println(sshRun("systemctl restart xray"))

        case "4":
            fmt.Println("Installing tools...")
            fmt.Println(sshRun("apt install jq curl unzip -y"))

        case "5":
            fmt.Print("Enter command to run: ")
            cmd, _ := reader.ReadString('\n')
            cmd = strings.TrimSpace(cmd)
            fmt.Println(sshRun(cmd))

        case "6":
            id := uuid()
            fmt.Println("Generated UUID:", id)
            fmt.Print("Enter username: ")
            username, _ := reader.ReadString('\n')
            username = strings.TrimSpace(username)
            fmt.Print("Enter expiry date (YYYY-MM-DD): ")
            expiry, _ := reader.ReadString('\n')
            expiry = strings.TrimSpace(expiry)
            cmd := fmt.Sprintf(`jq '.inbounds[0].settings.clients += [{"id":"%s","email":"%s","expiry":"%s"}]' /etc/xray/config.json > /tmp/v.json && mv /tmp/v.json /etc/xray/config.json && systemctl restart xray && echo "VLESS CREATED: %s"`, id, username, expiry, id)
            fmt.Println(sshRun(cmd))

        case "7":
            id := uuid()
            fmt.Println("Generated UUID:", id)
            fmt.Print("Enter username: ")
            username, _ := reader.ReadString('\n')
            username = strings.TrimSpace(username)
            fmt.Print("Enter expiry date (YYYY-MM-DD): ")
            expiry, _ := reader.ReadString('\n')
            expiry = strings.TrimSpace(expiry)
            cmd := fmt.Sprintf(`jq '.inbounds[0].settings.clients += [{"id":"%s","email":"%s","expiry":"%s"}]' /etc/xray/config.json > /tmp/m.json && mv /tmp/m.json /etc/xray/config.json && systemctl restart xray && echo "VMESS CREATED: %s"`, id, username, expiry, id)
            fmt.Println(sshRun(cmd))

        case "8":
            fmt.Print("Enter Trojan password: ")
            pass := readPassword("Password: ")
            fmt.Print("Enter username for Trojan: ")
            username, _ := reader.ReadString('\n')
            username = strings.TrimSpace(username)
            fmt.Print("Enter expiry date (YYYY-MM-DD): ")
            expiry, _ := reader.ReadString('\n')
            expiry = strings.TrimSpace(expiry)
            cmd := fmt.Sprintf(`jq '.inbounds[0].settings.clients += [{"password":"%s","email":"%s","expiry":"%s"}]' /etc/xray/config.json > /tmp/t.json && mv /tmp/t.json /etc/xray/config.json && systemctl restart xray && echo "TROJAN CREATED"`, pass, username, expiry)
            fmt.Println(sshRun(cmd))

        case "9":
            fmt.Println("Installing Hysteria2...")
            fmt.Println(sshRun("curl -fsSL https://get.hy2.sh | bash"))

        case "10":
            fmt.Println("Note about UDP:")
            fmt.Println("UDP note: requires badvpn or udp-relay setup.")

        case "11":
            fmt.Println("Installing dnsmasq...")
            fmt.Println(sshRun("apt install dnsmasq -y && systemctl enable dnsmasq && systemctl restart dnsmasq"))

        case "12":
            // Create Server
            fmt.Println("=== Create New Server ===")
            fmt.Print("Enter server name: ")
            serverName, _ := reader.ReadString('\n')
            serverName = strings.TrimSpace(serverName)

            fmt.Print("Enter server IP: ")
            serverIP, _ := reader.ReadString('\n')
            serverIP = strings.TrimSpace(serverIP)

            // Collect username, password, expiry for admin user
            fmt.Print("Enter admin username: ")
            adminUser, _ := reader.ReadString('\n')
            adminUser = strings.TrimSpace(adminUser)

            fmt.Print("Enter admin password (leave blank to generate): ")
            adminPass := readPassword("Password (leave blank to generate): ")
            if adminPass == "" {
                adminPass = generatePassword()
                fmt.Printf("Generated password: %s\n", adminPass)
            }

            fmt.Print("Enter expiry date for admin user (YYYY-MM-DD): ")
            expiry, _ := reader.ReadString('\n')
            expiry = strings.TrimSpace(expiry)

            // Create user with provided info
            createUserCmd := fmt.Sprintf(
                "useradd -m -s /bin/bash %s && echo '%s:%s' | chpasswd && chage -E %s %s",
                adminUser, adminUser, adminPass, expiry, adminUser,
            )

            fmt.Println("Executing server setup commands...")
            output := sshRun(createUserCmd)
            fmt.Println(output)

            fmt.Printf("Server '%s' with IP '%s' created.\n", serverName, serverIP)
            fmt.Printf("Admin user '%s' with expiry '%s' and password '%s' created.\n", adminUser, expiry, adminPass)

        case "0":
            fmt.Println("Exiting...")
            return
        default:
            fmt.Println("Invalid option, please try again.")
        }
    }
}

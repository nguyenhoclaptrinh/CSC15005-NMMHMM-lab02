package main

import (
    "fmt"
    "secure_notes/internal/client/crypto"
    "secure_notes/internal/client/utils"
    "os"
)

func main() {
    // Ví dụ CLI menu đơn giản
    for {
        fmt.Println("\nSecure Notes Client")
        fmt.Println("1. Register")
        fmt.Println("2. Login")
        fmt.Println("3. Upload Note")
        fmt.Println("4. List Notes")
        fmt.Println("5. Share Note")
        fmt.Println("6. Create Temp URL")
        fmt.Println("0. Exit")
        fmt.Print("Choose option: ")

        var choice int
        fmt.Scanln(&choice)

        switch choice {
        case 1:
            // Call register function
            utils.LogInfo("Register selected")
        case 2:
            // Call login function
            utils.LogInfo("Login selected")
        case 3:
            // Call upload note function
            utils.LogInfo("Upload note selected")
        case 4:
            // Call list notes function
            utils.LogInfo("List notes selected")
        case 5:
            // Call share note function
            utils.LogInfo("Share note selected")
        case 6:
            // Call create temp URL function
            utils.LogInfo("Create temp URL selected")
        case 0:
            os.Exit(0)
        default:
            fmt.Println("Invalid option")
        }
    }
}

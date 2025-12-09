package main

import (
	"fmt"
	"os"
	clientinternal "secure-notes-client/pkg"
)

func main() {
	// CLI menu with separate states for logged-out and logged-in
	loggedIn := clientinternal.IsLoggedIn()

	for {
		fmt.Println("\nSecure Notes Client")
		if !loggedIn {
			fmt.Println("1. Register")
			fmt.Println("2. Login")
			fmt.Println("0. Exit")
			fmt.Print("Choose option: ")

			var choice int
			fmt.Scanln(&choice)

			switch choice {
			case 1:
				clientinternal.Register()
				clientinternal.LogInfo("Register selected")
			case 2:
				clientinternal.Login()
				clientinternal.LogInfo("Login selected")
				loggedIn = clientinternal.IsLoggedIn()
			case 0:
				os.Exit(0)
			default:
				fmt.Println("Invalid option")
			}
		} else {
			fmt.Println("3. Upload Note")
			fmt.Println("4. List Notes")
			fmt.Println("5. Share Note")
			fmt.Println("6. Create Temp URL")
			fmt.Println("7. Logout")
			fmt.Println("0. Exit")
			fmt.Print("Choose option: ")

			var choice int
			fmt.Scanln(&choice)

			switch choice {
			case 3:
				clientinternal.UploadNote()
				clientinternal.LogInfo("Upload note selected")
			case 4:
				clientinternal.ListNotes()
				clientinternal.LogInfo("List notes selected")
			case 5:
				clientinternal.ShareNote()
				clientinternal.LogInfo("Share note selected")
			case 6:
				clientinternal.CreateTempURL()
				clientinternal.LogInfo("Create temp URL selected")
			case 7:
				clientinternal.Logout()
				clientinternal.LogInfo("Logout selected")
				loggedIn = clientinternal.IsLoggedIn()
			case 0:
				os.Exit(0)
			default:
				fmt.Println("Invalid option")
			}
		}
	}
}

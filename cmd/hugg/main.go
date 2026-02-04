package main

import (
	"fmt"
	"log"
	"os"

	"syscall"

	"github.com/nathfavour/huggkey/internal/storage"
	"github.com/nathfavour/huggkey/pkg/identity"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "hugg",
	Short: "Huggkey - The Sovereign Identity & Transport Mesh",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new HuggID identity",
	Run: func(cmd *cobra.Command, args []string) {
		vaultPath, err := storage.GetDefaultVaultPath()
		if err != nil {
			log.Fatalf("Failed to get vault path: %v", err)
		}

		if _, err := os.Stat(vaultPath); err == nil {
			fmt.Println("An identity already exists at", vaultPath)
			fmt.Print("Overwrite? (y/N): ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "Y" {
				fmt.Println("Initialization cancelled.")
				return
			}
		}

		fmt.Print("Enter passphrase to encrypt your vault: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Failed to read password: %v", err)
		}
		password := string(bytePassword)
		fmt.Println()

		mnemonic, id, err := identity.CreateNewIdentity()
		if err != nil {
			log.Fatalf("Failed to create identity: %v", err)
		}

		if err := storage.SaveVault(vaultPath, password, mnemonic); err != nil {
			log.Fatalf("Failed to save vault: %v", err)
		}

		fmt.Println("\nSuccessfully generated and saved a new HuggID!")
		fmt.Printf("HuggID: %s\n", id.HuggID)
		fmt.Printf("Vault saved to: %s\n", vaultPath)
		fmt.Println("\nIMPORTANT: Store your 24-word seed phrase securely. It is the only way to recover your identity.")
		fmt.Printf("\nSeed Phrase:\n%s\n", mnemonic)
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
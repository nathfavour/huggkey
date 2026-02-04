package main

import (
	"bufio"
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

var (
	passphrase string
)

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

		password := passphrase
		if password == "" {
			fmt.Print("Enter passphrase to encrypt your vault: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf("Failed to read password: %v", err)
			}
			password = string(bytePassword)
			fmt.Println()
		}

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

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display the current HuggID",
	Run: func(cmd *cobra.Command, args []string) {
		vaultPath, err := storage.GetDefaultVaultPath()
		if err != nil {
			log.Fatalf("Failed to get vault path: %v", err)
		}

		if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
			fmt.Println("No identity found. Run 'hugg init' to create one.")
			return
		}

		password := passphrase
		if password == "" {
			fmt.Print("Enter passphrase to unlock vault: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf("Failed to read password: %v", err)
			}
			password = string(bytePassword)
			fmt.Println()
		}

		mnemonic, err := storage.LoadVault(vaultPath, password)
		if err != nil {
			log.Fatalf("Failed to unlock vault: %v", err)
		}

		id, err := identity.FromMnemonic(mnemonic)
		if err != nil {
			log.Fatalf("Failed to derive identity: %v", err)
		}

		fmt.Printf("HuggID: %s\n", id.HuggID)
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore an identity from a seed phrase",
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
				fmt.Println("Restore cancelled.")
				return
			}
		}

		fmt.Println("Enter your 24-word seed phrase:")
		var mnemonic string
		// Read full mnemonic (could be multiple lines or spaces)
		// For simplicity in CLI, we'll ask for it in one go or word by word.
		// Using a simple Scanln for now, but in a real CLI we might want something better.
		// Let's use a scanner to read the full line.
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			mnemonic = scanner.Text()
		}

		if !identity.IsValidMnemonic(mnemonic) {
			log.Fatal("Invalid mnemonic phrase.")
		}

		fmt.Print("Enter passphrase to encrypt your new vault: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Failed to read password: %v", err)
		}
		password := string(bytePassword)
		fmt.Println()

		id, err := identity.FromMnemonic(mnemonic)
		if err != nil {
			log.Fatalf("Failed to derive identity: %v", err)
		}

		if err := storage.SaveVault(vaultPath, password, mnemonic); err != nil {
			log.Fatalf("Failed to save vault: %v", err)
		}

		fmt.Println("\nIdentity successfully restored!")
		fmt.Printf("HuggID: %s\n", id.HuggID)
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Show the seed phrase for the current identity",
	Run: func(cmd *cobra.Command, args []string) {
		vaultPath, err := storage.GetDefaultVaultPath()
		if err != nil {
			log.Fatalf("Failed to get vault path: %v", err)
		}

		if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
			fmt.Println("No identity found.")
			return
		}

		password := passphrase
		if password == "" {
			fmt.Print("Enter passphrase to unlock vault: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf("Failed to read password: %v", err)
			}
			password = string(bytePassword)
			fmt.Println()
		}

		mnemonic, err := storage.LoadVault(vaultPath, password)
		if err != nil {
			log.Fatalf("Failed to unlock vault: %v", err)
		}

		fmt.Println("Your 24-word seed phrase:")
		fmt.Println(mnemonic)
	},
}

var burnCmd = &cobra.Command{
	Use:   "burn",
	Short: "Permanently delete the local identity",
	Run: func(cmd *cobra.Command, args []string) {
		vaultPath, err := storage.GetDefaultVaultPath()
		if err != nil {
			log.Fatalf("Failed to get vault path: %v", err)
		}

		if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
			fmt.Println("No identity found.")
			return
		}

		fmt.Printf("WARNING: This will permanently delete your local identity vault at %s\n", vaultPath)
		fmt.Print("Are you sure you want to proceed? (type 'DELETE'): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "DELETE" {
			fmt.Println("Burn cancelled.")
			return
		}

		if err := os.Remove(vaultPath); err != nil {
			log.Fatalf("Failed to delete vault: %v", err)
		}

		fmt.Println("Identity successfully burned.")
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase for vault encryption/decryption")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(burnCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
/*
Copyright © 2023 Carter Ewen
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var date string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jrnl",
	Short: "simple journaling cli tool",
	Long:  `A simple journaling cli tool that allows you to create, read, update and delete journal entries all from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {

		vimStatus, _ := cmd.Flags().GetBool("vim")
		overwriteStatus, _ := cmd.Flags().GetBool("overwrite")

		if !validateDate(date) {
			fmt.Println("Invalid date format. Please use the format: mm-dd-yyyy")
		} else {
			run(date, args, vimStatus, overwriteStatus)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jrnl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&date, "date", "d", time.Now().Format("01-02-2006"), "date of journal entry to be edjited/added to")
	rootCmd.PersistentFlags().BoolP("vim", "v", false, "open vim to edit journal entry")
	rootCmd.PersistentFlags().BoolP("overwrite", "o", false, "delete an entry specified by using —date=MM-DD-YYYY (or -d=MM-DD-YYYY) and leaving your entry blank or start typing and your new entry will replace the old one")
}

func validateDate(dateString string) bool {

	layout := "01-02-2006"
	fmt.Println(dateString)

	_, err := time.Parse(layout, dateString)
	if err != nil {
		fmt.Println("Invalid date format:", err)
		return false
	} else {
		fmt.Println("Valid date format")
		return true
	}
}

func vimOpen(file string) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	absPath, error := filepath.Abs(dirname + "/.jrnl/" + file + ".txt")
	if error != nil {
		fmt.Println(error)
	}
	cmd := exec.Command("vim", absPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	erroring := cmd.Run()
	if erroring != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}

func run(date string, input []string, vim bool, overwrite bool) {
	if vim {
		vimOpen(date)
	} else {
		dailyJounal(input, date, overwrite)
	}
}

func dailyJounal(input []string, date string, overwrite bool) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	absPath, error := filepath.Abs(dirname + "/.jrnl/" + date + ".txt")
	if error != nil {
		fmt.Println(error)
	}

	if overwrite {
		if len(input) == 0 {
			if _, err := os.Stat(absPath); err == nil {
				os.Remove(absPath)
			}
		} else {
			if _, err := os.Stat(absPath); err == nil {
				os.Truncate(absPath, 0)
			}
			file, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
			} else {
				file.WriteString(strings.Join(input, " "))
			}
			file.Close()
		}
	} else {
		file, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
		} else {
			file.WriteString(strings.Join(input, " "))
		}
		file.Close()
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"

	"github.com/jasonlvhit/gocron"
)

func main() {
	fmt.Println("Starting jobs.")

	containerName := "postgres-pgbackrest"
	stanzaName := "pgbackrest_test_stanza_name"
	logLevel := "info"
	backupCommand := fmt.Sprintf("docker exec -u postgres -t %s pgbackrest backup --stanza='%s' --log-level-console=%s", containerName, stanzaName, logLevel)

	gocron.Every(60).Second().From(gocron.NextTick()).Do(executeMultiLineCmd, backupCommand)

	infoCommand := fmt.Sprintf("docker exec -u postgres -t %s pgbackrest info --stanza='%s' --log-level-console=%s", containerName, stanzaName, "warn")
	gocron.Every(15).Second().From(gocron.NextTick()).Do(executeMultiLineCmd, infoCommand)

	fullBackupCommand := fmt.Sprintf("docker exec -u postgres -t %s pgbackrest backup --type=full --stanza='%s' --log-level-console=%s", containerName, stanzaName, logLevel)

	gocron.Every(5*60).Second().Do(executeMultiLineCmd, fullBackupCommand)
	// Begin job immediately upon start
	// gocron.Every(10).Second().From(gocron.NextTick()).Do(cmdRun)
	// Do jobs with params
	// gocron.Every(1).Second().Do(taskWithParams, 1, "hello")

	// Begin job at a specific date/time
	// t := time.Date(2019, time.November, 10, 15, 0, 0, 0, time.Local)
	// gocron.Every(1).Hour().From(&t).Do(task)

	// NextRun gets the next running time
	// _, time := gocron.NextRun()
	// fmt.Println(time)

	// Remove a specific job
	// gocron.Remove(task)

	// Clear all scheduled jobs
	// gocron.Clear()

	// Start all the pending jobs
	<-gocron.Start()

}

func cmdRun() {
	fmt.Println("Executing cmdRun pgBackup..")
	// cmdStr := "docker exec -it postgres bin/bash"
	cmdStr := "docker exec -u postgres -t postgres-pgbackrest pgbackrest info --stanza='pgbackrest_test_stanza_name' --log-level-console=info"
	out, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	fmt.Printf("Command OUtput: %s, %s\n", out, err)
}

func executeMultiLineCmd(cmdStr string) {

	fmt.Println("-----------------------Scheduled command start.-----------------------")
	fmt.Println(cmdStr)
	// func runPgBackrest(containerNAme string, stanzaName string, logLevel string) {
	// fmt.Println("\nExecuting longCmdFullOpWait pgBackup..")
	// // docker exec -u postgres -t postgres-pgbackrest pgbackrest info --stanza='pgbackrest_test_stanza_name' --log-level-console=info
	// cmdStr := fmt.Sprintf("docker exec -u postgres -t %s pgbackrest backup --stanza='%s' --log-level-console=%s", containerNAme, stanzaName, logLevel)
	// The command you want to run along with the argument
	cmd := exec.Command("/bin/sh", "-c", cmdStr)

	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan struct{})

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	// Use the scanner to scan the output line by line and log it
	// It's running in a goroutine so that it doesn't block
	go func() {

		// Read line by line and process it
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}

		// We're all done, unblock the channel
		done <- struct{}{}

	}()

	// Start the command and check for errors
	// Start starts the specified command but does not wait for it to complete.
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	// The Wait method will return the exit code and release associated resources once the command exits.
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Couldnt complete the command successfully ", err)
		// log.Fatal(err)
	}
	fmt.Printf("-----------------------Scheduled command end.-----------------------\n\n")
}

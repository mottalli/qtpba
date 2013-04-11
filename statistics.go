package main

import (
    "time"
    "os/exec"
    "strconv"
)

func runStatsDaemon() {
    retryCount := 3
    for {
        // Me fijo las últimas horas
        timestamp := time.Now().Add(-2 * time.Hour).UTC().Unix()
        cmd := exec.Command("/home/marcelo/Documents/Programacion/qtpba/process_data.R", strconv.FormatInt(timestamp, 10))
        logger.Println("Running statistics with timestamp", timestamp, "...")

        if err := cmd.Run(); err != nil {
            logger.Println("Error running stats:", err)
            if retryCount--; retryCount > 0 {
                logger.Println("Retrying in 5 seconds...")
                time.Sleep(5 * time.Second)
                continue
            } else {
                logger.Println("Giving up this run.")
            }
        } else {
            logger.Println("Successfully ran stats")
        }

        retryCount = 3
        time.Sleep(30 * time.Minute)
    }
}

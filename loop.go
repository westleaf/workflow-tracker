package main

import (
	"fmt"

	"github.com/westleaf/workflow-tracker/internal/config"
)

func watchWorkflows(c *config.Config) {
	fmt.Println("Watching workflows with config:", c)

	for {
		err := checkWorkflows(c)
		if err != nil {
			fmt.Println("Error checking workflows:", err)
		}

		// Sleep or wait for a certain interval before checking again
	}

		// Sleep or wait for a certain interval before checking again
	}
}

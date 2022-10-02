package middlewares

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
	"treehole_next/config"
)

func Logger(c *fiber.Ctx) (err error) {
	errHandler := c.App().ErrorHandler
	var start, stop time.Time

	// Set latency start time
	if config.Debug {
		start = time.Now()
	}

	// Handle request, store err for logging
	chainErr := c.Next()

	// Manually call error handler
	if chainErr != nil {
		if err := errHandler(c, chainErr); err != nil {
			_ = c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// Set latency stop time
	if config.Debug {
		stop = time.Now()
	}

	// Format log
	fmt.Printf("%s %s %3d %s %7v \n",
		c.Method(),
		c.Path(),
		c.Response().StatusCode(),
		stop.Sub(start).Round(time.Millisecond),
		c.IP(),
	)

	if err != nil {
		fmt.Println(err.Error())
	}

	// End chain
	return nil
}

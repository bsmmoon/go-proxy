package seleniumwrapper

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bsmmoon/go-proxy/tool/logger"
	"github.com/tebeka/selenium"
)

// Options Options
type Options struct {
	SeleniumDriverPath string
	GeckoDriverPath    string
	Port               int
}

// Selenium selenium
func Selenium(options Options) {
	logger.INFO("Starting Selenium: %v", options)
	opts := []selenium.ServiceOption{
		// selenium.StartFrameBuffer(),        // Start an X frame buffer for the browser to run in. // xvfb not supported in MacOS?
		selenium.GeckoDriver(options.GeckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),                    // Output debug information to STDERR.
	}
	service, err := selenium.NewSeleniumService(options.SeleniumDriverPath, options.Port, opts...)
	if err != nil {
		logger.FATAL(err.Error()) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{
		"browserName": "firefox",
		"proxy": selenium.Proxy{
			Type: selenium.Manual,
			HTTP: "localhost:8089",
		},
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", options.Port))
	if err != nil {
		logger.FATAL(err.Error())
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("http://eu.httpbin.org/"); err != nil {
		panic(err)
	}

	logger.INFO("insert y value here: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	logger.INFO(input.Text())
}

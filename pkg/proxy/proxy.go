package proxy

import (
	"fmt"
	"os"

	"github.com/bsmmoon/go-proxy/tool/logger"
	"github.com/tebeka/selenium"
)

type ProxyOptions struct {
	SeleniumDriverPath string
	GeckoDriverPath    string
	Port               int
}

// HelloWorld hello world!
func HelloWorld(options ProxyOptions) {
	logger.INFO("HELLO WORLD")

	logger.INFO("\nPort: %v\nSelenium: %v\nGecko: %v", options.Port, options.GeckoDriverPath, options.SeleniumDriverPath)

	opts := []selenium.ServiceOption{
		// selenium.StartFrameBuffer(),        // Start an X frame buffer for the browser to run in. // xvfb not supported in MacOS?
		selenium.GeckoDriver(options.GeckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),                    // Output debug information to STDERR.
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(options.SeleniumDriverPath, options.Port, opts...)
	if err != nil {
		logger.FATAL(err.Error()) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", options.Port))
	if err != nil {
		logger.FATAL(err.Error())
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("http://play.golang.org/?simple=1"); err != nil {
		panic(err)
	}
}

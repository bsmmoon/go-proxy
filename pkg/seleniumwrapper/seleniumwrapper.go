package seleniumwrapper

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bsmmoon/go-proxy/tool/logger"
	"github.com/tebeka/selenium"
)

// Option Option
type Option struct {
	SeleniumDriverPath string
	GeckoDriverPath    string
	ChromeDriverPath   string
	Port               int
	ProxyPort          int
	Browser            string
}

// Selenium selenium
func Selenium(option Option) {
	logger.INFO("Starting Selenium: %v", option)
	service, err := makeFirefoxService(option)
	if err != nil {
		logger.FATAL(err.Error()) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	capabilities := selenium.Capabilities{
		"browserName": option.Browser,
		"proxy": selenium.Proxy{
			Type: selenium.Manual,
			HTTP: fmt.Sprintf("localhost:%v", option.ProxyPort),
		},
	}
	wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", option.Port))
	if err != nil {
		logger.FATAL(err.Error())
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("http://eu.httpbin.org/"); err != nil {
		panic(err)
	}

	// Quit upon user input
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

const FIREFOX = "firefox"

// CHROME make sure you add chromedriver to your PATH!!
const CHROME = "chrome"

func makeService(options Option) (*selenium.Service, error) {
	switch options.Browser {
	case FIREFOX:
		return makeFirefoxService(options)
	case CHROME:
		return makeChromeService(options)
	}
	return nil, nil
}

func makeFirefoxService(option Option) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(option.GeckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),                   // Output debug information to STDERR.
	}
	return selenium.NewSeleniumService(option.SeleniumDriverPath, option.Port, opts...)
}

func makeChromeService(option Option) (*selenium.Service, error) {
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(option.ChromeDriverPath), // Specify the path to GeckoDriver in order to use Chrome.
		selenium.Output(os.Stderr),                     // Output debug information to STDERR.
	}
	return selenium.NewChromeDriverService(option.SeleniumDriverPath, option.Port, opts...)
}

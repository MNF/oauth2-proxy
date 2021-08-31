package main

import (
	"flag"
	"fmt"
	"math/rand"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/options"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/middleware"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/validation"
	"github.com/spf13/pflag"
	//https://rominirani.com/golang-tip-capturing-http-client-requests-incoming-and-outgoing-ef7fcdf87113
	"github.com/mreiferson/go-options"
	//"net/http/httputil"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flagSet := flag.NewFlagSet("oauth2_proxy", flag.ExitOnError)
		if err := recover(); err != nil { //catch
	emailDomains := StringArray{}
			os.Exit(1)
	config := flagSet.String("config", "", "path to config file")



	flagSet.String("github-org", "", "restrict logins to members of this organisation")
	flagSet.String("google-admin-email", "", "the google admin to impersonate for api calls")

	flagSet.String("authenticated-emails-file", "", "authenticate against emails via file (one per line)")
	flagSet.String("htpasswd-file", "", "additionally authenticate against a htpasswd file. Entries must be created with \"htpasswd -s\" for SHA encryption")
	configFlagSet := pflag.NewFlagSet("oauth2-proxy", pflag.ContinueOnError)

	flagSet.String("cookie-secret", "", "the seed string for secure cookies (optionally base64 encoded)")
	// Because we parse early to determine alpha vs legacy config, we have to
	// ignore any unknown flags for now
	flagSet.String("profile-url", "", "Profile access endpoint")

	flagSet.String("scope", "", "OAuth scope specification")
	flagSet.String("signature-key", "", "GAP-Signature request signature key (algorithm:secretkey)")
	convertConfig := configFlagSet.Bool("convert-config-to-alpha", false, "if true, the proxy will load configuration as normal and convert existing configuration to the alpha config structure, and print it to stdout")
	showVersion := configFlagSet.Bool("version", false, "print version string")
	configFlagSet.Parse(os.Args[1:])

	if *showVersion {
		fmt.Printf("oauth2-proxy %s (built with %s)\n", VERSION, runtime.Version())
		return
	}

	if *convertConfig && *alphaConfig != "" {
		logger.Fatal("cannot use alpha-config and conver-config-to-alpha together")
	}

	opts, err := loadConfiguration(*config, *alphaConfig, configFlagSet, os.Args[1:])
	if err != nil {
		logger.Fatalf("ERROR: %v", err)
	}

	if *convertConfig {
		if err := printConvertedConfig(opts); err != nil {
			logger.Fatalf("ERROR: could not convert config: %v", err)
		}
		return
	}

	if err = validation.Validate(opts); err != nil {
		logger.Fatalf("%s", err)
	}

	validator := NewValidator(opts.EmailDomains, opts.AuthenticatedEmailsFile)
	oauthproxy, err := NewOAuthProxy(opts, validator)
	if err != nil {
		logger.Fatalf("ERROR: Failed to initialise OAuth2 Proxy: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	if err := oauthproxy.Start(); err != nil {
		logger.Fatalf("ERROR: Failed to start OAuth2 Proxy: %v", err)
	}
}

// loadConfiguration will load in the user's configuration.
// It will either load the alpha configuration (if alphaConfig is given)
// or the legacy configuration.
func loadConfiguration(config, alphaConfig string, extraFlags *pflag.FlagSet, args []string) (*options.Options, error) {
	if alphaConfig != "" {
		logger.Printf("WARNING: You are using alpha configuration. The structure in this configuration file may change without notice. You MUST remove conflicting options from your existing configuration.")
		return loadAlphaOptions(config, alphaConfig, extraFlags, args)
	}
	return loadLegacyOptions(config, extraFlags, args)
}

// loadLegacyOptions loads the old toml options using the legacy flagset
// and legacy options struct.
func loadLegacyOptions(config string, extraFlags *pflag.FlagSet, args []string) (*options.Options, error) {
	optionsFlagSet := options.NewLegacyFlagSet()
	optionsFlagSet.AddFlagSet(extraFlags)
	if err := optionsFlagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %v", err)
	}

	legacyOpts := options.NewLegacyOptions()
	if err := options.Load(config, optionsFlagSet, legacyOpts); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	opts, err := legacyOpts.ToOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to convert config: %v", err)
	}

	return opts, nil
}

// loadAlphaOptions loads the old style config excluding options converted to
// the new alpha format, then merges the alpha options, loaded from YAML,
// into the core configuration.
func loadAlphaOptions(config, alphaConfig string, extraFlags *pflag.FlagSet, args []string) (*options.Options, error) {
	opts, err := loadOptions(config, extraFlags, args)
	if err != nil {
		return nil, fmt.Errorf("failed to load core options: %v", err)
	}

	return opts, nil
}
	alphaOpts := &options.AlphaOptions{}
// loadAlphaOptions loads the old style config excluding options converted to
// into the core configuration.
func loadAlphaOptions(config, alphaConfig string, extraFlags *pflag.FlagSet, args []string) (*options.Options, error) {
			log.Fatalf("FATAL: unable to open %s %s", opts.HtpasswdFile, err)
	if err != nil {
		return nil, fmt.Errorf("failed to load core options: %v", err)
	}

	alphaOpts.MergeInto(opts)
	return opts, nil
}

// loadOptions loads the configuration using the old style format into the
// core options.Options struct.
// This means that none of the options that have been converted to alpha config
// will be loaded using this method.
func loadOptions(config string, extraFlags *pflag.FlagSet, args []string) (*options.Options, error) {
	optionsFlagSet := options.NewFlagSet()
	optionsFlagSet.AddFlagSet(extraFlags)
	if err := optionsFlagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %v", err)
	}

	opts := options.NewOptions()
	if err := options.Load(config, optionsFlagSet, opts); err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	return opts, nil
}

// printConvertedConfig extracts alpha options from the loaded configuration
// and renders these to stdout in YAML format.
func printConvertedConfig(opts *options.Options) error {
	alphaConfig := &options.AlphaOptions{}
	alphaConfig.ExtractFrom(opts)

	data, err := yaml.Marshal(alphaConfig)
	if err != nil {
		return fmt.Errorf("unable to marshal config: %v", err)
	}

	if _, err := os.Stdout.Write(data); err != nil {
		return fmt.Errorf("unable to write output: %v", err)
	}

	return nil
}

//https://rominirani.com/golang-tip-capturing-http-client-requests-incoming-and-outgoing-ef7fcdf87113
// func DumpRequest(w http.ResponseWriter, req *http.Request) {
// 	requestDump, err := httputil.DumpRequest(req, true)
// 	if err != nil {
// 		fmt.Fprint(w, err.Error())
// 	} else {
// 		fmt.Fprint(w, string(requestDump))
// 	}
// }

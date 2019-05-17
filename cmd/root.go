// Copyright © 2019 Vyacheslav Akhmetov <iam@itaddict.ru>

package cmd

import (
	"fmt"
	"github.com/akme/get-closer/loaders"
	"github.com/akme/get-closer/protocols"

	"net"
	"os"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/cheggaaa/pb.v1"
	"math"
	"sort"
)

var (
	cfgFile     string
	hostsFile   string
	timeout     uint // time.Duration( 10 * time.Second)
	concurrency uint
	limit       uint
	dnsResolver string
	dnsWarmUp   bool
	progressBar bool
	count       int
	verbose     bool
	sslEnable   bool
	tcpPort     int
	delay       int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "get-closer",
	Short: "Find out closest hosts in terms of network latency.",
	Long:  `Find out closest hosts in terms of network latency.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.get-closer.yaml)")
	rootCmd.PersistentFlags().StringVarP(&hostsFile, "from-file", "f", "", "path to file with hosts to check")
	rootCmd.MarkPersistentFlagRequired("from-file")
	rootCmd.PersistentFlags().UintVarP(&timeout, "timeout", "t", 60, "timeout for request")
	rootCmd.PersistentFlags().UintVarP(&concurrency, "concurrency", "c", 1, "concurrency")
	rootCmd.PersistentFlags().UintVarP(&limit, "limit", "l", 0, "number of hosts to return")
	rootCmd.PersistentFlags().StringVarP(&dnsResolver, "dns-resolver", "r", "", "use custom DNS resolver")
	rootCmd.PersistentFlags().BoolVarP(&dnsWarmUp, "dns-warm-up", "w", true, "warm up DNS cache before request")
	rootCmd.PersistentFlags().IntVarP(&count, "count", "", 1, "number of tests per host")
	rootCmd.PersistentFlags().BoolVarP(&progressBar, "progress-bar", "b", true, "show progress bar")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	rootCmd.PersistentFlags().IntVarP(&delay, "delay", "d", 3, "set delay between checks")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".get-closer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".get-closer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.BindPFlag("from-file", rootCmd.PersistentFlags().Lookup("from-file"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("concurrency", rootCmd.Flags().Lookup("concurrency"))
	viper.BindPFlag("limit", rootCmd.PersistentFlags().Lookup("limit"))
	viper.BindPFlag("dns-server", rootCmd.PersistentFlags().Lookup("dns-server"))
	viper.BindPFlag("dns-warm-up", rootCmd.PersistentFlags().Lookup("dns-warm-up"))

}

//durationSlice represents collection of measurements
type durationSlice []time.Duration

//Measurement struct for storing results
type Measurement struct {
	Host     string
	Duration durationSlice
}

//Measurements list
type Measurements []Measurement

func (m Measurements) Len() int           { return len(m) }
func (m Measurements) Less(i, j int) bool { return m[i].Duration.Avg() > m[j].Duration.Avg() }
func (m Measurements) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func startMeasurements(ccmd *cobra.Command, args []string) {
	hostsList, err := loaders.LoadHosts(hostsFile)
	if err != nil {
		fmt.Println(err)
	}
	var MeasurementsList Measurements

	for _, v := range hostsList {
		MeasurementsList = append(MeasurementsList, Measurement{Host: v})
	}

	listSize := len(MeasurementsList)
	bar := pb.New(listSize)

	if progressBar {
		bar = pb.StartNew(listSize).Prefix("Checking hosts")
		bar.ShowTimeLeft = false
	}
	if len(dnsResolver) != 0 {
		protocols.UseCustomDNS([]string{dnsResolver})
	}

	for k, v := range MeasurementsList {
		if verbose {
			fmt.Printf("Measuring [%d/%d] %s\n", k+1, listSize, v.Host)
		}
		if progressBar {
			bar.Increment()
		}
		var duration time.Duration

		if err != nil {
			fmt.Println("url.Parse:", err)
		}
		if dnsWarmUp {
			if verbose {
				fmt.Printf("\tWarming up DNS for %s\n", v.Host)
			}
			before := time.Now()
			addrs, err := net.LookupHost(v.Host)
			diff := time.Since(before)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if verbose {
				fmt.Printf("\tResolving %s to %s took %s\n", v.Host, addrs, diff)
			}
		}
		for i := 0; i < count; i++ {
			switch ccmd.Name() {
			case "http":
				duration = protocols.HTTPPing(v.Host, sslEnable)
			case "dns":
				duration = protocols.DNSPing(v.Host)
			case "icmp":
				duration = protocols.ICMPPing(v.Host)
			case "tcp":
				duration = protocols.TCPPing(v.Host, tcpPort)
			default:
			}

			if err != nil {
				fmt.Println(err)
			}
			if verbose {
				fmt.Printf("\t%s measurement for %s took %s\n", ccmd.Name(), v.Host, duration)
			}
			MeasurementsList[k].Duration = append(MeasurementsList[k].Duration, duration)

			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
	if progressBar {
		bar.FinishPrint("Closest hosts:")
	}
	if verbose {
		fmt.Println("Closest hosts:")
	}
	if limit == 0 {
		limit = uint(listSize)
	}
	sortedList := sortByDuration(MeasurementsList)
	for _, v := range sortedList[:limit] {
		fmt.Printf("\t%s %s\n", v.Host, v.Duration)
	}

}

func (p durationSlice) Len() int           { return len(p) }
func (p durationSlice) Less(i, j int) bool { return int64(p[i]) < int64(p[j]) }
func (p durationSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p durationSlice) Min() time.Duration {
	sort.Sort(p)
	return p[0]
}
func (p durationSlice) Max() time.Duration {
	sort.Sort(p)
	return p[p.Len()-1]
}
func (p durationSlice) Avg() time.Duration {
	var avg int64
	for i := 0; i < p.Len(); i++ {
		avg += int64(p[i])
	}
	if p.Len() == 0 {
		return time.Duration(0)
	}
	return time.Duration(avg / int64(p.Len()))
}
func (p durationSlice) Std() time.Duration {
	sqdifs := make(durationSlice, p.Len(), p.Len())
	avg := p.Avg()
	var avgsqdif int64
	for i := 0; i < p.Len(); i++ {
		sqdif := p[i] - avg
		sqdifs[i] = sqdif * sqdif
		avgsqdif += int64(sqdifs[i])
	}
	avgsqdif /= int64(sqdifs.Len())
	return time.Duration(math.Sqrt(float64(avgsqdif)))
}

func sortByDuration(measurements Measurements) Measurements {
	pl := make(Measurements, len(measurements))
	i := 0
	for _, v := range measurements {
		pl[i] = Measurement{v.Host, durationSlice{v.Duration.Avg()}}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

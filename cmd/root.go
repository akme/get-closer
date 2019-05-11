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
	rootCmd.PersistentFlags().StringVarP(&hostsFile, "from-file", "f", "", "Path to file with hosts to check")
	rootCmd.MarkPersistentFlagRequired("from-file")
	rootCmd.PersistentFlags().UintVarP(&timeout, "timeout", "t", 60, "Timeout for request")
	rootCmd.PersistentFlags().UintVarP(&concurrency, "concurrency", "c", 1, "Concurrency")
	rootCmd.PersistentFlags().UintVarP(&limit, "limit", "l", 0, "number of hosts to return")
	rootCmd.PersistentFlags().StringVarP(&dnsResolver, "dns-server", "", "", "use custom DNS resolver")
	rootCmd.PersistentFlags().BoolVarP(&dnsWarmUp, "dns-warm-up", "w", true, "warm up DNS cache before request")
	rootCmd.PersistentFlags().IntVarP(&count, "count", "", 1, "number of tests per host")
	rootCmd.PersistentFlags().BoolVarP(&progressBar, "progress-bar", "b", true, "show progress bar")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
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

//DurationSlice represents collection of measurements
type DurationSlice []time.Duration

//Measurement struct for storing results
type Measurement struct {
	Host     string
	Duration DurationSlice
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
	bar := pb.New(len(MeasurementsList))

	if progressBar {
		bar = pb.StartNew(len(MeasurementsList)).Prefix("Checking hosts")
		bar.ShowTimeLeft = false
	}
	if len(dnsResolver) != 0 {
		protocols.UseCustomDNS([]string{dnsResolver})
	}

	for k, v := range MeasurementsList {
		//fmt.Println("Measuring ", v.Host)
		if progressBar {
			bar.Increment()
		}
		var duration time.Duration

		if err != nil {
			fmt.Println("url.Parse:", err)
		}
		//fmt.Println("+++++", host.Hostname())
		if dnsWarmUp {
			//fmt.Println("Warming up DNS")
			//before := time.Now()
			_, err := net.LookupHost(v.Host)
			//after := time.Now()
			if err != nil {
				fmt.Println(err)
				continue
			}
			// diff := after.Sub(before)
			//fmt.Println(addrs, diff)
		}
		for i := 0; i < count; i++ {
			switch ccmd.Name() {
			case "http":
				duration = protocols.HTTPPing(v.Host)
			case "dns":
				duration = protocols.DNSPing(v.Host)
			case "icmp":
				duration = protocols.ICMPPing(v.Host)
			case "tcp":
				duration = protocols.TCPPing(v.Host)
			default:
			}

			if err != nil {
				fmt.Println(err)
			}
			MeasurementsList[k].Duration = append(MeasurementsList[k].Duration, duration)
		}
	}
	if progressBar {
		bar.FinishPrint("Closest hosts:")
	}
	for _, v := range sortByDuration(MeasurementsList) {
		fmt.Println(v.Host, v.Duration)
	}

}

// NOTE: This implements the sortable interface
func (p DurationSlice) Len() int           { return len(p) }
func (p DurationSlice) Less(i, j int) bool { return int64(p[i]) < int64(p[j]) }
func (p DurationSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// NOTE: Wasteful Convenience Functions
func (p DurationSlice) Min() time.Duration {
	sort.Sort(p)
	return p[0]
}
func (p DurationSlice) Max() time.Duration {
	sort.Sort(p)
	return p[p.Len()-1]
}
func (p DurationSlice) Avg() time.Duration {
	var avg int64
	for i := 0; i < p.Len(); i++ {
		avg += int64(p[i])
	}
	return time.Duration(avg / int64(p.Len()))
}
func (p DurationSlice) Std() time.Duration {
	sqdifs := make(DurationSlice, p.Len(), p.Len())
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
		pl[i] = Measurement{v.Host, DurationSlice{v.Duration.Avg()}}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

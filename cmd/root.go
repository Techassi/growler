package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/Techassi/growler/internal/queue"
	"github.com/Techassi/growler/internal/crawl"
	"github.com/Techassi/growler/internal/events"
	"github.com/Techassi/growler/internal/workerpool"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	urlFlag string
	queueFlag int
	workersFlag int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   	"growler",
	Short: 	"growler is a web crawler written in Go",
	Long: 	`growler crawls the web propagating from the given url (seed)
in a parallized manner with scalable queue and workers.`,
	Run: func(cmd *cobra.Command, args []string) {
		q, err := queue.NewQueue(queueFlag)
		q.URLJob(urlFlag)
		if err != nil {
			panic(err)
		}

		p, pool_err := workerpool.NewWorkerPool(workersFlag, q, crawl.Crawl)
		if pool_err != nil {
			panic(err)
		}

		err = p.On("worker:finish", events.WorkerProcess)
		if err != nil {
			panic(err)
		}

		p.Start()
	},
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
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&urlFlag, "url", "", "The URL used as an entry point")
	// workaround (see https://github.com/spf13/cobra/issues/921)
	cobra.MarkFlagRequired(rootCmd.PersistentFlags(), "url")
	rootCmd.PersistentFlags().IntVar(&queueFlag, "queue", 1000, "The max amount of items in queue at one time")
	rootCmd.PersistentFlags().IntVar(&workersFlag, "workers", 10, "The max amount of cuncurrent workers running")
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

		// Search config in home directory with name ".growler" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".growler")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

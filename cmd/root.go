package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/Techassi/growler/internal/queue"
	"github.com/Techassi/growler/internal/crawl"
	"github.com/Techassi/growler/internal/workerpool"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var urlFlag string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   	"growler",
	Short: 	"A brief description of your application",
	Long: 	`A longer description that spans multiple lines and likely contains
			examples and usage of using your application. For example:`,
	Run: func(cmd *cobra.Command, args []string) {
		q, err := queue.NewQueue(1000)
		q.URLJob(urlFlag)
		if err != nil {
			panic(err)
		}

		p := workerpool.NewWorkerPool(10, q, crawl.Crawl)
		// err = p.On("init", workerInit)
		// if err != nil {
		// 	panic(err)
		// }

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

func workerInit(pool *workerpool.WorkerPool) {
	fmt.Println(pool.Queue)
}

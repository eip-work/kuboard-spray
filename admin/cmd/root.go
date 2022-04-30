package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kuboard-spray-admin",
	Short: "Kuboard 命令行控制",
	Long: `Kuboard 命令行控制程序，用于执行重置 admin 的密码等操作。
kuboard-spray-admin reset-password	
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	logrus.SetFormatter(new(KuboardLogFormatter))
	value := os.Getenv("KUBOARD_SPRAY_ADMIN_LOGRUS_LEVEL")
	if value == "" {
		value = "info"
	}
	level, err := logrus.ParseLevel(value)
	if err == nil {
		fmt.Println("设置日志级别为 " + value)
		logrus.SetLevel(level)
	} else {
		fmt.Println("请检查 KUBOARD_SPRAY_ADMIN_LOGRUS_LEVEL 的值，可选的有 panic / fatal / error / warn / info / debug / trace ，当前为： " + value)
		logrus.SetLevel(logrus.InfoLevel)
	}

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(resetPasswordCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AutomaticEnv() // read in environment variables that match

}

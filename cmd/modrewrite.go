// Copyright © 2018 Steve Huff <steve.huff@asics.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/asicsdigital/safetyrail/util"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// modrewriteCmd represents the modrewrite command
var modrewriteCmd = &cobra.Command{
	Use:   "modrewrite",
	Short: "Reformat S3 website rewrite rules for mod_rewrite",
	Long:  `Reformat S3 website rewrite rules for mod_rewrite.`,
	Run: func(cmd *cobra.Command, args []string) {
		jww.INFO.Println("modrewrite called")

		redirects := []string{}

		var r util.XMLRoutingRules
		infile := viper.GetString("infile")

		if len(infile) > 0 {
			err := r.Parse(infile)

			if err != nil {
				jww.ERROR.Panic(err)
			}

			for i := range r.RoutingRules {
				rule := &r.RoutingRules[i]
				redirects = append(redirects, rule.ModRewrite())
			}
		}

		bucket := viper.GetString("bucket")

		if len(bucket) > 0 {
			redirects = append(redirects, util.GetS3Redirects(bucket, util.AwsInit())...)
		}

		for i := range redirects {
			fmt.Println(redirects[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(modrewriteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modrewriteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modrewriteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	viper.SetDefault("infile", "")
	modrewriteCmd.Flags().String("infile", "", "Path to XML rewrite rules")
	viper.BindPFlag("infile", modrewriteCmd.Flags().Lookup("infile"))

	viper.SetDefault("bucket", "")
	modrewriteCmd.Flags().String("bucket", "", "S3 bucket for redirect objects")
	viper.BindPFlag("bucket", modrewriteCmd.Flags().Lookup("bucket"))
}

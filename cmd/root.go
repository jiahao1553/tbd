package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gwenwindflower/tbd/internal"
	"github.com/gwenwindflower/tbd/sourcerer"
	"github.com/spf13/cobra"
)

type Elapsed struct {
	DbStart           time.Time
	ProcessingStart   time.Time
	DbElapsed         float64
	ProcessingElapsed float64
}

var (
	greenBold = color.New(color.FgMagenta).Add(color.Bold).SprintFunc()
	rootCmd   = &cobra.Command{
		Use:   "tbd",
		Short: "🏁 A sweet and speedy code generator for dbt projects. 🏎️✨",
		Long: fmt.Sprintf(`🏁  %s 🏎️✨
tbd uses your database schema to generate YAML configs and
SQL staging models, including tests and docs via LLM.

It's the easy button for starting a dbt project.`, greenBold("A sweet and speedy code generator for dbt projects.")),
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ps, err := internal.FetchDbtProfiles()
			if err != nil {
				log.Fatalf("Error fetching dbt profiles: %v\n", err)
			}
			fr, err := internal.Forms(ps)
			if err != nil {
				log.Fatalf("Error running form: %v\n", err)
			}
			if !fr.Confirm {
				log.Fatal("⛔ User cancelled.")
			}
			cd := internal.SetConnectionDetails(fr, ps)

			e := Elapsed{}
			e.DbStart = time.Now()

			bd := fr.BuildDir
			err = internal.PrepBuildDir(bd)
			if err != nil {
				log.Fatalf("Error preparing build directory: %v\n", err)
			}
			dbc, err := sourcerer.GetConn(cd)
			if err != nil {
				log.Fatalf("Error getting database connection: %v\n", err)
			}
			err = dbc.ConnectToDb(ctx)
			if err != nil {
				log.Fatalf("Error connecting to database: %v\n", err)
			}
			fmt.Println("Connected to database")
			ts, err := dbc.GetSourceTables(ctx)
			if err != nil {
				log.Fatalf("Error getting sources: %v\n", err)
			}
			fmt.Println("Got source tables")
			fmt.Println("Putting columns on tables...")
			err = sourcerer.PutColumnsOnTables(ctx, ts, dbc)
			if err != nil {
				log.Fatalf("Error putting columns on tables: %v\n", err)
			}

			e.DbElapsed = time.Since(e.DbStart).Seconds()
			// End of database interaction, start of processing
			e.ProcessingStart = time.Now()

			if fr.GenerateDescriptions {
				llm, err := internal.GetLlm(fr)
				if err != nil {
					// Using Printf instead of log.Fatalf since the program doesn't
					// need to totally fail if the API provider can't be fetched
					fmt.Printf("Error getting API provider: %v\n", err)
				}
				fmt.Println("Generating descriptions and tests...")
				internal.InferColumnFields(llm, ts)
				if err != nil {
					// Using Printf instead of log.Fatalf since the program
					// doesn't need to totally fail if there's an error in the column field inference
					fmt.Printf("Error inferring column fields: %v\n", err)
				}
			}
			fmt.Println("Writing files...")
			if fr.CreateProfile {
				internal.WriteProfile(cd, bd)
			}
			if fr.ScaffoldProject {
				// include DbtProfileName because currently, profile is pointing to connection type
				s, err := internal.WriteScaffoldProject(cd, bd, fr.ProjectName, fr.DbtProfileName)
				if err != nil {
					log.Fatalf("Error scaffolding project: %v\n", err)
				}
				bd = s
			}
			// include database and schema because both are missing in the _sources.yml, which are required
			err = internal.WriteFiles(ts, bd, fr.Prefix, &fr.Database, &fr.Schema)
			if err != nil {
				log.Fatalf("Error writing files: %v\n", err)
			}
			e.ProcessingElapsed = time.Since(e.ProcessingStart).Seconds()
			pinkUnderline := color.New(color.FgMagenta).Add(color.Bold, color.Underline).SprintFunc()
			fmt.Printf("\n🏁 Done in %.1fs fetching data and %.1fs writing files!\nYour YAML and SQL files are in the %s directory.", e.DbElapsed, e.ProcessingElapsed, pinkUnderline(fr.BuildDir))
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

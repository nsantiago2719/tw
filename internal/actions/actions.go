package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v3"

	"github.com/nsantiago2719/tw/internal/app"
	"github.com/nsantiago2719/tw/internal/terraform"
	"github.com/nsantiago2719/tw/pkg/utils"
)

// TODO: run terraform apply along with the var files passed if exist
func RunTerraform(ctx context.Context, cmd *cli.Command, cfg string) error {
	resourceName := cmd.StringArg("resource-name")
	if resourceName == "" {
		return fmt.Errorf("resource-name cannot be empty")
	}

	var resources []app.Resource
	config, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)
	if err != nil {
		return err
	}

	resourcePath, args := utils.GetDetails(resourceName, resources)
	if resourcePath == "" {
		return fmt.Errorf("the resource registered has an empty path")
	}

	if err := utils.RunInit(ctx, resourcePath); err != nil {
		return err
	}

	execApply := terraform.InitCmd("apply")
	err = execApply.CreateCmd(resourcePath, args...)
	if err != nil {
		return err
	}

	if cmd.Bool("auto-approve") {
		execApply.AddArg("-auto-approve")
	}

	if cmd.Bool("dry-run") {
		execApply.AddArg("-dry-run")
	}

	// Get output channel and stdin channels from command execution
	// The stdinRequestChan and stdinInputChan are used by handleCommandIO to manage interactive input
	execApplyOutput, stdinRequestChan, stdinInputChan, err := execApply.Exec(ctx)
	if err != nil {
		return err
	}

	// Handle both output and potential input requests
	// This function uses all three channels to manage command IO
	utils.HandleCommandIO(execApplyOutput, stdinRequestChan, stdinInputChan)

	return nil
}

// PlanTerraform runs terraform plan against the resource specified
func PlanTerraform(ctx context.Context, cmd *cli.Command, cfg string) error {
	resourceName := cmd.StringArg("resource-name")

	if resourceName == "" {
		return fmt.Errorf("resouce-name cannot be empty")
	}

	var resources []app.Resource
	config, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)
	if err != nil {
		return err
	}

	resourcePath, varFiles := utils.GetDetails(resourceName, resources)
	if resourcePath == "" {
		return fmt.Errorf("the resource registered has an empty path")
	}

	if err := utils.RunInit(ctx, resourcePath); err != nil {
		return err
	}

	execPlan := terraform.InitCmd("plan")
	err = execPlan.CreateCmd(resourcePath, varFiles...)
	if err != nil {
		return err
	}

	// Get output channel and stdin channels from command execution
	// The stdinRequestChan and stdinInputChan are used by handleCommandIO to manage interactive input
	execPlanOutput, stdinRequestChan, stdinInputChan, err := execPlan.Exec(ctx)
	if err != nil {
		return err
	}

	// Handle both output and potential input requests
	// This function uses all three channels to manage command IO
	utils.HandleCommandIO(execPlanOutput, stdinRequestChan, stdinInputChan)

	return nil
}

// TODO: set path as the current path if path flag is `.`
func RegisterResource(_ context.Context, cmd *cli.Command, cfg string) error {
	if cmd.String("name") == "" {
		return errors.New("Name must not be empty")
	}

	if cmd.String("path") == "" {
		return errors.New("Path must not be empty")
	}

	rs := app.Resource{
		Name:     cmd.String("name"),
		Path:     cmd.String("path"),
		VarFiles: cmd.StringSlice("var-files"),
	}

	file, err := os.OpenFile(cfg, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil
	}

	defer file.Close()

	resources := []app.Resource{}

	config, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)
	if err != nil {
		return nil
	}

	for resource := range resources {
		nameIndex := resource + 1
		baseName := cmd.String("name")
		if resources[resource].Name == rs.Name {
			rs.Name = fmt.Sprintf("%s-%d", baseName, nameIndex)
		}
	}

	resources = append(resources, rs)

	marshalResources, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = file.Write(marshalResources)
	if err != nil {
		return err
	}

	return nil
}

// Resources lists all registered resources in a table format
func Resources(_ context.Context, _ *cli.Command, cfg string) error {
	config, err := os.ReadFile(cfg)
	resources := []app.Resource{}
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &resources)

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Path", "Var Files")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, resource := range resources {
		tbl.AddRow(resource.Name, resource.Path, resource.VarFiles)
	}

	tbl.Print()
	return nil
}

// Init create a config.json file if the file does not exist
// else it would do nothing
func Init(_ context.Context, _ *cli.Command, cfg string) error {
	if _, err := os.Stat(cfg); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(cfg, []byte("[]"), 0o755)
		if err != nil {
			return err
		}
	}
	fmt.Println("tw initialized: created config.json file")
	return nil
}

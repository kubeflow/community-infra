package main

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/kubeflow/community-infra/pkg/api/v1alpha1"
	"github.com/kubeflow/community-infra/pkg/controllers"
	"github.com/kubeflow/community-infra/pkg/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/api/cloudresourcemanager/v1"
	crmV2 "google.golang.org/api/cloudresourcemanager/v2"
	"google.golang.org/api/servicemanagement/v1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplyOptions struct {
	Input string
}

type GenerateOptions struct {
	Output string
}

var (
	opts  = ApplyOptions{}
	gOpts = GenerateOptions{}

	rootCmd = &cobra.Command{}

	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply the specified file",
		Run: func(cmd *cobra.Command, args []string) {
			apply()
		},
	}

	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a bulk deleter spec.",
		Run: func(cmd *cobra.Command, args []string) {
			generate()
		},
	}

	log logr.Logger
)

func init() {
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(generateCmd)

	applyCmd.Flags().StringVarP(&opts.Input, "file", "f", "", "The input file to process")
	applyCmd.MarkFlagRequired("file")

	generateCmd.Flags().StringVarP(&gOpts.Output, "output", "o", "", "The output file to write")
	generateCmd.MarkFlagRequired("output")
}

func initLogger() {
	// TODO(jlewi): Make the verbosity level configurable.

	// Start with a production logger config.
	config := zap.NewProductionConfig()

	// TODO(jlewi): In development mode we should use the console encoder as opposed to json formatted logs.

	// Increment the logging level.
	// TODO(jlewi): Make this a flag.
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	zapLog, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Could not create zap instance (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	zap.ReplaceGlobals(zapLog)
}

func apply() {
	initLogger()

	if opts.Input == "" {
		log.Error(fmt.Errorf("No input file supplied"), "Input file must be supplied")
		return
	}

	b, err := ioutil.ReadFile(opts.Input)

	if err != nil {
		log.Error(err, "Could not read file", "file", opts.Input)
		return
	}

	kind, err := utils.GetObjectKind(opts.Input)

	if err != nil {
		log.Error(err, "Could not read kind from file", "file", opts.Input)
		return
	}

	log.Info("Got Object Kind", "kind", kind)


	ctx := context.Background()
	crmService, err := cloudresourcemanager.NewService(ctx)

	if err != nil {
		log.Error(err, "Could not initialize cloud resource manager client")
		return
	}

	crmServiceV2, err := crmV2.NewService(ctx)

	if err != nil {
		log.Error(err, "Could not initialize cloud resource manager v2 client")
		return
	}

	smService, err := servicemanagement.NewService(ctx)

	if err != nil {
		log.Error(err, "Could not create service management endpoint")
		return
	}

	orgHelper := &controllers.OrgHelper{
		Crm: crmService,
		CrmV2: crmServiceV2,
		Sm:    smService,
	}


	var applyErr error
	switch kind {
	case v1alpha1.BulkProjectDeleteKind:
		applier := &v1alpha1.BulkProjectDelete{}

		err = yaml.Unmarshal(b, applier)

		if err != nil {
			log.Error(err, "Could not unmarshal applier for kind", "kind", kind)
			return
		}

		applyErr = controllers.ApplyBulkDelete(applier, orgHelper)
		break
	case v1alpha1.BulkProjectMoveKind:
		applier := &v1alpha1.BulkProjectMove{}

		err = yaml.Unmarshal(b, applier)

		if err != nil {
			log.Error(err, "Could not unmarshal applier for kind", "kind", kind)
			return
		}

		applyErr = controllers.ApplyBulkMove(applier, orgHelper)
		break
		default:
		log.Error(fmt.Errorf("Unrecognized error"), "Unrecognized kind", "kind", kind)
	}

	if applyErr != nil {
		log.Error(applyErr, "Apply failed")
		return
	}
}

func generate() {
	initLogger()

	bulkDelete := v1alpha1.BulkProjectDelete{
		TypeMeta: metav1.TypeMeta{
			Kind: "BulkProjectDelete",
		},
		Spec: v1alpha1.BulkProjectDeleteSpec{
			Projects: []string{},
		},
	}
	for i := 4; i < 59; i = i + 1 {
		bulkDelete.Spec.Projects = append(bulkDelete.Spec.Projects, fmt.Sprintf("kf-load-test-project%v", i))
	}

	b, err := yaml.Marshal(bulkDelete)

	if err != nil {
		log.Error(err, "Could not marshal the bulk delete code")
		return
	}

	err = ioutil.WriteFile(gOpts.Output, b, 0777)

	if err != nil {
		log.Error(err, "Could not write file", "file", gOpts.Output)
	}
	return
}

func main() {
	rootCmd.Execute()
}

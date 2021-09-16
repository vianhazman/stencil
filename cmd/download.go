package cmd

import (
	"context"
	"errors"
	"os"

	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// DownloadCmd creates a new cobra command for download descriptor
func DownloadCmd() *cobra.Command {
	var host, filePath string
	var req stencilv1.DownloadDescriptorRequest
	var latest bool

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download filedescriptorset file",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if latest {
				if req.Version != "" {
					return errors.New("version and latest cannot be specified in the same query")
				}
				req.Channel = "latest"
			} else if req.Version == "" {
				return errors.New("need either version tag or latest flag")
			}

			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			res, err := client.DownloadDescriptor(context.Background(), &req)
			if err != nil {
				return err
			}
			err = os.WriteFile(filePath, res.Data, 0666)
			return err
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "stencil host address eg: localhost:8000")
	cmd.MarkFlagRequired("host")
	cmd.Flags().StringVar(&req.Namespace, "namespace", "", "provide namespace/group or entity name")
	cmd.MarkFlagRequired("namespace")
	cmd.Flags().StringVar(&req.Name, "name", "", "provide proto repo name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&req.Version, "version", "", "provide semantic version compatible value")
	cmd.Flags().BoolVar(&latest, "latest", false, "search snapshots marked with latest")

	cmd.Flags().StringVar(&filePath, "output", "", "write to file")
	cmd.Flags().StringSliceVar(&req.Fullnames, "fullnames", []string{}, "provide fully qualified proto full names. You can provide multiple names separated by \",\" Eg: google.protobuf.FileDescriptorProto,google.protobuf.FileDescriptorSet")
	return cmd
}

package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	stencilv1 "github.com/odpf/stencil/server/odpf/stencil/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// Snapshot creates a new cobra command to manage snapshot
func Snapshot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "list, update snapshot details",
		Annotations: map[string]string{
			"group:core": "true",
		},
	}
	cmd.PersistentFlags().String("host", "", "stencil host address eg: localhost:8000")
	cmd.MarkPersistentFlagRequired("host")
	cmd.AddCommand(listCmd())
	cmd.AddCommand(promoteCmd())
	cmd.AddCommand(printCmd())
	return cmd
}

func listCmd() *cobra.Command {
	var req stencilv1.ListSnapshotsRequest
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list snapshots with optional filters",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			host, _ := cmd.Flags().GetString("host")
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			res, err := client.ListSnapshots(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}
			data, _ := protojson.MarshalOptions{EmitUnpopulated: true, Multiline: true, Indent: "  "}.Marshal(res)
			fmt.Println(string(data))
			return nil
		},
	}
	cmd.Flags().StringVar(&req.Namespace, "namespace", "", "provide namespace/group or entity name")
	cmd.Flags().StringVar(&req.Name, "name", "", "provide proto repo name")
	cmd.Flags().StringVar(&req.Version, "version", "", "provide semantic version compatible value")
	cmd.Flags().BoolVar(&req.Latest, "latest", false, "mark as latest version")
	return cmd
}

func promoteCmd() *cobra.Command {
	var req stencilv1.PromoteSnapshotRequest
	cmd := &cobra.Command{
		Use:   "promote",
		Short: "promote specified snapshot to latest",
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"group:core": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			host, _ := cmd.Flags().GetString("host")
			conn, err := grpc.Dial(host, grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()
			client := stencilv1.NewStencilServiceClient(conn)
			res, err := client.PromoteSnapshot(context.Background(), &req)
			if err != nil {
				errStatus := status.Convert(err)
				return errors.New(errStatus.Message())
			}
			data, _ := protojson.MarshalOptions{EmitUnpopulated: true, Multiline: true, Indent: "  "}.Marshal(res)
			fmt.Println(string(data))
			return nil
		},
	}
	cmd.Flags().Int64Var(&req.Id, "id", 0, "snapshot id")
	cmd.MarkFlagRequired("id")
	return cmd
}

// printCmd creates a new cobra command for upload
func printCmd() *cobra.Command {
	var host string
	var latest bool
	var req stencilv1.DownloadDescriptorRequest
	cmd := &cobra.Command{
		Use:  "print",
		Args: cobra.NoArgs,
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

			fds := &descriptorpb.FileDescriptorSet{}
			if err := proto.Unmarshal(res.GetData(), fds); err != nil {
				return fmt.Errorf("descriptor set file is not valid. %w", err)
			}

			fd, err := desc.CreateFileDescriptorFromSet(fds)
			if err != nil {
				return err
			}

			return (&protoprint.Printer{}).PrintProtoFile(fd, os.Stdout)
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
	return cmd
}
